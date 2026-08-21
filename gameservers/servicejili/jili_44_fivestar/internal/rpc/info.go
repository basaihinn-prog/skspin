package rpc

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_44_fivestar/internal/message"
	"serve/servicejili/jili_44_fivestar/internal/models"
	"serve/servicejili/jiliut"
	"serve/servicejili/jiliut/AckType"
	"strings"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.Info] = reqInfo
}

var BetGrade = []int64{1, 1000, 2000, 3000, 5000, 8000, 10000, 20000, 50000, 100000, 200000, 300000, 400000, 500000, 700000, 1000000}

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
		var ack = &message.Server_InfoResponse{
			//	MaxOdd: proto.Float64(1000),
			Wallet: []*message.Server_Wallet{
				{
					CurrencyName:   proto.String(currencyName),
					CurrencySymbol: proto.String(currencySymbol),
					Bet:            ut.Float64ToInt64TwoDecimal(ut.FloatArrMul(info.Cs, currencyMulti)),
					Coin:           proto.Int64(int64(ut.Gold2Money(gold * 1000))),
					Decimal:        proto.Int32(2),
					Rate:           proto.Float64(1),
					Ratio:          proto.Float64(1),
					Unit:           proto.Float64(1000),
				},
			},
			SpinResponse: []*message.Server_SpinResponse{
				{
					AllPlate:  []byte{10, 35, 10, 5, 10, 3, 7, 7, 4, 10, 5, 10, 3, 6, 4, 1, 10, 5, 10, 3, 5, 5, 7, 10, 5, 10, 3, 1, 1, 4, 10, 5, 10, 3, 7, 1, 5},
					PostMoney: proto.Int64(int64(ut.Gold2Money(gold * 1000))),
					PreMoney:  proto.Int64(int64(ut.Gold2Money(gold * 1000))),
				},
			},
			JpUnlockBet: []float64{0, 0, 0, 0},
			Free: []*message.Server_FreeSpinData{
				{
					//	Remain: make([]*message.ServerFreeData, 0),
				},
			},
			//	FreeSpinType: proto.Int32(-1),
		}
		ret = ack
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
