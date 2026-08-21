package rpc

import (
	"fmt"
	"log/slog"
	"serve/comm/redisx"
	"strings"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_5_chilli/internal"
	"serve/servicejili/jili_5_chilli/internal/message"
	"serve/servicejili/jili_5_chilli/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{1,
	2,
	3,
	5,
	8,
	10,
	20,
	50,
	100,
	200,
	300,
	400,
	500,
	700,
	1000}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	var gameReq message.Chilli_GameReqData
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
		// aid := int32(pid)
		// var ack = message.Csh_LoginDataAck{
		// 	Aid: &aid,
		// }

		// encode, _ := proto.Marshal(&ack)

		// encode := []byte{10, 160, 1, 33, 174, 71, 225, 122, 148, 110, 168, 64, 42, 120, 0, 0, 0, 0, 0, 0, 224, 63, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 36, 64, 0, 0, 0, 0, 0, 0, 52, 64, 0, 0, 0, 0, 0, 0, 62, 64, 0, 0, 0, 0, 0, 0, 68, 64, 0, 0, 0, 0, 0, 0, 73, 64, 0, 0, 0, 0, 0, 0, 84, 64, 0, 0, 0, 0, 0, 0, 89, 64, 0, 0, 0, 0, 0, 0, 105, 64, 0, 0, 0, 0, 0, 64, 127, 64, 0, 0, 0, 0, 0, 64, 143, 64, 49, 0, 0, 0, 0, 0, 0, 240, 63, 57, 0, 0, 0, 0, 0, 0, 240, 63, 65, 0, 0, 0, 0, 0, 0, 240, 63, 96, 4, 16, 133, 209, 57, 24, 1, 32, 1, 42, 0, 58, 24, 5, 7, 5, 0, 2, 5, 0, 7, 3, 4, 4, 0, 3, 5, 8, 8, 3, 3, 1, 7, 9, 0, 5, 6}

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

		var ack = message.Chilli_GameInfoAck{
			ChilliSettingVec: []*message.Chilli_ChilliSetting{
				{
					PlateCount:    proto.Int32(4),
					ChilliCount:   proto.Int32(30),
					WildLockCount: proto.Int32(3),
					GameType:      proto.Int32(4),
				},
				{
					PlateCount:    proto.Int32(4),
					ChilliCount:   proto.Int32(14),
					WildLockCount: proto.Int32(2),
					GameType:      proto.Int32(3),
				},
				{
					PlateCount:    proto.Int32(3),
					ChilliCount:   proto.Int32(9),
					WildLockCount: proto.Int32(1),
					GameType:      proto.Int32(2),
				},
				{
					PlateCount: proto.Int32(2),
					GameType:   proto.Int32(1),
				},
			},
			ArcadeNo:      proto.Int32(1234843),
			CanIntoArcade: proto.Bool(true),
			//Plate:         []int32{5, 7, 5, 0, 2, 5, 0, 7, 3, 4, 4, 0, 3, 5, 8, 8, 3, 3, 1, 7, 9, 0, 5, 6},
			Prefer: []*message.Chilli_PreferRoundShow{},
			Test:   proto.Bool(true),
			WalletInfo: []*message.Chilli_WalletInfo{
				{
					CurrencyName:   proto.String(currencyName),
					CurrencySymbol: proto.String(currencySymbol),
					Bet:            ut.FloatArrMul(info.Cs, currencyMulti),
					Coin:           proto.Float64(ut.Gold2Money(gold)),
					Unit:           proto.Float64(1),
					Ratio:          proto.Float64(1),
					Rate:           proto.Float64(1),
					Decimal:        proto.Int32(2),
					// BetType
					//todo betTypeS 未找到定义
					//BetTypeS:
				},
			},
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Chilli_ResData{
			Type: proto.Int32(AckType["info"]),
			Data: []*message.Chilli_InfoData{
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
