package jdbcomm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
)

func hitPermillageArr(arr []int, i int) bool {
	if i < len(arr) {
		w := arr[i]
		if rand.IntN(1000) < w {
			return true
		}
	}
	return false
}

type NextPlayRespParam struct {
	Player             *Player
	Bet                int64
	SelfPoolGold       int64
	IsBuy              bool
	BuyFeatureType     string
	BuyMinAwardPercent int
	Combine            Combine
	NoAwardPercent     int
	HitBigAwardPercent []int
	AppStore           *redisx.AppStore
	OptionsSlots       int //可选插槽

	NextFunc               func(str string, witchBuckets int) (playResp *SimulateData, err error)
	SimulateByBucketIdFunc func(bucketId, witchBuckets int) (playResp *SimulateData, err error)
	BigRewardFunc          func(witchBuckets int) (playResp *SimulateData, err error)
	BigRewardFunc2_5       func(witchBuckets int) (playResp *SimulateData, err error)
	MinBuyFunc             func() (playResp *SimulateData, err error)
	NextBuyFunc            func(str string) (playResp *SimulateData, err error)
	NextSuperBuyFunc       func(str string) (playResp *SimulateData, err error)
	ControlNextData        func(nextData *redisx.NextMulti, witchBucket int) (playResp *SimulateData, err error)
}

type NextPlayRespParam2 struct {
	Player             *Player
	Bet                int64
	SelfPoolGold       int64
	IsBuy              bool
	BuyMinAwardPercent int
	Combine            Combine
	NoAwardPercent     int
	HitBigAwardPercent []int
	AppStore           *redisx.AppStore
	OptionsSlots       int //可选插槽
	OptionsSlots2      int //可选插槽2

	NextFunc               func(str string, witchBuckets, god int) (playResp *SimulateData, err error)
	SimulateByBucketIdFunc func(bucketId, witchBuckets, god int) (playResp *SimulateData, err error)
	BigRewardFunc          func(witchBuckets, god int) (playResp *SimulateData, err error)
	MinBuyFunc             func(god, ty int) (playResp *SimulateData, err error)
	NextBuyFunc            func(str string, god, ty int) (playResp *SimulateData, err error)
	ControlNextData        func(nextData *redisx.NextMulti, witchBucket, god int) (playResp *SimulateData, err error)
}

type NextPlayRespParamHasReSpin struct {
	Player             *Player
	Bet                int64
	SelfPoolGold       int64
	IsBuy              bool
	BuyMinAwardPercent int
	Combine            Combine
	NoAwardPercent     int
	HitBigAwardPercent []int
	AppStore           *redisx.AppStore
	OptionsSlots       int //可选插槽
	IsReSpin           bool

	NextFunc               func(str string, witchBuckets int) (playResp *SimulateData, err error)
	SimulateByBucketIdFunc func(bucketId, witchBuckets int) (playResp *SimulateData, err error)
	BigRewardFunc          func(witchBuckets int) (playResp *SimulateData, err error)
	MinBuyFunc             func() (playResp *SimulateData, err error)
	NextBuyFunc            func(str string) (playResp *SimulateData, err error)
}

func NextFunc(PPPlayers *Player, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	PLRRTP := (Bet*5 + PPPlayers.WinAmount) / (PPPlayers.BetAmount)

	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= PPPlayers.MoniterRTPNode && v.RangeMinValue <= PPPlayers.MoniterRTPNode {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(PPPlayers.SpinCount) {
				//新手概率
				if v.NewbieValue >= NewbieValue {
					if PLRRTP < RTP {
						return true
					}
				}
			} else {
				//老手概率
				if v.NotNewbieValue >= NewbieValue {
					if PLRRTP < RTP {
						return true
					}
				}
			}
		}
	}

	return false
}

