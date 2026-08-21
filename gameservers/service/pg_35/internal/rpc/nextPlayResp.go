package rpc

import (
	"context"
	"fmt"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/redisx"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"serve/service/pg_35/internal/gendata"
	"serve/service/pg_35/internal/models"
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

func nextPlayResp(player *models.Player, bet int64, selfPoolGold int64, isBuy bool, App *redisx.AppStore) (playResp *gendata.SimulateData, hitBigAward, forcedKill bool, err error) {
	if false {
		coll := db.Collection("simulate")
		match := bson.M{}
		if isBuy {
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
		// id, _ := primitive.ObjectIDFromHex("6604d1cce04f4526f342ac8a")
		// cursor, _ := coll.Find(context.TODO(), db.ID(id))
		if err != nil {
			return nil, false, false, err
		}

		// cursor.Next()
		var docs []*gendata.SimulateData
		cursor.All(context.TODO(), &docs)
		playResp = docs[0]
		return playResp, false, false, err
	}

	if isBuy {
		playResp, err = gendata.GCombineDataMng.NextBuy(fmt.Sprintf("%d", bet))
		return
	}

	if nextData := redisx.GetPlayerNextData(player.PID); nextData != nil {
		boundType := db.BoundType(0)
		playResp, err = gendata.GCombineDataMng.ControlNextData(nextData, boundType, App.GamePatten)
		return
	}

	if player.NoAwardPercent != 0 {
		//走用户设置
		if selfPoolGold < 0 && rand.IntN(1000) < player.NoAwardPercent {
			playResp, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
			forcedKill = true
			return
		}
	} else { //走平台设置
		if selfPoolGold < 0 && rand.IntN(1000) < redisx.LoadNoAwardPercent(player.AppID) {
			playResp, err = gendata.GCombineDataMng.SampleSimulate(0, App.GamePatten)
			forcedKill = true
			return
		}
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
				//兜底，极限情况只有一条数据
				for i := 0; i < 5; i++ {
					playResp, err = gendata.GCombineDataMng.GetBigReward(App.GamePatten)
					if sid := strings.Split(player.LastSid, "_"); len(sid) > 1 && sid[0] != playResp.Id.Hex() {
						hitBigAward = true
						return
					}
				}
			}
		}
	}

	for {
		playResp, err = gendata.GCombineDataMng.Next(fmt.Sprintf("%d", bet), App.GamePatten)
		//保底直接返回
		if App.MaxMultiple == 0 || App.MaxWinPoints*10000 == 0 {
			return
		}
		//当大于平台设置的倍数时丢弃，再取一条（不需要管++逻辑，就是丢弃）

		if playResp.Times < App.MaxMultiple && playResp.Times*float64(bet) <= App.MaxWinPoints*10000 {
			//fmt.Println("==============================")
			//fmt.Println(playResp.Times, len(playResp.DropPan), playResp.Id, playResp)
			return
			//break
		}
	}

}
