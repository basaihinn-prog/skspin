package mq

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

// /////////////////////////////////////////////////
var emptyMsgType = reflect.TypeOf(&nats.Msg{})

func Invoke( /*c *nats.Conn,*/ method string, v interface{}, vPtr interface{}) error {
	// c.Request(method, )
	lo.Must0(defaultNC != nil)
	c := defaultNC

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m, err := c.Request(method, b, 100*time.Second)
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
			return errors.New(ret.Err)
		}
	}
	return nil
}
