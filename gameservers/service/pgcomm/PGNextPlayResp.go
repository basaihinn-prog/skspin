package pgcomm

import (
	"fmt"
	"math"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"strings"
)

type PGNextPlayRespStruct struct {
	Player             *PGPlayer
	Bet                int64
	IsBuy              bool
	SelfPoolGold       int64
	HitBigAwardPercent int
	App                *redisx.AppStore
	BoundType          db.BoundType
}

type GetBigRewardStruct_1 struct {
	Next            func(string, int) (*slotsmongo.SimulateData, error)
	NextBuy         func(string) (*slotsmongo.SimulateData, error)
	ControlNextData func(*redisx.NextMulti, db.BoundType, int) (*slotsmongo.SimulateData, error)
	SampleSimulate  func(int, int) (*slotsmongo.SimulateData, error)
	GetBigReward    func(int) (*slotsmongo.SimulateData, error)
	GetBigReward2_5 func(int) (*slotsmongo.SimulateData, error)
	GetBuyMinData   func() (*slotsmongo.SimulateData, error)
}

type GetBigRewardStruct_2 struct {
	Next            func(string, db.BoundType, int) (*slotsmongo.SimulateData, error)
	NextBuy         func(string) (*slotsmongo.SimulateData, error)
	ControlNextData func(*redisx.NextMulti, db.BoundType, int) (*slotsmongo.SimulateData, error)
	SampleSimulate  func(int, int, db.BoundType) (*slotsmongo.SimulateData, error)
	GetBigReward    func(db.BoundType, int) (*slotsmongo.SimulateData, error)
	GetBigReward2_5 func(db.BoundType, int) (*slotsmongo.SimulateData, error)
	GetBuyMinData   func() (*slotsmongo.SimulateData, error)
}

func NextFunc(PGPlayers *PGPlayer, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
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

func NextFunc1(PGPlayers *PGPlayer, Bet, RTP, MoniterNewbieNum int64, MoniterRtpRangeValue []redisx.MoniterRtpRangeValue) bool {
	NewbieValue := float64(rand.IntN(10000)) / 100
	MoniterRTPNode := math.Abs(float64(PGPlayers.MoniterRTPNode))
	for _, v := range MoniterRtpRangeValue {
		//rtp 区间
		if v.RangeMaxValue >= int64(MoniterRTPNode) && v.RangeMinValue <= int64(MoniterRTPNode) {
			//监控新手旋转数量
			if MoniterNewbieNum >= int64(PGPlayers.SpinCount) {
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

func TypeFunc[T GetBigRewardStruct_1 | GetBigRewardStruct_2](GetBigRewardStruct T, FuncName string, PGNextPlayRespS PGNextPlayRespStruct) (playResp *slotsmongo.SimulateData, err error) {
	switch FuncName {
	case "Next":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.Next(fmt.Sprintf("%d", PGNextPlayRespS.Bet), PGNextPlayRespS.App.GamePatten)
		case GetBigRewardStruct_2:
			playResp, err = v.Next(fmt.Sprintf("%d", PGNextPlayRespS.Bet), PGNextPlayRespS.BoundType, PGNextPlayRespS.App.GamePatten)
		}
	case "NextBuy":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.NextBuy(fmt.Sprintf("%d", PGNextPlayRespS.Bet))
		case GetBigRewardStruct_2:
			playResp, err = v.NextBuy(fmt.Sprintf("%d", PGNextPlayRespS.Bet))
		}
	case "GetBuyMinData":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.GetBuyMinData()
		case GetBigRewardStruct_2:
			playResp, err = v.GetBuyMinData()
		}
	case "ControlNextData":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			nextData := redisx.GetPlayerNextData(PGNextPlayRespS.Player.PID)
			playResp, err = v.ControlNextData(nextData, db.BoundType(0), PGNextPlayRespS.App.GamePatten)
		case GetBigRewardStruct_2:
			nextData := redisx.GetPlayerNextData(PGNextPlayRespS.Player.PID)
			playResp, err = v.ControlNextData(nextData, db.BoundType(0), PGNextPlayRespS.App.GamePatten)
		}
	case "SampleSimulate":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.SampleSimulate(0, PGNextPlayRespS.App.GamePatten)
		case GetBigRewardStruct_2:
			playResp, err = v.SampleSimulate(0, PGNextPlayRespS.App.GamePatten, PGNextPlayRespS.BoundType)
		}
	case "GetBigReward":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.GetBigReward(PGNextPlayRespS.App.GamePatten)
		case GetBigRewardStruct_2:
			playResp, err = v.GetBigReward(PGNextPlayRespS.BoundType, PGNextPlayRespS.App.GamePatten)
		}
	case "GetBigReward2_5":
		switch v := any(GetBigRewardStruct).(type) {
		case GetBigRewardStruct_1:
			playResp, err = v.GetBigReward2_5(PGNextPlayRespS.App.GamePatten)
		case GetBigRewardStruct_2:
			playResp, err = v.GetBigReward2_5(PGNextPlayRespS.BoundType, PGNextPlayRespS.App.GamePatten)
		}
	}
	return
}

