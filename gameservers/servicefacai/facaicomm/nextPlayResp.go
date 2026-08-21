package facaicomm

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"serve/comm/redisx"
)

type NextPlayRespParam struct {
	Player             *Player
	Bet                int64
	SelfPoolGold       int64
	IsBuy              bool
	BuyMinAwardPercent int
	Combine            Combine
	NoAwardPercent     int
	HitBigAwardPercent []int
	AppStore           *redisx.AppStore
	Mod                int //模式购买相关
	OptionsSlots       int //可选插槽

	NextFunc               func(str string, witchBuckets int) (playResp *SimulateData, err error)
	SimulateByBucketIdFunc func(bucketId, witchBuckets int) (playResp *SimulateData, err error)
	BigRewardFunc          func(witchBuckets int) (playResp *SimulateData, err error)
	BigRewardFunc2_5       func(witchBuckets int) (playResp *SimulateData, err error)
	MinBuyFunc             func(mod int) (playResp *SimulateData, err error)
	NextBuyFunc            func(mod int, str string) (playResp *SimulateData, err error)
	ControlNextData        func(nextData *redisx.NextMulti, witchBucket int) (playResp *SimulateData, err error)
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
	if param.IsBuy {
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
			playResp, err = param.MinBuyFunc(param.Mod)
		} else {
			if param.NextBuyFunc == nil {
				err = errors.New("not found NextBuyFunc")
				return
			}
			playResp, err = param.NextBuyFunc(param.Mod, fmt.Sprintf("%d", param.Bet))
		}
		return
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
	for {
		playResp, err = param.NextFunc(fmt.Sprintf("%d", param.Bet), param.AppStore.GamePatten)
		//保底直接返回
		if param.AppStore.MaxMultiple == 0 || param.AppStore.MaxWinPoints*10000 == 0 {
			return
		}
		if playResp.Times < param.AppStore.MaxMultiple && playResp.Times*float64(param.Bet) <= param.AppStore.MaxWinPoints*10000 {
			return
		}
	}

}
