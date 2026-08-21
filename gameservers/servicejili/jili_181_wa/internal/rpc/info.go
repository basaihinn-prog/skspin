package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_181_wa/internal/gamedata"
	"serve/servicejili/jili_181_wa/internal/models"
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

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 400, 500, 700, 1000}

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
			MaxOdd: proto.Float64(10000),
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
			// ExtraInfo: jiliut.ProtoEncode(&message.Wa_GameInfoData{
			// 	Mul:       []float64{1, 1.5},
			// 	ShowExtra: proto.Bool(true),
			// 	ComboSetting: []*message.Wa_Column{
			// 		{
			// 			Col: []int32{2, 5, 8, 11},
			// 		},
			// 	},
			// }),
			ExtraInfo: []byte{8, 1, 18, 16, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 248, 63, 26, 6, 10, 4, 2, 5, 8, 11, 26, 6, 10, 4, 2, 5, 8, 11, 26, 6, 10, 4, 2, 3, 4, 5, 26, 6, 10, 4, 1, 1, 1, 1, 26, 6, 10, 4, 3, 5, 7, 10, 26, 10, 10, 8, 196, 19, 196, 19, 196, 19, 196, 19},
			FreeSpin: []*jiliUtMessage.ServerFreeSpinList{
				{
					Free: make([]*jiliUtMessage.ServerFreeData, 0),
				},
			},
			JpUnlockBet:  []float64{},
			FreeSpinType: proto.Int32(-1),
			Mall: []*jiliUtMessage.Server_MallInfo{
				{
					AlterID:  proto.Int32(50),
					MaxBet:   proto.Float64(1000 * currencyMulti),
					PriceOdd: proto.Float64(float64(gamedata.GetSettings().BuyMulti)),
					Show:     proto.Int32(2),
				},
			}}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
