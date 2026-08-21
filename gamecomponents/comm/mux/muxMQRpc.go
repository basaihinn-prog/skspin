package mux

import (
	"encoding/json"
	"fmt"
	"game/duck/lazy"
	"log"
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func validRpcPorc(handler interface{}) {
	t := reflect.TypeOf(handler)

	Assert(t.NumIn() == RpcParamNum)

	pstypeReply := t.In(RpcParamReply)
	Assert(pstypeReply.Kind() == reflect.Ptr)

	Assert(t.NumOut() == 1)
	u := reflect.TypeOf((*error)(nil)).Elem()
	Assert(t.Out(0).Implements(u))
}

func (m *Mux) RegRpc(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	validRpcPorc(handler)

	data := &PHandler{
		Path:         path,
		Handler:      handler,
		Desc:         desc,
		Kind:         kind,
		ParamsSample: ps,
		Class:        "rpc",
	}

	return m.Add(data)
}

func RegRpc(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	return DefaultRpcMux.RegRpc(path, desc, kind, handler, ps)
}

func RegistRpcToMQ(c *nats.Conn) {
	DefaultRpcMux.RegistRpcToMQ(c)
}

func (m *Mux) RegistRpcToMQ(c *nats.Conn) {
	for _, h := range m.ToArr() {
		if h.Class == "rpc" {
			// c.QueueSubscribe()
			// lo.Must(c.Subscribe(h.Path, rpcWrapper(c, h.Handler, h.Path, h.OnlyDev)))
			lo.Must(c.QueueSubscribe(h.Path, "queue", rpcWrapper(c, h.Handler, h.Path, h.OnlyDev)))

			slog.Info("mq QueueSubscribe",
				"path", h.Path,
				"kind", h.Kind,
				"desc", h.Desc,
				"class", h.Class,
			)
		}
	}
}

func rpcWrapper(c *nats.Conn, handler any, path string, onlyDev bool) func(*nats.Msg) {
	fn := reflect.ValueOf(handler)
	t := reflect.TypeOf(handler)
	f := func(msg *nats.Msg) {
		logUnit := NewReqLogUnit()
		logUnit.URI = path
		logUnit.ReqBody = string(msg.Data)

		var ret struct {
			Resp any
			Err  string
		}

		defer func() {
			if ret.Err != "" || lazy.CommCfg().IsDev {
				logUnit.Err = ret.Err
				logUnit.Result = ret.Resp
				logUnit.Print()
			}
		}()

		defer func() {
			data := lo.Must(json.Marshal(ret))
			c.Publish(msg.Reply, data)
		}()

		defer func() {
			if x := recover(); x != nil {
				os.Stdout.Write(debug.Stack())
				log.Println(x)

				// result = fmt.Errorf("panic %v", x)
				// logUnit.Err = fmt.Sprintf("panic %v", x)
				ret.Err = fmt.Sprintf("panic %v", x)
			}
		}()

		pstypeArgs := t.In(RpcParamArgs)
		paramsValue, err := bindwrap(pstypeArgs, func(i any) error {
			return json.Unmarshal(msg.Data, i)
		})
		if err != nil {
			logUnit.Params = json.RawMessage(msg.Data)
			ret.Err = err.Error()
			return
		}
		logUnit.Params = paramsValue.Interface()

		if onlyDev && !lazy.CommCfg().IsDev {
			ret.Err = "仅开发环境可以调用此方法!"
			return
		}

		pstypeReply := t.In(RpcParamReply)
		paramsReply := reflect.New(pstypeReply.Elem())

		in := []reflect.Value{paramsValue, paramsReply}
		out := fn.Call(in)

		if !out[0].IsNil() {
			err := out[0].Interface().(error)
			ret.Err = err.Error()
			return
		}
		ret.Resp = paramsReply.Interface()
	}
	return func(msg *nats.Msg) {
		go f(msg)
	}
}
