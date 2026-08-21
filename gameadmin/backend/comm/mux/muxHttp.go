package mux

import (
	"bytes"
	"errors"
	"game/duck/lazy"
	"io"
	"net/http"
	"reflect"

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
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"https://admin.ppslot.site":    true,
			"https://ppslot.site":          true,
			"https://h5.ppslot.site":       true,
			"https://api.cocodr.xyz":       true,
			"https://history.cocodr.xyz":   true,
			"https://pg.bxbet.asia":        true,
			"https://m-pg.bxbet.asia":      true,
			"https://static-pg.bxbet.asia": true,
			"http://localhost:4568":     true,
			"http://localhost:3000":     true,
		}
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://admin.ppslot.site")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			HttpReturn2(w, nil, err)
			logUnit.Err = err.Error()
			return
		}

		logUnit.ReqBody = string(body)

		r.Body = io.NopCloser(bytes.NewReader(body))

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
