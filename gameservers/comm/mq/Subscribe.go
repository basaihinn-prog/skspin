package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

var (
	defaultNC *nats.Conn
	// gobNC     *nats.EncodedConn
	JsonNC *nats.EncodedConn
)

// func GobNC() *nats.EncodedConn {
// 	return gobNC
// }

// func JsonNC() *nats.EncodedConn {
// 	return jsonNC
// }

func NC() *nats.Conn {
	return defaultNC
}

func PublishMsg(subj string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = defaultNC.Publish(subj, data)
	if err != nil {
		return err
	}
	// 添加Flush()以确保消息发送成功
	if err := defaultNC.Flush(); err != nil {
		return fmt.Errorf("failed to flush messages: %w", err)
	}
	return nil
}

func Subscribe[T any](subj string, fn func(something *T)) {
	// 订阅主题
	defaultNC.Subscribe(subj, func(msg *nats.Msg) {
		// 从字节数组中重新转换为结构体
		something := new(T)
		err := json.Unmarshal(msg.Data, &something)
		if err != nil {
			log.Fatalf("Failed to unmarshal message: %s", err)
		}
		fn(something)
	})
}

// go get github.com/nats-io/nats.go
func ConnectServerMust(addr string) (*nats.Conn, error) {
	// DO NOT DELETE THIS GC
	runtime.GC()

	// addr := lazy.GetAddr("proxy.mq")
	nc, err := nats.Connect("nats://" + addr)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}
	defaultNC = nc
	// gobNC, _ = nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	JsonNC = lo.Must(nats.NewEncodedConn(nc, nats.JSON_ENCODER))
	return nc, nil
}

// 发布一个消息（Gob）
// func Publish(subj string, data any) {
// 	err := gobNC.Publish(subj, data)
// 	if err != nil {
// 		slog.Info("Publish错误", "err", err)
// 	}
// }

// 订阅一个消息（Gob）
// func Subscribe[T any](subj string, fn func(msg *T)) {
// 	_, err := gobNC.Subscribe(subj, fn)
// 	if err != nil {
// 		slog.Info("Subscribe错误", "err", err)
// 	}
// }

// 使用nats请求回复模式。(Gob)
// func Call(subj string, req, resp any, timeout time.Duration) error {
// 	return gobNC.Request(subj, req, resp, timeout)
// }
