package rpc

import (
	"fmt"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_4_tks/internal"
	"serve/servicejili/jili_4_tks/internal/message"
	"serve/servicejili/jili_4_tks/internal/models"
	"serve/servicejili/jiliut"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 400, 500, 700, 1000}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Tks_GameReqData
	err = proto.Unmarshal(ps.Data, &gameReq)
	if err != nil {
		return
	}

	token := gameReq.GetToken()

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
		currencyName := ""
		currencySymbol := ""
		currencyMulti := 1.0
		operatorInfo := jiliut.GetOperatorInfo(plr.AppID)
		if operatorInfo.CurrencyManufactureVisibleOff != nil {
			plat := "jili"
			if strings.HasPrefix(lazy.ServiceName, "tada") {
				plat = "tada"
			}
			if _, ok := operatorInfo.CurrencyManufactureVisibleOff[plat]; ok && operatorInfo.CurrencyManufactureVisibleOff[plat] == 1 {
				item := lazy.GetCurrencyItem(operatorInfo.CurrencyKey)
				currencyName = item.Key
				currencySymbol = item.Symbol
				currencyMulti = item.Multi
			}
		}
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		var ack = message.Tks_GameInfoAck{
			WalletInfo: []*message.Tks_WalletInfo{
				{
					CurrencyName:   proto.String(currencyName),
					CurrencySymbol: proto.String(currencySymbol),
					Bet:            ut.FloatArrMul(info.Cs, currencyMulti),
					// BetType
					Coin:    proto.Float64(ut.Gold2Money(gold)),
					Decimal: proto.Int32(2),
					Ratio:   proto.Float64(1),
					Rate:    proto.Float64(1),
					Unit:    proto.Float64(1),
				},
			},
			ArcadeNo:      proto.Int32(1278994),
			CanIntoArcade: proto.Bool(true),
			Test:          proto.Bool(true),
			Prefer:        []*message.Tks_PreferRoundShow{},
			HighRateBet:   proto.Float64(3),
			Plate:         []int32{},
		}

		encode, _ := proto.Marshal(&ack)
		plr.SpinCountOfThisEnter = 0
		var resdata = message.Tks_ResData{
			Type: proto.Int32(AckType["info"]),
			Data: []*message.Tks_InfoData{
				{
					Encode: encode,
				},
			},
		}
		// resdata.

		ret, _ = proto.Marshal(&resdata)

		return nil
	})
	return
}
