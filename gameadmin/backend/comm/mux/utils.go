package mux

import (
	"encoding/json"
	"game/comm/define"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const ContentTypeJson = "application/json"

func HttpParseArgs(r *http.Request, p interface{}) error {
	if r.Method == "POST" {
		return json.NewDecoder(r.Body).Decode(p)
	} else {
		return HttpBindQuery(p, r)
	}
}

func HttpReturn2(w http.ResponseWriter, result interface{}, err error) {
	var ret Response

	if err != nil {
		ret.Error = err.Error()
		if ec, ok := err.(define.IErrcode); ok {
			ret.Code = ec.Code()
			// ret.Error = strings.
			i := strings.IndexByte(ret.Error, '#')
			if i != -1 {
				// 兼容
				ret.Error = ret.Error[i+1:]
			}
		}
		if ret.Code == 0 {
			ret.Code = 1
		}
	} else {
		ret.Data = result
	}

	buf, _ := json.Marshal(ret)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", ContentTypeJson)
	w.Write(buf)
}

const (
	HttpRpcParamRequest = iota
	HttpRpcParamArgs
	HttpRpcParamReply
	HttpRpcParamNum
)

const (
	RpcParamArgs = iota
	RpcParamReply
	RpcParamNum
)

// var (
// 	httpSecret string
// )

/*
func Handler_Layout(t reflect.Type) string {
	// t := reflect.TypeOf(h)
	numIn := t.NumIn()
	Assert(numIn == RpcParamNum || numIn == HttpRpcParamNum)

	u := reflect.TypeOf((*error)(nil)).Elem()
	Assert(t.Out(0).Implements(u))

	var fnLayout string

	if numIn == RpcParamNum {
		fnLayout = "rpc"
	} else {
		ps1type := t.In(HttpRpcParamRequest)
		isHttpReq := ps1type == reflect.TypeOf((*http.Request)(nil))
		// Assert(isHttpReq || ps1type.Kind() == reflect.Int64)

		if isHttpReq {
			fnLayout = "http"
		} else {
			fnLayout = "player"
		}

	}

	return fnLayout
}
*/

func Handler_Args_Reply(t reflect.Type) (pstypeArgs, pstypeReply reflect.Type) {
	numIn := t.NumIn()

	if numIn == 2 {
		pstypeArgs = t.In(RpcParamArgs)
		pstypeReply = t.In(RpcParamReply)
	} else if numIn == 3 {
		pstypeArgs = t.In(HttpRpcParamArgs)
		pstypeReply = t.In(HttpRpcParamReply)
	} else {
		log.Panicf("error numIn: %d", numIn)
	}

	return
}
func WrapBindParams(r *http.Request, paramsType reflect.Type) (v reflect.Value, err error) {
	if IsEmptyParamsType(paramsType) {
		v = EmptyParamsValue
		return
	}
	switch r.Method {
	case http.MethodPost:
		// Content-Type: application/x-www-form-urlencoded

		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {

			return bindwrap(paramsType, func(i interface{}) error {
				err := r.ParseForm()
				if err != nil {
					return err
				}

				return BindData(i, r.Form, "query")
			})
		}
		return bindwrap(paramsType, func(i interface{}) error {
			return json.NewDecoder(r.Body).Decode(i)
		})
	default:
		return bindwrap(paramsType, func(i interface{}) error {
			return HttpBindQuery(i, r)
		})
	}

}

func bindwrap(paramsType reflect.Type, fn func(i interface{}) error) (v reflect.Value, err error) {
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

func Bindwrap(paramsType reflect.Type, fn func(i interface{}) error) (v reflect.Value, err error) {
	return bindwrap(paramsType, fn)
}

func Assert(b bool, msg ...interface{}) {
	if !b {
		// debug.PrintStack()
		log.Panicln(msg...)
	}
}

func Indirect(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
