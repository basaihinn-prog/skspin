package rpc

import (
	"encoding/base64"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_302_mcp/internal/models"
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
		info.Cs = []float64{1, 5, 10, 50, 100, 500, 1000}
		// 你的 Base64 编码字符串
		extraInfoStr := "OgA6ADoAOkIKQAAAAAAAABRAAAAAAAAAJEAAAAAAAAAuQAAAAAAAADRAAAAAAAAAPkAAAAAAAABJQAAAAAAAAFlAAAAAAAAAaUA6QgpAAAAAAAAAFEAAAAAAAAAkQAAAAAAAAC5AAAAAAAAANEAAAAAAAAA+QAAAAAAAAElAAAAAAAAAWUAAAAAAAABpQDpCCkAAAAAAAAAkQAAAAAAAADRAAAAAAAAAPkAAAAAAAABJQAAAAAAAAFlAAAAAAAAAaUAAAAAAAEB/QAAAAAAAQI9AOkIKQAAAAAAAACRAAAAAAAAANEAAAAAAAAA+QAAAAAAAAElAAAAAAAAAWUAAAAAAAABpQAAAAAAAQH9AAAAAAABAj0A6QgpAAAAAAAAAJEAAAAAAAAA0QAAAAAAAAD5AAAAAAAAASUAAAAAAAABZQAAAAAAAAGlAAAAAAABAf0AAAAAAAECPQDpCCkAAAAAAAAAkQAAAAAAAADRAAAAAAAAAPkAAAAAAAABJQAAAAAAAAFlAAAAAAAAAaUAAAAAAAEB/QAAAAAAAQI9AOkIKQAAAAAAAACRAAAAAAAAANEAAAAAAAAA+QAAAAAAAAElAAAAAAAAAWUAAAAAAAABpQAAAAAAAQH9AAAAAAABAn0A6QgpAAAAAAAAAJEAAAAAAAAA0QAAAAAAAAD5AAAAAAAAASUAAAAAAAABZQAAAAAAAAGlAAAAAAABAf0AAAAAAAECfQDpCCkAAAAAAAAAkQAAAAAAAADRAAAAAAAAAPkAAAAAAAABJQAAAAAAAAFlAAAAAAAAAaUAAAAAAAEB/QAAAAAAAQJ9AOkIKQAAAAAAAACRAAAAAAAAANEAAAAAAAAA+QAAAAAAAAElAAAAAAAAAWUAAAAAAAABpQAAAAAAAQH9AAAAAAABAn0A6QgpAAAAAAAAAJEAAAAAAAAA0QAAAAAAAAD5AAAAAAAAASUAAAAAAAABZQAAAAAAAAGlAAAAAAABAf0AAAAAAAIizQDpCCkAAAAAAAAAkQAAAAAAAADRAAAAAAAAAPkAAAAAAAABJQAAAAAAAAFlAAAAAAAAAaUAAAAAAAEB/QAAAAAAAiMNAQniamZmZmZm5PwAAAAAAAOA/AAAAAAAA8D8AAAAAAAAUQAAAAAAAACRAAAAAAAAASUAAAAAAAABZQAAAAAAAQH9AAAAAAABAj0AAAAAAAIizQAAAAAAAiMNAAAAAAABq6EAAAAAAAGr4QAAAAACAhB5BAAAAAICELkE="
		// 解码 Base64 字符串
		extraInfo, _ := base64.StdEncoding.DecodeString(extraInfoStr)
		var ack = &serverOfficial.GameInfoAck{
			ExtraInfo: extraInfo,
			MaxOdd:    10100,
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					CurrencyName:   currencyName,
					CurrencySymbol: currencySymbol,
					Bet:            info.Cs,
					Coin:           ut.Gold2Money(gold),
					Decimal:        4,
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
