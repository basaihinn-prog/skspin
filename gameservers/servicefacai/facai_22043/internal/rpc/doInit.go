package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicefacai/facai_22043/internal"
	"serve/servicefacai/facaicomm"
	"time"
)

func init() {
	facaicomm.RegRpc("init", doInit)
}

var (
	sampleDoInit = `{"betRange":[40,80,120,200,400,800,1200,2000,4000,8000,12000,20000,40000,80000,100000],"nickName":"demo977","userPoint":200000,"userPls":0,"features":3,"defaultBet":40,"probData":{"extraBetRange":[60,120,180,300,600,1200,1800,3000,6000,12000,18000,30000,60000,120000,150000],"featureBuyBetRange":[2000,4000,6000,10000,20000,40000,60000,100000],"featureBuyExtraBetRange":[3000,6000,9000,15000,30000,60000,90000,150000],"initPicture":[[[3,5,6,7],[8,11,2,12,9],[5,10,16,6,7],[8,13,2,14,9],[3,5,6,7]],[[4,5,6,7],[8,18,2,11,9],[5,4,13,6,7],[8,12,2,16,7],[5,5,6,8]]],"wheel":[[[2,10,8,6,3,9,4,5,5,4,3,8,6,6,3,8,5,4],[2,7,6,9,4,7,5,7,9,5,7,9,6,6,4,7,9,10],[2,9,10,6,13,7,10,10,4,6,9,10,7,6,9,7,6,9],[2,16,7,10,7,5,9,10,9,10,7,6,8,4,7,10,10,10],[2,10,9,3,7,10,7,7,9,10,6,6,6,8,10,10,9,7]],[[2,3,8,8,8,5,5,5,10,10,6,6,8,8,8,9,8,3],[2,5,9,3,6,15,4,10,9,3,9,4,6,9,11,7,5,9],[2,10,10,8,9,7,5,10,10,8,9,9,11,3,10,7,8,8],[2,7,9,6,3,4,9,5,3,10,10,13,9,7,6,5,7,3],[2,6,10,7,10,6,4,5,7,9,10,4,10,8,6,10,7,9]]],"freeTimesMax":50,"bigWinLevelOdds":[7,15,30]},"timeStamp":1744700847093,"testMode":false,"purchaseFeature":true,"extraJpBet":false,"jpBet":0}`
)

func doInit(msg *nats.Msg) (ret []byte, err error) {
	pid, _, err := facaicomm.ParsePidSfsObject(msg.Data)
	if err != nil {
		return nil, err
	}

	err = db.CallWithPlayer(pid, func(plr *facaicomm.Player) error {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(sampleDoInit), &data); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return err
		}
		info, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		// 修改数据
		betLevels := []int{}
		extraBetRange := []int{}
		featureBuyBetRange := []int{}
		featureBuyExtraBetRange := []int{}
		for i := range info.Cs {
			betLevels = append(betLevels, int(info.Cs[i]))
			extraBetRange = append(extraBetRange, int(info.Cs[i]*internal.InitExMul))
			if i < internal.InitBuyRange {
				featureBuyBetRange = append(featureBuyBetRange, int(info.Cs[i])*internal.InitBuyMul)
				featureBuyExtraBetRange = append(featureBuyExtraBetRange, int(info.Cs[i])*internal.InitExBuyMul)
			}
		}
		data["ts"] = time.Now().UnixNano()
		data["betRange"] = betLevels
		data["extraBetRange"] = extraBetRange
		data["featureBuyBetRange"] = featureBuyBetRange
		data["featureBuyExtraBetRange"] = featureBuyExtraBetRange
		data["nickName"] = pid
		balance, _ := slotsmongo.GetBalance(pid)
		data["userPoint"] = ut.HackGold2Money(balance)
		modifiedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}
		ret = modifiedJSON

		plr.SpinCountOfThisEnter = 0
		return nil
	})

	return
}
