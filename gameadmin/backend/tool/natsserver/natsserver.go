package main

import (
	"errors"
	"game/comm/mq"
	"game/duck/lazy"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

type addPs struct {
	X, Y int
}
type addRet struct {
	Sum int
}

type rpc struct{}

func (rpc) Add(ps addPs, ret *addRet) (err error) {
	slog.Info("rpc Add", "ps", ps)
	ret.Sum = ps.X + ps.Y
	return nil
}

func (rpc) Div(ps *[2]float64, ret *float64) (err error) {
	slog.Info("rpc Div", "ps", ps)
	if ps[1] == 0 {
		err = errors.New("the ps y is 0")
		return
	}

	*ret = ps[0] / ps[1]
	return nil
}

func main() {
	lazy.Init("natsserver")
	mq.ConnectServerMust()

	mq.NC().Subscribe("test.hello", func(msg *nats.Msg) {
		msg.Respond([]byte("hello from" + time.Now().Format(time.RFC3339Nano)))
	})
	// mq.RegObj(lazy.ServiceName, rpc{})
	// mq.RegObj(lazy.ServiceName, rpc{})
	// mq.RegObj(lazy.ServiceName, rpc{})
	// mq.RegObj(lazy.ServiceName, rpc{})
	select {}
}
