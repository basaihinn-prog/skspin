package rpc

import (
	"fmt"
	"log/slog"
	"serve/comm/redisx"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_146_fb/internal"
	"serve/servicejili/jili_146_fb/internal/message"
	"serve/servicejili/jili_146_fb/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 400, 500, 700, 1000}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Fb_GameReqData
	err = proto.Unmarshal(ps.Data, &gameReq)
	if err != nil {
		return
	}

	// fmt.Println(login.String())

	token := gameReq.GetToken()
	slog.Info("login", "token", token)

	//TODO
	// pid := int64(123456)

	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	gold, err := slotsmongo.GetBalance(pid)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		info, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		var ack = message.Fb_GameInfoAck{
			WalletInfo: []*message.Fb_WalletInfo{
				{
					Bet: ut.FloatArrMul(info.Cs, curItem.Multi),
					// BetType
					Coin:    proto.Float64(ut.Gold2Money(gold)),
					Decimal: proto.Int32(2),
					Ratio:   proto.Float64(1),
					Rate:    proto.Float64(1),
					Unit:    proto.Float64(1),
				},
			},
			Football:      []float64{1, 2, 5, 10, 15, 25, 50, 100, 250},
			ArcadeNo:      proto.Int32(1234843),
			CanIntoArcade: proto.Bool(true),
			Test:          proto.Bool(true),
			Mul:           []float64{1, 1.5},
			Prefer:        []*message.Fb_PreferRoundShow{},
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Fb_ResData{
			Type: proto.Int32(AckType["info"]),
			Data: []*message.Fb_InfoData{
				{
					Encode: encode,
				},
			},
		}
		// resdata.

		ret, _ = proto.Marshal(&resdata)
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
