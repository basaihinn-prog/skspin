package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_76_nightclub/internal/gamedata"
	"serve/servicejili/jili_76_nightclub/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"serve/servicejili/jiliut/jiliUtMessage"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []float64{0.3, 0.5, 1, 1.5, 2, 3, 5, 10, 20, 30, 40, 50, 80, 100, 200}

func reqInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {

	gold, err := slotsmongo.GetBalance(pid)
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
		var ack = &jiliUtMessage.Server_GameInfoAck{
			MaxOdd: proto.Float64(1000),
			WalletInfo: []*jiliUtMessage.Server_WalletInfo{
				{
					CurrencyName:   proto.String(currencyName),
					CurrencySymbol: proto.String(currencySymbol),
					Bet:            ut.FloatArrMul(info.Cs, currencyMulti),
					Coin:           proto.Float64(ut.Gold2Money(gold)),
					Decimal:        proto.Int32(2),
					Rate:           proto.Float64(1),
					Ratio:          proto.Float64(1),
					Unit:           proto.Float64(1),
				},
			},
			FreeSpin: []*jiliUtMessage.ServerFreeSpinList{
				{
					Free: make([]*jiliUtMessage.ServerFreeData, 0),
				},
			},
			FreeSpinType: proto.Int32(-1),
			Mall: []*jiliUtMessage.Server_MallInfo{
				{
					AlterID:  proto.Int32(50),
					MaxBet:   proto.Float64(200 * currencyMulti),
					PriceOdd: proto.Float64(float64(gamedata.Settings.BuyMulti)),
					Show:     proto.Int32(2),
				},
			},
			GameMall: []*jiliUtMessage.Server_GameMallInfo{
				{
					MaxBet: proto.Float64(200 * currencyMulti),
					PriceOdd: []float64{
						float64(gamedata.Settings.BuyMulti),
					},
					AlterID: []int32{
						50,
					},
				},
			},
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