func NextFunc1(PPPlayers *Player, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	MoniterRTPNode := math.Abs(float64(PPPlayers.MoniterRTPNode))
	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= int64(MoniterRTPNode) && v.RangeMinValue <= int64(MoniterRTPNode) {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(PPPlayers.SpinCount) {
				//新手概率
				if v.NewbieValue > NewbieValue {
					return true
				}
			} else {
				//老手概率
				if v.NotNewbieValue > NewbieValue {
					return true
				}
			}
		}
	}

	return false
}

func NextPlayResp(param *NextPlayRespParam) (playResp *SimulateData, hitBigAward, forcedKill, buyKill bool, err error) {
	defer func(p *NextPlayRespParam) {
		if playResp.Times != 0 {
			slotsmongo.UpdateGamePlayerMulti(p.Player.PID, playResp.Times)
		}
	}(param)

	if param.IsBuy {

		if param.BuyFeatureType == "BuyFeature_02" {
			if param.NextSuperBuyFunc == nil {
				err = errors.New("not found NextBuyFunc")
				return
			}
			playResp, err = param.NextSuperBuyFunc(fmt.Sprintf("%d", param.Bet))
			return
		}
		minAwardPercent := param.Player.BuyMinAwardPercent
		fmt.Println("player:BuyMinAwardPercent", minAwardPercent)
		if minAwardPercent == 0 {
			minAwardPercent = param.BuyMinAwardPercent
			fmt.Println("App:BuyMinAwardPercent", minAwardPercent)
		}
		if rand.IntN(1000) < minAwardPercent {
			if param.MinBuyFunc == nil {
				err = errors.New("not found minBuyFunc")
				return
			}
			playResp, err = param.MinBuyFunc()
		} else {
			if param.NextBuyFunc == nil {
				err = errors.New("not found NextBuyFunc")
				return
			}
			playResp, err = param.NextBuyFunc(fmt.Sprintf("%d", param.Bet))
		}
		return
	}

	//限制
	if param.Player.RestrictionsStatus == 1 {
		if param.Player.Win >= param.Player.RestrictionsMaxWin*10000 || param.Player.Multi >= param.Player.RestrictionsMaxMulti {
			playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
			forcedKill = true
			return
		}
	}

	if nextData := redisx.GetPlayerNextData(param.Player.PID); nextData != nil {
		if param.ControlNextData == nil {
			err = errors.New("not found ControlNextData")
			return
		}
		playResp, err = param.ControlNextData(nextData, param.AppStore.GamePatten)
		return
	}

	if param.AppStore.PersonalMoniterConfig.IsMoniter == 1 || param.AppStore.MoniterConfig.IsMoniter == 1 {
		RTP := param.Player.TargetRTP
		if RTP == 0 {
			RTP = param.AppStore.RTP
		}

		param.Player.MoniterCountNode++
		if param.AppStore.PersonalMoniterConfig.IsMoniter == 1 { //个人监控
			//监控周期
			if int64(param.Player.MoniterCountNode) >= param.AppStore.PersonalMoniterConfig.MoniterNumCycle {
				param.Player.MoniterCountNode = 0
				//监控RTP节点
				param.Player.MoniterRTPNode = 0
				if param.Player.BetAmount != 0 {
					Tmp := float64(param.Player.WinAmount) / float64(param.Player.BetAmount)
					param.Player.MoniterRTPNode = int64(RTP) - int64(Tmp*100)
				}
			}

			//增加
			if param.Player.MoniterRTPNode > 0 {
				ok := NextFunc(param.Player, param.Bet, int64(RTP), param.AppStore.PersonalMoniterConfig.MoniterNewbieNum, param.AppStore.PersonalMoniterConfig.MoniterAddRTPRangeValue)
				fmt.Println("增", ok)
				if ok {
					playResp, err = param.BigRewardFunc2_5(param.AppStore.GamePatten)
					hitBigAward = true
					return
				}
			}

			//减少
			if param.Player.MoniterRTPNode < 0 {
				ok := NextFunc1(param.Player, param.Bet, int64(RTP), param.AppStore.PersonalMoniterConfig.MoniterNewbieNum, param.AppStore.PersonalMoniterConfig.MoniterReduceRTPRangeValue)
				fmt.Println("减", ok)
				if ok {
					playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
					forcedKill = true
					return
				}
			}
		} else if param.AppStore.MoniterConfig.IsMoniter == 1 { //游戏监控
			//监控周期
			if int64(param.Player.MoniterCountNode) == param.AppStore.MoniterConfig.MoniterNumCycle {
				param.Player.MoniterCountNode = 0
				//监控RTP节点
				param.Player.MoniterRTPNode = 0
				if param.Player.BetAmount != 0 {
					Tmp := float64(param.Player.WinAmount) / float64(param.Player.BetAmount)
					param.Player.MoniterRTPNode = int64(RTP) - int64(Tmp*100)
				}
			}

			//增加
			if param.Player.MoniterRTPNode > 0 {
				ok := NextFunc(param.Player, param.Bet, int64(RTP), param.AppStore.MoniterConfig.MoniterNewbieNum, param.AppStore.MoniterConfig.MoniterAddRTPRangeValue)
				if ok {
					playResp, err = param.BigRewardFunc2_5(param.AppStore.GamePatten)
					hitBigAward = true
					return
				}
			}

			//减少
			if param.Player.MoniterRTPNode < 0 {
				ok := NextFunc1(param.Player, param.Bet, int64(RTP), param.AppStore.MoniterConfig.MoniterNewbieNum, param.AppStore.MoniterConfig.MoniterReduceRTPRangeValue)
				if ok {
					playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
					forcedKill = true
					return
				}
			}
		}
	}

	noAwardPercent := param.Player.NoAwardPercent
	if noAwardPercent == 0 {
		noAwardPercent = redisx.LoadNoAwardPercent(param.Player.AppID)
	}
	if param.SelfPoolGold < 0 && rand.IntN(1000) < noAwardPercent {
		if param.SimulateByBucketIdFunc == nil {
			err = errors.New("not found SimulateByBucketIdFunc")
			return
		}
		playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
		forcedKill = true
		return
	}
	//百分之35机率触发
	if rand.IntN(1000) < HitBigAwardPercent {
		fmt.Println("enter")
		// 同时满足以下2个条件
		// 1.个人奖池金额>=最近200把投注均值不包含购买小游戏*10
		// 2.个人奖池金额>=当前投注额*40
		// 可以爆个人奖池为20-40倍的一个小游戏奖励，该奖励不消耗库存
		if param.Bet*10 < param.SelfPoolGold {
			avg := param.Player.CaclAvgBet()
			if avg*10 < param.SelfPoolGold {
				if param.BigRewardFunc == nil {
					err = errors.New("not found BigRewardFunc")
					return
				}
				playResp, err = param.BigRewardFunc(param.AppStore.GamePatten)
				hitBigAward = true
				return
			}
		}
	}
	if param.NextFunc == nil {
		err = errors.New("not found NextFunc")
		return
	}

	if param.Player.PersonWinMaxMult != 0 && param.Player.PersonWinMaxScore != 0 {
		for {
			playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten)
			if playResp.Times < float64(param.Player.PersonWinMaxMult) && playResp.Times*float64(param.Bet) <= float64(param.Player.PersonWinMaxScore)*10000 {
				return
			}
		}
	}

	if param.AppStore.MaxMultiple != 0 && param.AppStore.MaxWinPoints != 0 {
		for {
			playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten)
			if playResp.Times < param.AppStore.MaxMultiple && playResp.Times*float64(param.Bet) <= param.AppStore.MaxWinPoints*10000 {
				return
			}
		}
	}

	playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten)
	return
}

