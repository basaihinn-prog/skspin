package rpc

import (
	"errors"
	"fmt"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_44_fivestar/internal"
	"serve/servicejili/jili_44_fivestar/internal/message"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/req", internal.GameShortName), req)
}

var reqMux = map[int32]func(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error){}

func req(ps *nats.Msg) (ret []byte, err error) {
	// ps.Header
	var req message.Server_Request

	err = proto.Unmarshal(ps.Data, &req)
	if err != nil {
		return
	}

	token := ps.Header.Get("Token")
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	//slog.Info("req", "pid", pid, "token", token, "ack", req.GetAck())

	var (
		retData protoreflect.ProtoMessage
	)

	fn, ok := reqMux[req.GetAck()]
	if !ok {
		err = errors.New("error method")
		return
	}

	retData, err = fn(pid, req.GetReq(), ps)
	if err != nil {
		return
	}

	var resp = &message.Server_Response{
		Ack:      proto.Int32(req.GetAck()),
		Response: jiliut.ProtoEncode(retData),
	}

	ret = jiliut.ProtoEncode(resp)
	return
}
