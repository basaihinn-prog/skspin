package main

import (
	"context"
	"game/duck/etcd"
	"game/duck/exit"
	"game/duck/logger"
	"game/duck/rpc1"
	"game/duck/rpc1/discovery"
	"game/duck/rpc1/test/pb/pb"
)

type MailServer struct {
}

func (MailServer) SendMail(context.Context, *pb.SendMailReq) (*pb.SendMailResp, error) {
	logger.Info("123123123")
	return &pb.SendMailResp{}, nil
}

func main() {
	e, err := etcd.NewEtcd("")
	if err != nil {
		logger.Err(err)
		return
	}
	exit.Close("", e)

	d := discovery.NewEtcdRegiser(e)
	p := discovery.NewSystemPortProvider()

	server := rpc1.NewServer(&rpc1.ServerOption{Name: "server"}, p, d)
	exit.Close("", server)

	pb.RegisterMailServer(server, &MailServer{})

	err = server.Serve()
	if err != nil {
		logger.Err(err)
	}

	exit.Exit()
}
