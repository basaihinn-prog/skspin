package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_8003/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":0,"singleBetCombinations":{"10_10_5_NoExtraBet":500,"10_1_5_NoExtraBet":50,"10_2_5_NoExtraBet":100,"10_3_5_NoExtraBet":150,"10_5_5_NoExtraBet":250},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v2_8003_05_01_001","baseGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":13,"maxBetLine":0,"betSpec":{"waysBetList":[1,2,3,5,10],"waysBetColumnList":[5],"extraBetTypeList":["NoExtraBet"]},"specialFeatureCount":2,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,0,0,0],[0,0,50,200,1000],[0,0,25,150,400],[0,0,25,150,400],[0,0,20,75,200],[0,0,20,75,200],[0,0,10,50,150],[0,0,10,50,150],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[],"symbolAttribute":["Wild","BonusGame","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"WaysGame","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_05","specialHitInfo":["freeGame_01"],"basePay":2},{"specialHitPattern":"HP_06","specialHitInfo":["bonusGame_01"],"basePay":0}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":112,"noWinIndex":[0],"wheelData":[12,2,9,6,8,3,10,7,12,6,9,6,11,7,9,7,11,2,9,9,9,3,11,9,12,3,10,12,6,8,8,2,11,9,8,5,8,11,12,4,9,7,11,4,11,5,9,6,9,7,8,12,11,4,9,8,11,4,11,4,11,9,6,8,10,5,9,5,8,11,3,9,8,2,11,12,8,6,11,11,11,7,8,12,11,5,10,8,11,2,9,9,12,3,11,8,12,6,11,11,4,8,11,3,11,9,12,5,11,8,12,5]},{"wheelLength":112,"noWinIndex":[0],"wheelData":[11,4,8,7,12,0,10,5,10,12,8,7,10,2,10,5,10,7,8,8,8,5,9,7,0,6,10,7,11,3,8,5,8,7,12,6,10,3,12,7,9,5,10,10,10,4,12,7,11,4,10,8,12,2,9,8,8,3,10,7,11,5,10,3,11,0,10,4,12,5,10,8,11,5,10,4,10,5,9,7,10,2,10,8,11,6,8,7,11,3,8,4,10,7,10,6,11,2,8,8,10,7,9,7,10,12,5,10,3,11,10,2]},{"wheelLength":111,"noWinIndex":[0],"wheelData":[11,0,6,6,11,4,10,6,12,7,0,6,10,1,9,7,10,6,6,5,9,2,8,6,10,5,10,3,12,3,10,6,11,2,10,5,12,7,9,6,12,2,10,7,10,12,5,12,10,2,9,11,10,6,6,4,10,1,10,3,12,3,10,6,10,2,12,6,12,4,10,4,12,6,10,6,12,8,5,11,12,1,10,7,11,3,10,6,10,2,8,5,10,5,10,6,6,11,4,10,6,8,1,10,7,11,2,8,6,10,4]},{"wheelLength":115,"noWinIndex":[0],"wheelData":[12,5,9,6,9,1,12,5,9,5,9,5,11,7,7,12,6,9,11,1,8,6,9,7,8,3,9,8,7,7,12,9,1,12,7,9,9,9,0,8,7,10,3,9,7,10,3,11,4,12,7,10,4,9,7,11,5,5,12,7,9,4,12,6,11,4,9,3,12,3,9,7,11,4,12,5,9,5,12,7,9,1,12,9,5,9,7,8,4,9,7,11,5,12,12,12,7,11,5,9,7,12,5,9,7,9,5,12,6,9,1,9,3,12,5]},{"wheelLength":114,"noWinIndex":[0],"wheelData":[11,4,4,8,5,10,0,9,4,11,5,9,5,10,3,3,8,5,11,1,8,7,12,6,8,6,10,3,8,3,10,4,8,6,10,4,9,11,1,12,10,6,11,4,8,5,11,7,8,4,10,10,10,3,8,5,8,7,10,1,8,10,7,11,7,8,6,11,3,9,6,11,3,9,6,8,7,8,3,10,4,8,4,12,3,8,6,10,11,1,9,6,10,3,8,7,10,4,9,5,10,4,10,3,10,0,9,6,10,5,9,4,10,6]}]]},"freeGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":13,"maxBetLine":0,"specialFeatureCount":1,"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,0,0,0],[0,0,50,200,1000],[0,0,25,150,400],[0,0,25,150,400],[0,0,20,75,200],[0,0,20,75,200],[0,0,10,50,150],[0,0,10,50,150],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[],"symbolAttribute":["Wild","BonusGame","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"WaysGame","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_05","specialHitInfo":["reSpin_01"],"basePay":2}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":120,"noWinIndex":[7],"wheelData":[11,4,11,6,11,2,10,11,8,8,8,6,11,2,11,8,10,11,6,7,5,4,11,8,2,10,7,10,12,11,10,8,12,10,4,10,6,7,5,4,12,6,10,11,12,6,10,8,12,10,11,8,12,10,9,12,10,4,11,12,11,7,12,10,11,4,11,6,10,7,10,11,9,10,11,2,9,5,11,11,8,10,12,7,11,8,12,7,8,8,11,10,8,12,11,10,12,3,11,7,8,12,11,2,12,4,8,11,12,8,11,3,8,8,10,11,9,8,8,11]},{"wheelLength":120,"noWinIndex":[12],"wheelData":[10,10,6,7,5,4,11,2,12,7,12,2,12,10,10,12,11,10,12,12,10,2,12,3,10,2,12,7,10,12,12,11,12,10,10,12,11,2,10,3,12,7,10,8,11,12,12,12,10,9,10,12,10,10,12,11,3,12,10,9,12,10,10,10,3,12,10,12,11,10,12,10,7,10,6,12,4,10,6,10,4,12,6,10,9,12,0,10,12,10,9,10,8,6,10,12,3,10,6,12,12,3,10,5,12,3,12,6,5,9,12,10,12,3,10,2,12,12,10,6]},{"wheelLength":130,"noWinIndex":[0],"wheelData":[11,2,8,5,9,8,8,6,9,2,9,8,9,10,11,8,12,2,8,4,11,5,8,2,9,9,8,12,11,8,9,8,9,4,8,5,11,4,8,5,9,4,8,12,9,11,8,5,10,2,9,8,9,11,8,9,8,3,9,5,8,4,9,8,3,9,5,8,5,11,7,9,4,11,9,5,10,11,9,0,8,9,9,8,4,9,11,9,2,8,7,9,11,9,8,6,8,8,8,5,9,8,8,11,9,8,9,4,9,8,10,6,9,8,12,9,8,8,8,11,9,9,8,4,9,8,9,7,9,11]},{"wheelLength":260,"noWinIndex":[0],"wheelData":[9,11,11,11,9,4,10,6,8,7,12,3,8,4,12,12,12,10,10,9,11,11,5,10,9,9,9,6,10,10,9,3,12,9,11,6,9,9,10,5,10,10,9,11,10,3,12,12,9,7,11,5,10,10,10,12,8,12,11,9,12,10,5,9,8,9,11,7,8,9,9,7,11,11,11,9,9,9,4,12,12,9,7,10,11,8,3,11,10,9,12,7,10,8,9,10,11,4,3,6,10,11,11,10,9,10,9,12,11,6,12,10,8,11,11,9,10,10,8,12,10,6,9,3,10,6,7,3,10,8,9,11,11,11,9,10,4,6,7,9,12,8,8,4,12,12,12,10,10,9,11,11,5,10,9,9,9,8,10,10,3,9,12,9,11,10,9,9,10,11,10,10,9,5,10,3,12,12,9,11,11,5,10,10,10,12,10,12,11,9,12,10,5,9,10,9,11,11,8,9,9,7,11,11,11,9,9,9,4,12,12,9,8,10,12,8,0,11,10,9,12,12,10,8,9,11,7,4,3,9,10,11,11,10,9,10,12,9,10,11,12,8,10,6,11,7,10,10,8,12,10,11,9,10,9,11,6,9,10,8]},{"wheelLength":240,"noWinIndex":[0],"wheelData":[9,11,11,11,7,8,8,8,11,0,10,9,12,5,12,11,3,12,12,12,9,11,11,11,12,9,4,10,12,3,9,5,11,3,9,4,10,6,11,5,12,9,6,12,11,12,12,10,5,11,12,9,9,12,9,12,9,6,11,9,9,12,12,12,5,11,9,11,9,4,9,12,11,12,12,12,9,9,11,11,8,6,9,8,9,11,12,7,9,9,10,12,12,9,11,9,11,3,12,11,7,10,5,11,4,9,7,9,8,4,12,3,12,6,9,11,7,12,12,12,9,11,11,11,7,8,8,8,11,3,10,9,12,5,12,11,3,12,12,12,9,11,11,11,12,9,4,10,12,3,9,5,11,4,12,6,10,4,11,12,5,3,6,12,11,9,12,10,5,11,12,9,9,12,9,12,9,6,11,9,9,12,12,12,5,11,9,11,9,4,9,12,11,12,12,12,9,9,11,11,8,6,9,8,9,11,12,7,9,9,10,12,7,3,5,9,11,10,12,11,12,11,9,11,8,4,3,7,9,12,9,4,12,6,9,11,7,12,12,12]}]],"freeGameExtendSetting":{"faceChangeProbability":[0.111111,0.111111,0.111111,0.333333,0.333334],"defaultRound":10,"addRoundPerHit":10,"maxRound":150}},"bonusGameSetting":{},"doubleGameSetting":{"spinTimeLimit":5,"spinBetLimit":1000000000,"rtp":0.96,"tieRate":0.1}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":-1,"betCombinations":{"10_5_NoExtraBet":500,"1_5_NoExtraBet":50,"2_5_NoExtraBet":100,"3_5_NoExtraBet":150,"5_5_NoExtraBet":250},"gambleLimit":50000000,"defaultWaysBetColumnIdx":0}`
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
