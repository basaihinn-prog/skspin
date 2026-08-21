package define

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type PGParams struct {
	Path     string
	TraceId  string
	Form     url.Values
	Pid      int64
	Header   map[string]any
	BuyEX    int
	BuyBonus bool
	Count    int
	GameId   string
	GameName string
	UName    string
}

// func (ps *PGParams) GetPid() int64 {
// 	return ps.Pid
// }

func (ps PGParams) GetPid() int64 {
	return ps.Pid
}

func (ps *PGParams) Get(key string) string {
	return ps.Form.Get(key)
}

func (ps *PGParams) GetInt(key string) (ans int) {
	v := ps.Form.Get(key)
	if v == "" {
		return
	}

	ans, _ = strconv.Atoi(v)
	return
}

func (ps *PGParams) GetFloat(key string) (ans float64) {
	v := ps.Form.Get(key)
	if v == "" {
		return
	}

	ans, _ = strconv.ParseFloat(v, 64)
	return
}

type M = map[string]any
type D = []any

func M2Values(m M) url.Values {
	ans := make(url.Values, len(m))
	for k, v := range m {
		ans.Set(k, fmt.Sprint(v))
	}
	return ans
}

type PGError struct {
	Cd  string `json:"cd"`
	Msg string `json:"msg"`
	Tid string `json:"tid"`
}

func (e *PGError) Error() string {
	return e.Msg
}

type PGRetWrapper struct {
	Dt  *json.RawMessage `json:"dt"`
	Err *PGError         `json:"err"`
}
