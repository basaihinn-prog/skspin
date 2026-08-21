package rpc

import (
	"fmt"
	"log/slog"
	"serve/comm/redisx"
	"serve/servicejili/jili_45_cbt/internal"
	"strings"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_45_cbt/internal/message"
	"serve/servicejili/jili_45_cbt/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{0.5, 1, 2, 3, 5, 10, 20, 30, 40, 50, 80, 100, 200, 500, 1000}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Cbt_GameReqData
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
		var ack = message.Cbt_GameInfoAck{

			ArcadeNo:      proto.Int32(1320487),
			CanIntoArcade: proto.Bool(true),
			OddList:       []float64{1, 3, 5, 8, 10, 40, 135, 625},
			// Plate:         []int32{5, 7, 5, 0, 2, 5, 0, 7, 3, 4, 4, 0, 3, 5, 8, 8, 3, 3, 1, 7, 9, 0, 5, 6},
			Prefer: []*message.Cbt_PreferRoundShow{},
			Test:   proto.Bool(true),
			WalletInfo: []*message.Cbt_WalletInfo{
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
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Cbt_ResData{
			Type: proto.Int32(AckType["info"]),
			Data: []*message.Cbt_InfoData{
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
