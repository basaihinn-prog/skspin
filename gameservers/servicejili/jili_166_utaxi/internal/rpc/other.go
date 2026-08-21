package rpc

import (
	"fmt"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_166_utaxi/internal"
	"serve/servicejili/jili_166_utaxi/internal/message"
	"serve/servicejili/jili_166_utaxi/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/mall/checkmallproto", internal.GameShortName), checkmallproto)

	jiliut.RegRpc("/fulljp/JPInfoProto", unionjp_JPBlockProto)

	jiliut.RegRpc(fmt.Sprintf("/%s/account/heart", internal.GameShortName), accountHeart)

}

// /csh/account/heart
func accountHeart(ps *nats.Msg) (ret []byte, err error) {
	var ack = message.Utaxi_HeartAck{
		// Aid: &aid,
		Message: proto.String("testtest"),
	}

	encode, _ := proto.Marshal(&ack)

	var resdata = message.Utaxi_ResData{
		Type: proto.Int32(AckType["heartBeat"]),
		Data: []*message.Utaxi_InfoData{
			{
				Encode: encode,
			},
		},
	}
	// resdata.

	ret, _ = proto.Marshal(&resdata)
	return
}

// /csh/mall/checkmallproto
func checkmallproto(ps *nats.Msg) (ret []byte, err error) {
	//jili_166没有购买小游戏

	// var req message.Gaia_CheckMallReq
	// err = proto.Unmarshal(ps.Data, &req)
	// if err != nil {
	// 	return
	// }

	// token := req.GetToken()
	// pid, err := jwtutil.ParseToken(token)
	// if err != nil {
	// 	// err = define.NewErrCode("Invalid player session", 1302)
	// 	return
	// }
	// err = db.CallWithPlayer(pid, func(plr *models.Player) error {
	// 	curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
	// 	var resdata = message.Gaia_GaiaResponse{
	// 		Type:  proto.Int32(AckType["buyBonus"]),
	// 		Token: req.Token,
	// 		Data: jiliut.ProtoEncode(&message.Gaia_CheckMallAck{
	// 			Show: proto.Int32(2),
	// 			Settings: []*message.Gaia_GameMallSetting{
	// 				{
	// 					APIID:    req.Apiid,
	// 					AlterID:  proto.Int32(50),
	// 					Event:    proto.Bool(true),
	// 					ForSale:  proto.Bool(true),
	// 					GameID:   proto.Int32(internal.GameNo),
	// 					MaxBet:   proto.Float64(1000 * curItem.Multi),
	// 					PriceOdd: proto.Float64(36.5),
	// 					BetVec:   []float64{0.5, 1, 2, 3, 5, 10, 20, 30, 40, 50, 80, 100, 200, 500, 1000},
	// 				},
	// 			},
	// 		}),
	// 	}

	// 	ret = jiliut.ProtoEncode(&resdata)
	// 	return nil
	// })
	return
}

func unionjp_JPBlockProto(ps *nats.Msg) (ret []byte, err error) {
	var preq message.Gaia_FullJPInfoReq

	err = proto.Unmarshal(ps.Data, &preq)
	if err != nil {
		return
	}

	token := preq.GetToken()
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		aid := int32(pid)
		var ack = message.Utaxi_LoginDataAck{
			Aid: &aid,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Utaxi_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Utaxi_InfoData{
				{
					Encode: encode,
				},
			},
		}
		// resdata.

		ret, _ = proto.Marshal(&resdata)

		return err
	})

	return
}
