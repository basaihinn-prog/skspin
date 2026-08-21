package ppcomm

import (
	"context"
	"fmt"
	"log/slog"
	"serve/comm/redisx"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"
	"serve/comm/slotsmongo"
	"serve/comm/slotspool"

	"github.com/samber/lo"
)

type Player struct {
	db.DocPlayer `bson:"inline"`
	LastData     Variables
	IsBuy        bool
	C            float64
	CostGold     int64
	BigReward    int64 // 大奖
	IsCollect    bool
	LastSid      string //最近一次的游戏结果id，SimulateData表id
	LastIndex    int    //最近一次的游戏结果下标，SimulateData中DropPan下标
	IsEnd        bool   //最近一次的游戏是否为最后一轮

	EnterBalance int64 //用户进入游戏的时候的钱包金额

	FRBData *FRBData
}

type FRBData struct {
	Config      *FRBPlayer
	Fra         float64
	Frn         int
	ReplenishEV bool
}

func (d *FRBData) IsValid() bool {
	if d == nil {
		return false
	}

	if d.Frn <= 0 {
		return false
	}

	// if d.Config.ExpirationDate < time.Now().Unix() {
	// 	return false
	// }

	return true
}

func (d *FRBData) EVFinish(line int) string {
	lo.Must0(line != 0)
	ev := "FR1~%.2f,%d,%.2f,,"

	frb := d.Config
	return fmt.Sprintf(ev, frb.TotalBet/float64(line), line, d.Fra)
}

func (d *FRBData) EVStart(line int) string {
	lo.Must0(line != 0)
	ev := "FR0~%.2f,%d,%d,0,0,%d,1,,"

	frb := d.Config
	return fmt.Sprintf(ev, frb.TotalBet/float64(line), line, d.Frn, frb.ExpirationDate)
	// return d.Config.EVStart(line)
}

type UpdatePoolParam struct {
	Bet                  int64
	SelfPoolGold         int64
	ToSelfAwardPool      int64
	ForcedKill           bool
	IsBuy                bool
	RewardPercentLess100 int
	App                  *redisx.AppStore
}

// func (p *Player) GetFRBData( /*totalBet float64*/ ) *FRBData {
// 	if !p.FRBData.IsValid() {
// 		p.FRBData = nil
// 	}

// 	return p.FRBData
// }

func (p *Player) UpdatePool(param *UpdatePoolParam) {
	if param == nil {
		return
	}
	if !param.IsBuy {
		if param.SelfPoolGold < 0 && param.ForcedKill {
			param.ToSelfAwardPool += param.Bet //个人奖池为负的时候，直接增加下注额，需求如此
		} else {
			rotateCount := p.SpinCount
			betMul := redisx.LoadAwardPercent(p.AppID) + slotspool.GetSlotsPool(p.AppID).Value
			if p.RewardPercent != 0 {
				//走用户的
				betMul = p.RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			}
			//源码为配置文件类型
			//betMul := slotspool.GetGameSlotsPool(lazy.ServiceName).RewardPercent + slotspool.GetSlotsPool(p.AppID).Value
			if param.App.IsProtection == 1 && rotateCount < param.App.ProtectionRotateCount {
				betMul += param.App.ProtectionRewardPercentLess
			}
			param.ToSelfAwardPool += param.Bet * int64(betMul) / 1000
		}
	}

	slotsmongo.IncSelfSlotsPool(p.PID, param.ToSelfAwardPool)
}

func (plr *Player) IsEndO() (isEnd bool, params []string) {
	var lastId string
	if plr.LastData != nil {
		lastId = plr.LastData.Str("gid")
	}
	if lastId != "" {
		params = strings.Split(lastId, "_")
		if params[1] == params[2] && plr.LastData.Str("na") == "s" {
			isEnd = true
		}
	} else {
		isEnd = true
	}

	return
}

// 删除用户的上一局， todo 应该把betlog 相关的也删掉
func (plr *Player) RewriteLastData() {
	var vm Variables
	_, err := db.Collection("players").UpdateOne(context.TODO(), bson.M{"_id": plr.PID}, bson.M{"$set": bson.M{"lastdata": vm}})
	if err != nil {
		slog.Error(err.Error())
		return

	}

}
