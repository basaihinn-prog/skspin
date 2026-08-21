package rpc

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicejili/jili_226_witch/internal"
	"serve/servicejili/jili_226_witch/internal/gendata"
	"serve/servicejili/jili_226_witch/internal/models"

	"serve/comm/db"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NextFunc(JILIPlayers *models.Player, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	PLRRTP := (Bet*5 + JILIPlayers.WinAmount) / JILIPlayers.BetAmount

	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= JILIPlayers.MoniterRTPNode && v.RangeMinValue <= JILIPlayers.MoniterRTPNode {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(JILIPlayers.SpinCount) {
				//新手概率
				if v.NewbieValue > NewbieValue {
					if PLRRTP < RTP {
						return true
					}
				}
			} else {
				//老手概率
				if v.NotNewbieValue > NewbieValue {
					if PLRRTP < RTP {
						return true
					}
				}
			}
		}
	}

	return false
}

func NextFunc1(JILIPlayers *models.Player, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	MoniterRTPNode := math.Abs(float64(JILIPlayers.MoniterRTPNode))
	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= int64(MoniterRTPNode) && v.RangeMinValue <= int64(MoniterRTPNode) {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(JILIPlayers.SpinCount) {
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
func hitPermillageArr(arr []int, i int) bool {
	if i < len(arr) {
		w := arr[i]
		if rand.IntN(1000) < w {
			return true
		}
	}
	return false
}

func nextPlayResp(player *models.Player, bet int64, selfPoolGold int64, ty int, App *redisx.AppStore) (playResp *models.RawSpin, hitBigAward, forcedKill, buyKill bool, err error) {
	defer func(p *models.Player) {
		if playResp.Times != 0 {
			slotsmongo.UpdateGamePlayerMulti(p.PID, playResp.Times)
		}
	}(player)

	if false {
		coll := db.Collection("rawSpinData")
		match := bson.M{}
		if ty == internal.GameTypeGame {
			match["type"] = 1
		}
		if ty == internal.GameTypeExtra {
			match["type"] = 2
		}
		if ty == internal.GameTypeNormal {
			match["type"] = 0
		}

		cur := lo.Must(coll.Aggregate(context.TODO(), mongo.Pipeline{
			bson.D{{"$match", match}},
			db.D("$sample", db.D("size", 1)),
		}))
		var spindocs []*models.RawSpin
		lo.Must0(cur.All(context.TODO(), &spindocs))
		lo.Must0(len(spindocs) == 1)
		spindoc := spindocs[0]
		return spindoc, false, false, false, err
	}

	if ty == internal.GameTypeGame {
		minAwardPercent := player.BuyMinAwardPercent
		if minAwardPercent == 0 {
			minAwardPercent = App.BuyMinAwardPercent
		}
		if rand.IntN(1000) < minAwardPercent {
			playResp, err = gendata.GCombineDataMng.GetBuyMinData()
			buyKill = true
		} else {
			playResp, err = gendata.GCombineDataMng.NextBuy(fmt.Sprintf("%d", bet))
		}
		return
	}

	//限制
	if player.RestrictionsStatus == 1 {
		if player.Win >= player.RestrictionsMaxWin || player.Multi >= player.RestrictionsMaxMulti {
			playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
			forcedKill = true
			return
		}
	}

	if ty == internal.GameTypeNormal {
		if nextData := redisx.GetPlayerNextData(player.PID); nextData != nil {
			playResp, err = gendata.GCombineDataMng.ControlNextData(nextData, ty, App.GamePatten)
			return
		}

		if App.PersonalMoniterConfig.IsMoniter == 1 || App.MoniterConfig.IsMoniter == 1 {
			RTP := player.TargetRTP
			if RTP == 0 {
				RTP = App.RTP
			}
			player.MoniterCountNode++

			fmt.Println("bet", bet)
			fmt.Println("player.BetAmount", player.BetAmount)
			fmt.Println("player.WinAmount", player.WinAmount)

			if App.PersonalMoniterConfig.IsMoniter == 1 { //个人监控
				//监控周期
				if int64(player.MoniterCountNode) == App.PersonalMoniterConfig.MoniterNumCycle {
					player.MoniterCountNode = 0
					//监控RTP节点
					player.MoniterRTPNode = 0
					if player.BetAmount != 0 {
						TMP := float64(player.WinAmount) / float64(player.BetAmount)
						player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
					}
				}
				//增加
				if player.MoniterRTPNode > 0 {
					ok := NextFunc(player, bet, int64(RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterAddRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.GetBigReward2_5(ty, App.GamePatten)
						hitBigAward = true
						return
					}
				}

				//减少
				if player.MoniterRTPNode < 0 {
					ok := NextFunc1(player, bet, int64(RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterReduceRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
						forcedKill = true
						return
					}
				}
			} else if App.MoniterConfig.IsMoniter == 1 { //游戏监控
				//监控周期
				if int64(player.MoniterCountNode) == App.MoniterConfig.MoniterNumCycle {
					player.MoniterCountNode = 0
					//监控RTP节点
					player.MoniterRTPNode = 0
					if player.BetAmount != 0 {
						TMP := float64(player.WinAmount) / float64(player.BetAmount)
						player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
					}
				}
				//增加
				if player.MoniterRTPNode > 0 {
					ok := NextFunc(player, bet, int64(RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterAddRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.GetBigReward2_5(ty, App.GamePatten)
						hitBigAward = true
						return
					}
				}
				//减少
				if player.MoniterRTPNode < 0 {
					ok := NextFunc1(player, bet, int64(RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterReduceRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
						forcedKill = true
						return
					}
				}
			}
		}

		noAwardPercent := player.NoAwardPercent
		if noAwardPercent == 0 {
			noAwardPercent = redisx.LoadNoAwardPercent(player.AppID)
		}
		if selfPoolGold < 0 && rand.IntN(1000) < noAwardPercent {
			playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
			forcedKill = true
			return
		}

		if rand.IntN(1000) < gendata.HitBigAwardPercent {
			fmt.Println("enter")
			// 同时满足以下2个条件
			// 1.个人奖池金额>=最近200把投注均值不包含购买小游戏*10
			// 2.个人奖池金额>=当前投注额*40
			// 可以爆个人奖池为20-40倍的一个小游戏奖励，该奖励不消耗库存
			if bet*10 < selfPoolGold {
				avg := player.CaclAvgBet()
				if avg*10 < selfPoolGold {
					playResp, err = gendata.GCombineDataMng.GetBigReward(ty, App.GamePatten)
					hitBigAward = true
					return
				}
			}
		}
	}
	if ty == internal.GameTypeExtra {
		if nextData := redisx.GetPlayerNextData(player.PID); nextData != nil {
			playResp, err = gendata.GCombineDataMng.ControlNextData(nextData, ty, App.GamePatten)
			return
		}

		if App.PersonalMoniterConfig.IsMoniter == 1 || App.MoniterConfig.IsMoniter == 1 {
			RTP := player.TargetRTP
			if RTP == 0 {
				RTP = App.RTP
			}
			player.MoniterCountNode++

			if App.PersonalMoniterConfig.IsMoniter == 1 { //个人监控
				//监控周期
				if int64(player.MoniterCountNode) == App.PersonalMoniterConfig.MoniterNumCycle {
					player.MoniterCountNode = 0
					//监控RTP节点
					player.MoniterRTPNode = 0
					if player.BetAmount != 0 {
						TMP := float64(player.WinAmount) / float64(player.BetAmount)
						player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
					}
				}
				//增加
				if player.MoniterRTPNode > 0 {
					ok := NextFunc(player, bet, int64(RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterAddRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.GetBigReward2_5(ty, App.GamePatten)
						hitBigAward = true
						return
					}
				}

				//减少
				if player.MoniterRTPNode < 0 {
					ok := NextFunc1(player, bet, int64(RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterReduceRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
						forcedKill = true
						return
					}
				}
			} else if App.MoniterConfig.IsMoniter == 1 { //游戏监控
				//监控周期
				if int64(player.MoniterCountNode) == App.MoniterConfig.MoniterNumCycle {
					player.MoniterCountNode = 0
					//监控RTP节点
					player.MoniterRTPNode = 0
					if player.BetAmount != 0 {
						TMP := float64(player.WinAmount) / float64(player.BetAmount)
						player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
					}
				}
				//增加
				if player.MoniterRTPNode > 0 {
					ok := NextFunc(player, bet, int64(RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterAddRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.GetBigReward2_5(ty, App.GamePatten)
						hitBigAward = true
						return
					}
				}
				//减少
				if player.MoniterRTPNode < 0 {
					ok := NextFunc1(player, bet, int64(RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterReduceRTPRangeValue)
					if ok {
						playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
						forcedKill = true
						return
					}
				}
			}
		}

		noAwardPercent := player.NoAwardPercent
		if noAwardPercent == 0 {
			noAwardPercent = redisx.LoadNoAwardPercent(player.AppID)
		}
		if selfPoolGold < 0 && rand.IntN(1000) < noAwardPercent {
			playResp, err = gendata.GCombineDataMng.SampleForceSimulate(ty, App.GamePatten)
			forcedKill = true
			return
		}

		if rand.IntN(1000) < gendata.HitBigAwardPercent {
			fmt.Println("enter")
			// 同时满足以下2个条件
			// 1.个人奖池金额>=最近200把投注均值不包含购买小游戏*10
			// 2.个人奖池金额>=当前投注额*40
			// 可以爆个人奖池为20-40倍的一个小游戏奖励，该奖励不消耗库存
			if bet*10 < selfPoolGold {
				avg := player.CaclAvgBet()
				if avg*10 < selfPoolGold {
					playResp, err = gendata.GCombineDataMng.GetBigReward(ty, App.GamePatten)
					hitBigAward = true
					return
				}
			}
		}
	}

	//个人
	if player.PersonWinMaxMult != 0 && player.PersonWinMaxScore != 0 {
		for {
			playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), ty, App.GamePatten)
			if playResp.Times < float64(player.PersonWinMaxMult) && playResp.Times*float64(bet) <= float64(player.PersonWinMaxScore)*10000 {
				return
			}
		}
	}

	//机台
	if App.MaxMultiple != 0 && App.MaxWinPoints != 0 {
		for {
			playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), ty, App.GamePatten)
			if playResp.Times < App.MaxMultiple && playResp.Times*float64(bet) <= App.MaxWinPoints*10000 {
				return
			}
		}
	}

	playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), ty, App.GamePatten)
	return
}
