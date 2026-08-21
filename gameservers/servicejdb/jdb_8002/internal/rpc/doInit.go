package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_8002/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":-1,"singleBetCombinations":{"10_10_50_NoExtraBet":500,"10_1_50_NoExtraBet":50,"10_2_50_NoExtraBet":100,"10_3_50_NoExtraBet":150,"10_5_50_NoExtraBet":250},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":0,"gameFeatureCount":3,"executeSetting":{"settingId":"v2_8002_05_00_001","baseGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":11,"maxBetLine":50,"betSpec":{"lineBetList":[1,2,3,5,10],"betLineList":[50],"extraBetTypeList":["NoExtraBet"]},"specialFeatureCount":3,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,50,500,1000],[0,0,50,300,800],[0,0,30,150,600],[0,0,15,50,300],[0,0,15,50,200],[0,0,10,30,100],[0,0,10,30,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[[1,1,1,1,1],[2,2,2,2,2],[0,0,0,0,0],[3,3,3,3,3],[1,2,3,2,1],[2,1,0,1,2],[0,0,1,2,3],[3,3,2,1,0],[1,0,0,0,1],[2,3,3,3,2],[0,1,2,3,3],[3,2,1,0,0],[1,0,1,2,1],[2,3,2,1,2],[0,1,0,1,0],[3,2,3,2,3],[1,2,1,0,1],[2,1,2,3,2],[0,1,1,1,0],[3,2,2,2,3],[1,1,2,3,3],[2,2,1,0,0],[1,1,0,1,1],[2,2,3,2,2],[1,2,2,2,3],[2,1,1,1,0],[0,0,1,0,0],[3,3,2,3,3],[0,1,2,2,3],[3,2,1,1,0],[0,0,0,1,2],[3,3,3,2,1],[1,0,0,1,2],[2,3,3,2,1],[0,1,1,2,3],[3,2,2,1,0],[1,0,1,2,3],[2,3,2,1,0],[0,1,2,3,2],[3,2,1,0,1],[1,0,1,0,1],[2,3,2,3,2],[0,1,0,1,2],[3,2,3,2,1],[2,1,0,0,1],[1,2,3,3,2],[2,1,0,0,0],[1,2,3,3,3],[0,0,1,2,2],[3,3,2,1,1]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"LeftToRight","wheelUsePattern":"Independent","specialHitInfo":[{"specialHitPattern":"HP_01","specialHitInfo":["freeGame_01"],"basePay":0},{"specialHitPattern":"HP_02","specialHitInfo":["freeGame_02"],"basePay":0},{"specialHitPattern":"HP_03","specialHitInfo":["freeGame_03"],"basePay":0}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":29,"noWinIndex":[0],"wheelData":[5,7,8,9,10,4,7,7,10,9,4,9,10,5,9,3,8,6,3,8,9,4,9,6,10,2,9,8,6]},{"wheelLength":28,"noWinIndex":[0],"wheelData":[8,2,8,7,7,8,9,6,0,3,4,9,2,10,10,5,7,7,8,10,6,8,6,5,8,5,4,3]},{"wheelLength":26,"noWinIndex":[0],"wheelData":[5,7,8,3,6,5,6,4,0,5,9,4,7,0,8,10,10,2,3,7,8,9,6,9,1,6]},{"wheelLength":34,"noWinIndex":[0],"wheelData":[10,3,7,5,6,6,6,0,6,7,10,2,6,8,5,9,7,5,5,6,4,10,6,0,6,3,9,6,10,0,5,3,6,6]},{"wheelLength":31,"noWinIndex":[0],"wheelData":[3,5,0,10,6,6,8,7,6,6,9,6,7,6,9,6,10,6,6,8,8,10,5,5,2,9,6,4,5,6,5]}]]},"freeGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":11,"maxBetLine":50,"specialFeatureCount":0,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,50,500,1000],[0,0,50,300,800],[0,0,30,150,600],[0,0,15,50,300],[0,0,15,50,200],[0,0,10,30,100],[0,0,10,30,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[[1,1,1,1,1],[2,2,2,2,2],[0,0,0,0,0],[3,3,3,3,3],[1,2,3,2,1],[2,1,0,1,2],[0,0,1,2,3],[3,3,2,1,0],[1,0,0,0,1],[2,3,3,3,2],[0,1,2,3,3],[3,2,1,0,0],[1,0,1,2,1],[2,3,2,1,2],[0,1,0,1,0],[3,2,3,2,3],[1,2,1,0,1],[2,1,2,3,2],[0,1,1,1,0],[3,2,2,2,3],[1,1,2,3,3],[2,2,1,0,0],[1,1,0,1,1],[2,2,3,2,2],[1,2,2,2,3],[2,1,1,1,0],[0,0,1,0,0],[3,3,2,3,3],[0,1,2,2,3],[3,2,1,1,0],[0,0,0,1,2],[3,3,3,2,1],[1,0,0,1,2],[2,3,3,2,1],[0,1,1,2,3],[3,2,2,1,0],[1,0,1,2,3],[2,3,2,1,0],[0,1,2,3,2],[3,2,1,0,1],[1,0,1,0,1],[2,3,2,3,2],[0,1,0,1,2],[3,2,3,2,1],[2,1,0,0,1],[1,2,3,3,2],[2,1,0,0,0],[1,2,3,3,3],[0,0,1,2,2],[3,3,2,1,1]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"LeftToRight","wheelUsePattern":"Independent","specialHitInfo":[],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":40,"noWinIndex":[0,1,2,3],"wheelData":[5,7,5,5,10,4,7,7,7,9,6,9,4,10,9,3,5,6,10,8,9,6,8,3,10,2,9,10,9,6,10,6,5,8,10,8,9,8,7,7]},{"wheelLength":44,"noWinIndex":[0,1,2,3],"wheelData":[8,2,8,7,7,9,9,6,0,3,4,9,9,10,10,5,7,7,2,10,6,10,6,5,8,5,3,8,4,0,3,4,9,9,8,10,5,7,7,2,8,6,10,6]},{"wheelLength":31,"noWinIndex":[0,1,2,3],"wheelData":[5,7,8,3,6,10,9,4,0,2,9,4,4,8,10,10,5,3,7,8,6,9,7,9,3,3,10,6,7,8,5]},{"wheelLength":33,"noWinIndex":[0,1,2,3],"wheelData":[10,4,7,5,9,6,6,7,10,7,10,2,7,5,8,9,3,8,5,9,4,8,0,8,9,3,8,4,10,0,8,3,6]},{"wheelLength":44,"noWinIndex":[0,1,2,3],"wheelData":[3,5,0,10,10,6,8,7,6,6,9,10,7,9,7,7,10,4,4,8,8,7,5,8,2,9,10,4,8,6,7,6,3,9,10,7,9,9,9,10,6,5,8,8]}]],"freeGameExtendSetting":{"triggerCount5C1":25,"triggerCount4C1":15,"triggerCount3C1":10,"maxRound":25}},"bonusGameSetting":{},"doubleGameSetting":{"spinTimeLimit":5,"spinBetLimit":1000000000,"rtp":0.96,"tieRate":0.1}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":0,"betCombinations":{"10_50_NoExtraBet":500,"1_50_NoExtraBet":50,"2_50_NoExtraBet":100,"3_50_NoExtraBet":150,"5_50_NoExtraBet":250},"gambleLimit":50000000,"defaultWaysBetColumnIdx":-1}`
)

