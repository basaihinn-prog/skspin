package mq

import (
	"encoding/json"
	"errors"
	"game/comm/define"
	"log"
	"log/slog"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

const (
	RpcParamArgs = iota + 1
	RpcParamReply
	RpcParamNum
)

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

var e = reflect.TypeOf((*error)(nil)).Elem()

func validMethod(t reflect.Type) bool {
	pstypeReply := t.In(RpcParamReply)

	ok := t.NumIn() == RpcParamNum &&
		t.NumOut() == 1 &&
		pstypeReply.Kind() == reflect.Pointer &&
		t.Out(0).Implements(e)

	return ok
}

func RegObj(ns string, rcvr any) {
	lo.Must0(defaultNC != nil)

	c := defaultNC

	v := reflect.ValueOf(rcvr)

	typ := v.Type()
	nmethod := typ.NumMethod()

	sname := reflect.Indirect(v).Type().Name()

	if sname == "" {
		s := "rpc.Register: no service name for type " + typ.String()
		log.Panic(s)
	}

	ns = ns + "." + sname

	for i := 0; i < nmethod; i++ {
		method := typ.Method(i)
		methodv := v.Method(i)
		t := method.Type

		if !validMethod(t) {
			continue
		}

		path := ns + "." + method.Name

		_, err := c.QueueSubscribe(path, "queue", func(msg *nats.Msg) {
			var ret struct {
				Resp any
				Err  string
			}
			defer func() {
				data := lo.Must(json.Marshal(ret))
				c.Publish(msg.Reply, data)
			}()

			pstypeArgs := t.In(RpcParamArgs)
			paramsValue, err := bindwrap(pstypeArgs, func(i any) error {
				return json.Unmarshal(msg.Data, i)
			})

			pstypeReply := t.In(RpcParamReply)
			paramsReply := reflect.New(pstypeReply.Elem())

			if err != nil {
				ret.Err = err.Error()
				return
			}

			in := []reflect.Value{paramsValue, paramsReply}
			out := methodv.Call(in)

			if !out[0].IsNil() {
				err := out[0].Interface().(error)
				ret.Err = err.Error()
				return
			}
			ret.Resp = paramsReply.Interface()
		})

		slog.Info("Subscribe",
			"numin", t.NumIn(),
			"path", path,
			"error", err,
		)
	}
}

// /////////////////////////////////////////////////
var emptyMsgType = reflect.TypeOf(&nats.Msg{})

func Invoke(method string, v interface{}, vPtr interface{}) error {
	// c.Request(method, )
	lo.Must0(defaultNC != nil)
	c := defaultNC

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m, err := c.Request(method, b, 60*time.Second)
	if err != nil {
		return err
	}
	if reflect.TypeOf(vPtr) == emptyMsgType {
		mPtr := vPtr.(*nats.Msg)
		*mPtr = *m
	} else {
		var ret struct {
			Resp any
			Err  string
		}
		ret.Resp = vPtr
		err = json.Unmarshal(m.Data, &ret)
		if err != nil {
			return err
		}

		if ret.Err != "" {
			if err, ok := define.NewErrCodeFromStr(ret.Err); ok {
				return err
			}
			return errors.New(ret.Err)
		}
	}
	return nil
}
