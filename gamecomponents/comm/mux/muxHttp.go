package mux

import (
	"bytes"
	"errors"
	"fmt"
	"game/comm/mq"
	"game/duck/lazy"
	"game/duck/ut2/httputil"
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

const Host = "127.0.0.1"

func checkIp(ip string) (allow bool, err error) {
	var netAllow struct {
		Allow bool
	}
	fmt.Println(ip)
	if ip == Host {
		fmt.Println(ip, "进入")
		netAllow.Allow = true
		return netAllow.Allow, nil
	}
	err = mq.Invoke("/AdminInfo/Interior/whiteList_exists", map[string]any{
		"IP": ip,
	}, &netAllow)
	fmt.Println("invoke:", ip)
	if err != nil {
		return
	}
	allow = netAllow.Allow
	return
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
		ip := httputil.GetIPFromRequest(r)
		//loc := string(iploc.Country(net.ParseIP(ip)))
		netAllow, err := checkIp(ip)
		fmt.Println("httpRpcWrapper", "netAllow", netAllow, "err", err, "ip", ip)
		if err != nil {
			http.Error(w, "{\"code\":6009,\"error\":\"checkIp error\",\"data\":null}", http.StatusForbidden)
			return
		}
		// 进行 IP 检查
		if !netAllow {
			http.Error(w, "{\"code\":6010,\"error\":\"Access denied\",\"data\":null}", http.StatusForbidden)
			return
		}
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
