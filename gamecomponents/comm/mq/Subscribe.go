package mq

import (
	"game/duck/lazy"
	"game/duck/logger"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

var (
	defaultNC *nats.Conn
	gobNC     *nats.EncodedConn
	JsonNC    *nats.EncodedConn
)

func NC() *nats.Conn {
	return defaultNC
}

func GobNC() *nats.EncodedConn {
	return gobNC
}

// go get github.com/nats-io/nats.go
func ConnectServerMust() (*nats.Conn, error) {
	addr := lazy.GetAddr("proxy.mq")
	nc, err := nats.Connect("nats://" + addr)
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	defaultNC = nc
	gobNC, _ = nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	JsonNC = lo.Must(nats.NewEncodedConn(nc, nats.JSON_ENCODER))
	return nc, nil
}

// 发布一个消息（Gob）
func Publish(subj string, data any) {
	err := gobNC.Publish(subj, data)
	if err != nil {
		logger.Info("Publish错误", err)
	}
}

// 订阅一个消息（Gob）
func Subscribe[T any](subj string, fn func(msg *T)) {
	_, err := gobNC.Subscribe(subj, fn)
	if err != nil {
		logger.Info("Subscribe错误", err)
	}
}

// 使用nats请求回复模式。(Gob)
// func Call(subj string, req, resp any, timeout time.Duration) error {
// 	return gobNC.Request(subj, req, resp, timeout)
// }
