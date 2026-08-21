package slotsdb

import (
	"bytes"
	"compress/gzip"
	"game/comm/xbin"
	"game/duck/logger"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	"google.golang.org/protobuf/encoding/protowire"
)

type DB struct {
	filename string
	rawfile  *os.File
	colls    map[string]*Coll
	hashColl map[int]*Coll
	mutex    sync.Mutex
	collID   int
	writer   *DBWriter
}

type BulkCache struct {
	Header []byte
	Body   []byte
}

func NewDB(filename string) (*DB, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	f.Seek(0, os.SEEK_END) // 要在末尾写入。

	sf := &DB{
		filename: filename,
		rawfile:  f,
		colls:    map[string]*Coll{},
		hashColl: map[int]*Coll{},
		writer:   NewDBWriter(f),
	}

	logger.Info("read file:", filename)

	err = sf.init(filename)
	if err != nil {
		logger.Err("sf.init() err:", filename, err)
		return nil, err
	}

	return sf, nil
}

func (sf *DB) init(filename string) error {
	zipbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Err("ioutil.ReadFile(filename) err:", filename, err)
		return err
	}

	if len(zipbytes) == 0 {
		return nil
	}

	// 前10个字节是gzip的文件头: 1f8b 0800 0000 0000 00ff
	// 只有10个字节说明没有业务数据。
	if len(zipbytes) == 10 {
		return nil
	}

	gzReader, err := gzip.NewReader(bytes.NewReader(zipbytes))
	if err != nil {
		logger.Err("gzip.NewReader(bytes.NewReader(zipbytes))", err, len(zipbytes))
		return err
	}

	bytes, err := ioutil.ReadAll(gzReader)
	if err != nil {
		logger.Err("ioutil.ReadAll(gzReader) err:", len(bytes), err)
		return err
	}

	logger.Info("Init()", len(bytes))

	idx := 0
	for idx < len(bytes) {
		end := idx + 8
		if end > len(bytes) {
			end = len(bytes)
		}
		t := bytes[idx:end]
		size, n := protowire.ConsumeVarint(t)
		idx += n
		dataBytes := bytes[idx : idx+int(size)]
		idx += len(dataBytes)

		var root DataRoot
		err := xbin.Unmarshal(dataBytes, &root)
		if err != nil {
			logger.Err("xbin.Unmarshal(dataBytes, &root) err:", err)
			return err
		}

		var bodyBytes []byte
		if root.DataSize > 0 {
			bodyBytes = bytes[idx : idx+root.DataSize]
			idx += root.DataSize
		}

		switch root.Type {
		case SubDateType_TableDefine:
			var tablename = string(bodyBytes)
			coll := NewColl(root.TableID, tablename, sf)
			sf.colls[tablename] = coll
			sf.hashColl[root.TableID] = coll

			if sf.collID < root.TableID {
				sf.collID = root.TableID
			}

			logger.Info("TableDefine:", tablename)
		case SubDateType_Drop:
			collName := sf.hashColl[root.TableID].Name
			delete(sf.hashColl, root.TableID)
			delete(sf.colls, collName)
		default:
			coll := sf.hashColl[root.TableID]
			if coll == nil {
				panic("coll == nil")
			}
			err := sf.hashColl[root.TableID].load(&root, bodyBytes)
			if err != nil {
				logger.Err("xbin.Unmarshal(root.SubData, &td) err:", err)
				return err
			}
		}
	}

	logger.Info("Init() finish.")

	return nil
}

func (sf *DB) writeData(tableID int, rowID int, dataTipe SubDataType, data []byte) error {
	bytes, err := xbin.Marshal(&DataRoot{
		Type:     dataTipe,
		TableID:  tableID,
		RowID:    rowID,
		DataSize: len(data),
	})
	if err != nil {
		return err
	}

	sf.writer.addItem(bytes, data)

	return err
}

func (sf *DB) Commit() error {
	return sf.writer.Commit()
}

func (sf *DB) C(tablename string) *Coll {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	coll := sf.colls[tablename]
	if coll == nil {
		logger.Info("create table:", sf.collID, tablename)
		sf.collID++
		coll = NewColl(sf.collID, tablename, sf)
		sf.colls[tablename] = coll
		sf.hashColl[coll.collID] = coll
		sf.writeData(sf.collID, 0, SubDateType_TableDefine, []byte(tablename))
	}

	return coll
}

func (sf *DB) GetAllTableNames() []string {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	var arr = []string{}

	for name := range sf.colls {
		arr = append(arr, name)
	}

	sort.Strings(arr)

	return arr
}

func (sf *DB) Close() {
	sf.writer.Commit()
	sf.rawfile.Close()
}

type DBWriter struct {
	rawfile *os.File
	cache   []*BulkCache
	mutex   sync.Mutex
}

func NewDBWriter(file *os.File) *DBWriter {
	return &DBWriter{
		rawfile: file,
		cache:   []*BulkCache{},
	}
}

func (sf *DBWriter) addItem(header []byte, body []byte) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.cache = append(sf.cache, &BulkCache{
		Header: header,
		Body:   body,
	})
}

func (sf *DBWriter) Commit() error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	// gzWriter不要使用成员变量的方式，当程序中断后进行gzip读取的时候出现unexpected err，原因是gzWriter有缓存没有写完，即使我使用了gzWriter.Flush()也没有解决。改成局部变量的方式可以完美解决这个问题。
	gzWriter := gzip.NewWriter(sf.rawfile)

	// 使用cache的方式，可以每秒钟写入5万条记录左右。如果不使用cache，每条记录新建一个gzWriter，是每秒钟写入2500条记录左右。
	for _, item := range sf.cache {
		_, err := sf.write(gzWriter, item.Header, item.Body)
		if err != nil {
			return err
		}
	}

	gzWriter.Close()

	sf.cache = []*BulkCache{}

	return nil
}

func (sf *DBWriter) write(gzWriter *gzip.Writer, p []byte, data []byte) (n int, err error) {
	total := 0

	// 写入p的长度
	size, err := gzWriter.Write(protowire.AppendVarint(nil, uint64(len(p))))
	if err != nil {
		logger.Err("Db.write err:", err)
		return 0, err
	}

	total += size

	// 写入p的数据
	size, err = gzWriter.Write(p)
	if err != nil {
		logger.Err("Db.write err:", err)
		return total, err
	}

	total += size

	// 写入数据
	if len(data) > 0 {
		size, err = gzWriter.Write(data)
		total += size
	}

	return total, err
}
