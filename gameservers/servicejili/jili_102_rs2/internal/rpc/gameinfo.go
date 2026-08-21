package rpc

import (
	"fmt"
	"log/slog"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_102_rs2/internal"
	"serve/servicejili/jili_102_rs2/internal/message"
	"serve/servicejili/jili_102_rs2/internal/models"
	"serve/servicejili/jiliut"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{1, 1.2, 3, 6, 9, 15, 30, 45, 60, 90, 150, 210, 300, 600, 900, 1200}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Rs2_GameReqData
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
		var ack = message.Rs2_GameInfoAck{
			WalletInfo: []*message.Rs2_WalletInfo{
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
			Mul:           []float64{},
			Prefer:        []*message.Rs2_PreferRoundShow{},
		}

		encode, _ := proto.Marshal(&ack)
		var resdata = message.Rs2_ResData{
			Type: proto.Int32(AckType["info"]),
			Data: []*message.Rs2_InfoData{
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
