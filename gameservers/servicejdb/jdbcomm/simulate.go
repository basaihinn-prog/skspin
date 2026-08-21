package jdbcomm

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"reflect"
	"strconv"
	"time"
)

type SimulateData struct {
	Id                   primitive.ObjectID `bson:"_id"`
	DropPan              Variables          `bson:"droppan"` //自行解析的数据
	HasGame              bool               `bson:"hasgame"`
	Times                float64            `bson:"times"`
	BucketId             int                `bson:"bucketid"`
	Type                 int                `bson:"type"`
	Selected             bool               `bson:"selected"`
	RoundID              int                `bson:"RoundID"`     //数据ID
	TurnIndex            int                `bson:"TurnIndex"`   //轮次，他跟数据ID是属于一组的数据
	FreeFlag             int                `bson:"FreeFlag"`    //旋转类型   0 就是普通旋转, 1就是买免费  2  超级购买   1010代表10次购买   1012代表12次  2010  也是代表10次购买
	GroupFlag            int                `bson:"GroupFlag"`   //普通下注，双倍下注，部分游戏有双倍下注
	QueryString          Variables          `bson:"QueryString"` //数据值
	BucketHeartBeat      int                `bson:"BucketHeartBeat"`
	BucketWave           int                `bson:"BucketWave"`
	BucketGov            int                `bson:"BucketGov"`
	BucketMix            int                `bson:"BucketMix"`
	BucketStable         int                `bson:"BucketStable"`
	BucketHighAward      int                `bson:"BucketHighAward"`
	BucketSuperHighAward int                `bson:"BucketSuperHighAward"`
}

//	func (sd *SimulateData) Deal(c float64, line float64, balance float64) Variables {
//		bet := c * line
//		//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())
//
//		pan := sd.DropPan[0]
//
//		rid := strings.Split(pan.Str("gid"), "_")[0]
//		pan.Set("rid", rid)
//		pan.SetFloat("c", c)
//		pan.MKMulFloat("tw", bet)
//		pan.MKMulFloat("tmb_win", bet)
//		pan.MKMulFloat("w", bet)
//		pan.MKMulFloat("wp", bet)
//		pan.MKMulFloat("fsres", bet)
//		pan.MKMulFloat("fswin", bet)
//		pan.MKMulFloat("fsres_total", bet)
//		pan.MKMulFloat("fswin_total", bet)
//		pan.MKMulFloat("tmb_res", bet)
//		pan.MKMulFloat("rs_iw", bet)
//		if pan.Get("rs_iw") != "" {
//			pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
//		}
//		pan.MKMulFloat("rs_win", bet)
//		if pan.Get("rs_win") != "" {
//			pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
//		}
//		pan.MKMulFloat("mo_tw", bet)
//		pan.MKMulFloat("apwa", bet)
//		pan.MKMulFloat("rw", bet)
//		psym := pan.Get("psym")
//		split := strings.Split(psym, "~")
//		if len(split) == 3 {
//			ret, _ := strconv.ParseFloat(split[1], 64)
//			p := message.NewPrinter(language.English)
//			split[1] = p.Sprintf("%.2f", ret*bet)
//			psym = strings.Join(split, "~")
//			pan.Set("psym", psym)
//		}
//		gsfA := pan.Get("gsf_a")
//		if gsfA != "" {
//			split = strings.Split(gsfA, ";")
//			for i, cell := range split {
//				cellsplit := strings.Split(cell, "~")
//				if len(cellsplit) != 2 {
//					break
//				}
//				ret, _ := strconv.ParseFloat(cellsplit[1], 64)
//				p := message.NewPrinter(language.English)
//				cellsplit[1] = p.Sprintf("%.2f", ret*bet)
//				split[i] = strings.Join(cellsplit, "~")
//			}
//			gsfA = strings.Join(split, ";")
//			pan.Set("gsf_a", gsfA)
//		}
//
//		// 消除
//		for k := range pan {
//			// 正则表达式，使用捕获组
//			r := regexp.MustCompile(`^(l\d+)`)
//			// 查找匹配
//			match := r.FindStringSubmatch(k)
//			// 判断是否匹配并获取结果
//			if len(match) > 0 {
//				lstrs := strings.Split(pan[k], "~")
//				lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//				score, err := strconv.ParseFloat(lstrs[1], 64)
//				if err != nil {
//					panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
//				}
//				p := message.NewPrinter(language.English)
//				lstrs[1] = p.Sprintf("%.2f", score*bet)
//				pan[k] = strings.Join(lstrs, "~")
//			}
//		}
//		//for i := 0; i < pan.Int("l"); i++ {
//		//	l := fmt.Sprintf("l%d", i)
//		//	if lstr, ok := pan[l]; ok {
//		//		lstrs := strings.Split(lstr, "~")
//		//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//		//		score, err := strconv.ParseFloat(lstrs[1], 64)
//		//		if err != nil {
//		//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
//		//		}
//		//		p := message.NewPrinter(language.English)
//		//		lstrs[1] = p.Sprintf("%.2f", score*bet)
//		//		pan[l] = strings.Join(lstrs, "~")
//		//	}
//		//}
//
//		pan.SetCurrency("balance", balance)
//		pan.SetCurrency("balance_cash", balance)
//
//		return pan
//	}
func parseMongoArray(raw interface{}) []map[string]interface{} {
	// 尝试断言为 bson.A，即 []interface{}
	array, ok := raw.(bson.A)
	if !ok {
		fmt.Println("类型断言失败，raw 不是 bson.A")
		return nil
	}

	var result []map[string]interface{}
	for i, item := range array {
		doc, ok := item.(bson.M) // bson.M == map[string]interface{}
		if !ok {
			fmt.Printf("第 %d 个元素不是 bson.M\n", i)
			continue
		}
		result = append(result, doc)
	}
	return result
}

