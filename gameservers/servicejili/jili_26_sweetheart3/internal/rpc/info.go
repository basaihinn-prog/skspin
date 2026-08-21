package rpc

import (
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_26_sweetheart3/internal/message"
	"serve/servicejili/jili_26_sweetheart3/internal/models"
	"serve/servicejili/jiliut/AckType"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 400, 500, 700, 1000}

func reqInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	var inforeq message.Server_InfoReq
	err = proto.Unmarshal(data, &inforeq)
	if err != nil {
		return
	}
	gold, err := slotsmongo.GetBalance(pid)
	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		_, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		var ack = &serverOfficial.GameInfoAck{
			WalletInfo: []*serverOfficial.WalletInfo{
				{
					Bet:     BetGrade,
					Coin:    ut.Gold2Money(gold * 1000),
					Decimal: 4,
					Rate:    0.43,
					Ratio:   1,
					Unit:    1,
				},
			},
			FreeSpin:     &serverOfficial.FreeSpinList{Free: []*serverOfficial.FreeData{}},
			FreeSpinType: -1,
			JpUnlockBet:  []float64{},
			//Free: []*message.Server_FreeSpinData{
			//	{},
			//},
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
