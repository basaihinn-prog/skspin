package rpc

import (
	"fmt"

	"serve/servicejili/jili_176_cny/internal"
	"serve/servicejili/jili_176_cny/internal/message"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// /fd/mall/checkmallproto

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
		Type:  proto.Int32(AckType["buyBonus"]),
		Token: req.Token,
		Data: jiliut.ProtoEncode(&message.Gaia_CheckMallAck{
			Show: proto.Int32(2),
			Settings: []*message.Gaia_GameMallSetting{
				{
					APIID:    req.Apiid,
					AlterID:  proto.Int32(50),
					Event:    proto.Bool(true),
					ForSale:  proto.Bool(true),
					GameID:   proto.Int32(internal.GameNo),
					MaxBet:   proto.Float64(1000),
					PriceOdd: proto.Float64(36.5),
					BetVec:   []float64{0.6, 1.2, 1.8, 2.4, 3, 6, 12, 18, 24, 30, 60, 120, 180, 300, 600},
				},
			},
		}),
	}

	ret = jiliut.ProtoEncode(&resdata)
	return
}