func (sd *SimulateData) Deal4(c, balance, multiply float64) Variables {
	pan := sd.DropPan
	spinResult := pan["spinResult"].(Variables)
	boardDisplayResult := spinResult["boardDisplayResult"].(Variables)
	boardDisplayResult["displayBet"] = c
	gameStateResult, _ := convertMongoData(spinResult["gameStateResult"])
	for i := range gameStateResult {
		gameStateResult[i].MKMulFloat("stateWin", multiply)
		roundResult, _ := convertMongoData(gameStateResult[i]["roundResult"])
		for i2 := range roundResult {
			//正常游戏
			roundResult[i2].MKMulFloat("roundWin", multiply)
			gameResult := convertToMapSafely(roundResult[i2]["gameResult"])
			gameResult.MKMulFloat("playerWin", multiply)
			lineWinResult, _ := convertMongoData(gameResult["lineWinResult"])
			for i3 := range lineWinResult {
				lineWinResult[i3].MKMulFloat("lineWin", multiply)
			}

			quantityWinResult, _ := convertMongoData(gameResult["quantityWinResult"])
			for item := range quantityWinResult {
				quantityWinResult[item].MKMulFloat("quantityWin", multiply)
			}

			// scatterValues * multiply 处理
			extendGameStateResult := convertToMapSafely(roundResult[i2]["extendGameStateResult"])
			scatterValuesHandle(extendGameStateResult, multiply)

			//spinResult.gameStateResult[1].roundResult[0].extendGameStateResult.holdAndSpinResult
			holdAndSpinResultList, err := convertMongoData(extendGameStateResult["holdAndSpinResult"])
			if err == nil && len(holdAndSpinResultList) > 0 {
				for _, holdAndSpin := range holdAndSpinResultList {
					scatterScoreRecordHandle(holdAndSpin, "value", multiply)
				}

			}

			//spinResult.gameStateResult[1].roundResult[0].extendGameStateResult.scatterTwoScoreBeforeBS
			slice, err := convertToSlice(extendGameStateResult["scatterTwoScoreBeforeBS"])
			if err == nil && len(slice) > 0 {
				for i, a := range slice {
					slice[i] = a.(float64) * multiply
				}
				extendGameStateResult["scatterTwoScoreBeforeBS"] = slice
			}

			wayGameWinScoreBefore, err := convertToSlice(extendGameStateResult["wayGameWinScoreBefore"])
			if err == nil && len(wayGameWinScoreBefore) > 0 {
				for i, a := range wayGameWinScoreBefore {
					wayGameWinScoreBefore[i] = a.(float64) * multiply
				}
				extendGameStateResult["wayGameWinScoreBefore"] = wayGameWinScoreBefore
			}

			scatterScoreRecordHandle(extendGameStateResult, "scatterScoreRecord", multiply)

			extendGameStateResult.MKMulFloat("extendWin", multiply)

			// 111
			gameResultWayGameResult := convertToMapSafely(gameResult["wayGameResult"])

			convertMongoDataHandle(gameResultWayGameResult, "wayWinResult", "symbolWin", multiply)

			// symbolWin * multiply 处理

			convertMongoDataHandle(roundResult[i2], "specialFeatureResult", "specialScreenWin", multiply)
			convertMongoDataHandle(gameResult, "wayWinResult", "symbolWin", multiply)

			// playerWin * multiply 处理
			convertToMapSafelyHandle(gameResult, "wayGameResult", "playerWin", multiply)

			// 处理连环销

			continuousWinHandle(gameResult, multiply)
			holdAndSpinResultHandle(extendGameStateResult, multiply)
			scatterTotalScoreHandle(extendGameStateResult)

			// reSpin * multiply 处理
			reSpinMutiplyHandle(convertToMapSafely(roundResult[i2]["gameResult"]), multiply)
			//fg相关
			displayResult := convertToMapSafely(roundResult[i2]["displayResult"])
			accumulateWinResult := convertToMapSafely(displayResult["accumulateWinResult"])
			accumulateWinResult.MKMulFloat("beforeSpinFirstStateOnlyBasePayAccWin", multiply)
			accumulateWinResult.MKMulFloat("afterSpinFirstStateOnlyBasePayAccWin", multiply)
			accumulateWinResult.MKMulFloat("beforeSpinAccWin", multiply)
			accumulateWinResult.MKMulFloat("afterSpinAccWin", multiply)
		}
	}
	spinResult.MKMulFloat("totalWin", multiply)
	pan.SetCurrency("balance", balance)
	pan.SetStr("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	return pan
}

func (sd *SimulateData) Deal3(c float64, balance float64, multiply float64) Variables {
	pan := sd.DropPan
	SpinResult := pan["spinResult"].(Variables)
	SpinResult.MKMulFloat("playerTotalWin", multiply)
	baseGameResult := convertToMapSafely(SpinResult["baseGameResult"])
	freeGameResult := convertToMapSafely(SpinResult["freeGameResult"])
	if freeGameResult != nil {
		freeGameResult.MKMulFloat("freeGameTotalWin", multiply)
	}

	freeGameOneRoundResult, err := convertMongoData(freeGameResult["freeGameOneRoundResult"])

	if err == nil && len(freeGameOneRoundResult) > 0 {

		for item := range freeGameOneRoundResult {
			lineGameResult := convertToMapSafely(freeGameOneRoundResult[item]["lineGameResult"])
			lineGameResult.MKMulFloat("playerWin", multiply)
			freeGameOneRoundResult[item].MKMulFloat("playerWin", multiply)

			displayLogicInfo := convertToMapSafely(freeGameOneRoundResult[item]["displayLogicInfo"])
			displayLogicInfo.MKMulFloat("afterAccumulateWinWithBaseGameWin", multiply)
			displayLogicInfo.MKMulFloat("afterAccumulateWinWithoutBaseGameWin", multiply)
			displayLogicInfo.MKMulFloat("beforeAccumulateWinWithBaseGameWin", multiply)
			displayLogicInfo.MKMulFloat("beforeAccumulateWinWithoutBaseGameWin", multiply)

			lineResult, _ := convertMongoData(lineGameResult["lineResult"])
			for line := range lineResult {
				lineResult[line].MKMulFloat("lineWin_LR", multiply)
				lineResult[line].MKMulFloat("lineWin_RL", multiply)
			}

			convertToMapSafelyHandle(freeGameOneRoundResult[item], "extendInfoForFreeGameResult", "extendPlayerWin", multiply)

			freeGameRoundsWaysGameResult := convertToMapSafely(freeGameOneRoundResult[item]["waysGameResult"])
			if freeGameRoundsWaysGameResult != nil {
				convertMongoDataHandle(freeGameRoundsWaysGameResult, "waysResult", "symbolWin", multiply)
			}

			convertToMapSafelyHandle(freeGameOneRoundResult[item], "waysGameResult", "playerWin", multiply)

		}
	}

	specialFeatureResult, _ := convertMongoData(baseGameResult["specialFeatureResult"])
	for item := range specialFeatureResult {
		specialFeatureResult[item].MKMulFloat("specialScreenWin", multiply)
	}

	waysGameResult := convertToMapSafely(baseGameResult["waysGameResult"])

	waysResult, _ := convertMongoData(waysGameResult["waysResult"])
	for item := range waysResult {
		waysResult[item].MKMulFloat("symbolWin", multiply)
	}

	baseGameResult.MKMulFloat("baseGameTotalWin", multiply)

	bonusGameResult := convertToMapSafely(SpinResult["bonusGameResult"])
	if bonusGameResult != nil {
		bonusGameResult.MKMulFloat("bonusGameTotalWin", multiply)
	}
	convertToMapSafelyHandle(baseGameResult, "extendInfoForbaseGameResult", "extendPlayerWin", multiply)

	lineGameResult := convertToMapSafely(baseGameResult["lineGameResult"])
	lineGameResult.MKMulFloat("playerWin", multiply)

	convertMongoDataHandle(lineGameResult, "lineResult", "lineWin_RL", multiply)
	convertMongoDataHandle(lineGameResult, "lineResult", "lineWin_LR", multiply)
	pan.SetCurrency("balance", balance)
	pan.SetStr("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	return pan
}

func truncateFloat(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(num*shift) / shift
}

func (sd *SimulateData) Deal2(c float64, balance float64) Variables {
	pan := sd.DropPan
	spinResult := pan["spinResult"].(Variables)
	boardDisplayResult := spinResult["boardDisplayResult"].(Variables)
	originC := boardDisplayResult.Float("displayBet")
	boardDisplayResult["displayBet"] = c
	multiply := c / originC // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数
	gameStateResult, _ := convertMongoData(spinResult["gameStateResult"])
	for i := range gameStateResult {
		gameStateResult[i].MKMulFloat("stateWin", multiply)
		roundResult, _ := convertMongoData(gameStateResult[i]["roundResult"])
		for i2 := range roundResult {
			//正常游戏

			roundResult[i2].MKMulFloat("roundWin", multiply)
			gameResult := convertToMapSafely(roundResult[i2]["gameResult"])
			gameResult.MKMulFloat("playerWin", multiply)
			lineWinResult, _ := convertMongoData(gameResult["lineWinResult"])
			for i3 := range lineWinResult {
				lineWinResult[i3].MKMulFloat("lineWin", multiply)
			}

			quantityWinResult, _ := convertMongoData(gameResult["quantityWinResult"])
			for item := range quantityWinResult {
				quantityWinResult[item].MKMulFloat("quantityWin", multiply)
			}
			// scatterValues * multiply 处理
			extendGameStateResult := convertToMapSafely(roundResult[i2]["extendGameStateResult"])
			extendGameStateResult.MKMulFloat("holdAndSpinWin", multiply)
			extendGameStateResult.MKMulFloat("fullExtendWin", multiply)
			extendGameStateResult.MKMulFloat("fullWin", multiply)
			extendGameStateResult.MKMulFloat("holdAndSpinExtendWin", multiply)
			fullResult := convertToMapSafely(extendGameStateResult["fullResult"])
			lineGameResult := convertToMapSafely(fullResult["lineGameResult"])
			convertMongoDataHandle(lineGameResult, "lineWinResult", "lineWin", multiply)
			lineGameResult.MKMulFloat("playerWin", multiply)

			scatterValuesHandle(extendGameStateResult, multiply)

			if extendGameStateResult != nil {

				winC2(extendGameStateResult, multiply)
				//gameDescriptor(extendGameStateResult, multiply)

				screenWinsInfoList, err := convertMongoData(extendGameStateResult["screenWinsInfo"])

				if len(screenWinsInfoList) > 0 && err == nil {
					for _, screenWins := range screenWinsInfoList {
						screenval := convertToMapSafely(screenWins)
						screenval.MKMulFloat("playerWin", multiply)

						quantityWinResults, _ := convertToSlice(screenval["quantityWinResult"])

						for _, quantityv := range quantityWinResults {
							quantityval := convertToMapSafely(quantityv)
							quantityval.MKMulFloat("quantityWin", multiply)
						}
					}
				}

			}
			quantityGameResult(gameResult, multiply)

			//spinResult.gameStateResult[1].roundResult[0].extendGameStateResult.holdAndSpinResult
			//holdAndSpinResultList, err := convertMongoData(extendGameStateResult["holdAndSpinResult"])
			//if err == nil && len(holdAndSpinResultList) > 0 {
			//	for _, holdAndSpin := range holdAndSpinResultList {
			//		scatterScoreRecordHandle(holdAndSpin, "value", multiply)
			//	}
			//
			//}

			//spinResult.gameStateResult[1].roundResult[0].extendGameStateResult.scatterTwoScoreBeforeBS
			slice, err := convertToSlice(extendGameStateResult["scatterTwoScoreBeforeBS"])
			if err == nil && len(slice) > 0 {
				for i, a := range slice {
					slice[i] = a.(float64) * multiply
				}
				extendGameStateResult["scatterTwoScoreBeforeBS"] = slice
			}

			wayGameWinScoreBefore, err := convertToSlice(extendGameStateResult["wayGameWinScoreBefore"])
			if err == nil && len(wayGameWinScoreBefore) > 0 {
				for i, a := range wayGameWinScoreBefore {
					wayGameWinScoreBefore[i] = a.(float64) * multiply
				}
				extendGameStateResult["wayGameWinScoreBefore"] = wayGameWinScoreBefore
			}
			scatterScoreRecordHandle(extendGameStateResult, "scatterScoreRecord", multiply)
			scatterScoreRecordHandle(extendGameStateResult, "eachPositionPoints", multiply)

			scatterScoreRecordHandle(extendGameStateResult, "lanternScoreArray", multiply)
			scatterScoreRecordHandle(extendGameStateResult, "beforeLanternScoreArray", multiply)
			scatterScoreRecordHandle(extendGameStateResult, "afterLanternScoreArray", multiply)

			extendGameStateResult.MKMulFloat("extendWin", multiply)

			// 111
			gameResultWayGameResult := convertToMapSafely(gameResult["wayGameResult"])

			convertMongoDataHandle(gameResultWayGameResult, "wayWinResult", "symbolWin", multiply)

			// symbolWin * multiply 处理

			convertMongoDataHandle(roundResult[i2], "specialFeatureResult", "specialScreenWin", multiply)
			convertMongoDataHandle(gameResult, "wayWinResult", "symbolWin", multiply)

			// playerWin * multiply 处理
			convertToMapSafelyHandle(gameResult, "wayGameResult", "playerWin", multiply)

			// 处理连环销

			continuousWinHandle(gameResult, multiply)
			holdAndSpinResultHandle(extendGameStateResult, multiply)
			scatterTotalScoreHandle(extendGameStateResult)

			// reSpin * multiply 处理
			reSpinMutiplyHandle(convertToMapSafely(roundResult[i2]["gameResult"]), multiply)
			//fg相关
			displayResult := convertToMapSafely(roundResult[i2]["displayResult"])
			accumulateWinResult := convertToMapSafely(displayResult["accumulateWinResult"])
			accumulateWinResult.MKMulFloat("beforeSpinFirstStateOnlyBasePayAccWin", multiply)
			accumulateWinResult.MKMulFloat("afterSpinFirstStateOnlyBasePayAccWin", multiply)
			accumulateWinResult.MKMulFloat("beforeSpinAccWin", multiply)
			accumulateWinResult.MKMulFloat("afterSpinAccWin", multiply)

		}
	}
	spinResult.MKMulFloat("totalWin", multiply)
	pan.SetCurrency("balance", balance)
	pan.SetStr("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	return pan
}

//
//func (sd *SimulateData) DealSP(c float64, line float64, balance float64, originC float64) Variables {
//	bet := c * line
//	//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())
//	pan := sd.DropPan[0]
//	//originC := pan.Float("c")
//	multiply := c / originC // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数
//	bet = multiply
//	rid := strings.Split(pan.Str("gid"), "_")[0]
//	pan.Set("rid", rid)
//	pan.SetFloat("c", c)
//	pan.MKMulFloat("tw", bet)
//	pan.MKMulFloat("tmb_win", bet)
//	pan.MKMulFloat("w", bet)
//	pan.MKMulFloat("wp", bet)
//	pan.MKMulFloat("fsres", bet)
//	pan.MKMulFloat("fswin", bet)
//	pan.MKMulFloat("fsres_total", bet)
//	pan.MKMulFloat("fswin_total", bet)
//	pan.MKMulFloat("tmb_res", bet)
//	pan.MKMulFloat("pw", bet)
//	pan.MKMulFloat("rs_iw", bet)
//	if pan.Get("rs_iw") != "" {
//		pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
//	}
//	pan.MKMulFloat("rs_win", bet)
//	if pan.Get("rs_win") != "" {
//		pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
//	}
//	pan.MKMulFloat("mo_tw", bet)
//	//针对pp_vswayschilheat的处理
//	pan.MKMulmo_twInG("mo_tw", bet)
//	//pan.MKMulFloat("apwa", bet)
//	if pan.Get("apwa") != "" {
//		splitAPWA := strings.Split(pan.Get("apwa"), ",")
//		for k, v := range splitAPWA {
//			ret, _ := strconv.ParseFloat(v, 64)
//			p := message.NewPrinter(language.English)
//			splitAPWA[k] = p.Sprintf("%.2f", ret*bet)
//		}
//		pan.Set("apwa", strings.Join(splitAPWA, ","))
//	}
//	pan.MKMulFloat("rw", bet)
//	psym := pan.Get("psym")
//	split := strings.Split(psym, "~")
//	if len(split) == 3 {
//		ret, _ := strconv.ParseFloat(split[1], 64)
//		p := message.NewPrinter(language.English)
//		split[1] = p.Sprintf("%.2f", ret*bet)
//		psym = strings.Join(split, "~")
//		pan.Set("psym", psym)
//	}
//	gsfA := pan.Get("gsf_a")
//	if gsfA != "" {
//		split = strings.Split(gsfA, ";")
//		for i, cell := range split {
//			cellsplit := strings.Split(cell, "~")
//			if len(cellsplit) != 2 {
//				break
//			}
//			ret, _ := strconv.ParseFloat(cellsplit[1], 64)
//			p := message.NewPrinter(language.English)
//			cellsplit[1] = p.Sprintf("%.2f", ret*bet)
//			split[i] = strings.Join(cellsplit, "~")
//		}
//		gsfA = strings.Join(split, ";")
//		pan.Set("gsf_a", gsfA)
//	}
//	//ReplaceEmbedScoreField 兴许可以使用
//	//ReplaceScoreField2(pan, bet, "wlc_v", '~')
//	if _, ok := pan["wlc_v"]; ok {
//		wlc_vStr := pan.Get("wlc_v")
//		wlc_v := strings.Split(wlc_vStr, ";")
//		for i := range wlc_v {
//			temp := strings.Split(wlc_v[i], "~")
//			if len(temp) > 2 {
//				temp[1] = strings.ReplaceAll(temp[1], ",", "")
//				ret, _ := strconv.ParseFloat(temp[1], 64)
//				p := message.NewPrinter(language.English)
//				temp[1] = p.Sprintf("%.2f", ret*bet)
//				wlc_v[i] = strings.Join(temp, "~")
//				fmt.Println(wlc_v[i])
//			}
//		}
//		wlc_vStr = strings.Join(wlc_v, ";")
//		pan.Set("wlc_v", wlc_vStr)
//	}
//	// 消除
//	for k := range pan {
//		// 正则表达式，使用捕获组
//		r := regexp.MustCompile(`^(l\d+)`)
//		// 查找匹配
//		match := r.FindStringSubmatch(k)
//		// 判断是否匹配并获取结果
//		if len(match) > 0 {
//			lstrs := strings.Split(pan[k], "~")
//			lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//			score, err := strconv.ParseFloat(lstrs[1], 64)
//			if err != nil {
//				panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
//			}
//			p := message.NewPrinter(language.English)
//			lstrs[1] = p.Sprintf("%.2f", score*bet)
//			pan[k] = strings.Join(lstrs, "~")
//		}
//
//	}
//	//for i := 0; i < pan.Int("l"); i++ {
//	//	l := fmt.Sprintf("l%d", i)
//	//	if lstr, ok := pan[l]; ok {
//	//		lstrs := strings.Split(lstr, "~")
//	//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//	//		score, err := strconv.ParseFloat(lstrs[1], 64)
//	//		if err != nil {
//	//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
//	//		}
//	//		p := message.NewPrinter(language.English)
//	//		lstrs[1] = p.Sprintf("%.2f", score*bet)
//	//		pan[l] = strings.Join(lstrs, "~")
//	//	}
//	//}
//
//	pan.SetCurrency("balance", balance)
//	pan.SetCurrency("balance_cash", balance)
//
//	return pan
//}
//
//// pp的数据进行解析成返回前端的格式
//func (sd *SimulateData) ParsingData(bet int64, balance float64, multiply int, enterBalance float64) Variables {
//	pan := sd.QueryString
//
//	//处理嵌入字段的分数*倍数
//	ReplaceEmbedScoreField(pan, multiply)
//
//	// 添加余额
//	AppendBalance(pan, bet, balance, multiply, enterBalance)
//
//	// 替换其他字段
//	ReplaceAllScoreField(pan, multiply)
//
//	pan.SetCurrency("balance", balance)
//	pan.SetCurrency("balance_cash", balance)
//	pan.SetCurrency("balance_bonus", 0.00)
//
//	pan.SetInt("stime", int(time.Now().UnixMilli()))
//
//	return pan
//}
//
//func GetBet(isBuyFree bool, coin float64, lines int, bl string, hasDouble int, FreeFlag int, TurnIndex int) int64 {
//	if isBuyFree {
//		bet := coin * float64(lines) // 0.01 * 20 = 0.20
//		// 待观察
//		bet *= 100
//		return int64(bet)
//	}
//
//	if FreeFlag == 0 && TurnIndex == 0 {
//		if hasDouble == 1 && bl == "1" { // 加倍模式
//			return int64(coin * float64(lines) * 1.25)
//		}
//		return int64(coin * float64(lines))
//	}
//
//	return 0
//}
//
//func parseToMap(query string) (Variables, error) {
//	// 解析查询字符串
//	values, err := url.ParseQuery(query)
//	if err != nil {
//		return nil, err
//	}
//
//	// 转换为 map[string]string
//	result := make(map[string]string)
//	for key, val := range values {
//		// url.ParseQuery 会返回 []string，即使每个键只有一个值
//		if len(val) > 0 {
//			result[key] = val[0]
//		}
//	}
//
//	return result, nil
//}

func convertToSlice(data any) ([]any, error) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, fmt.Errorf("input is not a slice or array")
	}

	result := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}
	return result, nil
}

