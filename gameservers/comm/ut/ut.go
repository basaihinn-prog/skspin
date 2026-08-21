package ut

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func Time2YMD(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatGold(g int) string {
	return fmt.Sprintf("%.2f", float64(g)/1e4)
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func ToInt(s string) int {
	r, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0
	}
	return int(r)
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func ClipInt(n, min, max int) int {
	lo.Must0(min <= max)
	if n < min {
		return min
	}

	if max < n {
		return max
	}

	return n
}

func MaxInt(x ...int) int {
	var max = int(math.MinInt64)
	for _, v := range x {
		if v > max {
			max = v
		}
	}

	return max
}

// 随机 [low, high]
func RandomInt(low int, high int) int {
	if low == high {
		return low
	}
	if low > high || high == 0 {
		return 0
	}

	return rand.IntN(high-low+1) + low
}

// rand num is in [0,max)
func UtilRand(max int64) int {
	return RandomInt(0, int(max-1))
}

func RandomFromArray(arr []int) int {
	idx := rand.IntN(len(arr))
	return arr[idx]
}

// var UnmarshalJsonWithComment = ut2.UnmarshalJsonWithComment

func RandomFromWeightArr(arr []int) int {
	sum := 0
	for _, n := range arr {
		if n < 0 {
			panic("n < 0")
		}
		sum += n
	}

	if sum == 0 {
		panic("sum == 0")
	}

	rd := rand.IntN(sum)

	for i, n := range arr {
		if rd < n {
			return i
		}

		rd -= n
	}

	panic("no reachable!!")
}

func randomFromWeightArrNopanic(arr []int) (int, bool) {
	sum := 0
	for _, n := range arr {
		if n < 0 {
			panic("n < 0")
		}
		sum += n
	}

	if sum == 0 {
		// panic("sum == 0")
		return 0, false
	}

	rd := rand.IntN(sum)

	for i, n := range arr {
		if rd < n {
			return i, true
		}

		rd -= n
	}

	panic("no reachable!!")
}

func SampleByWeightsPred[T any](arr []T, getW func(item T) (weight int)) (ans T, ok bool) {
	weights := make([]int, len(arr))
	for i, item := range arr {
		weights[i] = getW(item)
	}

	choiceIdx, ok := randomFromWeightArrNopanic(weights)
	if !ok {
		return
	}
	ans = arr[choiceIdx]
	return
}

func PrintJson(i interface{}) {
	// buf, _ := json.MarshalIndent(i, "", "    ")
	// fmt.Println(string(buf))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(i)
}

func PopN[T any](arr *[]T, n int) []T {
	mid := len(*arr) - n
	if mid < 0 {
		mid = 0
	}
	result := (*arr)[mid:]

	*arr = (*arr)[:mid]
	return result
}

func Pop[T any](arr *[]T) T {
	mid := len(*arr) - 1
	result := (*arr)[mid]

	*arr = (*arr)[:mid]
	return result
}

// 随机出栈一个元素
func PopRand[T any](arr *[]T) T {
	n := len(*arr)
	i := rand.IntN(n)
	if i != n-1 {
		(*arr)[i], (*arr)[n-1] = (*arr)[n-1], (*arr)[i]
	}
	return Pop(arr)
}

func JoinStr(sep rune, elems ...string) string {
	return strings.Join(elems, string(sep))
}

func ErrString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// 初始化字段
// i 必须是想 *struct
// map, slice , *struct 类型的字段如果是nil的话将被初始化
func InitNilFields(i interface{}) {
	t := reflect.TypeOf(i)
	lo.Must0(t.Kind() == reflect.Ptr)
	t = t.Elem()
	lo.Must0(t.Kind() == reflect.Struct)

	v := reflect.ValueOf(i).Elem()

	propertyNums := t.NumField()

	for i := 0; i < propertyNums; i++ {
		stField := t.Field(i)
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		switch stField.Type.Kind() {
		case reflect.Map:
			if field.IsNil() {
				field.Set(reflect.MakeMap(stField.Type))
			}
		case reflect.Slice:
			if field.IsNil() {
				field.Set(reflect.MakeSlice(stField.Type, 0, 4))
			}
			// case reflect.Ptr:
			// 	if field.IsNil() && stField.Type.Elem().Kind() == reflect.Struct {
			// 		field.Set(reflect.New(stField.Type.Elem()))
			// 	}
		}
	}
}

func SetField(ptr any, fieldname string, val any) {
	v := reflect.ValueOf(ptr).Elem()
	field := v.FieldByName(fieldname)
	field.Set(reflect.ValueOf(val))
}

func MaxLen[T any](collection []T, maxcount int) []T {
	if n := len(collection); maxcount < n {
		return collection[n-maxcount:]
	}

	return collection
}

func IsSameDate(t1, t2 time.Time) bool {
	lo.Must0(t1.Location() == t1.Location())
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func Bindwrap(paramsType reflect.Type, fn func(i interface{}) error) (v reflect.Value, err error) {
	isPtr := paramsType.Kind() == reflect.Ptr
	if isPtr {
		paramsType = paramsType.Elem()
	}
	paramsValue := reflect.New(paramsType)

	err = fn(paramsValue.Interface())
	if err != nil {
		return
	}

	if isPtr {
		v = paramsValue
	} else {
		v = paramsValue.Elem()
	}

	return

}

func Money2Gold(money float64) int64 {
	//避免精度缺失，math.round四舍五入
	return int64(math.Round(money * 10000))
}

func Gold2Money(gold int64) float64 {
	return float64(gold) / 10000
}

// 不给小数点
func HackGold2Money(gold int64) int64 {
	return gold / 100
}
func JdbBet2Money(bet float64) float64 {
	return bet / 100
}
func FaCaiBet2Money(bet int64) float64 {
	return float64(bet / 100)
}
func FTrunc(x float64) float64 {
	return math.Floor(x*100) / 100
}

func TruncMsg(msg []byte, maxlen int) string {
	if len(msg) < maxlen {
		return string(msg)
	}

	return fmt.Sprintf("%s +[%d]more ", msg[:maxlen], len(msg))
}

func GetJsonRaw(i any) (raw json.RawMessage, err error) {
	buf, err := json.Marshal(i)
	if err != nil {
		return
	}

	raw = json.RawMessage(buf)
	return
}

func MKMul(m bson.M, k string, mul float64) {
	if v := m[k]; v != nil {
		fv := GetFloat(v)
		m[k] = Round6(fv * mul)
	}
}

func GetFloat(v interface{}) float64 {
	rv := reflect.ValueOf(v)
	if rv.CanInt() {
		return (float64)(rv.Int())
	}
	return rv.Float()
}

func GetInt(v interface{}) int {
	rv := reflect.ValueOf(v)
	return int(rv.Int())
}

func MK2Mul(m bson.M, k1, k2 string, mul float64) {
	if m1 := m[k1]; m1 != nil {
		if _, ok := m1.(bson.M); ok {
			MKMul(m1.(bson.M), k2, mul)
		}
	}
}

func MK3Mul(m bson.M, k1, k2, k3 string, mul float64) {
	if m1 := m[k1]; m1 != nil {
		if _, ok := m1.(bson.M); ok {
			MK2Mul(m1.(bson.M), k2, k3, mul)
		}
	}
}

func Round6(fv float64) float64 {
	return math.Round(fv*1e6) / 1e6
}

func FloatEQ(f1, f2 float64) bool {
	return math.Abs(f1-f2) < 1e-6
}

func FloatPtrMul(pf *float64, mul float64) {
	if pf != nil {
		*pf = Round6(*pf * mul)
	}
}

func formatNumber(num float64) string {
	// 如果是整数,直接返回整数部分
	if num == math.Trunc(num) {
		return fmt.Sprintf("%.0f", num)
	}

	// 如果是小数,保留2位小数
	return strconv.FormatFloat(num, 'f', 2, 64)
}

func formatNumber1(num float64) string {
	// 如果是整数,直接返回整数部分
	if num == math.Trunc(num) {
		return fmt.Sprintf("%.0f", num)
	}

	// 如果是小数,保留2位小数
	return strconv.FormatFloat(num, 'f', 1, 64)
}

func FloatStrPtrMul(pf *string, mul float64) {
	if pf == nil || *pf == "" || *pf == "0" {
		return
	}

	v := lo.Must(strconv.ParseFloat(*pf, 64))

	//s := strconv.FormatFloat(v*mul, 'f', -1, 64)
	s := formatNumber(v * mul)
	*pf = s
}
func FloatStrPtrMul1(pf *string, mul float64) {
	if pf == nil || *pf == "" || *pf == "0" {
		return
	}

	v := lo.Must(strconv.ParseFloat(*pf, 64))

	//s := strconv.FormatFloat(v*mul, 'f', -1, 64)
	s := formatNumber1(v * mul)
	*pf = s
}
func FloatStrPtrDiv(pf *string, mul float64) {
	if pf == nil || *pf == "" || *pf == "0" {
		return
	}

	v := lo.Must(strconv.ParseFloat(*pf, 64))

	// s := strconv.FormatFloat(v/mul, 'f', -1, 64)
	s := formatNumber(v / mul)
	*pf = s
}

func Ftoa(f float64) string {
	f = Round6(f)
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func MK1Div(m bson.M, k string, div float64) {
	if v := m[k]; v != nil {
		fv := GetFloat(v)
		m[k] = Round6(fv / div)
	}
}

func MK2Div(m bson.M, k1, k2 string, div float64) {
	if m1 := m[k1]; m1 != nil {
		if _, ok := m1.(bson.M); ok {
			MK1Div(m1.(bson.M), k2, div)
		}
	}
}

func MK3Div(m bson.M, k1, k2, k3 string, div float64) {
	if m1 := m[k1]; m1 != nil {
		if _, ok := m1.(bson.M); ok {
			MK2Div(m1.(bson.M), k2, k3, div)
		}
	}
}

func FloatPtrDiv(pf *float64, mul float64) {
	if pf != nil {
		*pf = Round6(*pf / mul)
	}
}

func ReverseString(s string) string {
	a := []byte(s)
	slices.Reverse(a)
	return string(a)
}

// fivestar专用
func Float64ToInt64TwoDecimal(floats []float64) []int64 {
	ints := make([]int64, len(floats))
	for i, f := range floats {
		// 处理潜在溢出
		scaled := f * 100
		if scaled > math.MaxInt64 || scaled < math.MinInt64 {
			fmt.Errorf("value %f scaled to %f, out of range for int64", f, scaled)
		}

		// 四舍五入处理
		rounded := math.Round(scaled)

		// 转换回int64
		ints[i] = int64(rounded) / 100 * 1000

	}
	return ints
}

func FloatArrMul(arr []float64, mul float64) []float64 {
	ret := make([]float64, 0, len(arr))
	for i := range arr {
		ret = append(ret, arr[i]*mul)
	}
	return ret
}

// bool: equal
func FloatInArr(arr []float64, anim float64) bool {
	for i := range arr {
		if FloatEQ(arr[i], anim) {
			return true
		}
	}
	return false
}
func Int64ArrMul(arr []int64, mul float64) []int64 {
	ret := make([]int64, 0, len(arr))
	for i := range arr {
		ret = append(ret, int64(float64(arr[i])*mul))
	}
	return ret
}

// bool: equal
func Int64InArr(arr []int64, anim int64) bool {
	for i := range arr {
		if arr[i] == anim {
			return true
		}
	}
	return false
}

// bool: equal
func Int64InArrCount(arr []int64, anim int64) int64 {
	count := int64(0)
	for i := range arr {
		if arr[i] == anim {
			count++
		}
	}
	return count
}

// map转字符串
func Map2Str(v map[string]any) string {
	if len(v) == 0 {
		return ""
	}

	buf := BP.Get()
	defer BP.Put(buf)

	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		vs := v[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		switch vs.(type) { //Reduce the use of fmt.Sprint for basic types.
		case string:
			buf.WriteString(vs.(string))
		case int:
			buf.WriteString(strconv.FormatInt(int64(vs.(int)), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(int64(vs.(int64)), 10))
		case uint64:
			buf.WriteString(strconv.FormatUint(uint64(vs.(uint64)), 10))
		case float64:
			buf.WriteString(strconv.FormatFloat(float64(vs.(float64)), 'f', 4, 64))
		case bool:
			buf.WriteString(strconv.FormatBool(vs.(bool)))
		case []string:
			buf.WriteString(strings.Join(vs.([]string), ","))
		case []int:
			arr := make([]string, len(vs.([]int)))
			for i, vv := range vs.([]int) {
				arr[i] = strconv.Itoa(vv)
			}
			buf.WriteString(strings.Join(arr, ","))
		case []float64:
			arr := make([]string, len(vs.([]float64)))
			for i, vv := range vs.([]float64) {
				arr[i] = strconv.FormatFloat(vv, 'f', 4, 64)
			}
			buf.WriteString(strings.Join(arr, ","))
		default:
			str := fmt.Sprint(vs)
			//array
			if len(str) > 2 && str[0] == '[' && str[len(str)-1] == ']' {
				arr := strings.Split(str[1:len(str)-1], " ")
				buf.WriteString(strings.Join(arr, ","))
			} else {
				buf.WriteString(str)
			}
		}
	}
	buf.WriteByte('&')

	return buf.String()
}

func ToMD5Str(s string) string {
	data := md5.Sum([]byte(s))
	return hex.EncodeToString(data[:])
}

// 服务名置换厂商名
func ServiceNameToManufacturerName(s string) string {
	return strings.Split(s, "_")[0]
}

func DoHttpReq(client *http.Client, req *http.Request) (body []byte, resp *http.Response, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		err = errors.New(resp.Status)
		return
	}

	body, err = io.ReadAll(resp.Body)
	return
}

func GetShortUUID() string {
	newUUID := uuid.New()

	// 截取前 8 位作为短 UUID
	shortUUID := newUUID.String()[:8]
	return shortUUID
}

func Diff(sli []string, val string) (newSli []string) {
	for _, v := range sli {
		if v != val {
			newSli = append(newSli, v)
		}
	}
	return newSli
}
