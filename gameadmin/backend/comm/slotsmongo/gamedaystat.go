package slotsmongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SlotsPlayerGameDayStatRecord struct {
	Game           string
	Pid            int64
	Date           string
	Flag           string //特殊的Flag就是针对All，代表
	Flow           int64  //押注
	Win            int64  //赢分
	YeJi           int64  //业绩
	NetProfit      int64  //净输赢
	Fee            int64  //抽水
	SelfPoolReward int64  //个人奖池爆奖励
}

func IncAll(req *SlotsPlayerGameDayStatRecord) (err error) {
	req.Flag = "All"
	err = IncFlag(req)
	return
}

func IncBoth(req *SlotsPlayerGameDayStatRecord) (err error) {
	err = IncFlag(req)
	if err != nil {
		return
	}

	if req.Flag == "All" {
		return
	}

	req.Flag = "All"
	err = IncFlag(req)
	return
}

func IncFlag(req *SlotsPlayerGameDayStatRecord) (err error) {
	if req == nil {
		return
	}
	if req.Game == "" {
		err = fmt.Errorf("game is empty")
		return
	}
	if req.Pid == 0 {
		err = fmt.Errorf("Pid is zero")
		return
	}
	if req.Date == "" {
		err = fmt.Errorf("date is empty")
		return
	}
	if req.Flag == "" {
		err = fmt.Errorf("flag is empty")
		return
	}

	_, err = SlotsCollection("PlayerGameDayStat").UpdateOne(context.TODO(),
		bson.M{"Game": req.Game, "Pid": req.Pid, "Date": req.Date, "Flag": req.Flag},
		bson.M{"$inc": bson.M{
			"Flow":           req.Flow,
			"Win":            req.Win,
			"YeJi":           req.YeJi,
			"NetProfit":      req.NetProfit,
			"Fee":            req.Fee,
			"SelfPoolReward": req.SelfPoolReward,
		}},
		options.Update().SetUpsert(true))
	return
}
