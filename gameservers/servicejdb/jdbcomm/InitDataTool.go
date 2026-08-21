package jdbcomm

//
//import (
//	"fmt"
//	"math"
//	"regexp"
//	"strconv"
//	"strings"
//	"sync"
//)
//
//var (
//	RoundID   int
//	TurnIndex int
//	lock      sync.RWMutex
//)
//
//func init() {
//	RoundID = 0
//	TurnIndex = 0
//}
//
//func GetRoundID() int {
//	lock.Lock()
//	defer lock.Unlock()
//	return RoundID
//}
//
//func SetRoundID(id int) {
//	lock.Lock()
//	defer lock.Unlock()
//	RoundID = id
//}
//
//func GetTurnIndex() int {
//	lock.Lock()
//	defer lock.Unlock()
//	return TurnIndex
//}
//
//func SetTurnIndex(index int) {
//	lock.Lock()
//	defer lock.Unlock()
//	TurnIndex = index
//}
//
//// ReplaceEmbedScoreField 替换嵌入的得分字段
//func ReplaceEmbedScoreField(query Variables, multiply int) Variables {
//	// 遍历替换中奖线
//	for i := 0; i <= 99; i++ {
//		field := "l" + strconv.Itoa(i)
//		if !ReplaceScoreField(query, multiply, field, '~') {
//			break
//		}
//	}
//
//	ReplaceScoreField(query, multiply, "wlc_v", '~')
//	ReplaceScoreField(query, multiply, "psym", '~')
//	ReplaceScoreFieldRegex(query, multiply, "trail", regexp.MustCompile(`totmul~\d+;nmwin~(\d+\.\d{2})`))
//
//	return query
//}
//
//// ReplaceScoreField 使用正则表达式替换字段值
//func ReplaceScoreFieldRegex(query Variables, multiply int, key string, r *regexp.Regexp) bool {
//	// 获取对应键的值
//	value, exists := query[key]
//	if !exists || value == "" {
//		return false
//	}
//
//	isChanged := false
//
//	// 替换逻辑
//	newValue := r.ReplaceAllStringFunc(value, func(match string) string {
//		submatches := r.FindStringSubmatch(match)
//		if len(submatches) > 1 {
//			v := submatches[1] // 捕获组中的值
//			if result, err := strconv.ParseFloat(v, 64); err == nil && result != 0.00 {
//				isChanged = true
//				replacement := fmt.Sprintf("%.2f", result*float64(multiply))
//				return match[:len(match)-len(v)] + replacement
//			}
//		}
//		return match
//	})
//
//	// 更新 query 值
//	if isChanged {
//		query[key] = newValue
//	}
//
//	return isChanged
//}
//
//// ReplaceScoreField 按指定标识分割并将分数乘以倍数
//func ReplaceScoreField(query Variables, multiply int, key string, splitor rune) bool {
//	// 获取键对应的值
//	value, exists := query[key]
//	if !exists || !strings.ContainsRune(value, splitor) {
//		return false
//	}
//
//	isChanged := false
//
//	// 分割字符串
//	array := strings.Split(value, string(splitor))
//	for i, val := range array {
//		if strings.Contains(val, ".") {
//			if result, err := strconv.ParseFloat(val, 64); err == nil && result != 0.00 {
//				// 替换分数
//				array[i] = fmt.Sprintf("%.2f", result*float64(multiply))
//				isChanged = true
//				break // 只处理第一个匹配到的值
//			}
//		}
//	}
//
//	// 如果有修改，则更新 query
//	if isChanged {
//		query[key] = strings.Join(array, string(splitor))
//	}
//	return isChanged
//}
//
//// ReplaceScoreField 按指定标识分割并将分数乘以倍数
//func ReplaceScoreField2(query Variables, multiply float64, key string, splitor rune) bool {
//	// 获取键对应的值
//	value, exists := query[key]
//	if !exists || !strings.ContainsRune(value, splitor) {
//		return false
//	}
//
//	isChanged := false
//
//	// 分割字符串
//	array := strings.Split(value, string(splitor))
//	for i, val := range array {
//		if strings.Contains(val, ".") {
//			if result, err := strconv.ParseFloat(val, 64); err == nil && result != 0.00 {
//				// 替换分数
//				array[i] = fmt.Sprintf("%.2f", result*multiply)
//				isChanged = true
//				break // 只处理第一个匹配到的值
//			}
//		}
//	}
//
//	// 如果有修改，则更新 query
//	if isChanged {
//		query[key] = strings.Join(array, string(splitor))
//	}
//	return isChanged
//}
//
//// AppendBalance 更新查询参数并计算余额
//func AppendBalance(query Variables, bet int64, balance float64, multiply int, enterBalance float64) {
//	// 获取 "w" 的值
//	w, exists := query["w"]
//	if exists {
//		if win, err := strconv.ParseFloat(w, 64); err == nil {
//			balance -= float64(bet)
//
//			//这个值要找从哪里来的
//			changedBalance := 1
//			changedBalance += int(math.Round(win * float64(multiply) * 100))
//		}
//	}
//	// 添加 "ntp"
//	query["ntp"] = fmt.Sprintf("%.2f", balance-enterBalance)
//
//	// 添加 "balance"、"balance_cash" 和 "balance_bonus"
//	b := fmt.Sprintf("%.2f", balance/100)
//	query["balance"] = b
//	query["balance_cash"] = b
//	query["balance_bonus"] = "0.00"
//}
//
//// ReplaceAllScoreField 替换查询参数中指定字段的分数乘以倍数
//func ReplaceAllScoreField(query Variables, multiply int) {
//	// 指定需要处理的字段列表
//	fields := []string{
//		"w", "tw", "rs_win", "fswin", "fswin_total", "fsres", "rs_iw",
//		"fsres_total", "mo_tw", "tmb_win", "tmb_res", "pw", "apwa",
//	}
//
//	// 遍历字段并更新分数
//	for _, field := range fields {
//		if value, exists := query[field]; exists {
//			// 尝试将字段值解析为浮点数
//			if val, err := strconv.ParseFloat(value, 64); err == nil {
//				// 将值乘以倍数并格式化为两位小数
//				query[field] = fmt.Sprintf("%.2f", val*float64(multiply))
//			}
//		}
//	}
//}