func scatterValuesHandle(extendGameStateResult Variables, multiply float64) (err error) {
	scatterValuesList := extendGameStateResult["scatterValues"]

	if scatterValuesList != nil {

		slice, err := convertToSlice(scatterValuesList)

		if err == nil && len(slice) > 0 {
			arr1 := make([][]float64, 0)
			if len(slice) > 0 {

				for _, a := range slice {
					sliceHom, err := convertToSlice(a)
					if err != nil {
						return nil
					}
					arr2 := make([]float64, 0)
					for _, a := range sliceHom {

						arr2 = append(arr2, a.(float64)*multiply)
					}
					arr1 = append(arr1, arr2)
				}
			}
			extendGameStateResult["scatterValues"] = arr1
		}
	}

	//screenWinsInfo := extendGameStateResult["screenWinsInfo"]
	//
	//if screenWinsInfo != nil {
	//
	//	screenWinsInfos, _ := convertToSlice(screenWinsInfo)
	//
	//	for _, screenv := range screenWinsInfos {
	//
	//		screenval := convertToMapSafely(screenv)
	//		screenval.MKMulFloat("playerWin", 100)
	//
	//		quantityWinResults, _ := convertToSlice(screenval["quantityWinResult"])
	//
	//		for _, quantityv := range quantityWinResults {
	//			quantityval := convertToMapSafely(quantityv)
	//			quantityval.MKMulFloat("hitOdds", 100)
	//			quantityval.MKMulFloat("quantityWin", 100)
	//		}
	//	}
	//
	//}

	return err
}

