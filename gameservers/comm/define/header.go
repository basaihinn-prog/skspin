package define

import (
	"net/url"
	"strconv"
)

type Values map[string][]string

func (h Values) Add(key, value string) {
	h[key] = append(h[key], value)
}

func (h Values) Set(key, value string) {
	h[key] = []string{value}
}

const (
	_EMPTY_ = ""
)

func (h Values) Get(key string) string {
	if h == nil {
		return _EMPTY_
	}
	if v := h[key]; v != nil {
		return v[0]
	}
	return _EMPTY_
}

func (h Values) Values(key string) []string {
	return h[key]
}

func (h Values) Del(key string) {
	delete(h, key)
}

func (h Values) GetInt(key string) (ans int) {
	v := h.Get(key)
	if v == "" {
		return
	}

	ans, _ = strconv.Atoi(v)
	return
}

func (h Values) GetFloat(key string) (ans float64) {
	v := h.Get(key)
	if v == "" {
		return
	}

	ans, _ = strconv.ParseFloat(v, 64)
	return
}

func (h Values) Has(key string) bool {
	_, ok := h[key]
	return ok
}

func (v Values) Encode() string {
	return url.Values(v).Encode()
}

func ParseQuery(query string) (Values, error) {
	vals, err := url.ParseQuery(query)
	if err != nil {
		return Values{}, err
	}

	return Values(vals), nil
}
