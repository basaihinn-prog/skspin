package define

import (
	"fmt"
	"net/url"
	"strconv"
)

type PGParams struct {
	Path    string
	TraceId string
	Form    url.Values
	Pid     int64
	GameId  string
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