func winC2(extendGameStateResult Variables, multiply float64) {
	winC2List := extendGameStateResult["winC2"]
	if winC2List != nil {
		slice, _ := convertToSlice(winC2List)
		arr1 := make([][]float64, 0)
		for _, v := range slice {
			slices, _ := convertToSlice(v)
			arr2 := make([]float64, 0)
			for _, val := range slices {
				arr2 = append(arr2, val.(float64)*multiply)
			}
			arr1 = append(arr1, arr2)
		}
		extendGameStateResult["winC2"] = arr1
	}

}
func gameDescriptor(extendGameStateResult Variables, multiply float64) {
	gameDes := extendGameStateResult["gameDescriptor"]
	if gameDes != nil {
		gameDesCom := convertToMapSafely(gameDes)
		slice, _ := convertToSlice(gameDesCom["component"])
		for _, v := range slice {
			slices, _ := convertToSlice(v)
			for _, val := range slices {
				valPlace := convertToMapSafely(val)
				if valPlace["placeholders"] != nil {
					slices1, _ := convertToSlice(valPlace["placeholders"])
					for _, valPlaceholder := range slices1 {
						valPlaceholderMap := convertToMapSafely(valPlaceholder)
						if valPlaceholderMap["type"] == "score" {
							valPlaceholderMap["value"] = MKMulFloat(valPlaceholderMap["value"].(string), multiply)
						}
					}
				}
			}
		}
		extendGameStateResult["gameDescriptor"] = gameDesCom
	}
}