func NextPlayResp2(param *NextPlayRespParam2, isRand bool) (playResp *SimulateData, hitBigAward, forcedKill, buyKill bool, err error) {
	if isRand {
		coll := db.Collection("simulate")
		match := bson.M{}
		//id, _ := primitive.ObjectIDFromHex("66e3e561295a1ba68b5e8fc2")
		//match["_id"] = id
		if param.IsBuy {
			match["type"] = 1
		} else {
			match["type"] = 0
		}
		cursor, err := coll.Aggregate(context.TODO(), bson.A{
			bson.D{{"$match", match}},
			bson.D{
				{"$sample", bson.D{{"size", 1}}},
			},
		})
		if err != nil {
			return nil, false, false, false, err
		}

		// cursor.Next()
		var docs []*SimulateData
		cursor.All(context.TODO(), &docs)
		playResp = docs[0]
		return playResp, false, false, false, err
	}
	if param.IsBuy {
		minAwardPercent := param.Player.BuyMinAwardPercent
		if minAwardPercent == 0 {
			minAwardPercent = param.BuyMinAwardPercent
		}
		if rand.IntN(1000) < minAwardPercent {
			if param.MinBuyFunc == nil {
				err = errors.New("not found minBuyFunc")
				return
			}
			playResp, err = param.MinBuyFunc(param.OptionsSlots, param.OptionsSlots2)
		} else {
			if param.NextBuyFunc == nil {
				err = errors.New("not found NextBuyFunc")
				return
			}
			playResp, err = param.NextBuyFunc(fmt.Sprintf("%d", param.Bet), param.OptionsSlots, param.OptionsSlots2)
		}
		return
	}
	if nextData := redisx.GetPlayerNextData(param.Player.PID); nextData != nil {
		if param.ControlNextData == nil {
			err = errors.New("not found ControlNextData")
			return
		}
		playResp, err = param.ControlNextData(nextData, param.AppStore.GamePatten, param.OptionsSlots)
		return
	}
	noAwardPercent := param.Player.NoAwardPercent
	if noAwardPercent == 0 {
		noAwardPercent = redisx.LoadNoAwardPercent(param.Player.AppID)
	}
	if param.SelfPoolGold < 0 && rand.IntN(1000) < noAwardPercent {
		if param.SimulateByBucketIdFunc == nil {
			err = errors.New("not found SimulateByBucketIdFunc")
			return
		}
		playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten, param.OptionsSlots)
		forcedKill = true
		return
	}
	//if param.SelfPoolGold < 0 && rand.IntN(1000) < param.NoAwardPercent {
	//	if param.SimulateByBucketIdFunc == nil {
	//		err = errors.New("not found SimulateByBucketIdFunc")
	//		return
	//	}
	//	playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
	//	forcedKill = true
	//	return
	//}

	if hitPermillageArr(param.HitBigAwardPercent, param.Player.SpinCountOfThisEnter) {
		// 同时满足以下2个条件
		// 1.个人奖池金额>=最近200把投注均值不包含购买小游戏*10
		// 2.个人奖池金额>=当前投注额*40
		// 可以爆个人奖池为20-40倍的一个小游戏奖励，该奖励不消耗库存

		// bet := Grades[grade]

		if param.Bet*40 < param.SelfPoolGold {
			// var avg int64
			// sum := lo.Sum(player.BetHistory)
			// if len(player.BetHistory) > 0 {
			// 	avg = sum / int64(len(player.BetHistory))
			// }

			avg := param.Player.CaclAvgBet()
			if avg*40 < param.SelfPoolGold {
				if param.BigRewardFunc == nil {
					err = errors.New("not found BigRewardFunc")
					return
				}
				playResp, err = param.BigRewardFunc(param.AppStore.GamePatten, param.OptionsSlots)
				hitBigAward = true
				return
			}
		}
	}
	if param.NextFunc == nil {
		err = errors.New("not found NextFunc")
		return
	}
	for {
		playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten, param.OptionsSlots)
		if param.AppStore.MaxMultiple == 0 || param.AppStore.MaxWinPoints*10000 == 0 {
			return
		}
		if playResp.Times < param.AppStore.MaxMultiple && playResp.Times*float64(param.Bet) <= param.AppStore.MaxWinPoints*10000 {
			return
		}
	}

}

