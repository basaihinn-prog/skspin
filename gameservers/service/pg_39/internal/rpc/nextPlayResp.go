package rpc

import (
	"fmt"
	"math"
	"math/rand/v2"
	"serve/comm/redisx"
	"serve/service/pg_39/internal/gendata"
	"serve/service/pg_39/internal/models"

	"go.mongodb.org/mongo-driver/bson"
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

var bigRewardBucketId int

func init() {
	// bigRewardBucketId = gendata.GetBucketIds(10, 20, true)[0]
	bigRewardBucketId = 6
}

func NextFunc(PGPlayers *models.Player, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	PLRRTP := ((float64(Bet*5) + float64(PGPlayers.WinAmount)) / float64(PGPlayers.BetAmount)) * 100

	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= PGPlayers.MoniterRTPNode && v.RangeMinValue <= PGPlayers.MoniterRTPNode {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(PGPlayers.SpinCount) {
				//新手概率
				if v.NewbieValue >= NewbieValue {
					if int64(PLRRTP) < RTP {
						return true
					}
				}
			} else {
				//老手概率
				if v.NotNewbieValue >= NewbieValue {
					if int64(PLRRTP) < RTP {
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

func nextPlayResp(player *models.Player, bet int64, selfPoolGold int64, App *redisx.AppStore) (playResp *bson.M, hitBigAward, forcedKill bool, err error) {
	// primitive.NewObjectIDFromTimestamp()

	nextData := redisx.GetPlayerNextData(player.PID)
	if nextData != nil {
		playResp, err = gendata.GCombineDataMng.ControlNextData(nextData, App.GamePatten)
		return
	}

	//监控
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
				ok := NextFunc(player, bet, int64(App.RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterAddRTPRangeValue)
				if ok {
					playResp, err = gendata.GCombineDataMng.GetBigReward2_5(App.GamePatten)
					hitBigAward = true
					return
				}
			}

			//减少
			if player.MoniterRTPNode < 0 {
				ok := NextFunc(player, bet, int64(App.RTP), App.PersonalMoniterConfig.MoniterNewbieNum, App.PersonalMoniterConfig.MoniterReduceRTPRangeValue)
				if ok {
					playResp, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
					forcedKill = true
					return
				}
			}
		} else if App.MoniterConfig.IsMoniter == 1 { //游戏监控
			//增加
			if int64(player.MoniterCountNode) == App.MoniterConfig.MoniterNumCycle {
				player.MoniterCountNode = 0
				//监控RTP节点
				player.MoniterRTPNode = 0

				if player.BetAmount != 0 {
					TMP := float64(player.WinAmount) / float64(player.BetAmount)
					player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
				}
			}

			if player.MoniterRTPNode > 0 {
				ok := NextFunc(player, bet, int64(App.RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterAddRTPRangeValue)
				if ok {
					playResp, err = gendata.GCombineDataMng.GetBigReward2_5(App.GamePatten)
					hitBigAward = true
					return
				}
			}
			//减少
			if player.MoniterRTPNode < 0 {
				ok := NextFunc(player, bet, int64(App.RTP), App.MoniterConfig.MoniterNewbieNum, App.MoniterConfig.MoniterReduceRTPRangeValue)
				if ok {
					playResp, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
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
		playResp, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
		forcedKill = true
		return
	}
	if rand.IntN(1000) < gendata.HitBigAwardPercent {
		fmt.Println("enter")
		// 同时满足以下2个条件
		// 1.个人奖池金额>=最近200把投注均值不包含购买小游戏*10
		// 2.个人奖池金额>=当前投注额*40
		// 可以爆个人奖池为20-40倍的一个小游戏奖励，该奖励不消耗库存

		// bet := Grades[grade]
		if bet*10 < selfPoolGold {
			// var avg int64
			// sum := lo.Sum(player.BetHistory)
			// if len(player.BetHistory) > 0 {
			// 	avg = sum / int64(len(player.BetHistory))
			// }
			avg := player.CaclAvgBet()
			if avg*10 < selfPoolGold {
				playResp, err = gendata.GCombineDataMng.SampleSimulate(bigRewardBucketId, App.GamePatten)
				hitBigAward = true
				return
			}
		}
	}

	if player.PersonWinMaxMult != 0 || player.PersonWinMaxScore != 0 {
		playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), App.GamePatten)
		if (*playResp)["aw"].(float64) < float64(player.PersonWinMaxMult) && (*playResp)["aw"].(float64)*float64(bet) <= float64(player.PersonWinMaxScore)*10000 {
			return
		}
	}

	if App.MaxMultiple != 0 || App.MaxWinPoints != 0 {
		playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), App.GamePatten)
		if (*playResp)["aw"].(float64) < App.MaxMultiple && (*playResp)["aw"].(float64)*float64(bet) <= App.MaxWinPoints*10000 {
			return
		}
	}

	playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), App.GamePatten)
	return
}