func quantityGameResult(quantityGameResult Variables, multiply float64) {
	quantityGameResultMap := convertToMapSafely(quantityGameResult["quantityGameResult"])
	quantityGameResultMap.MKMulFloat("playerWin", multiply)
	quantityGameResultSlice, _ := convertToSlice(quantityGameResultMap["quantityWinResult"])
	for _, v := range quantityGameResultSlice {
		val := convertToMapSafely(v)
		val.MKMulFloat("quantityWin", multiply)
	}

	lineGameResultMap := convertToMapSafely(quantityGameResult["lineGameResult"])
	lineGameResultMap.MKMulFloat("playerWin", multiply)
	lineWinResultSlice, _ := convertToSlice(lineGameResultMap["lineWinResult"])
	for _, v := range lineWinResultSlice {
		val := convertToMapSafely(v)
		val.MKMulFloat("lineWin", multiply)
	}

	cascadeEliminateResultSlice, _ := convertToSlice(quantityGameResult["cascadeEliminateResult"])
	for _, v := range cascadeEliminateResultSlice {
		val := convertToMapSafely(v)
		/////////////////////////////////////////////////////
		vval := convertToMapSafely(val["quantityGameResult"])
		vval.MKMulFloat("playerWin", multiply)
		vvval, _ := convertToSlice(vval["quantityWinResult"])
		for _, vv := range vvval {
			vvv := convertToMapSafely(vv)
			vvv.MKMulFloat("quantityWin", multiply)
		}
		/////////////////////////////////////////////////////
		sval := convertToMapSafely(val["lineGameResult"])
		sval.MKMulFloat("playerWin", multiply)
		ssval, _ := convertToSlice(sval["lineWinResult"])
		for _, ss := range ssval {
			sss := convertToMapSafely(ss)
			sss.MKMulFloat("lineWin", multiply)
		}
	}
}

