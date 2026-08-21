package rpc

import (
	"context"
	"fmt"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/redisx"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"serve/service/pg_1312883/internal/gendata"
	"serve/service/pg_1312883/internal/models"
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

func nextPlayResp(player *models.Player, bet int64, selfPoolGold int64, isBuy bool, App *redisx.AppStore) (playResp *gendata.SimulateData, hitBigAward, forcedKill, buyKill bool, err error) {
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
		if err != nil {
			return nil, false, false, false, err
		}

		// cursor.Next()
		var docs []*gendata.SimulateData
		cursor.All(context.TODO(), &docs)
		playResp = docs[0]
		return playResp, false, false, false, err
	}
	if isBuy {
		minAwardPercent := player.BuyMinAwardPercent
		if minAwardPercent == 0 {
			minAwardPercent = App.BuyMinAwardPercent
		}
		if rand.IntN(1000) < minAwardPercent {
			playResp, err = gendata.GCombineDataMng.GetBuyMinData()
		} else {
			playResp, err = gendata.GCombineDataMng.NextBuy(fmt.Sprintf("%d", bet))
		}
		return
	}

	if nextData := redisx.GetPlayerNextData(player.PID); nextData != nil {
		playResp, err = gendata.GCombineDataMng.ControlNextData(nextData, db.BoundType(0), App.GamePatten)
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
		if bet*10 < selfPoolGold {
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