func NextPlayRespHasReSpin(param *NextPlayRespParamHasReSpin, isRand bool) (playResp *SimulateData, hitBigAward, forcedKill, buyKill bool, err error) {
	if isRand {
		coll := db.Collection("simulate")
		match := bson.M{}
		//id, _ := primitive.ObjectIDFromHex("66e3e561295a1ba68b5e8fc2")
		//match["_id"] = id
		if param.IsBuy {
			match["type"] = 1
		} else {
			match["type"] = 0
		}
		cursor, err := coll.Aggregate(context.TODO(), bson.A{
			bson.D{{"$match", match}},
			bson.D{
				{"$sample", bson.D{{"size", 1}}},
			},
		})
		if err != nil {
			return nil, false, false, false, err
		}

		// cursor.Next()
		var docs []*SimulateData
		cursor.All(context.TODO(), &docs)
		playResp = docs[0]
		return playResp, false, false, false, err
	}
	//if param.IsReSpin {
	//	if param.IsBuy {
	//		if param.NextBuyFunc == nil {
	//			err = errors.New("not found NextBuyFunc")
	//			return
	//		}
	//		playResp, err = param.NextBuyFunc(fmt.Sprintf("%d", param.Bet))
	//		return
	//	} else {
	//		match := bson.M{}
	//		match["hasgame"] = true
	//		cursor, err := db.Collection("simulate").Aggregate(context.TODO(), bson.A{
	//			bson.D{{"$match", match}},
	//			bson.D{
	//				{"$sample", bson.D{{"size", 1}}},
	//			},
	//		})
	//		if err != nil {
	//			return nil, false, false, false, err
	//		}
	//
	//		// cursor.Next()
	//		var docs []*SimulateData
	//		cursor.All(context.TODO(), &docs)
	//		playResp = docs[0]
	//		return playResp, false, false, false, err
	//	}
	//}
	if param.IsBuy || param.IsReSpin {
		if param.SelfPoolGold < 0 && rand.IntN(1000) < param.BuyMinAwardPercent {
			// 玩家奖池为负，购买小游戏强杀
			if param.MinBuyFunc == nil {
				err = errors.New("not found minBuyFunc")
				return
			}
			playResp, err = param.MinBuyFunc()
			buyKill = true
		} else {
			if param.NextBuyFunc == nil {
				err = errors.New("not found NextBuyFunc")
				return
			}
			playResp, err = param.NextBuyFunc(fmt.Sprintf("%d", param.Bet))
		}
		return
	}
	if param.Player.NoAwardPercent != 0 {
		if param.SelfPoolGold < 0 && rand.IntN(1000) < param.Player.NoAwardPercent {
			if param.SimulateByBucketIdFunc == nil {
				err = errors.New("not found SimulateByBucketIdFunc")
				return
			}
			playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
			forcedKill = true
			return
		}
	} else { //走平台设置
		if param.SelfPoolGold < 0 && rand.IntN(1000) < redisx.LoadNoAwardPercent(param.Player.AppID) {
			if param.SimulateByBucketIdFunc == nil {
				err = errors.New("not found SimulateByBucketIdFunc")
				return
			}
			playResp, err = param.SimulateByBucketIdFunc(0, param.AppStore.GamePatten)
			forcedKill = true
			return
		}
	}
	if hitPermillageArr(param.HitBigAwardPercent, param.Player.SpinCountOfThisEnter) {
		if param.Bet*40 < param.SelfPoolGold {
			avg := param.Player.CaclAvgBet()
			if avg*40 < param.SelfPoolGold {
				if param.BigRewardFunc == nil {
					err = errors.New("not found BigRewardFunc")
					return
				}
				playResp, err = param.BigRewardFunc(param.AppStore.GamePatten)
				hitBigAward = true
				return
			}
		}
	}
	if param.NextFunc == nil {
		err = errors.New("not found NextFunc")
		return
	}
	for {
		playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten)
		return
	}

}
