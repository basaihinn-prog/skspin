package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/servicejdb/jdb_8015/internal"
	"serve/servicejdb/jdbcomm"
)

func init() {
	jdbcomm.RegRpc("h5.init", doInit)
}

var (
	sampleDoInit = `{"maxBet":9223372036854775807,"defaultWaysBetIdx":-1,"singleBetCombinations":{"10_10_40_NoExtraBet":400,"10_15_40_NoExtraBet":600,"10_1_40_NoExtraBet":40,"10_2_40_NoExtraBet":80,"10_5_40_NoExtraBet":200},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":0,"gameFeatureCount":2,"executeSetting":{"settingId":"v2_8015_05_02_001","baseGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":12,"maxBetLine":40,"betSpec":{"lineBetList":[1,2,5,10,15],"betLineList":[40],"extraBetTypeList":["NoExtraBet"]},"specialFeatureCount":1,"payTable":[[0,0,30,150,1000],[0,0,0,0,0],[0,0,25,100,400],[0,0,25,100,400],[0,0,15,50,200],[0,0,15,50,150],[0,0,10,40,150],[0,0,10,30,125],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[[1,1,1,1,1],[2,2,2,2,2],[0,0,0,0,0],[3,3,3,3,3],[1,2,3,2,1],[2,1,0,1,2],[0,0,1,2,3],[3,3,2,1,0],[1,0,0,0,1],[2,3,3,3,2],[0,1,2,3,3],[3,2,1,0,0],[1,0,1,2,1],[2,3,2,1,2],[0,1,0,1,0],[3,2,3,2,3],[1,2,1,0,1],[2,1,2,3,2],[0,1,1,1,0],[3,2,2,2,3],[1,1,2,3,3],[2,2,1,0,0],[1,1,0,1,1],[2,2,3,2,2],[1,2,2,2,3],[2,1,1,1,0],[0,0,1,0,0],[3,3,2,3,3],[0,1,2,2,3],[3,2,1,1,0],[0,0,0,1,2],[3,3,3,2,1],[1,0,0,1,2],[2,3,3,2,1],[0,1,1,2,3],[3,2,2,1,0],[1,0,1,2,3],[2,3,2,1,0],[0,1,2,3,2],[3,2,1,0,1]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"LeftToRight","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_11","specialHitInfo":["freeGame_01"],"basePay":2}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":65,"noWinIndex":[0],"wheelData":[6,2,8,7,10,5,7,8,6,3,10,7,8,5,6,8,7,6,5,10,9,7,4,7,6,9,8,6,0,0,0,0,0,8,11,6,2,8,6,7,10,5,7,11,8,3,10,7,10,5,10,6,8,3,9,7,8,2,8,6,10,3,6,11,7]},{"wheelLength":59,"noWinIndex":[0],"wheelData":[5,11,10,11,4,10,8,2,9,11,3,11,5,11,7,0,0,0,0,8,11,3,8,5,6,11,1,9,10,9,1,11,9,2,11,9,5,6,11,9,4,9,6,4,6,9,5,6,11,4,6,11,1,7,11,9,3,6,9]},{"wheelLength":63,"noWinIndex":[0],"wheelData":[4,9,3,8,5,10,2,9,3,11,3,6,9,10,2,9,1,11,10,3,9,2,11,7,9,1,10,2,9,4,10,5,9,4,8,0,0,0,0,9,6,5,9,8,7,1,8,9,2,8,7,5,10,4,6,9,5,8,9,4,6,3,9]},{"wheelLength":59,"noWinIndex":[0],"wheelData":[5,10,4,9,3,11,9,5,11,6,4,6,2,10,3,7,2,6,9,1,8,5,8,5,8,1,11,5,6,9,4,10,8,3,7,3,10,1,7,3,9,3,10,2,11,0,0,0,0,0,0,0,6,5,7,2,8,4,11]},{"wheelLength":50,"noWinIndex":[0],"wheelData":[3,8,7,2,9,4,11,7,2,11,8,3,0,0,0,0,0,6,3,11,3,2,4,5,4,10,2,6,5,6,9,3,10,7,4,6,7,5,8,5,11,7,2,7,3,6,5,9,10,4]}]]},"freeGameSetting":{"screenRow":4,"screenColumn":5,"symbolCount":12,"maxBetLine":40,"specialFeatureCount":1,"payTable":[[0,0,30,150,1000],[0,0,0,0,0],[0,0,25,100,400],[0,0,25,100,400],[0,0,15,50,200],[0,0,15,50,150],[0,0,10,40,150],[0,0,10,30,125],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100],[0,0,5,20,100]],"lineTable":[[1,1,1,1,1],[2,2,2,2,2],[0,0,0,0,0],[3,3,3,3,3],[1,2,3,2,1],[2,1,0,1,2],[0,0,1,2,3],[3,3,2,1,0],[1,0,0,0,1],[2,3,3,3,2],[0,1,2,3,3],[3,2,1,0,0],[1,0,1,2,1],[2,3,2,1,2],[0,1,0,1,0],[3,2,3,2,3],[1,2,1,0,1],[2,1,2,3,2],[0,1,1,1,0],[3,2,2,2,3],[1,1,2,3,3],[2,2,1,0,0],[1,1,0,1,1],[2,2,3,2,2],[1,2,2,2,3],[2,1,1,1,0],[0,0,1,0,0],[3,3,2,3,3],[0,1,2,2,3],[3,2,1,1,0],[0,0,0,1,2],[3,3,3,2,1],[1,0,0,1,2],[2,3,3,2,1],[0,1,1,2,3],[3,2,2,1,0],[1,0,1,2,3],[2,3,2,1,0],[0,1,2,3,2],[3,2,1,0,1]],"symbolAttribute":["Wild","FreeGame","Base","Base","Base","Base","Base","Base","Base","Base","Base","Base"],"gameHitPattern":"LeftToRight","wheelUsePattern":"Dependent","specialHitInfo":[{"specialHitPattern":"HP_11","specialHitInfo":["reSpin_01"],"basePay":2}],"mixGroupCount":0,"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":200,"noWinIndex":[0],"wheelData":[4,10,9,6,5,11,9,4,11,10,9,5,11,7,8,4,9,7,11,5,10,6,3,8,9,2,11,10,6,5,6,7,3,11,8,10,3,7,11,10,3,9,6,8,5,11,6,9,5,8,10,11,4,6,9,11,3,7,9,6,4,9,8,7,4,8,11,10,3,7,6,8,0,0,0,0,7,6,8,3,6,11,2,8,10,6,3,7,6,3,9,10,6,5,9,7,6,4,10,6,8,6,3,4,10,9,6,5,10,11,9,4,11,10,9,5,11,7,8,4,9,7,11,5,10,6,3,8,9,2,11,10,6,5,6,7,3,11,8,10,3,7,11,10,3,9,6,8,5,11,6,9,5,8,10,11,4,6,9,11,3,7,9,6,4,9,8,7,4,8,10,3,7,6,8,4,6,10,6,3,7,6,8,3,6,11,2,8,10,6,3,7,6,3,9,10,6,5,7,6]},{"wheelLength":33,"noWinIndex":[0],"wheelData":[3,8,9,6,1,8,6,7,3,10,6,4,9,11,0,0,0,0,0,0,7,4,11,10,2,9,1,11,8,5,7,10,11]},{"wheelLength":54,"noWinIndex":[0],"wheelData":[5,8,11,2,7,9,6,3,8,10,11,0,0,0,0,9,10,3,6,11,5,9,7,3,9,8,7,3,9,10,4,6,10,3,7,9,11,2,6,8,3,7,6,7,4,11,8,9,5,7,6,1,7,10]},{"wheelLength":30,"noWinIndex":[0],"wheelData":[7,10,2,11,8,10,0,0,0,0,8,10,9,3,6,10,9,1,6,10,11,6,5,7,11,8,3,9,11,4]},{"wheelLength":200,"noWinIndex":[0],"wheelData":[3,8,9,7,3,9,11,10,2,11,8,4,10,9,5,7,11,10,3,11,8,5,10,8,3,10,7,6,4,6,9,3,11,9,4,6,7,5,11,5,8,9,2,7,6,4,11,9,5,8,10,9,4,5,6,9,3,10,7,4,6,7,5,10,3,8,7,2,7,6,4,11,8,9,0,0,0,0,6,5,8,9,2,9,3,11,10,2,11,8,4,10,8,5,3,11,10,2,11,8,4,10,8,4,10,2,6,5,6,9,3,11,9,6,5,7,10,11,5,8,9,2,7,6,11,5,2,11,8,10,2,6,5,6,9,3,10,7,4,6,7,3,10,5,8,7,2,7,4,11,8,9,7,6,4,6,9,3,11,9,4,6,7,5,11,5,8,9,2,7,6,4,11,9,5,8,10,9,4,5,6,9,3,10,7,4,6,7,5,10,3,8,7,2,7,6,4,11,8,9]}]],"freeGameExtendSetting":{"defaultRound":10,"addRoundPerHit":10,"maxRound":20}},"bonusGameSetting":{},"doubleGameSetting":{"spinTimeLimit":5,"spinBetLimit":1000000000,"rtp":0.96,"tieRate":0.1}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":0,"betCombinations":{"10_40_NoExtraBet":400,"15_40_NoExtraBet":600,"1_40_NoExtraBet":40,"2_40_NoExtraBet":80,"5_40_NoExtraBet":200},"gambleLimit":50000000,"defaultWaysBetColumnIdx":-1}`
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