//pp.PutByteArray("entity", []byte(`{"maxBet":9223372036854775807,"defaultWaysBetIdx":-1,"singleBetCombinations":{"10_10_9_NoExtraBet":90,"10_1_9_NoExtraBet":9,"10_20_9_NoExtraBet":180,"10_30_9_NoExtraBet":270,"10_40_9_NoExtraBet":360,"10_50_9_NoExtraBet":450,"10_5_9_NoExtraBet":45},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":0,"defaultConnectBetIdx":-1,"defaultQuantityBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v3_14027_05_03_002","betSpecSetting":{"paymentType":"PT_001","extraBetTypeList":["NoExtraBet"],"betSpecification":{"lineBetList":[1,5,10,20,30,40,50],"betLineList":[9],"betType":"LineGame"}},"gameStateSetting":[{"gameStateType":"GS_003","frameSetting":{"screenColumn":3,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":55,"noWinIndex":[0],"wheelData":[3,7,5,0,7,6,2,8,6,1,9,6,9,5,0,8,0,2,6,1,9,4,7,7,2,6,0,3,7,4,8,1,9,3,8,5,0,9,8,3,7,5,1,4,8,2,7,8,0,9,6,3,8,9,1]},{"wheelLength":67,"noWinIndex":[0],"wheelData":[6,4,8,2,5,3,7,0,9,4,6,2,5,8,1,4,7,2,9,5,0,6,8,3,9,4,7,1,6,8,3,9,7,2,5,9,3,8,6,0,9,5,8,3,6,4,9,1,8,6,4,9,2,8,5,3,9,7,0,8,6,1,3,7,9,2,5]},{"wheelLength":62,"noWinIndex":[0],"wheelData":[8,4,3,8,6,4,1,5,7,2,6,4,9,0,7,5,7,2,9,4,5,9,1,6,9,4,3,9,3,0,9,4,7,2,8,5,9,1,8,4,5,2,7,7,0,5,7,3,9,4,1,9,3,6,7,2,5,9,0,5,4,8]}]]},"symbolSetting":{"symbolCount":10,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","M7","M8"],"payTable":[[0,0,100],[0,0,0],[0,0,50],[0,0,25],[0,0,15],[0,0,12],[0,0,8],[0,0,8],[0,0,3],[0,0,3]],"mixGroupCount":0,"mixGroupSetting":[]},"lineSetting":{"maxBetLine":9,"lineTable":[[1,1,1],[0,0,0],[2,2,2],[0,1,2],[2,1,0],[2,1,2],[0,1,0],[1,0,1],[1,2,1]]},"gameHitPatternSetting":{"gameHitPattern":"LineGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":1,"specialHitInfo":[{"specialHitPattern":"HP_05","triggerEvent":"Trigger_01","basePay":0}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":1,"addRound":0,"maxRound":1}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_01"]}}},{"gameStateType":"GS_069","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":48,"noWinIndex":[0],"wheelData":[5,3,7,6,5,2,8,7,9,4,8,7,9,0,9,7,8,5,9,7,3,8,7,5,8,6,9,2,8,7,4,9,7,3,6,9,5,7,9,0,6,7,4,5,9,7,8,8]},{"wheelLength":39,"noWinIndex":[0],"wheelData":[8,6,4,7,3,8,9,5,7,9,6,7,4,8,6,3,7,5,8,4,7,9,6,0,9,7,5,8,9,7,4,9,5,8,2,9,5,8,6]},{"wheelLength":44,"noWinIndex":[0],"wheelData":[7,9,2,6,1,7,7,0,6,5,8,1,7,9,4,6,5,7,8,4,9,1,7,4,8,0,9,5,1,8,6,3,7,9,3,9,5,7,9,1,8,4,9,6]},{"wheelLength":51,"noWinIndex":[0],"wheelData":[8,3,1,5,4,6,2,7,6,3,7,5,1,8,4,6,7,0,6,5,8,4,9,6,1,7,8,4,9,6,2,8,6,4,9,5,1,7,9,4,6,8,0,5,9,3,6,4,9,6,4]},{"wheelLength":58,"noWinIndex":[0],"wheelData":[6,5,9,4,0,6,5,7,1,8,4,9,6,2,8,5,3,6,9,5,8,4,9,5,8,0,6,4,7,3,9,5,6,1,8,4,9,6,2,8,6,4,7,9,1,9,6,3,8,5,9,2,6,9,3,6,4,6]}]]},"symbolSetting":{"symbolCount":10,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","M7","M8"],"payTable":[[0,0,100,200,1000],[0,0,0,0,0],[0,0,50,100,500],[0,0,25,50,200],[0,0,15,25,100],[0,0,12,25,100],[0,0,8,15,50],[0,0,8,15,50],[0,0,3,8,30],[0,0,3,5,30]],"mixGroupCount":0,"mixGroupSetting":[]},"lineSetting":{"maxBetLine":25,"lineTable":[[1,1,1,1,1],[0,0,0,0,0],[2,2,2,2,2],[0,1,2,1,0],[2,1,0,1,2],[2,1,2,2,2],[0,1,0,0,0],[1,0,1,2,2],[1,2,1,0,0],[0,1,1,1,2],[2,1,1,1,0],[1,1,0,1,1],[1,1,2,1,1],[1,0,0,0,1],[1,2,2,2,1],[0,0,1,2,1],[2,2,1,0,1],[0,0,1,2,2],[2,2,1,0,0],[0,0,0,1,2],[2,2,2,1,0],[2,1,0,0,0],[0,1,2,2,2],[0,1,2,1,2],[2,1,0,1,0]]},"gameHitPatternSetting":{"gameHitPattern":"LineGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":1,"specialHitInfo":[{"specialHitPattern":"HP_07","triggerEvent":"ReTrigger_01","basePay":0}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":5,"addRound":5,"maxRound":25}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_06"]}},"extendSetting":{"roundOddsRadix":2,"startPower":0,"maxRoundOdds":64}}],"doubleGameSetting":{"doubleRoundUpperLimit":5,"doubleBetUpperLimit":1000000000,"rtp":0.96,"tieRate":0.1},"boardDisplaySetting":{"winRankSetting":{"BigWin":38,"MegaWin":246,"UltraWin":550}},"gameFlowSetting":{"conditionTableWithoutBoardEnd":[["CD_False","CD_True","CD_False"],["CD_False","CD_False","CD_01"],["CD_False","CD_False","CD_False"]]}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":0,"betCombinations":{"10_9_NoExtraBet":90,"1_9_NoExtraBet":9,"20_9_NoExtraBet":180,"30_9_NoExtraBet":270,"40_9_NoExtraBet":360,"50_9_NoExtraBet":450,"5_9_NoExtraBet":45},"gambleLimit":0,"buyFeatureLimit":2147483647,"buyFeature":true,"defaultWaysBetColumnIdx":-1}`))

func doInit(msg *nats.Msg) (ret []byte, err error) {
	pid, _, err := jdbcomm.ParsePidSfsObject(msg.Data)
	if err != nil {
		return nil, err
	}

	err = db.CallWithPlayer(pid, func(plr *jdbcomm.Player) error {
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
		if betCombinations, ok := data["singleBetCombinations"].(map[string]interface{}); ok {
			betCombinations = make(map[string]interface{})
			for i := range info.Cs {
				key := fmt.Sprintf("10_%d_%d_NoExtraBet", int32(info.Cs[i]), internal.Line)
				betCombinations[key] = info.Cs[i] * internal.Line
			}
			data["singleBetCombinations"] = betCombinations
		}
		modifiedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}
		ret = modifiedJSON
		//curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		//todo 后台没写
		//c, err := redisx.GetPlayerCs(plr.AppID, plr.PID, true)
		//if err != nil {
		//	return err
		//}

		plr.SpinCountOfThisEnter = 0
		return nil
	})

	return
}
