package comm

import (
	"context"
	"errors"
	"fmt"
	"game/comm/db"
	"game/duck/lazy"
	"game/duck/logger"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"os"
	"strings"
)

const statdir = "statistics"

func OpenStatFile(file string, flag int) (f *os.File, err error) {
	os.Mkdir(statdir, 0755)
	filename := fmt.Sprintf("%v/%v_%v.stat", statdir, lazy.ServiceName, file)
	f, err = os.OpenFile(filename, flag|os.O_RDWR, 0644)
	if err != nil {
		logger.Err(err)
	}
	return
}

func DIV[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](num1, num2 T) T {
	if num2 == 0 {
		return 0
	}
	return num1 / num2
}

func RandPassword() string {
	baseStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*+-_="

	length := 8
	bytes := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		bytes = append(bytes, baseStr[rand.Int()%len(baseStr)])
	}
	return string(bytes)
}

// 千分位
func ThousandSeparator(s string) string {
	numStrs := strings.Split(s, ".")
	numStr := numStrs[0]
	length := len(numStr)
	if length < 4 {
		return s
	}
	cnt := (length - 1) / 3
	for i := 0; i < cnt; i++ {
		numStr = numStr[:length-(i+1)*3] + "," + numStr[length-(i+1)*3:]
	}
	if len(numStrs) == 1 {
		return numStr
	}
	return fmt.Sprintf("%s.%s", numStr, numStrs[1])
}

func VerifyMaxRTP(setRTP float64) (err error, passThrough bool) {
	// 验证RTP可以设置的最大值
	sys := SystemConfig{}
	CollConfig := db.Collection2("GameAdmin", "SystemConfig")
	err = CollConfig.FindOne(context.TODO(), bson.M{"AppId": "admin", "MaxRTP": bson.M{"$gte": setRTP}}).Decode(&sys)

	if err != nil {

		err = errors.New("超出RTP最大值")
		return err, false
	}

	if &sys == nil {

		err = errors.New("超出RTP最大值")
		return err, false
	}
	return err, true
}
