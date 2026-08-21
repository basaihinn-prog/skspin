package rpc

import (
	"encoding/base64"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_51_mc/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 500, 800, 1000, 2000}

func reqInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	gold, err := slotsmongo.GetBalance(pid)
	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		info, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		currencyName := ""
		currencySymbol := ""
		//currencyMulti := 1.0
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
				//currencyMulti = item.Multi
			}
		}
		//{
		//	"WalletInfo":[{"coin":99992.69,"bet":[1,5,10,50,100],"unit":1,"ratio":1,"rate":0.43,"decimal":4}],
		//	"extraInfo":"",
		//	"MaxOdd":10000,
		//	"freeSpin":{},
		//	"freeSpinType":-1
		//}
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		info.Cs = []float64{1, 5, 10, 50, 100}
		extraInfoStr := "CkIKQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKQgpAAAAAAAAAOUAAAAAAAABJQAAAAAAAwFJAAAAAAAAAWUAAAAAAAMBiQAAAAAAAQG9AAAAAAABAf0AAAAAAAECPQApCCkAAAAAAAABJQAAAAAAAAFlAAAAAAADAYkAAAAAAAABpQAAAAAAAwHJAAAAAAABAf0AAAAAAAECPQAAAAAAAQJ9ACkIKQAAAAAAAQH9AAAAAAABAj0AAAAAAAHCXQAAAAAAAiKNAAAAAAACIs0AAAAAAAIjDQAAAAAAAathAAAAAAABq6EAKQgpAAAAAAABAj0AAAAAAAECfQAAAAAAAcKdAAAAAAACIs0AAAAAAAIjDQAAAAAAAiNNAAAAAAABq6EAAAAAAAGr4QA=="
		// 解码 Base64 字符串
		extraInfo, _ := base64.StdEncoding.DecodeString(extraInfoStr)
		var ack = &serverOfficial.GameInfoAck{
			MaxOdd:    100000,
			ExtraInfo: extraInfo,
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					CurrencyName:   currencyName,
					CurrencySymbol: currencySymbol,
					Bet:            info.Cs,
					Coin:           ut.Gold2Money(gold),
					Decimal:        2,
					Rate:           1,
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
