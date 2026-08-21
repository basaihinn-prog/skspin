package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_14068/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":0,"singleBetCombinations":{"10_10_5_NoExtraBet":550,"10_1_5_NoExtraBet":55,"10_2_5_NoExtraBet":110,"10_3_5_NoExtraBet":165,"10_5_5_NoExtraBet":275},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":-1,"defaultConnectBetIdx":-1,"defaultQuantityBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v3_14068_05_00_001","betSpecSetting":{"paymentType":"PT_027","extraBetTypeList":["NoExtraBet"],"betSpecification":{"wayBetList":[1,2,3,5,10],"betColumnList":[5],"betType":"WayGame"}},"gameStateSetting":[{"gameStateType":"GS_098","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":41,"noWinIndex":[11],"wheelData":[3,8,12,6,11,3,3,3,4,12,10,4,4,4,8,1,2,2,2,1,2,7,4,9,5,8,3,3,6,11,3,8,3,5,7,6,10,9,3,5,1]},{"wheelLength":51,"noWinIndex":[30],"wheelData":[4,6,5,1,2,2,2,1,4,4,4,2,9,0,10,6,4,4,8,6,3,3,3,0,3,12,4,11,2,8,1,5,5,9,6,11,2,1,7,10,5,0,12,7,6,8,5,4,7,5,8]},{"wheelLength":47,"noWinIndex":[38],"wheelData":[9,4,0,10,7,1,5,1,5,3,3,3,12,0,11,4,4,4,1,8,4,4,11,6,8,3,7,12,4,7,4,0,10,3,4,9,3,1,2,2,2,1,6,7,4,8,6]},{"wheelLength":50,"noWinIndex":[14],"wheelData":[4,11,10,12,1,2,1,8,10,0,11,9,4,7,6,6,6,12,11,7,3,12,4,0,8,4,11,3,3,3,0,3,12,1,2,2,2,1,7,8,4,7,6,10,5,12,10,3,12,9]},{"wheelLength":63,"noWinIndex":[6],"wheelData":[12,4,9,7,4,9,3,3,3,12,11,1,2,2,2,1,2,10,6,11,9,4,12,7,3,10,8,2,11,9,12,7,2,8,11,3,9,8,7,4,4,4,7,6,8,10,6,7,2,10,12,6,10,8,5,8,10,6,12,5,3,9,6]}]]},"symbolSetting":{"symbolCount":13,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","A","K","Q","J","TE","NI"],"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,80,160,800],[0,0,50,100,500],[0,0,25,50,250],[0,0,20,40,200],[0,0,10,20,100],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50],[0,0,5,10,50]],"mixGroupCount":0},"gameHitPatternSetting":{"gameHitPattern":"WayGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":3,"specialHitInfo":[{"specialHitPattern":"HP_16","triggerEvent":"Trigger_01","basePay":50},{"specialHitPattern":"HP_17","triggerEvent":"Trigger_02","basePay":10},{"specialHitPattern":"HP_18","triggerEvent":"Trigger_03","basePay":5}]},"progressSetting":{"triggerLimitType":"NoLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":1,"addRound":0,"maxRound":1}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_03"]}},"extendSetting":{"ifBaseGame":true}},{"gameStateType":"GS_098","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":27,"noWinIndex":[18],"wheelData":[3,3,3,5,4,4,1,3,2,5,3,3,3,4,4,1,2,2,2,4,2,6,5,3,4,4,6]},{"wheelLength":25,"noWinIndex":[10],"wheelData":[4,1,4,0,4,1,3,3,3,1,3,0,3,4,4,4,3,5,1,2,2,2,1,5,3]},{"wheelLength":36,"noWinIndex":[5],"wheelData":[5,6,4,0,6,5,1,5,3,3,4,0,4,4,5,1,4,4,4,5,1,0,4,5,3,3,3,5,1,2,2,2,1,3,3,5]},{"wheelLength":48,"noWinIndex":[10],"wheelData":[4,3,2,5,6,4,4,2,2,5,3,0,3,4,4,1,2,2,2,1,5,4,4,3,3,3,5,0,2,2,1,5,3,3,5,4,6,3,3,6,4,4,4,3,5,6,3,5]},{"wheelLength":45,"noWinIndex":[6],"wheelData":[4,5,4,6,5,6,2,4,2,5,4,4,2,2,2,4,4,4,6,5,2,2,5,4,6,5,1,4,4,4,6,6,2,5,4,5,3,3,5,4,2,2,6,3,1]}]]},"symbolSetting":{"symbolCount":7,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5"],"payTable":[[0,0,0,0,0],[0,0,0,0,0],[0,0,80,160,800],[0,0,50,100,500],[0,0,25,50,250],[0,0,20,40,200],[0,0,10,20,100]],"mixGroupCount":0},"gameHitPatternSetting":{"gameHitPattern":"WayGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":3,"specialHitInfo":[{"specialHitPattern":"HP_16","triggerEvent":"Trigger_01","basePay":50},{"specialHitPattern":"HP_17","triggerEvent":"Trigger_02","basePay":10},{"specialHitPattern":"HP_18","triggerEvent":"Trigger_03","basePay":5}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":10,"addRound":10,"maxRound":100}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"MaxTriggerLimit","readyHandCount":1,"readyHandType":["ReadyHand_03"]}},"extendSetting":{"ifBaseGame":false}}],"doubleGameSetting":{"doubleRoundUpperLimit":5,"doubleBetUpperLimit":1000000000,"rtp":0.96,"tieRate":0.1},"boardDisplaySetting":{"winRankSetting":{"BigWin":20,"MegaWin":100,"UltraWin":300}},"gameFlowSetting":{"conditionTableWithoutBoardEnd":[["CD_False","CD_True","CD_False"],["CD_False","CD_False","CD_12"],["CD_False","CD_False","CD_False"]]}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":-1,"betCombinations":{"10_5_NoExtraBet":550,"1_5_NoExtraBet":55,"2_5_NoExtraBet":110,"3_5_NoExtraBet":165,"5_5_NoExtraBet":275},"gambleLimit":0,"buyFeatureLimit":2147483647,"buyFeature":true,"defaultWaysBetColumnIdx":0}`
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
