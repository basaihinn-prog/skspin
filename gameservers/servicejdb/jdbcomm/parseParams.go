package jdbcomm

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Variables map[string]any

func convertToMapSafely(obj any) Variables {
	// 直接断言为 map[string]interface{}
	if m, ok := obj.(Variables); ok {
		return m
	}

	// 使用反射进行转换
	val := reflect.ValueOf(obj)

	// 如果是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 处理 map 类型
	if val.Kind() == reflect.Map {
		result := make(Variables)
		for _, key := range val.MapKeys() {
			result[key.String()] = val.MapIndex(key).Interface()
		}
		return result
	}

	// 处理 struct 类型
	if val.Kind() == reflect.Struct {
		result := make(Variables)
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldName := typ.Field(i).Name

			// 获取字段的 tag（如果需要）
			tag := typ.Field(i).Tag.Get("bson")
			if tag != "" {
				fieldName = tag
			}

			result[fieldName] = field.Interface()
		}
		return result
	}

	return nil
}

func convertMongoData(data interface{}) ([]Variables, error) {
	// 处理 []interface{} 类型
	if sliceData, ok := data.([]interface{}); ok {
		result := make([]Variables, 0, len(sliceData))
		for _, item := range sliceData {
			// 直接断言为 map[string]interface{}
			if mapItem, ok := item.(Variables); ok {
				result = append(result, mapItem)
			} else {
				// 尝试使用反射转换
				val := reflect.ValueOf(item)
				if val.Kind() == reflect.Map {
					convertedMap := make(Variables)
					for _, key := range val.MapKeys() {
						convertedMap[key.String()] = val.MapIndex(key).Interface()
					}
					result = append(result, convertedMap)
				} else {
					return nil, fmt.Errorf("item is not a map: %v (type: %T)", item, item)
				}
			}
		}
		return result, nil
	}

	// 处理 primary.A 类型（如果使用了自定义类型）
	if sliceData, ok := data.(primitive.A); ok {
		result := make([]Variables, 0, len(sliceData))
		for _, item := range sliceData {
			// 直接断言为 map[string]interface{}
			if mapItem, ok := item.(Variables); ok {
				result = append(result, mapItem)
			} else {
				// 尝试使用反射转换
				val := reflect.ValueOf(item)
				if val.Kind() == reflect.Map {
					convertedMap := make(Variables)
					for _, key := range val.MapKeys() {
						convertedMap[key.String()] = val.MapIndex(key).Interface()
					}
					result = append(result, convertedMap)
				} else if val.Kind() == reflect.Slice {

					convertedMap := make(Variables)
					for i, datum := range sliceData {
						convertedMap[strconv.Itoa(i)] = datum
					}
					result = append(result, convertedMap)
				} else {
					return nil, fmt.Errorf("item is not a map: %v (type: %T)", item, item)
				}
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("unexpected type: %T", data)
}

// 安全地获取嵌套 map 的辅助函数
func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}

// 安全地获取特定类型值的辅助函数
func getValue[T any](m map[string]interface{}, key string) (T, bool) {
	val, ok := m[key].(T)
	return val, ok
}

// 获取变量的值，返回 any 和是否存在
func (v Variables) get(key string) (any, bool) {
	value, exists := v[key]
	return value, exists
}

// GetGameSt 根据 rsp 字段判断当前状态
func (rspMap Variables) GetGameSt() int {
	gameSt := 0
	fsmul, _ := rspMap.get("fsmul")

	if tw, exists := rspMap["tw"]; !exists || tw == "" {
		if rspMap.Str("na") == "s" {
			return 0
		}
		return 1
	}

	twValue := rspMap.Float("tw")
	wValue := rspMap.Float("w")

	if twValue == 0 && wValue == 0 && fsmul == "" {
		gameSt = 1
	} else if twValue > 0 && wValue > 0 && rspMap.Str("na") == "s" && fsmul == "" {
		gameSt = 2
	} else if twValue > 0 && rspMap.Str("na") == "c" && fsmul == "" {
		gameSt = 3
	} else if twValue > 0 && rspMap.Int("fsend_total") == 1 {
		gameSt = 5
	} else {
		gameSt = 4
	}
	return gameSt
}

// Int 获取整数值
func (ps Variables) Int(k string) int {
	if v, exists := ps[k]; exists {
		switch v := v.(type) {
		case string:
			ret, _ := strconv.Atoi(v)
			return ret
		case int:
			return v
		}
	}
	return 0
}

// SetInt 设置整数值
func (ps Variables) SetInt(k string, v int) {
	ps[k] = v
}

// Float 获取浮点数值
func (ps Variables) Float(k string) float64 {
	if v, exists := ps[k]; exists {
		switch v := v.(type) {
		case string:
			ret, _ := strconv.ParseFloat(v, 64)
			return ret
		case float64:
			return v
		}
	}
	return 0
}

// SetFloat 设置浮点数值
func (ps Variables) SetFloat(k string, v float64) {
	ps[k] = v
}

// 乘法运算
func (ps Variables) MKMulFloat(k string, mul float64) {
	if v, exists := ps[k]; exists {
		switch v := v.(type) {
		case string:
			// 尝试解析字符串为浮点数
			if value, err := strconv.ParseFloat(v, 64); err == nil {
				ps.SetFloat(k, value*mul)
			}
		case float64:
			// 如果已经是浮点数，直接乘法
			a := int(v * mul)
			ps.SetFloat(k, float64(a))
		case int:
			// 如果是整数，转换为浮点数计算
			ps.SetFloat(k, float64(v)*mul)
		}
	}
}

// 设置浮点数组并转换为字符串存储
func (ps Variables) SetFloatArr(k string, v []float64) {
	var str strings.Builder
	for i, val := range v {
		if i != 0 {
			str.WriteString(",")
		}
		str.WriteString(fmt.Sprintf("%.2f", val))
	}
	ps[k] = str.String()
}

// 货币格式化
func (ps Variables) Currency(k string) float64 {
	v := ps.Str(k)
	v = strings.ReplaceAll(v, ",", "")
	ret, _ := strconv.ParseFloat(v, 64)
	return ret
}

// 货币存储
func (ps Variables) SetCurrency(k string, v float64) {
	//p := message.NewPrinter(language.English)

	ps[k] = fmt.Sprintf("%.2f", v)
}

// 获取字符串
func (ps Variables) Str(k string) string {
	if v, exists := ps[k]; exists {
		return fmt.Sprintf("%v", v) // 确保转换为字符串
	}
	return ""
}

// 删除键
func (ps Variables) Delete(k string) {
	delete(ps, k)
}

// 是否存在键
func (ps Variables) Exist(k string) bool {
	_, exists := ps[k]
	return exists
}

// 设置字符串
func (ps Variables) SetStr(k string, v string) {
	ps[k] = v
}

// JSON 解析
func (ps Variables) JsonUnmarshal(k string, recv any) error {
	v := ps.Str(k)
	return json.Unmarshal([]byte(v), recv)
}

// JSON 序列化
func (ps Variables) SetJson(k string, v any) error {
	s, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ps[k] = string(s)
	return nil
}

// 编码为 URL 参数格式
func (ps Variables) Encode() string {
	keys := lo.Keys(ps)
	sort.Strings(keys)
	var sb strings.Builder

	for i, k := range keys {
		if i != 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(ps.Str(k))
	}
	return sb.String()
}

// 格式化输出
func (ps Variables) PrettyString() string {
	keys := lo.Keys(ps)
	sort.Strings(keys)
	var sb strings.Builder

	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(ps.Str(k))
		sb.WriteByte('\n')
	}
	return sb.String()
}

// 获取字节
func (ps Variables) Bytes() []byte {
	ps.Delete("rid")
	return []byte(ps.Encode())
}

// 解析 URL 参数格式的字符串
func ParseVariables(s string) Variables {
	vars := Variables{}
	pairs := strings.Split(s, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		k, v := kv[0], ""
		if len(kv) > 1 {
			v = kv[1]
		}
		vars[k] = v
	}
	return vars
}

// 打印格式化变量
func PrettyPrintVas(w io.Writer, s string) {
	io.WriteString(w, ParseVariables(s).PrettyString())
	io.WriteString(w, "\n===\n\n")
}

// 乘法修改 mo_tw 字段
func (ps Variables) MKMulmo_twInG(k string, mul float64) {
	g := ps.Str("g")

	re := regexp.MustCompile(`mo_tw:"([\d.]+)"`)
	g = re.ReplaceAllStringFunc(g, func(match string) string {
		split := strings.Split(match, ":")
		if len(split) != 2 {
			return match
		}
		value := MKMulFloat(strings.Trim(split[1], `"`), mul)
		return fmt.Sprintf(`mo_tw:"%v"`, value)
	})

	ps["g"] = g
}

// 乘法转换
func MKMulFloat(k string, mul float64) string {
	k = strings.ReplaceAll(k, ",", "")
	ret, _ := strconv.ParseFloat(k, 64)
	return fmt.Sprintf("%.2f", ret*mul)
}
