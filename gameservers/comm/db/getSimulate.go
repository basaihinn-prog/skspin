package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"unsafe"

	"github.com/golang/snappy"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSimulate() (docs []Playitem, err error) {
	docStr := "data"
	filePathStr := fmt.Sprintf("%v/%v.cache", docStr, GetMongoDbName())
	//在根目录下的data中获取文件中的数据
	// 打开文件
	if !fileExists(filePathStr) {
		{ //从mongo中查询数据
			coll := Collection("simulate")
			filter := D("selected", true)
			//cur, _ := coll.Find(context.TODO(), filter, options.Find().SetProjection(D("bucketid", 1, "type", 1, "BucketMix", 1, "BucketStable", 1, "BucketHeartBeat", 1, "times", 1)))
			cur, _ := coll.Find(context.TODO(), filter, options.Find().SetProjection(D("bucketid", 1, "type", 1, "BucketHeartBeat", 1, "BucketWave", 1, "BucketGov", 1, "BucketMix", 1, "BucketStable", 1, "BucketHighAward", 1, "BucketSuperHighAward", 1, "times", 1)))
			err = cur.All(context.TODO(), &docs)
			if err != nil {
				slog.Error("cur.All(context.TODO(), &docs)::err: ", err)
				return docs, err
			}
			if len(docs) == 0 {
				coll = Collection("rawSpinData")
				cur, _ := coll.Find(context.TODO(), filter, options.Find().SetProjection(D("bucketid", 1, "type", 1, "BucketHeartBeat", 1, "BucketWave", 1, "BucketGov", 1, "BucketMix", 1, "BucketStable", 1, "BucketHighAward", 1, "BucketSuperHighAward", 1, "times", 1)))
				err = cur.All(context.TODO(), &docs)
				if err != nil {
					slog.Error("cur.All(context.TODO(), &docs)::err: ", err)
					return docs, err
				}
			}
			//获取游戏id，将查出的数据存到根目录下的data
			binaryData := encodeWithBinary(docs)
			//创建文件，并写入
			// 创建文件夹（如果不存在）
			err = os.MkdirAll(docStr, os.ModePerm)
			if err != nil {
				slog.Error("Error creating directory:", err)
				return docs, err
			}
			// 创建文件
			var fileIn *os.File
			fileIn, createErr := os.Create(filePathStr)
			if createErr != nil {
				slog.Error("Error creating file:", createErr)
				return docs, createErr
			}
			defer fileIn.Close()
			_, err = fileIn.Write(binaryData)
			if err != nil {
				slog.Error("Error writing to file:", err)
				return docs, err
			}
		}
	} else {
		docs, err = readBinaryData(filePathStr)
	}
	return docs, nil
}

func GetSimulateByTableNames(tableName string) (docs []Playitem, err error) {
	docStr := "data"
	filePathStr := fmt.Sprintf("%v/%v.cache", docStr, GetMongoDbName())

	//在根目录下的data中获取文件中的数据
	if !fileExists(filePathStr) {
		{ //从mongo中查询数据
			coll := Collection(tableName)
			filter := D("selected", true)
			//cur, _ := coll.Find(context.TODO(), filter, options.Find().SetProjection(D("type", 1, "data.gametype", 1, "bucketid", 1, "BucketMix", 1, "BucketStable", 1, "BucketHeartBeat", 1)))
			cur, _ := coll.Find(context.TODO(), filter, options.Find().SetProjection(D("type", 1, "data.gametype", 1, "bucketid", 1, "BucketHeartBeat", 1, "BucketWave", 1, "BucketGov", 1, "BucketMix", 1, "BucketStable", 1, "BucketHighAward", 1, "BucketSuperHighAward", 1)))
			err = cur.All(context.TODO(), &docs)
			if err != nil {
				slog.Error("cur.All(context.TODO(), &docs)::err: ", err)
				return docs, err
			}
			//获取游戏id，将查出的数据存到根目录下的data
			binaryData := encodeWithBinary(docs)
			//创建文件，并写入
			// 创建文件夹（如果不存在）
			err = os.MkdirAll(docStr, os.ModePerm)
			if err != nil {
				slog.Error("Error creating directory:", err)
				return docs, err
			}
			// 创建文件
			var fileIn *os.File
			fileIn, createErr := os.Create(filePathStr)
			if createErr != nil {
				slog.Error("Error creating file:", createErr)
				return docs, createErr
			}
			defer fileIn.Close()
			_, err = fileIn.Write(binaryData)
			if err != nil {
				slog.Error("Error writing to file:", err)
				return docs, err
			}
		}
	} else {
		docs, err = readBinaryData(filePathStr)
	}
	return docs, nil
}
func readBinaryData(filePath string) (docs []Playitem, err error) {
	var file *os.File
	if runtime.GOOS == "linux" {
		shmFilePathStr := fmt.Sprintf("/dev/shm/%v.cache", GetMongoDbName())
		// 1. 复制文件到 /dev/shm
		err = copyFile(filePath, shmFilePathStr)
		if err != nil {
			slog.Error(fmt.Sprintf("复制文件到 /dev/shm 失败: %v\n", err))
			return
		}

		// 2. 从 /dev/shm 打开文件
		file, err = os.Open(shmFilePathStr)
		if err != nil {
			slog.Error(fmt.Sprintf("打开文件失败: %v\n", err))
			return
		}

		// 3. 使用完成后删除 /dev/shm 中的文件
		defer func() {
			err := os.Remove(shmFilePathStr)
			if err != nil {
				slog.Error(fmt.Sprintf("删除文件失败: %v\n", err))
			}
		}()
	} else {
		// 非 Linux 系统，直接读取原始文件
		file, err = os.Open(filePath)
		if err != nil {
			slog.Error("打开文件失败: %v\n", err)
			return
		}
	}
	defer file.Close()

	// 读取文件内容
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		slog.Error("Error:", err)
		return docs, err
	}

	// 获取文件大小
	fileSize := fileInfo.Size()
	// 预分配内存
	content := make([]byte, fileSize)

	// 读取文件内容到预分配的切片中
	_, err = file.Read(content)
	if err != nil {
		slog.Error("Error:", err)
		return docs, err
	}
	if len(content) > 0 {
		docs = decodeWithBinary(content)
	}
	return
}
func encodeWithBinary(data []Playitem) []byte {
	// 计算每个结构体的大小
	structSize := int(unsafe.Sizeof(Playitem{}))
	buf := make([]byte, 0, len(data)*structSize)

	// 将每个结构体转换为字节切片并追加到缓冲区
	for _, s := range data {
		// 使用 unsafe 将结构体转换为字节切片
		bytes := (*[unsafe.Sizeof(Playitem{})]byte)(unsafe.Pointer(&s))[:]
		buf = append(buf, bytes...)
	}

	return snappy.Encode(nil, buf)
}

func decodeWithBinary(data []byte) []Playitem {
	decompressed, _ := snappy.Decode(nil, data)
	// 计算每个结构体的大小
	structSize := int(unsafe.Sizeof(Playitem{}))
	result := make([]Playitem, len(decompressed)/structSize)

	// 将字节切片转换为结构体
	for i := range result {
		// 使用 unsafe 将字节切片转换为结构体
		result[i] = *(*Playitem)(unsafe.Pointer(&decompressed[i*structSize]))
	}
	return result
}
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // 文件存在
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // 文件不存在
	}
	return false // 其他错误
}
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	return nil
}
