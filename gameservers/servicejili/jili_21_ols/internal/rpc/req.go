package rpc

import (
	"errors"
	"fmt"
	"serve/servicejili/jili_21_ols/internal/models"

	"serve/comm/jwtutil"
	"serve/servicejili/jili_21_ols/internal"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/jiliUtMessage"

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
	var req jiliUtMessage.Server_Request

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

	var retData protoreflect.ProtoMessage

	fn, ok := reqMux[req.GetAck()]
	if !ok {
		err = errors.New("error method")
		return
	}

	retData, err = fn(pid, req.GetReq(), ps)
	if err != nil {
		return
	}
	data, err := jiliut.ProtoEncryption(token, retData)
	if err != nil {
		return
	}

	var resp = &jiliUtMessage.Server_GaiaResponse{
		Type: proto.Int32(req.GetAck()),
		Data: data,
	}

	ret = jiliut.ProtoEncode(resp)
	return
}

var reqMuxRtp = map[int32]func(plr *models.Player, pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error){}

func reqMuxRtpFunc(ps *nats.Msg) (ret []byte, err error) {
	// ps.Header
	var req jiliUtMessage.Server_Request

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

	var retData protoreflect.ProtoMessage

	fn, ok := reqMux[req.GetAck()]
	if !ok {
		err = errors.New("error method")
		return
	}

	retData, err = fn(pid, req.GetReq(), ps)
	if err != nil {
		return
	}

	data, err := jiliut.ProtoEncryption(token, retData)
	if err != nil {
		return
	}

	var resp = &jiliUtMessage.Server_GaiaResponse{
		Type: proto.Int32(req.GetAck()),
		Data: data,
	}

	ret = jiliut.ProtoEncode(resp)
	return
}