func scatterScoreRecordHandle(extendGameStateResult Variables, handleKey string, multiply float64) (err error) {
	scatterValuesList := extendGameStateResult[handleKey]

	if scatterValuesList != nil {

		slice, err := convertToSlice(scatterValuesList)

		if err == nil && len(slice) > 0 {
			arr1 := make([][]float64, 0)
			if len(slice) > 0 {

				for _, a := range slice {
					sliceHom, err := convertToSlice(a)
					if err != nil {
						return nil
					}
					arr2 := make([]float64, 0)
					for _, a := range sliceHom {

						arr2 = append(arr2, a.(float64)*multiply)
					}
					arr1 = append(arr1, arr2)
				}
			}
			extendGameStateResult[handleKey] = arr1
		}
	}
	return err
}

func convertMongoDataHandle(parentVariables Variables, currentKey, handleString string, multiply float64) {
	wayWinResultList, err := convertMongoData(parentVariables[currentKey])

	if err == nil && len(wayWinResultList) > 0 {

		for _, wayWin := range wayWinResultList {
			wayWin.MKMulFloat(handleString, multiply)
		}

	}
}

func convertToMapSafelyHandle(parentVariables Variables, currentKey, handleString string, multiply float64) {
	wayWinResultList := convertToMapSafely(parentVariables[currentKey])

	if wayWinResultList != nil {
		wayWinResultList.MKMulFloat(handleString, multiply)

	}
}

