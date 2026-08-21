package jdbcomm

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"log"
	"os"
	"runtime/debug"
	"serve/comm/lazy"
)

type RpcFunc func(*nats.Msg) ([]byte, error)

var rpcMux = map[string]RpcFunc{}

func RegRpc(path string, fn RpcFunc) {
	lo.Must0(rpcMux[path] == nil)
	rpcMux[path] = fn
}

func RegistRpcToMQ(c *nats.Conn) {
	for pth, h := range rpcMux {
		subj := fmt.Sprintf("%s.%s", lazy.ServiceName, pth)
		fmt.Println("RegistRpcToMQ", subj)
		c.QueueSubscribe(subj, "queue", rpcWrapper(h, c))
	}
}

func rpcWrapper(h RpcFunc, c *nats.Conn) func(*nats.Msg) {
	var f = func(msg *nats.Msg) {
		defer func() {
			if x := recover(); x != nil {
				os.Stdout.Write(debug.Stack())
				log.Println(x)

				// result = fmt.Errorf("panic %v", x)
				// logUnit.Err = fmt.Sprintf("panic %v", x)
				// ret.Err = fmt.Sprintf("panic %v", x)
			}
		}()
		data, err := h(msg)
		respM := &nats.Msg{
			Subject: msg.Reply,
			Data:    data,
			Header:  nats.Header{},
		}

		if err != nil {
			respM.Header.Set("error", err.Error())
		}

		c.PublishMsg(respM)
	}

	return func(m *nats.Msg) {
		go f(m)
	}
}
