package mux

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type ApiUnit struct {
	Path    string
	Params  string
	Desc    string
	Kind    string
	Class   string
	OnlyDev bool
}

func marshal(t reflect.Type) string {
	v := reflect.New(t)
	i := v.Interface()
	buf, _ := json.Marshal(i)
	return string(buf)
}

func (m *Mux) list_api(w http.ResponseWriter, r *http.Request) {
	apiarr := m.ToArr()
	arr := make([]ApiUnit, len(apiarr))
	for i, api := range apiarr {
		var u ApiUnit
		ShallowCopy(&u, api)
		t := reflect.TypeOf(api.Handler)

		if api.ParamsSample != nil {
			json, _ := json.Marshal(api.ParamsSample)
			u.Params = string(json)
		} else {
			pstypeArgs, _ := Handler_Args_Reply(t)
			u.Params = marshal(pstypeArgs)
		}

		arr[i] = u
	}

	HttpReturn2(w, arr, nil)
}
