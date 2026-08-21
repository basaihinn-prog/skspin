package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_15001/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":0,"singleBetCombinations":{"10_10_5_NoExtraBet":500,"10_1_5_NoExtraBet":50,"10_2_5_NoExtraBet":100,"10_3_5_NoExtraBet":150,"10_5_5_NoExtraBet":250},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v2_15001_05_01_001","baseGameSetting":{"screenRow":3,"screenColumn":5,"symbolCount":13,"maxBetLine":25,"betSpec":{"waysBetList":[1,2,3,5,10],"waysBetColumnList":[5],"extraBetTypeList":["NoExtraBet"]},"specialFeatureCount":3,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,100,200,1000],[0,0,50,100,500],[0,0,40,80,400],[0,0,25,50,250],[0,0,10,20,100],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50]],"lineTable":[[1,1,1,1,1],[0,0,0,0,0],[2,2,2,2,2],[0,1,2,1,0],[2,1,0,1,2],[0,0,1,0,0],[2,2,1,2,2],[1,2,2,2,1],[1,0,0,0,1],[0,1,1,1,0],[2,1,1,1,2],[0,1,0,1,0],[2,1,2,1,2],[1,0,1,0,1],[1,2,1,2,1],[1,1,0,1,1],[1,1,2,1,1],[0,2,0,2,0],[2,0,2,0,2],[1,0,2,0,1],[1,2,0,2,1],[0,0,2,0,0],[2,2,0,2,2],[0,2,2,2,0],[2,0,0,0,2]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"WaysGame","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_12","specialHitInfo":["freeGame_01"],"basePay":50},{"specialHitPattern":"HP_13","specialHitInfo":["freeGame_02"],"basePay":10},{"specialHitPattern":"HP_14","specialHitInfo":["freeGame_03"],"basePay":5}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":63,"noWinIndex":[0],"wheelData":[2,2,2,7,7,1,11,11,5,5,8,3,3,12,10,10,6,6,6,7,7,5,5,5,12,12,4,4,4,9,9,6,8,8,5,5,10,10,3,3,8,8,4,4,9,9,5,5,10,10,1,9,9,3,7,7,5,5,8,8,4,7,7]},{"wheelLength":63,"noWinIndex":[0],"wheelData":[3,3,3,9,9,1,11,11,6,8,8,2,2,2,10,4,4,4,0,5,5,0,12,12,5,5,5,7,0,6,6,6,11,11,4,4,10,10,0,7,7,0,12,12,5,5,10,10,4,4,10,10,5,5,12,12,6,8,8,0,6,9,9]},{"wheelLength":63,"noWinIndex":[0],"wheelData":[4,4,4,10,10,1,11,11,5,5,8,8,6,6,6,12,12,3,3,3,0,4,4,0,7,3,3,9,9,5,5,5,12,12,3,3,9,9,0,7,7,0,11,11,5,5,8,8,2,2,2,7,7,6,9,9,0,6,8,8,6,10,10]},{"wheelLength":82,"noWinIndex":[0],"wheelData":[5,5,5,8,8,1,11,11,6,9,9,4,4,4,7,0,3,3,0,6,6,6,7,7,2,2,2,8,8,3,3,3,10,10,5,12,12,6,10,10,6,12,12,6,11,11,5,12,12,6,10,10,5,11,11,6,10,10,3,12,12,5,11,11,4,10,10,6,12,12,3,11,11,5,9,9,6,11,11,4,8,8]},{"wheelLength":56,"noWinIndex":[0],"wheelData":[6,6,6,11,11,1,12,12,5,5,5,8,8,4,4,10,10,6,12,12,6,7,7,6,9,9,6,7,7,3,3,3,9,9,4,4,4,11,11,2,2,2,2,7,6,11,11,6,9,9,6,8,8,6,11,11]}]]},"freeGameSetting":{"screenRow":3,"screenColumn":5,"symbolCount":13,"maxBetLine":25,"specialFeatureCount":3,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,100,200,1000],[0,0,50,100,500],[0,0,40,80,400],[0,0,25,50,250],[0,0,10,20,100],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50]],"lineTable":[[1,1,1,1,1],[0,0,0,0,0],[2,2,2,2,2],[0,1,2,1,0],[2,1,0,1,2],[0,0,1,0,0],[2,2,1,2,2],[1,2,2,2,1],[1,0,0,0,1],[0,1,1,1,0],[2,1,1,1,2],[0,1,0,1,0],[2,1,2,1,2],[1,0,1,0,1],[1,2,1,2,1],[1,1,0,1,1],[1,1,2,1,1],[0,2,0,2,0],[2,0,2,0,2],[1,0,2,0,1],[1,2,0,2,1],[0,0,2,0,0],[2,2,0,2,2],[0,2,2,2,0],[2,0,0,0,2]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"WaysGame","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_12","specialHitInfo":["reSpin_01"],"basePay":50},{"specialHitPattern":"HP_13","specialHitInfo":["reSpin_02"],"basePay":10},{"specialHitPattern":"HP_14","specialHitInfo":["reSpin_03"],"basePay":5}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":80,"noWinIndex":[0],"wheelData":[2,2,1,5,5,5,1,3,3,3,4,4,4,6,6,6,4,5,5,5,4,3,3,4,5,5,5,4,3,3,3,4,1,5,5,5,4,3,3,4,5,5,5,4,3,5,5,5,3,3,4,1,3,3,5,5,5,4,3,3,4,5,5,5,4,3,3,1,4,3,3,5,5,5,4,3,3,4,1,2]},{"wheelLength":80,"noWinIndex":[7],"wheelData":[3,3,1,2,2,2,1,4,4,4,6,6,6,4,3,6,4,0,5,5,5,0,3,3,1,5,5,5,0,3,3,6,4,3,3,0,5,5,5,4,2,2,2,4,5,5,5,4,3,3,4,6,3,3,6,4,5,5,5,4,6,6,4,5,5,5,4,6,3,3,4,6,3,3,4,5,5,5,4,3]},{"wheelLength":80,"noWinIndex":[7],"wheelData":[4,4,1,2,2,2,1,3,3,3,6,6,6,4,3,3,6,0,5,5,5,0,6,6,4,3,3,6,4,5,5,5,6,4,3,3,4,6,5,5,5,4,6,3,3,4,6,3,3,4,6,5,5,5,6,6,4,4,4,6,6,3,3,4,6,6,4,3,3,4,6,6,4,3,3,4,6,6,3,4]},{"wheelLength":100,"noWinIndex":[7],"wheelData":[5,5,1,4,4,4,1,6,6,6,3,3,3,4,6,3,3,0,5,5,5,0,6,6,4,3,3,6,5,5,5,6,6,6,2,2,2,6,5,5,5,6,4,3,3,3,4,6,2,2,2,6,4,3,3,3,4,6,5,5,5,4,6,2,2,2,4,6,3,3,3,6,4,2,2,2,4,1,2,2,2,1,4,2,2,2,4,6,6,2,2,2,4,6,2,2,2,6,4,5]},{"wheelLength":91,"noWinIndex":[17],"wheelData":[6,6,1,2,2,2,1,3,3,3,4,6,6,6,4,4,4,5,5,5,4,6,6,6,4,5,5,5,4,6,6,6,4,2,2,2,4,6,6,6,4,5,5,5,4,6,6,6,4,5,5,5,4,6,6,6,1,2,2,2,1,6,6,6,1,2,2,2,1,6,6,6,1,2,2,2,1,6,6,6,1,2,2,2,1,6,6,6,1,4,6]}]],"freeGameExtendSetting":{"defaultRound":10,"addRoundPerHit":10,"maxRound":200}},"bonusGameSetting":{},"doubleGameSetting":{"spinTimeLimit":5,"spinBetLimit":1000000000,"rtp":0.96,"tieRate":0.1}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":-1,"betCombinations":{"10_5_NoExtraBet":500,"1_5_NoExtraBet":50,"2_5_NoExtraBet":100,"3_5_NoExtraBet":150,"5_5_NoExtraBet":250},"gambleLimit":50000000,"defaultWaysBetColumnIdx":0}`
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
				key := fmt.Sprintf("10_%.0f_%d_NoExtraBet", info.Cs[i], internal.Line)
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
