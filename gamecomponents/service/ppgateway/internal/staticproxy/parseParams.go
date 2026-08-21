package staticproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Variables map[string]string

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
	v, err := strconv.ParseFloat(ps[k], 64)
	if err != nil {
		panic(fmt.Sprintf("cannot parse to float, k:%s", k))
	}
	ps[k] = fmt.Sprintf("%.2f", v*mul)
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
