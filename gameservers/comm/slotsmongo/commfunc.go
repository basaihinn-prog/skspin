package slotsmongo

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func GenerateRandStr(length int) string {
	rand.Seed(time.Now().UnixNano()) // 必须设置随机种子

	buf := make([]byte, length)
	for i := range buf {
		buf[i] = charset[rand.Intn(len(charset))]
	}

	return string(buf)
}

func GenerateSha512(str string) (hexStr, truncatedHex string, decimalValue int64, err error) {
	// 输入数据（按需替换）
	input := []byte(str)

	// 生成 SHA512 哈希
	hash := sha512.Sum512(input)
	hexStr = hex.EncodeToString(hash[:])

	// 截取前13位十六进制
	truncatedHex = hexStr[:13]

	// 十六进制转十进制（处理奇数位）
	value, err := strconv.ParseUint(truncatedHex, 16, 64)
	if err != nil {
		err = fmt.Errorf("十六进制转十进制失败: %v", err)
	}
	decimalValue = int64(value)
	return
}

func GenerateSha256(str string) (hexStr string) {
	// 创建一个新的 SHA256 哈希对象
	hash := sha256.New()

	// 写入数据
	hash.Write([]byte(str))

	// 计算哈希值
	hashed := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashed)
}
