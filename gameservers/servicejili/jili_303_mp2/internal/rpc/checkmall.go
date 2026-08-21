package rpc

import (
	"fmt"

	"serve/servicejili/jili_303_mp2/internal"
	"serve/servicejili/jili_303_mp2/internal/message"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// /csh/mall/checkmallproto

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/mall/checkmallproto", internal.GameShortName), checkmall)
}

func checkmall(ps *nats.Msg) (ret []byte, err error) {
	var req message.Mp2_CheckMallReq
	err = proto.Unmarshal(ps.Data, &req)
	if err != nil {
		return
	}

	var resdata = message.Gaia_GaiaResponse{
		Type: proto.Int32(AckType["buyBonus"]),
		Data: jiliut.ProtoEncode(&message.Mp2_CheckMallAck{
			Settings: []*message.Mp2_GameMallSetting{},
		}),
	}

	ret = jiliut.ProtoEncode(&resdata)
	return
}
