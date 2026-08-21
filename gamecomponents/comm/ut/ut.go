package ut

import (
	"encoding/json"
	"errors"
	"fmt"
	"game/duck/ut2"
	"github.com/google/uuid"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

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

var UnmarshalJsonWithComment = ut2.UnmarshalJsonWithComment

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

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "    ")
	// enc.Encode(i)

	os.Stdout.Write(lo.Must(json.Marshal(i)))
}

func WriteJsonFile(name string, i any) (err error) {
	file, err := os.Create(name)
	if err != nil {
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	enc.Encode(i)
	return
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
	return int64(money * 10000)
}

func Gold2Money(gold int64) float64 {
	return float64(gold) / 10000
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

func GetJsonRawMust(i any) (raw json.RawMessage) {
	buf := lo.Must(json.Marshal(i))

	raw = json.RawMessage(buf)
	return
}

func MKMul(m bson.M, k string, mul float64) {
	if v := m[k]; v != nil {
		fv := GetFloat(v)
		m[k] = fv * mul
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
		MKMul(m1.(bson.M), k2, mul)
	}
}

func MK3Mul(m bson.M, k1, k2, k3 string, mul float64) {
	if m1 := m[k1]; m1 != nil {
		MK2Mul(m1.(bson.M), k2, k3, mul)
	}
}

func IsPGGame(gameID string) bool {
	return strings.HasPrefix(gameID, "pg_")
}

func IsPPGame(gameID string) bool {
	return strings.HasPrefix(gameID, "pp_")
}

type FnFundTransfer = func() bool

func FundTransferWithRetryAndInterval(transfer FnFundTransfer, interval time.Duration, maxcount int) (retry bool) {
	retry = transfer()

	for i := 0; i < maxcount && retry; i++ {
		time.Sleep(interval)
		retry = transfer()
	}

	return retry
}

func HasPrefix(s string, prefixs ...string) bool {
	return slices.ContainsFunc(prefixs, func(p string) bool {
		return strings.HasPrefix(s, p)
	})
}

func HttpGetBody(u string) (body []byte, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	body, err = io.ReadAll(resp.Body)
	return
}

// mkdir dir -p if isnot exist
func Writefile(file string, data []byte) {
	dir := path.Dir(file)
	if e := os.MkdirAll(dir, 0755); e == nil || os.IsExist(e) {
		os.WriteFile(file, data, 0644)
	}
}

func HttpRequestJson(r *http.Request, ptr any) (err error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(payload, ptr)
	return
}

func HttpReturnJson(w http.ResponseWriter, data any) {
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func GetShortUUID() string {
	newUUID := uuid.New()

	// 截取前 8 位作为短 UUID
	shortUUID := newUUID.String()[:8]
	return shortUUID
}
