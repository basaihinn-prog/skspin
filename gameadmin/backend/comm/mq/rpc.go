package mq

import (
	"encoding/json"
	"errors"
	"game/comm/define"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

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