func PGNextPlayRespFunc[T GetBigRewardStruct_1 | GetBigRewardStruct_2](PGNextPlayRespS PGNextPlayRespStruct, GetBigRewardStruct T) (playResp *slotsmongo.SimulateData, hitBigAward, forcedKill bool, err error) {
	defer func(p PGNextPlayRespStruct) {
		if playResp.Times != 0 {
			slotsmongo.UpdateGamePlayerMulti(p.Player.PID, playResp.Times)
		}
	}(PGNextPlayRespS)

	if PGNextPlayRespS.IsBuy {
		minAwardPercent := PGNextPlayRespS.Player.BuyMinAwardPercent
		if minAwardPercent == 0 {
			minAwardPercent = PGNextPlayRespS.App.BuyMinAwardPercent
		}
		if rand.IntN(1000) < minAwardPercent {
			playResp, err = TypeFunc(GetBigRewardStruct, "GetBuyMinData", PGNextPlayRespS)
		} else {
			playResp, err = TypeFunc(GetBigRewardStruct, "NextBuy", PGNextPlayRespS)
		}
		return
	}

	//限制
	if PGNextPlayRespS.Player.RestrictionsStatus == 1 {
		if PGNextPlayRespS.Player.Win >= PGNextPlayRespS.Player.RestrictionsMaxWin*10000 || PGNextPlayRespS.Player.Multi >= PGNextPlayRespS.Player.RestrictionsMaxMulti {
			playResp, err = TypeFunc(GetBigRewardStruct, "SampleSimulate", PGNextPlayRespS)
			forcedKill = true
			return
		}
	}

	nextData := redisx.GetPlayerNextData(PGNextPlayRespS.Player.PID)
	if nextData != nil {
		playResp, err = TypeFunc(GetBigRewardStruct, "ControlNextData", PGNextPlayRespS)
		return
	}

	//监控
	if PGNextPlayRespS.App.PersonalMoniterConfig.IsMoniter == 1 || PGNextPlayRespS.App.MoniterConfig.IsMoniter == 1 {
		RTP := PGNextPlayRespS.Player.TargetRTP
		if RTP == 0 {
			RTP = PGNextPlayRespS.App.RTP
		}

		PGNextPlayRespS.Player.MoniterCountNode++
		if PGNextPlayRespS.App.PersonalMoniterConfig.IsMoniter == 1 { //个人监控
			//监控周期
			if int64(PGNextPlayRespS.Player.MoniterCountNode) == PGNextPlayRespS.App.PersonalMoniterConfig.MoniterNumCycle {
				PGNextPlayRespS.Player.MoniterCountNode = 0
				//监控RTP节点
				PGNextPlayRespS.Player.MoniterRTPNode = 0

				if PGNextPlayRespS.Player.BetAmount != 0 {
					TMP := float64(PGNextPlayRespS.Player.WinAmount) / float64(PGNextPlayRespS.Player.BetAmount)
					PGNextPlayRespS.Player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
				}
			}

			//增加
			if PGNextPlayRespS.Player.MoniterRTPNode > 0 {
				ok := NextFunc(PGNextPlayRespS.Player, PGNextPlayRespS.Bet, int64(RTP), PGNextPlayRespS.App.PersonalMoniterConfig.MoniterNewbieNum, PGNextPlayRespS.App.PersonalMoniterConfig.MoniterAddRTPRangeValue)
				fmt.Println("增", ok)
				if ok {
					playResp, err = TypeFunc(GetBigRewardStruct, "GetBigReward2_5", PGNextPlayRespS)
					if sid := strings.Split(PGNextPlayRespS.Player.LastSid, "_"); len(sid) > 1 && sid[0] != playResp.Id.Hex() {
						hitBigAward = true
						return
					}
				}
			}

			//减少
			if PGNextPlayRespS.Player.MoniterRTPNode < 0 {
				ok := NextFunc1(PGNextPlayRespS.Player, PGNextPlayRespS.Bet, int64(RTP), PGNextPlayRespS.App.PersonalMoniterConfig.MoniterNewbieNum, PGNextPlayRespS.App.PersonalMoniterConfig.MoniterReduceRTPRangeValue)
				fmt.Println("减", ok)
				if ok {
					playResp, err = TypeFunc(GetBigRewardStruct, "SampleSimulate", PGNextPlayRespS)
					forcedKill = true
					return
				}
			}
		} else if PGNextPlayRespS.App.MoniterConfig.IsMoniter == 1 { //游戏监控
			//监控周期
			if int64(PGNextPlayRespS.Player.MoniterCountNode) >= PGNextPlayRespS.App.MoniterConfig.MoniterNumCycle {
				PGNextPlayRespS.Player.MoniterCountNode = 0
				//监控RTP节点
				PGNextPlayRespS.Player.MoniterRTPNode = 0

				if PGNextPlayRespS.Player.BetAmount != 0 {
					TMP := float64(PGNextPlayRespS.Player.WinAmount) / float64(PGNextPlayRespS.Player.BetAmount)
					PGNextPlayRespS.Player.MoniterRTPNode = int64(RTP) - int64(TMP*100)
				}
			}

			//增加
			if PGNextPlayRespS.Player.MoniterRTPNode > 0 {
				ok := NextFunc(PGNextPlayRespS.Player, PGNextPlayRespS.Bet, int64(RTP), PGNextPlayRespS.App.MoniterConfig.MoniterNewbieNum, PGNextPlayRespS.App.MoniterConfig.MoniterAddRTPRangeValue)
				fmt.Println("增", ok)
				if ok {
					playResp, err = TypeFunc(GetBigRewardStruct, "GetBigReward2_5", PGNextPlayRespS)
					if sid := strings.Split(PGNextPlayRespS.Player.LastSid, "_"); len(sid) > 1 && sid[0] != playResp.Id.Hex() {
						hitBigAward = true
						return
					}
				}
			}

			//减少
			if PGNextPlayRespS.Player.MoniterRTPNode < 0 {
				ok := NextFunc1(PGNextPlayRespS.Player, PGNextPlayRespS.Bet, int64(RTP), PGNextPlayRespS.App.MoniterConfig.MoniterNewbieNum, PGNextPlayRespS.App.MoniterConfig.MoniterReduceRTPRangeValue)
				fmt.Println("减", ok)
				if ok {
					playResp, err = TypeFunc(GetBigRewardStruct, "SampleSimulate", PGNextPlayRespS)
					forcedKill = true
					return
				}
			}
		}
	}

	//强杀
	if PGNextPlayRespS.Player.NoAwardPercent != 0 {
		//走用户设置
		if PGNextPlayRespS.SelfPoolGold < 0 && rand.IntN(1000) < PGNextPlayRespS.Player.NoAwardPercent {
			playResp, err = TypeFunc(GetBigRewardStruct, "SampleSimulate", PGNextPlayRespS)
			forcedKill = true
			return
		}
	} else {
		//走平台设置
		if PGNextPlayRespS.SelfPoolGold < 0 && rand.IntN(1000) < redisx.LoadNoAwardPercent(PGNextPlayRespS.Player.AppID) {
			playResp, err = TypeFunc(GetBigRewardStruct, "SampleSimulate", PGNextPlayRespS)
			forcedKill = true
			return
		}
	}

	//大奖
	if rand.IntN(1000) < PGNextPlayRespS.HitBigAwardPercent {
		if PGNextPlayRespS.Bet*10 < PGNextPlayRespS.SelfPoolGold {
			avg := PGNextPlayRespS.Player.CaclAvgBet()
			if avg*10 < PGNextPlayRespS.SelfPoolGold {
				for i := 0; i < 5; i++ {
					playResp, err = TypeFunc(GetBigRewardStruct, "GetBigReward", PGNextPlayRespS)
					if sid := strings.Split(PGNextPlayRespS.Player.LastSid, "_"); len(sid) > 1 && sid[0] != playResp.Id.Hex() {
						hitBigAward = true
						return
					}
				}
			}
		}
	}

	//个人限制
	if PGNextPlayRespS.Player.PersonWinMaxMult != 0 && PGNextPlayRespS.Player.PersonWinMaxScore != 0 {
		for {
			playResp, err = TypeFunc(GetBigRewardStruct, "Next", PGNextPlayRespS)
			if playResp.Times < float64(PGNextPlayRespS.Player.PersonWinMaxMult) && playResp.Times*float64(PGNextPlayRespS.Bet) <= float64(PGNextPlayRespS.Player.PersonWinMaxScore)*10000 {
				fmt.Println("个人限制", playResp.Times)
				return
			}
		}
	}

	//机台限制
	if PGNextPlayRespS.App.MaxMultiple != 0 && PGNextPlayRespS.App.MaxWinPoints != 0 {
		for {
			playResp, err = TypeFunc(GetBigRewardStruct, "Next", PGNextPlayRespS)
			if playResp.Times < PGNextPlayRespS.App.MaxMultiple && playResp.Times*float64(PGNextPlayRespS.Bet) <= PGNextPlayRespS.App.MaxWinPoints*10000 {
				fmt.Println("机台限制", playResp.Times)
				return
			}
		}
	}

	playResp, err = TypeFunc(GetBigRewardStruct, "Next", PGNextPlayRespS)
	return

}
