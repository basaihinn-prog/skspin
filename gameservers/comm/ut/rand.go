package ut

import (
	"math/rand/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.N(len(charset))]
	}
	return string(result)
}

const charset2 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成指定长度的随机字符串
func GenerateRandomString2(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset2[rand.N(len(charset2))]
	}
	return string(result)
}
