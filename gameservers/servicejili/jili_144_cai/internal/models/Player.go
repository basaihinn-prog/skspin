package models

import (
	"context"
	"log/slog"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/slotspool"
	"serve/servicejili/jili_144_cai/internal"

	"go.mongodb.org/mongo-driver/bson"
)

type Player struct {
	db.DocPlayer `bson:"inline"`

	// LS     string
	// LastID string
}

func (p *Player) UpdatePool(bet int64, selfPoolGold int64, toSelfAwardPool int64, forcedKill bool, ty int, app *redisx.AppStore) {
	if ty == internal.GameTypeNormal {
		if selfPoolGold < 0 && forcedKill {
			toSelfAwardPool += bet //个人奖池为负的时候，直接增加下注额，需求如此
		} else {
			betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
			if p.RewardPercent != 0 {
				//走用户的
				betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			}

			rotateCount := p.SpinCount
			if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
				betMul += app.ProtectionRewardPercentLess
			}
			toSelfAwardPool += bet * int64(betMul) / 1000
		}
	}
	if ty == internal.GameTypeExtra {
		if selfPoolGold < 0 && forcedKill {
			toSelfAwardPool += bet //个人奖池为负的时候，直接增加下注额，需求如此
		} else {
			betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
			if p.RewardPercent != 0 {
				//走用户的
				betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			}

			rotateCount := p.SpinCount
			if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
				betMul += app.ProtectionRewardPercentLess
			}
			toSelfAwardPool += bet * int64(betMul) / 1000
		}
	}
	if ty == internal.GameTypeExtraPlus {
		if selfPoolGold < 0 && forcedKill {
			toSelfAwardPool += bet //个人奖池为负的时候，直接增加下注额，需求如此
		} else {
			betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
			if p.RewardPercent != 0 {
				//走用户的
				betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			}

			rotateCount := p.SpinCount
			if app.IsProtection == 1 && rotateCount < app.ProtectionRotateCount {
				betMul += app.ProtectionRewardPercentLess
			}
			toSelfAwardPool += bet * int64(betMul) / 1000
		}
	}

	slotsmongo.IncSelfSlotsPool(p.PID, toSelfAwardPool)
}

// 删除用户的上一局， todo 应该把betlog 相关的也删掉
func (plr *Player) RewriteLastData() {
	var vm map[string]string
	_, err := db.Collection("players").UpdateOne(context.TODO(), bson.M{"_id": plr.PID}, bson.M{"$set": bson.M{"lastdata": vm}})
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
