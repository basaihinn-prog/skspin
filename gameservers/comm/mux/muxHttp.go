package mux

import (
	"errors"
	"net/http"
	"reflect"

	"serve/comm/lazy"

	"github.com/samber/lo"
)

func RegHttpWithSample(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	return DefaultRpcMux.RegHttpWithSample(path, desc, kind, handler, ps)
}

func (m *Mux) RegHttpWithSample(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	data := &PHandler{
		Path:         path,
		Handler:      handler,
		Desc:         desc,
		Kind:         kind,
		ParamsSample: ps,
		Class:        "http",
	}

	return m.Add(data)
}

func httpRpcWrapper(h *PHandler) func(w http.ResponseWriter, r *http.Request) {
	handler := h.Handler
	fn := reflect.ValueOf(handler)
	t := reflect.TypeOf(handler)
	// fnLayout := Handler_Layout(t)

	pstypeArgs, pstypeReply := Handler_Args_Reply(t)

	var isHttpReq bool
	if t.NumIn() == 3 {
		isHttpReq = t.In(0) == reflect.TypeOf((*http.Request)(nil))
		lo.Must0(isHttpReq || h.GetArg0 != nil)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			// w.Header().Set("Content-Type", ContentTypeJson)
			return
		}

		logUnit := NewReqLogUnit()
		logUnit.URI = r.RequestURI
		logUnit.Header = r.Header

		defer func() {
			logUnit.Print()
		}()

		// secret := r.Header.Get("secret")
		// if secret != httpSecret {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		paramsValue, err := WrapBindParams(r, pstypeArgs)
		if err != nil {
			HttpReturn2(w, nil, err)
			logUnit.Err = err.Error()
			return
		}
		logUnit.Params = paramsValue.Interface()

		if h.OnlyDev && !lazy.CommCfg().IsDev {
			err = errors.New("仅开发环境可以调用此方法!")
			HttpReturn2(w, nil, err)
			logUnit.Err = err.Error()
			return
		}

		var in []reflect.Value

		switch t.NumIn() {
		case 2:

		case 3:
			if isHttpReq {
				in = append(in, reflect.ValueOf(r))
			} else {
				arg0, err := h.GetArg0(r)
				if err != nil {
					HttpReturn2(w, nil, err)
					logUnit.Err = err.Error()
					return
				}

				in = append(in, reflect.ValueOf(arg0))
			}
		}

		in = append(in, paramsValue)

		paramsReply := reflect.New(pstypeReply.Elem())
		in = append(in, paramsReply)

		out := fn.Call(in)

		if !out[0].IsNil() {
			err = out[0].Interface().(error)
		}

		if err != nil {
			logUnit.Err = err.Error()
		}

		reply := paramsReply.Interface()
		logUnit.Result = reply
		HttpReturn2(w, reply, err)
	}
}
