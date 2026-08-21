package models

import (
	"context"
	"log/slog"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/slotspool"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// 删除用户的上一局， todo 应该把betlog 相关的也删掉
func (plr *Player) RewriteLastData() {
	_, err := db.Collection("players").UpdateOne(context.TODO(), bson.M{"_id": plr.PID}, bson.M{"$set": bson.M{"lastsid": "", "ls": ""}})
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

type Player struct {
	db.DocPlayer `bson:"inline"`

	LS string

	// SpinCountOfThisEnter int
	// SpinCount            int
	// BonusGameCount       int
	// BetHistory           []int64
	// BetAmount            int64
	// WinAmount            int64

	Cs        float64
	Ml        float64
	LastSid   string           //上一次的Sid,判断参数的正确性
	BdRecords []map[string]any //转完了后需要增加历史纪律
	IsBuy     bool             //是否是购买小游戏
	BigReward int64            //是否是大奖
	FpChoose  []string         //记录玩家小游戏的选择顺序
}

func (p *Player) UpdatePool(bet int64, selfPoolGold int64, toSelfAwardPool int64, forcedKill bool, app *redisx.AppStore) {
	if selfPoolGold < 0 && forcedKill {
		toSelfAwardPool += bet //个人奖池为负的时候，直接增加下注额，需求如此
	} else {
		rotateCount := p.SpinCount

		betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
		if p.RewardPercent != 0 {
			//走用户的
			betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
		}

		if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
			betMul += app.ProtectionRewardPercentLess
		}
		toSelfAwardPool += bet * int64(betMul) / 1000
	}
	slotsmongo.IncSelfSlotsPool(p.PID, toSelfAwardPool)
}

func (plr *Player) IsEndO() (isEnd bool, params []string) {
	if plr.LastSid != "" {
		params = strings.Split(plr.LastSid, "_")
		if params[1] == params[2] {
			isEnd = true
		}
	} else {
		isEnd = true
	}
	return
}
