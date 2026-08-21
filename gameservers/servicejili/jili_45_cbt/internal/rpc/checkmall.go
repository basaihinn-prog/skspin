package rpc

import (
	"fmt"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/servicejili/jili_45_cbt/internal/models"

	"serve/servicejili/jili_45_cbt/internal"
	"serve/servicejili/jili_45_cbt/internal/message"
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

	token := req.GetToken()
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
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
						MaxBet:   proto.Float64(1000 * curItem.Multi),
						PriceOdd: proto.Float64(36.5),
						BetVec:   []float64{0.5, 1, 2, 3, 5, 10, 20, 30, 40, 50, 80, 100, 200, 500, 1000},
					},
				},
			}),
		}

		ret = jiliut.ProtoEncode(&resdata)
		return nil
	})
	return
}
