package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_8050/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":0,"singleBetCombinations":{"10_10_5_NoExtraBet":500,"10_1_5_NoExtraBet":50,"10_2_5_NoExtraBet":100,"10_3_5_NoExtraBet":150,"10_6_5_NoExtraBet":300},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":-1,"defaultConnectBetIdx":-1,"defaultQuantityBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v3_8050_05_01_001","betSpecSetting":{"paymentType":"PT_004","extraBetTypeList":["NoExtraBet"],"betSpecification":{"wayBetList":[1,2,3,6,10],"betColumnList":[5],"betType":"WayGame"}},"gameStateSetting":[{"gameStateType":"GS_081","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":71,"noWinIndex":[0],"wheelData":[4,5,5,3,3,5,5,2,2,2,4,6,7,5,5,4,7,1,7,4,3,3,3,5,5,7,7,2,2,2,3,3,3,4,4,3,7,7,5,3,3,5,5,7,7,1,4,4,5,5,5,1,3,3,5,7,7,4,5,5,4,3,3,4,7,1,7,4,5,5,4]},{"wheelLength":71,"noWinIndex":[0],"wheelData":[6,3,4,4,3,3,6,6,4,4,5,5,4,6,6,1,3,3,3,4,4,6,6,2,2,2,5,5,1,5,3,5,0,2,5,5,6,6,4,4,5,6,6,4,5,0,5,4,3,6,1,6,3,3,0,5,3,6,6,3,1,4,4,5,3,3,7,6,3,6,1]},{"wheelLength":69,"noWinIndex":[0],"wheelData":[3,5,5,7,7,5,5,0,3,5,5,3,1,7,7,1,7,4,2,2,4,5,5,6,6,4,2,2,4,4,6,7,0,5,5,6,6,4,4,3,3,7,7,3,1,4,4,5,5,6,6,4,4,3,3,2,2,4,7,1,7,4,1,4,7,4,5,3,3]},{"wheelLength":69,"noWinIndex":[0],"wheelData":[3,3,5,5,0,4,7,3,3,5,6,7,5,4,4,2,5,7,7,5,0,7,5,4,6,6,3,3,4,4,3,6,6,3,2,2,6,6,5,7,7,5,2,0,6,4,4,2,2,2,7,7,5,3,3,5,0,4,3,3,4,2,2,5,6,0,6,5,3]},{"wheelLength":84,"noWinIndex":[0],"wheelData":[2,4,6,2,2,6,6,2,2,7,6,5,7,7,2,2,6,6,2,2,2,4,6,6,2,7,7,4,4,7,7,2,2,4,4,4,3,3,6,6,2,4,5,5,4,2,2,7,7,2,2,6,6,2,2,7,7,3,3,3,5,5,6,6,2,2,2,6,6,5,5,6,6,4,2,2,5,5,2,2,6,6,2,2]}]]},"symbolSetting":{"symbolCount":14,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","A","K","Q","J","TE","NI"],"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,60,200,350],[0,0,35,120,200],[0,0,35,100,150],[0,0,25,100,150],[0,0,25,80,125],[0,0,25,80,125],[0,0,20,60,100],[0,0,20,60,100],[0,0,15,40,100],[0,0,15,30,60],[0,0,15,30,60],[0,0,15,30,60]],"mixGroupCount":0},"lineSetting":{"maxBetLine":0},"gameHitPatternSetting":{"gameHitPattern":"WayGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":1,"specialHitInfo":[{"specialHitPattern":"HP_05","triggerEvent":"Trigger_01","basePay":5}]},"progressSetting":{"triggerLimitType":"NoLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":1,"addRound":0,"maxRound":1}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_01"]}},"extendSetting":{"lightReadyHandOdds":20}},{"gameStateType":"GS_036","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":2,"tableHitProbability":[0.5,0.5],"wheelData":[[{"wheelLength":36,"noWinIndex":[0],"wheelData":[0,0,0,2,2,7,0,0,0,4,6,0,0,0,0,7,3,1,5,0,0,0,0,2,6,0,0,0,0,2,5,0,0,0,7,4]},{"wheelLength":72,"noWinIndex":[0],"wheelData":[1,5,5,6,4,4,6,2,2,2,4,5,4,6,6,3,3,7,7,4,1,5,5,4,4,3,2,6,6,4,3,3,4,5,6,6,5,4,1,5,6,6,5,5,6,6,2,2,2,4,5,5,3,6,6,1,5,3,3,3,5,2,2,6,7,4,5,5,4,6,6,4]},{"wheelLength":72,"noWinIndex":[0],"wheelData":[1,3,4,5,3,7,5,3,4,5,7,7,5,3,5,2,7,3,1,7,3,5,7,3,4,7,3,4,4,7,3,2,2,3,5,5,1,4,5,3,3,6,5,5,7,7,5,7,7,5,3,3,3,7,1,3,7,5,3,7,4,3,7,4,5,7,2,2,2,3,5,4]},{"wheelLength":72,"noWinIndex":[0],"wheelData":[7,6,4,4,6,5,4,7,7,4,6,6,5,3,4,6,6,4,6,6,4,5,6,7,4,4,5,6,6,5,3,4,7,7,3,4,7,7,4,6,6,5,5,7,6,4,6,6,3,3,3,6,6,4,6,6,4,2,2,2,4,4,6,5,5,6,4,5,7,7,3,4]},{"wheelLength":36,"noWinIndex":[0],"wheelData":[0,0,0,6,7,0,0,0,0,7,5,2,0,0,0,0,5,3,0,0,0,4,4,3,0,0,0,0,4,4,0,0,0,3,5,4]}],[{"wheelLength":10,"noWinIndex":[0],"wheelData":[0,0,0,0,0,0,0,0,0,0]},{"wheelLength":69,"noWinIndex":[0],"wheelData":[2,2,2,5,5,5,4,4,4,6,6,3,3,5,5,7,7,4,4,3,3,6,6,5,5,7,7,2,2,2,6,7,4,4,3,3,6,7,4,4,5,5,3,3,3,4,4,4,2,2,2,4,4,4,3,3,3,4,4,4,3,3,3,4,4,4,5,5,5]},{"wheelLength":69,"noWinIndex":[0],"wheelData":[2,2,2,5,5,5,4,4,4,3,3,3,4,4,4,5,5,5,4,4,4,3,3,3,4,4,4,2,2,2,5,5,5,3,3,3,5,5,5,7,7,7,5,5,5,3,3,3,4,4,4,7,7,7,4,4,4,6,6,6,5,5,5,6,6,6,4,4,4]},{"wheelLength":69,"noWinIndex":[0],"wheelData":[2,2,2,5,5,5,3,3,3,4,4,4,3,3,3,5,5,5,6,6,6,7,7,7,3,3,3,2,2,2,5,5,5,3,3,3,5,5,5,4,4,4,5,5,5,3,3,3,4,4,4,6,6,6,3,3,3,5,5,5,7,7,7,4,4,4,3,3,3]},{"wheelLength":10,"noWinIndex":[0],"wheelData":[0,0,0,0,0,0,0,0,0,0]}]]},"symbolSetting":{"symbolCount":14,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","A","K","Q","J","TE","NI"],"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,60,200,350],[0,0,35,120,200],[0,0,35,100,150],[0,0,25,100,150],[0,0,25,80,125],[0,0,25,80,125],[0,0,20,60,100],[0,0,20,60,100],[0,0,15,40,100],[0,0,15,30,60],[0,0,15,30,60],[0,0,15,30,60]],"mixGroupCount":0},"lineSetting":{"maxBetLine":0},"gameHitPatternSetting":{"gameHitPattern":"WayGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":2,"specialHitInfo":[{"specialHitPattern":"HP_05","triggerEvent":"ReTrigger_01","basePay":5},{"specialHitPattern":"HP_75","triggerEvent":"ReTrigger_05","basePay":0}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":0,"addRound":0,"maxRound":999}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"MaxTriggerLimit","readyHandCount":2,"readyHandType":["ReadyHand_01","ReadyHand_24"]}},"extendSetting":{"defaultRound":{"Trigger_01":8},"addRound":{"ReTrigger_01":8,"ReTrigger_05":0},"reSpingRoundSetting":{"defaultRound":1,"addRound":0,"maxRound":1},"lightReadyHandOdds":40}}],"doubleGameSetting":{"doubleRoundUpperLimit":5,"doubleBetUpperLimit":1000000000,"rtp":0.96,"tieRate":0.1},"boardDisplaySetting":{"winRankSetting":{"BigWin":20,"MegaWin":40,"UltraWin":100}},"gameFlowSetting":{"conditionTableWithoutBoardEnd":[["CD_False","CD_True","CD_False"],["CD_False","CD_False","CD_01"],["CD_False","CD_False","CD_False"]]}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":-1,"betCombinations":{"10_5_NoExtraBet":500,"1_5_NoExtraBet":50,"2_5_NoExtraBet":100,"3_5_NoExtraBet":150,"6_5_NoExtraBet":300},"gambleLimit":0,"buyFeatureLimit":2147483647,"buyFeature":true,"defaultWaysBetColumnIdx":0}`
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
				key := fmt.Sprintf("10_%f_%d_NoExtraBet", info.Cs[i], internal.Line)
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
