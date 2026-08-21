package mq

import (
	"context"
	"encoding/json"
	"errors"
	"game/duck/lazy"
	"game/duck/rpc1"
	"time"

	"game/duck/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcReq struct {
	Service string
	Method  string
	Meta    metadata.MD
	Data    []byte // json
}

type grpcResp struct {
	Data  []byte // json
	Error string
}

func getDesc(service, method string) (*rpc1.DescAndImpl, *grpc.MethodDesc) {

	sdesc, ok := lazy.GrpcServer.ServiceDescMap.Load(service)
	if !ok {
		return nil, nil
	}

	for _, v := range sdesc.Desc.Methods {
		if v.MethodName == method {
			return sdesc, &v
		}
	}

	return nil, nil
}

// 通过mq调用别的服务器的grpc
func CallGrpc(svr, service, method string, meta metadata.MD, reqBuf []byte) ([]byte, error) {
	req := grpcReq{
		Service: service,
		Method:  method,
		Meta:    meta,
		Data:    reqBuf,
	}

	resp := grpcResp{}
	err := gobNC.Request("rpc."+svr, &req, &resp, 5*time.Second)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.Data, nil
}

// 接收来自mq的grpc
func RecvGrpcFromMQ() {
	_, err := gobNC.Subscribe("rpc."+lazy.ServiceName, func(subject, reply string, req *grpcReq) {

		serviceImpl, methodImpl := getDesc(req.Service, req.Method)
		if serviceImpl == nil {
			gobNC.Publish(reply, &grpcResp{Error: "没找到对应的grpc"})
			return
		}

		df := func(v interface{}) error {
			return rpc1.JsonBytesCodec{}.Unmarshal(req.Data, v)
		}

		ctx := metadata.NewIncomingContext(context.TODO(), req.Meta)

		resp, err := methodImpl.Handler(serviceImpl.Impl, ctx, df, rpc1.DefaultUnaryInterceptor)

		if err != nil {
			gobNC.Publish(reply, &grpcResp{Error: err.Error()})
		} else {
			buf, err := json.Marshal(resp)
			if err != nil {
				logger.Info("错误了", err)
				gobNC.Publish(reply, &grpcResp{Error: err.Error()})
				return
			}

			gobNC.Publish(reply, &grpcResp{Data: buf})
		}
	})
	if err != nil {
		logger.Info("订阅Grpc消息错误", err)
	}
}
