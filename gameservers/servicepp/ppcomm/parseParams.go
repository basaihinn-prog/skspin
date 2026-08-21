package ppcomm

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Variables map[string]string

// GetGameSt 根据rsp字段判断这次是什么状态，gameSt 1.bg未中奖，2.bg连转中，3.bg结束，4.fg连转中，5fg连转结束 0 doCollect
func (rspMap Variables) GetGameSt() int {
	gameSt := 0
	fsmul := rspMap["fsmul"]
	if rspMap["tw"] == "" && rspMap["na"] == "s" {
		gameSt = 0
		return gameSt
	}
	if rspMap["tw"] == "" {
		gameSt = 1
		return gameSt
	}
	if rspMap.Float("tw") == 0 && rspMap.Float("w") == 0 && fsmul == "" {
		gameSt = 1
	} else if rspMap.Float("tw") > 0 && rspMap.Float("w") > 0 && rspMap["na"] == "s" && fsmul == "" {
		gameSt = 2
	} else if rspMap.Float("tw") > 0 && rspMap["na"] == "c" && fsmul == "" {
		gameSt = 3
	} else if rspMap.Float("tw") > 0 && rspMap.Int("fsend_total") == 1 {
		gameSt = 5
	} else {
		gameSt = 4
	}
	return gameSt
}

func (ps Variables) Int(k string) int {
	v := ps[k]
	ret, _ := strconv.Atoi(v)
	return ret
}
func (ps Variables) SetInt(k string, v int) {
	s := strconv.Itoa(v)
	ps[k] = s
}

func (ps Variables) Float(k string) float64 {
	v := ps[k]
	ret, _ := strconv.ParseFloat(v, 64)
	return ret
}
func (ps Variables) SetFloat(k string, v float64) {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	ps[k] = s
}

func (ps Variables) MKMulFloat(k string, mul float64) {
	_, ok := ps[k]
	if !ok {
		return
	}
	v := ps.Currency(k)
	ps.SetCurrency(k, v*mul)
}

func (ps Variables) SetFloatArr(k string, v []float64) {
	var str strings.Builder
	for i := 0; i < len(v); i++ {
		if i != 0 {
			str.WriteByte(',')
		}
		str.WriteString(fmt.Sprintf("%.2f", v[i]))
	}
	ps[k] = str.String()
}

func (ps Variables) Currency(k string) float64 {
	v := ps[k]
	v = strings.ReplaceAll(v, ",", "")
	ret, _ := strconv.ParseFloat(v, 64)
	return ret
}
func (ps Variables) SetCurrency(k string, v float64) {
	p := message.NewPrinter(language.English)
	s := p.Sprintf("%.2f", v)
	ps[k] = s
}

func (ps Variables) Str(k string) string {
	v := ps[k]
	return v
}

func (ps Variables) Get(k string) string {
	v := ps[k]
	return v
}

func (ps Variables) Delete(k string) {
	_, ok := ps[k]
	if ok {
		delete(ps, k)
	}
}

func (ps Variables) Exist(k string) bool {
	_, ok := ps[k]
	return ok
}

func (ps Variables) SetStr(k string, v string) {
	ps[k] = v
}
func (ps Variables) Set(k string, v string) {
	ps[k] = v
}

func (ps Variables) JsonUnmarshal(k string, recv any) error {
	v := ps[k]
	return json.Unmarshal([]byte(v), recv)
}
func (ps Variables) SetJson(k string, v any) error {
	s, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ps[k] = string(s)
	return nil
}

func (ps Variables) Encode() string {
	keys := lo.Keys(ps)
	sort.Strings(keys)
	var sb strings.Builder
	cap := 0
	for k, v := range ps {
		cap += len(k) + len(v) + 2
	}

	sb.Grow(cap)
	for i, k := range keys {
		if i != 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(ps[k])
	}
	return sb.String()
}

func (ps Variables) PrettyString() string {
	keys := lo.Keys(ps)
	sort.Strings(keys)
	var sb strings.Builder
	cap := 0
	for k, v := range ps {
		cap += len(k) + len(v) + 2
	}

	sb.Grow(cap)
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(ps[k])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (ps Variables) Bytes() []byte {
	//处理rid
	ps.Delete("rid")
	return []byte(ps.Encode())
}

func ParseVariables(s string) Variables {
	/* build.js
	var variables = {};
	var pairs = String(data).split("&");
	for (var i = 0; i < pairs.length; ++i) {
	    var pair = pairs[i].split("=");
	    variables[pair[0]] = pair[1] || ""
	}
	*/

	var variables = Variables{}
	pairs := strings.Split(s, "&")
	for _, v := range pairs {
		pair := strings.SplitN(v, "=", 2)
		k, v := pair[0], ""
		if len(pair) > 1 {
			v = pair[1]
		}
		variables[k] = v
	}

	return variables
}

func PrettyPrintVas(w io.Writer, s string) {
	io.WriteString(w, ParseVariables(s).PrettyString())
	io.WriteString(w, "\n===\n\n")
}

func (ps Variables) MKMulmo_twInG(k string, mul float64) {
	g, ok := ps["g"]
	if !ok {
		return
	}
	// 正则表达式，匹配 mo_tw:"1.60"
	re := regexp.MustCompile(`mo_tw:"\d+\.\d+"`)
	mo_tw := re.FindString(g)
	if mo_tw == "" {
		return
	}
	split := strings.Split(mo_tw, ":")
	if len(split) != 2 {
		return
	}
	split[1] = MKMulFloat(split[1], mul)
	g = re.ReplaceAllString(g, fmt.Sprintf(`mo_tw:"%v"`, split[1]))

	//temp := make(map[string]interface{})
	//err := json.Unmarshal([]byte(g), &temp)
	//if err != nil {
	//}
	//for key := range temp {
	//	if strings.Contains(key, "collect") { //发现拥有该字段，直接用正则找g下对应的字符串修改
	//		marshal, _ := json.Marshal(temp[key])
	//		collect := make(map[string]interface{})
	//		json.Unmarshal(marshal, &collect)
	//		if _, ok := collect[k]; ok {
	//			oMotw := collect[k].(string)
	//			oMotw = MKMulFloat(oMotw, mul)
	//			// 正则表达式，匹配 mo_tw:"1.60"
	//			re := regexp.MustCompile(`mo_tw:"\d+\.\d+"`)
	//			g = re.ReplaceAllString(g, fmt.Sprintf(`mo_tw:"%v"`, oMotw))
	//			break
	//		} else {
	//			return
	//		}
	//	}
	//}
	ps["g"] = g
}

func MKMulFloat(k string, mul float64) string {
	k = strings.ReplaceAll(k, ",", "")
	k = strings.ReplaceAll(k, `"`, "")
	ret, _ := strconv.ParseFloat(k, 64)
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", ret*mul)
}
