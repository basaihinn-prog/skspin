package slotsmongo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"serve/comm/ut"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PPSimulateData struct {
	Id       primitive.ObjectID `bson:"_id"`
	DropPan  []bson.M           `bson:"droppan"`
	HasGame  bool               `bson:"hasgame"`
	Times    float64            `bson:"times"`
	BucketId int                `json:"bucketid"`
	Type     int                `bson:"type"`
	Selected bool               `bson:"selected"`
}

func processLine(line string, bet float64) string {
	values := strings.Split(line, "~")
	prize, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		fmt.Errorf("invalid prize value: %s", values[1])
		return ""
	}

	// 将奖金乘以2
	newPrize := prize * bet

	// 将新的奖金值格式化为字符串，保留两位小数
	values[1] = fmt.Sprintf("%.2f", newPrize)

	// 重新组合字符串
	newLine := fmt.Sprintf("%s", strings.Join(values, "~"))
	return newLine
}
func extractLLines(data map[string]interface{}) []string {
	var result []string
	for key, value := range data {
		if strings.HasPrefix(key, "l") {
			if num, err := strconv.Atoi(key[1:]); err == nil {
				result = append(result, fmt.Sprintf("l%d:%v", num, value))
			}
		}
	}
	return result
}

// 最小下注是1
func (sd *PPSimulateData) Deal(num, all int, balance, c, l, buyMul float64, isBuy bool) map[string]any {
	//todo 字段待替换成pp字段
	bet := c * l
	pan := sd.DropPan[0]

	delete(pan, "orignid")
	pan["c"] = c
	pan["l"] = l
	//mul := float64(1)
	//if isBuy {
	//	mul = buyMul
	//}

	//ut.MKMul(pan, "aw", bet)
	//ut.MKMul(pan, "actw", bet)

	ut.MKMul(pan, "tmb_win", bet)     // 连消的累计奖金
	ut.MKMul(pan, "tw", bet)          // 累计奖金
	ut.MKMul(pan, "w", bet)           // 本次奖金
	ut.MKMul(pan, "fswin", bet)       // FG中累计获取奖励
	ut.MKMul(pan, "fswin_total", bet) // FG结束累计获取奖励
	ut.MKMul(pan, "fsres", bet)       // FG中累计获取奖励
	ut.MKMul(pan, "rw", bet)          // EPIC LINK的奖金
	ut.MKMul(pan, "pw", bet)          // EPIC LINK的奖金
	//ut.MKMul(pan, "balance_bonus", bet) // 特殊字符奖金
	//ut.MKMul(pan, "mbv", bet)     		// 当前SPIN的倍数
	winLine := extractLLines(pan)
	for _, line := range winLine {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			fmt.Errorf("invalid line format: %s", line)
		}
		key := parts[0]
		newLine := processLine(parts[1], bet)
		pan[key] = newLine
	}

	//if num == all-1 {
	//	pan["balance"] = balance + ut.Round6(ut.GetFloat(pan["tw"]))      //SPIN之后剩余金币数
	//	pan["balance_cash"] = balance + ut.Round6(ut.GetFloat(pan["tw"])) //SPIN之后剩余金币数
	//} else {
	pan["balance"] = balance
	pan["balance_cash"] = balance
	//}
	//todo 返回中的时间戳也需要修改
	pan["stime"] = time.Now().UnixMilli()

	return pan
}

// 针对最小下注不是1的
func (sd *PPSimulateData) Deal2(num, all int, balance, c, l, line, buyMul, originMin float64, isBuy bool) map[string]any {
	bet := c * l
	pan := sd.DropPan[0]

	delete(pan, "orignid")
	pan["c"] = c
	pan["l"] = l
	bet = bet / originMin

	ut.MKMul(pan, "tmb_win", bet)       // 连消的累计奖金
	ut.MKMul(pan, "tw", bet)            // 累计奖金
	ut.MKMul(pan, "w", bet)             // 本次奖金
	ut.MKMul(pan, "fswin", bet)         // FG中累计获取奖励
	ut.MKMul(pan, "fswin_total", bet)   // FG结束累计获取奖励
	ut.MKMul(pan, "fsres", bet)         // FG中累计获取奖励
	ut.MKMul(pan, "rw", bet)            // EPIC LINK的奖金
	ut.MKMul(pan, "pw", bet)            // EPIC LINK的奖金
	ut.MKMul(pan, "balance_bonus", bet) // 特殊字符奖金
	winLine := extractLLines(pan)
	for _, line := range winLine {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			fmt.Errorf("invalid line format: %s", line)
		}
		key := parts[0]
		newLine := processLine(parts[1], bet)
		pan[key] = newLine
	}

	//if num == all-1 {
	//	pan["balance"] = balance + ut.Round6(ut.GetFloat(pan["tw"]))      //SPIN之后剩余金币数
	//	pan["balance_cash"] = balance + ut.Round6(ut.GetFloat(pan["tw"])) //SPIN之后剩余金币数
	//} else {
	pan["balance"] = balance
	pan["balance_cash"] = balance
	//}
	pan["stime"] = time.Now().UnixMilli()

	return pan
}
