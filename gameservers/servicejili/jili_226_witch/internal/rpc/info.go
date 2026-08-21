package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_226_witch/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []float64{0.5, 1, 2, 3, 5, 10, 20, 30, 40, 50, 80, 100, 200, 500, 1000}

func reqInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {

	gold, err := slotsmongo.GetBalance(pid)
	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		info, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		//currencyName := ""
		//currencySymbol := ""
		currencyMulti := 1.0
		operatorInfo := jiliut.GetOperatorInfo(plr.AppID)
		if operatorInfo.CurrencyManufactureVisibleOff != nil {
			plat := "jili"
			if strings.HasPrefix(lazy.ServiceName, "tada") {
				plat = "tada"
			}
			if _, ok := operatorInfo.CurrencyManufactureVisibleOff[plat]; ok && operatorInfo.CurrencyManufactureVisibleOff[plat] == 1 {
				item := lazy.GetCurrencyItem(operatorInfo.CurrencyKey)
				//currencyName = item.Key
				//currencySymbol = item.Symbol
				currencyMulti = item.Multi
			}
		}
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		var ack = &serverOfficial.GameInfoAck{
			MaxOdd: 8000,
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					//CurrencyName:   currencyName,
					//CurrencySymbol: currencySymbol,
					Coin:    ut.Gold2Money(gold),
					Bet:     ut.FloatArrMul(info.Cs, currencyMulti),
					Unit:    1,
					Ratio:   1,
					Rate:    5.75,
					Decimal: 4,
				},
			},
			Mall: &serverOfficial.MallInfo{
				PriceOdd: 50,
				MaxBet:   100,
				AlterID:  50,
				Show:     2,
			},
			GameMall: &serverOfficial.GameMallInfo{
				MaxBet:   100,
				PriceOdd: []float64{50},
				AlterID:  []int32{50},
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
