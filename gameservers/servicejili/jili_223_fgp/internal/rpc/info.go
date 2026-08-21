package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_223_fgp/internal/message"
	"serve/servicejili/jili_223_fgp/internal/models"
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
			MaxOdd: proto.Float64(15025),
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
			ExtraInfo: jiliut.ProtoEncode(&message.Fgp_GameInfoData{
				Mul:       []float64{1, 1.5},
				WheelOdds: []float64{1000, 10, 1, 30, 5, 50, 200, 8, 100, 3, 20, 15},
			}),
			FreeSpin: []*jiliUtMessage.ServerFreeSpinList{
				{
					Free: make([]*jiliUtMessage.ServerFreeData, 0),
				},
			},
			FreeSpinType: proto.Int32(-1)}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
