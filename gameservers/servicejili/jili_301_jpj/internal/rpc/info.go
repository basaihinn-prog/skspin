package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_301_jpj/internal/gamedata"
	"serve/servicejili/jili_301_jpj/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"strings"

	"github.com/nats-io/nats.go"
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
		var ack = &serverOfficial.GameInfoAck{
			MaxOdd: 15000,
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					CurrencyName:   currencyName,
					CurrencySymbol: currencySymbol,
					Bet:            ut.FloatArrMul(info.Cs, currencyMulti),
					Coin:           ut.Gold2Money(gold),
					Decimal:        4,
					Rate:           1,
					Ratio:          1,
					Unit:           1,
				},
			},
			ExtraInfo: []byte{
				10, 16, 0, 0, 0, 0, 0, 0, 36, 64, 0, 0, 0, 0, 0, 0, 36, 64, 18, 18, 9, 0, 0, 0, 0, 0, 0, 20, 64, 17, 0, 0, 0, 0, 0, 0, 36, 64, 18, 18, 9, 0, 0, 0, 0, 0, 0, 46, 64, 17, 0, 0, 0, 0, 0, 0, 62, 64, 18, 18, 9, 0, 0, 0, 0, 0, 0, 73, 64, 17, 0, 0, 0, 0, 0, 0, 105, 64, 18, 18, 9, 0, 0, 0, 0, 0, 64, 143, 64, 17, 0, 0, 0, 0, 0, 64, 143, 64, 26, 16, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 248, 63,
			},
			FreeSpin: &serverOfficial.FreeSpinList{},
			Mall: &serverOfficial.MallInfo{
				AlterID:  50,
				MaxBet:   1000 * currencyMulti,
				PriceOdd: float64(gamedata.GetSettings().BuyMulti),
				//PriceOdd: proto.Float64(33.5),
				Show: 2,
			},
			FreeSpinType: -1,
			ApiType:      8,
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