func continuousWinHandle(parentVariables Variables, multiply float64) {
	connectWinResult := convertToMapSafely(parentVariables["connectGameResult"])
	if connectWinResult != nil {
		connectWinResult.MKMulFloat("playerWin", multiply)
		convertMongoDataHandle(connectWinResult, "connectWinResult", "connectWin", multiply)
	}
	//convertToMapSafelyHandle(parentVariables, "connectWinResult", "playerWin", multiply)

	cascadeEliminateResultList, err := convertMongoData(parentVariables["cascadeEliminateResult"])
	if err == nil && len(cascadeEliminateResultList) > 0 {
		for _, cascadeEliminate := range cascadeEliminateResultList {

			connectGameResult := convertToMapSafely(cascadeEliminate["connectGameResult"])
			if connectGameResult != nil {
				connectGameResult.MKMulFloat("playerWin", multiply)
			}

			convertMongoDataHandle(connectGameResult, "connectWinResult", "connectWin", multiply)
			convertMongoDataHandle(connectGameResult, "specialFeatureResult", "specialScreenWin", multiply)

			cascadeEliminate.MKMulFloat("eliminateWin", multiply)
			convertToMapSafelyHandle(cascadeEliminate, "wayGameResult", "playerWin", multiply)
			wayGameResult := convertToMapSafely(cascadeEliminate["wayGameResult"])
			data, err := convertMongoData(wayGameResult["wayWinResult"])

			if err == nil && len(data) > 0 {
				for _, wayWin := range data {
					wayWin.MKMulFloat("symbolWin", multiply)
				}
			}

		}
	}
}

