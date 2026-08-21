package models

import (
	"context"
	"serve/comm/redisx"
	"serve/comm/slotspool"
	"strings"
	"time"

	"log/slog"
	"serve/comm/db"
	"serve/comm/slotsmongo"
	"serve/comm/ut"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (p *Player) UpdatePool(bet int64, selfPoolGold int64, toSelfAwardPool int64, forcedKill, isBuy, hitBigAward bool, allWin int64, app *redisx.AppStore) {
	if !isBuy {
		if selfPoolGold < 0 && forcedKill {
			toSelfAwardPool += bet //个人奖池为负的时候，直接增加下注额，需求如此
		} else {
			rotateCount := p.SpinCount
			//走平台的
			//修改start
			betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
			if p.RewardPercent != 0 {
				//走用户的
				betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			}

			//修改end

			if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
				betMul += app.ProtectionRewardPercentLess
			}
			toSelfAwardPool += bet * int64(betMul) / 1000
		}
	}

	after, err := slotsmongo.IncSelfSlotsPool(p.PID, toSelfAwardPool)

	if false {
		now := time.Now()
		hisdoc := db.D(
			"_id", primitive.NewObjectIDFromTimestamp(now),
			"pid", p.PID,
			"bet", bet,
			"allwin", allWin,
			"forcedkill", forcedKill,
			"hitBigAward", hitBigAward,
			"befer", selfPoolGold,
			"change", toSelfAwardPool,
			"after", after,
			"time", now,
			"func", "UpdatePool",
			"err", ut.ErrString(err),
		)

		his := db.Collection("PoolHistory")
		his.InsertOne(context.TODO(), hisdoc)
	}
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
