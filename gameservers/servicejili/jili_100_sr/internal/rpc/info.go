package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_100_sr/internal/models"
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

var BetGrade = []float64{1, 2, 3, 5, 10, 15, 20, 30, 40, 60, 80, 120, 160, 240, 320, 480, 640, 960, 1280, 1920}

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
			MaxOdd: proto.Float64(100000),
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
			JpUnlockBet:  []float64{},
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
