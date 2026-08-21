package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_17_sweetheart/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []int64{1000, 2000, 3000, 5000, 8000, 10000, 20000, 50000, 100000, 200000, 300000, 400000, 500000, 700000, 1000000}

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
		var ack = &serverOfficial.GameInfoAck{
			MaxOdd: 1000,
			GameMall: &serverOfficial.GameMallInfo{
				MaxBet:   100,
				PriceOdd: []float64{29.5},
				AlterID:  []int32{50},
			},
			Mall: &serverOfficial.MallInfo{
				AlterID:  50,
				MaxBet:   1000 * currencyMulti,
				PriceOdd: 29.5,
				Show:     2,
			},
			ExtraInfo: make([]byte, 0),
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					CurrencyName:   currencyName,
					CurrencySymbol: currencySymbol,
					Bet:            ut.FloatArrMul(info.Cs, currencyMulti),
					Coin:           ut.Gold2Money(gold),
					Decimal:        2,
					Rate:           5.75,
					Ratio:          1,
					Unit:           1,
				},
			},
			FreeSpin:     &serverOfficial.FreeSpinList{},
			FreeSpinType: -1,
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