func holdAndSpinResultHandle(parentVariables Variables, multiply float64) {
	holdAndSpinResultList, err := convertMongoData(parentVariables["holdAndSpinResult"])
	if err == nil && len(holdAndSpinResultList) > 0 {
		for _, holdAndSpinResult := range holdAndSpinResultList {
			values, _ := convertToSlice(holdAndSpinResult["value"])
			preValues, _ := convertToSlice(holdAndSpinResult["preValue"])
			merge := convertToMapSafely(holdAndSpinResult["merge"])
			multiplier := convertToMapSafely(holdAndSpinResult["multiplier"])
			var datas, preDatas [][]float64
			mergeDateMap := make(map[string][][]float64)
			multiplierDateMap := make(map[string][][]float64)
			for _, val := range values {
				var data []float64
				val, _ := convertToSlice(val)
				for _, v := range val {
					data = append(data, v.(float64)*multiply)
				}
				datas = append(datas, data)
			}
			for _, preVal := range preValues {
				var preData []float64
				preVal, _ := convertToSlice(preVal)
				for _, v := range preVal {
					preData = append(preData, v.(float64)*multiply)
				}
				preDatas = append(preDatas, preData)
			}
			for key, val := range merge {
				var mergeDates [][]float64
				val, _ := convertToSlice(val)
				for _, v := range val {
					var mergeDate []float64
					v, _ := convertToSlice(v)
					for _, vv := range v {
						mergeDate = append(mergeDate, vv.(float64)*multiply)
					}
					mergeDates = append(mergeDates, mergeDate)
				}
				mergeDateMap[key] = mergeDates
			}
			for key, val := range multiplier {
				var multiplierDates [][]float64
				val, _ := convertToSlice(val)
				for _, v := range val {
					var multiplierDate []float64
					v, _ := convertToSlice(v)
					for _, vv := range v {
						multiplierDate = append(multiplierDate, vv.(float64)*multiply)
					}
					multiplierDates = append(multiplierDates, multiplierDate)
				}
				multiplierDateMap[key] = multiplierDates
			}

			holdAndSpinResult_lineGameResult := convertToMapSafely(holdAndSpinResult["lineGameResult"])
			if holdAndSpinResult_lineGameResult != nil {

				convertMongoDataHandle(holdAndSpinResult_lineGameResult, "lineWinResult", "lineWin", multiply)
			}

			holdAndSpinResult["value"] = datas
			holdAndSpinResult["preValue"] = preDatas
			holdAndSpinResult["merge"] = mergeDateMap
			holdAndSpinResult["multiplier"] = multiplierDateMap
		}
	}
}

func scatterTotalScoreHandle(extendGameState Variables) {
	total := 0.0
	scatterScoreRecord, err := convertToSlice(extendGameState["scatterScoreRecord"])
	if err == nil && len(scatterScoreRecord) > 0 {
		for _, scatterScoreArray := range scatterScoreRecord {
			array, _ := convertToSlice(scatterScoreArray)
			for _, item := range array {
				total += item.(float64)
			}
		}
	}
	_, ok := extendGameState["scatterTotalScore"]
	if ok {
		extendGameState["scatterTotalScore"] = total
		fmt.Println(extendGameState["scatterTotalScore"])
	}
}

func reSpinMutiplyHandle(gameResult Variables, multiply float64) {
	cascadeEliminateList, err := convertMongoData(gameResult["cascadeEliminateResult"])
	if err == nil && len(cascadeEliminateList) > 0 {
		for _, element := range cascadeEliminateList {
			convertMongoDataHandle(element, "specialFeatureResult", "specialScreenWin", multiply)
		}
	}
}
