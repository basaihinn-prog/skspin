package main

import (
	"context"
	"game/duck/etcd"
	"game/duck/exit"
	"game/duck/logger"
	"game/duck/rpc1"
	"game/duck/rpc1/discovery"
	"game/duck/rpc1/test/pb/pb"
	"time"
)

func main() {
	e, err := etcd.NewEtcd("")
	if err != nil {
		panic(err)
	}

	client := rpc1.NewClientManager(discovery.NewEtcdDiscovery(e))
	exit.Close("grpc client", client)

	ch := make(chan struct{}, 1)

	exit.Callback("", func() {
		ch <- struct{}{}
	})

	t := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ch:
			goto EXIT

		case <-t.C:
			mailClient, err := client.GetClient("server", rpc1.RoundRobin())
			if err != nil {
				panic(err)
			}

			rpc := pb.NewMailClient(mailClient)

			resp, err := rpc.SendMail(context.TODO(), &pb.SendMailReq{})
			if err != nil {
				logger.Info(err)
			} else {
				logger.Info(resp)
			}
		}
	}
EXIT:
	exit.Exit()
}
