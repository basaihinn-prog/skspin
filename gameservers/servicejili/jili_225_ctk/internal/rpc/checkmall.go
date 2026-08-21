package rpc

import (
	"fmt"

	"serve/servicejili/jili_225_ctk/internal"
	"serve/servicejili/jili_225_ctk/internal/message"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// /csh/mall/checkmallproto

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/mall/checkmallproto", internal.GameShortName), checkmall)
}

func checkmall(ps *nats.Msg) (ret []byte, err error) {
	var req message.Gaia_CheckMallReq
	err = proto.Unmarshal(ps.Data, &req)
	if err != nil {
		return
	}

	var resdata = message.Gaia_GaiaResponse{
		/*Type:  proto.Int32(AckType["buyBonus"]),
		Token: req.Token,
		Data: jiliut.ProtoEncode(&message.Gaia_CheckMallAck{
			Show:     proto.Int32(2),
			Settings: []*message.Gaia_GameMallSetting{},
		}),*/
	}

	ret = jiliut.ProtoEncode(&resdata)
	return
}
