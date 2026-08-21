package report

import (
	"fmt"
	"os"
	"serve/fish_comm/rng"
	"serve/service_fish/application/lottery"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	PSFM_00008_97_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-97-1"
	"serve/service_fish/models"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

var fishList = map[string][]int32{
	models.PSF_ON_00001: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
	models.PSF_ON_00002: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
	models.PSF_ON_00003: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 303},
	models.PSF_ON_00004: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400},
	models.PSF_ON_00005: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 501, 502, 503, 504, 505},
	models.PSF_ON_00007: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400},
}
var triggerIconList = map[string]int{
	models.PSF_ON_00001: 31, models.PSF_ON_00002: 31, models.PSF_ON_00003: -1,
	models.PSF_ON_00004: -1, models.PSF_ON_00005: 102, models.PSF_ON_00007: 102,
}
var BonusIconList = map[string]map[int32]string{
	models.PSF_ON_00001: {
		26: PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.Drill.RTP8.ID,
		27: PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID,
	},
	models.PSF_ON_00002: {
		26: PSFM_00003_98_1.RTP98BS.PSF_ON_00002_1_BsMath.Icons.Drill.RTP8.ID,
		27: PSFM_00003_98_1.RTP98BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID,
	},
	models.PSF_ON_00003: {
		12: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp,
		13: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie2Drill.UseRtp,
		14: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie3Drill.UseRtp,
		15: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie4Drill.UseRtp,
		16: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie5Drill.UseRtp,
		17: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie6Drill.UseRtp,
		31: PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Bullets.UseRtp,
	},
	models.PSF_ON_00004: {
		201: PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath.Icons.MachineGun.UseRTP,
		202: PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath.Icons.SuperMachineGun.UseRTP,
	},
	models.PSF_ON_00005: {
		201: PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.MachineGun.UseRTP,
		202: PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.SuperMachineGun.UseRTP,
	},
	models.PSF_ON_00007: {
		201: PSFM_00008_97_1.RTP97BS.PSF_ON_00007_1_BsMath.Icons.MachineGun.UseRTP,
		202: PSFM_00008_97_1.RTP97BS.PSF_ON_00007_1_BsMath.Icons.SuperMachineGun.UseRTP,
	},
}

var mathModuleId string
var secWebSocketKey = "jerry"
var reportFileName string

/*
*
參數說明 :
 1. 遊戲代號 : Ex: PSF-ON-00007
 2. 數學代號 : 95、96、97、98
 3. 總次數 : 萬次(單位)
 4. 指定魚種 : 如要隨機魚種使用"-1"，用 "," 進行區隔。 Ex : 0,5,8
 5. 計算RTP的數量 : 每X次就記錄一次目前的RTP率

EX :

	XXX.exe PSF-ON-00001 95 100 0,1,2 10
*/
func TestService(t *testing.T) {
	if len(os.Args) != 6 {
		fmt.Println("參數說明:遊戲代號 數學代號 總次數 指定魚種 計算RTP的數量")
		return
	}

	gameId := os.Args[1]

	// 數學代號
	mathModuleId = chooseMath(gameId, os.Args[2])
	if mathModuleId == "" {
		fmt.Println("數學代號不存在")
		return
	}
	if mathModuleId == " " {
		fmt.Println("遊戲代號不存在")
		return
	}

	// 總次數
	runTimes, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("總次數輸入錯誤")
	}

	// 指定魚種
	var selectFishList []int32
	if os.Args[4] == "-1" {
		selectFishList = fishList[gameId]
	} else {
		inputFishIdList := strings.Split(os.Args[4], ",")
		for _, inputFishId := range inputFishIdList {
			if !isContain(fishList[os.Args[1]], inputFishId) {
				fmt.Println("魚種不存在")
				return
			}
		}

		for _, inputFishStr := range inputFishIdList {
			fishId, _ := strconv.Atoi(inputFishStr)
			selectFishList = append(selectFishList, int32(fishId))
		}
	}

	// 計算RTP的數量
	mileStone, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Println("計算RTP的數量錯誤")
		return
	}

	reportFileName = gameId + "_" + mathModuleId + "_Report.xlsx"

	var times, wins, envelopes, bullets uint64 = 0, 0, 0, 0
	reportMap := make(map[int]string, 0)
	var timesList = make(map[int32]uint64)
	var hitsList = make(map[int32]uint64)

	var totalBet uint64 = 0
	reportBetMap := make(map[int]uint64, 0)
	reportWinMap := make(map[int]uint64, 0)

	for i := 1; i <= runTimes*10000; i++ {
		hitFish := fish.Fish{}
		hitBullet := bullet.Bullet{}

		fishId := rngFishId(selectFishList)
		_, result, _, _, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId, mathModuleId, "", "", "", "",
			-1, fishId, -1,
			1, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0,
		)
		totalBet += 1

		hit, win, envelope, resultBullet := processResult(gameId, fishId, result)
		times++
		timesList[fishId]++
		hitsList[fishId] += hit
		wins += win
		envelopes += envelope
		bullets += resultBullet

		// Process Mercenary Bullet
		mercenaryWin, leftBullet := mercenaryBulletProcess(gameId, bullets, selectFishList)
		wins += mercenaryWin
		bullets = leftBullet

		// Record
		if i%mileStone == 0 {
			reportMap[i] = getRate(wins+envelopes+bullets, times)
			reportBetMap[i] = totalBet
			reportWinMap[i] = wins + envelopes + bullets
		}
	}

	// Print Data
	keyList := sortMapByKey(reportMap)

	//fmt.Println("Times", "RtpRate")
	fmt.Println("Times", "Bet", "Wins")
	for _, k := range keyList {
		//fmt.Println(k, reportMap[k])
		fmt.Println(k, reportBetMap[k], reportWinMap[k])
	}
	fmt.Println("Total", totalBet, wins+envelopes+bullets)

	fmt.Println("FishId", "Times")
	for _, v := range fishList[gameId] {
		fmt.Println(v, timesList[v], hitsList[v])
	}

	//// Export Excel File
	//f, index := reportCreateExcel("RTP_" + os.Args[2])
	//for i := 0; i < len(keyList); i++ {
	//	reportSetCellValue(f, "RTP_" + os.Args[2], i + 2, keyList[i], reportMap[keyList[i]])
	//}
	//f.SetActiveSheet(index)
	//
	//if err := f.SaveAs(reportFileName); err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//fmt.Println("Excel File Created. FileName: " + reportFileName)
}

func reportCreateExcel(sheetName string) (file *excelize.File, sheetIndex int) {
	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Header
	headers := &[]string{"次數", "RtpRate"}
	f.SetSheetRow(sheetName, "A1", headers)

	return f, index
}

func reportSetCellValue(file *excelize.File, sheetName string, cellIndex int, mileStone int, rtpRate string) {
	//strCell := strconv.Itoa(cellIndex)
	values := []interface{}{mileStone, rtpRate}

	for index, value := range values {
		reportGetCol(file, sheetName, index, strconv.Itoa(cellIndex), value)
	}
}

func reportGetCol(file *excelize.File, sheetName string, rowIndex int, cellIndex string, value interface{}) {
	switch value.(type) {
	case int:
		file.SetCellValue(sheetName, getAxis(rowIndex, cellIndex), value.(int))
	case uint64:
		file.SetCellValue(sheetName, getAxis(rowIndex, cellIndex), value.(uint64))
	case string:
		if strings.HasPrefix(value.(string), "=") {
			file.SetCellFormula(sheetName, getAxis(rowIndex, cellIndex), value.(string)[1:len(value.(string))])
		} else {
			file.SetCellValue(sheetName, getAxis(rowIndex, cellIndex), value.(string))
		}
	}
}

func getAxis(rowIndex int, cellIndex string) string {
	row := ""

	switch rowIndex {
	case 0:
		row = "A" + cellIndex
	case 1:
		row = "B" + cellIndex
	case 2:
		row = "C" + cellIndex
	case 3:
		row = "D" + cellIndex
	case 4:
		row = "E" + cellIndex
	case 5:
		row = "F" + cellIndex
	case 6:
		row = "G" + cellIndex
	case 7:
		row = "H" + cellIndex
	case 8:
		row = "I" + cellIndex
	case 9:
		row = "J" + cellIndex
	case 10:
		row = "K" + cellIndex
	case 11:
		row = "L" + cellIndex
	case 12:
		row = "M" + cellIndex
	case 13:
		row = "N" + cellIndex
	case 14:
		row = "O" + cellIndex
	}

	return row
}

func mercenaryBulletProcess(gameId string, bullets uint64, selectFishList []int32) (bulletWin, leftBullet uint64) {
	bulletWin = 0
	leftBullet = 0

	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	for {
		times := bullets / 60
		leftBullet = bullets - (times * 60)

		for i := 0; uint64(i) < (times * 60); i++ {
			fishId := rngFishId(selectFishList)

			_, result, _, _, _, _, _, _ := lottery.Service.MathProcess(
				secWebSocketKey, gameId, mathModuleId, "", "", "", BonusIconList[gameId][31],
				1, fishId, -1,
				1, 1, 0,
				&hitBullet, &hitFish,
				lottery.PLAYER, "",
				0, 0,
				true, 1,
			)

			if result.Pay > 0 {
				bulletWin += uint64(result.Pay * result.Multiplier)

				// Slot
				if fishId == 6 || fishId == 7 || fishId == 8 ||
					fishId == 9 || fishId == 10 || fishId == 11 {
					bulletWin += uint64(result.ExtraData[0].(int))
				}

				// Drill
				if fishId == 12 || fishId == 13 || fishId == 14 ||
					fishId == 15 || fishId == 16 || fishId == 17 {
					pay, drillBullet := shotDrill(gameId, BonusIconList[gameId][fishId], result.BonusPayload.(int), 200)
					bulletWin += pay
					leftBullet += drillBullet
				}
			}

			if result.Bullet > 0 {
				leftBullet += uint64(result.Bullet)
			}
		}

		if leftBullet < 60 {
			break
		}
	}

	return bulletWin, leftBullet
}

func sortMapByKey(mapData map[int]string) (sortKey []int) {
	var keys []int
	for k := range mapData {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}

func processResult(gameId string, fishId int32, result *probability.Probability) (hits, wins, envelopes, bullet uint64) {
	hits = 0
	wins = 0
	envelopes = 0
	bullet = 0

	if result.Pay > 0 {
		hits++
		wins += uint64(result.Pay * result.Multiplier)
	}

	// Mercenary Bullet
	if result.Bullet > 0 {
		bullet = uint64(result.Bullet)
	}

	if gameId != models.PSF_ON_00004 && gameId != models.PSF_ON_00003 {
		if result.TriggerIconId == triggerIconList[gameId] {
			envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId]))
		}
	}

	switch gameId {
	case models.PSF_ON_00001, models.PSF_ON_00002:
		// Drill
		if fishId == 26 && result.Pay > 0 {
			p, _ := shotDrill(gameId, BonusIconList[gameId][fishId], result.BonusPayload.(int), fishId)
			//p := shot_26(gameId, mathModule_id, BonusIconList[gameId][fishId], result.BonusPayload.(int))
			wins += p
		}

		// MachineGun
		if fishId == 27 && result.Pay > 0 {
			p, r := machineGunShot(gameId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId,
			)
			wins += p
			envelopes += r
		}

		// Red Envelope
		if fishId == 28 && result.TriggerIconId == 28 {
			hits++
			wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		}

		// Slot
		if fishId == 29 && result.TriggerIconId == 29 {
			hits++
			wins += uint64(result.ExtraData[0].(int))
		}

	default:
		// Red Envelope
		if fishId == 100 && result.TriggerIconId == 100 {
			hits++
			wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		}

		// Slot
		if gameId == models.PSF_ON_00003 {
			if (fishId == 6 || fishId == 7 || fishId == 8 ||
				fishId == 9 || fishId == 10 || fishId == 11) && result.TriggerIconId == 101 {
				hits++
				wins += uint64(result.ExtraData[0].(int))
			}
		} else {
			if fishId == 101 && result.TriggerIconId == 101 {
				hits++
				wins += uint64(result.ExtraData[0].(int))
			}
		}

		// MachineGun
		if fishId == 201 && result.Pay > 0 {
			p, r := machineGunShot(gameId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId,
			)
			wins += p
			envelopes += r
		}

		// SuperMachineGun
		if fishId == 202 && result.Pay > 0 {
			p, r := machineGunShot(gameId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId,
			)
			wins += p
			envelopes += r
		}

		// Drill
		if gameId == models.PSF_ON_00003 {
			if fishId == 12 || fishId == 13 || fishId == 14 ||
				fishId == 15 || fishId == 16 || fishId == 17 {
				p, b := shotDrill(gameId, BonusIconList[gameId][fishId], result.BonusPayload.(int), fishId)
				wins += p
				bullet += b
			}
		} else {
			// TODO
		}
	}

	// Fruit Dish
	if fishId == 300 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId]))
	}

	// A pack of beer
	if fishId == 301 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId]))
	}

	// XIAO LONG BAO
	if fishId == 302 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId]))
	}

	return hits, wins, envelopes, bullet
}

func rngFishId(fishList []int32) int32 {
	options := make([]rng.Option, 0)

	for _, v := range fishList {
		options = append(options, rng.Option{Weight: 1, Item: v})
	}

	return rng.New(options).Item.(int32)
}

func isContain(itemList []int32, fishIdStr string) bool {
	fishId, err := strconv.Atoi(fishIdStr)
	if err != nil {
		return false
	}

	for _, item := range itemList {
		if int32(fishId) == item {
			return true
		}
	}
	return false
}

func chooseMath(gameId, inputMath string) string {
	switch gameId {
	case models.PSF_ON_00001:
		switch inputMath {
		case "95":
			return models.PSFM_00002_95_1
		case "96":
			return models.PSFM_00002_96_1
		case "97":
			return models.PSFM_00002_97_1
		case "98":
			return models.PSFM_00002_98_1
		default:
			return ""
		}

	case models.PSF_ON_00002:
		switch inputMath {
		case "95":
			return models.PSFM_00003_95_1
		case "96":
			return models.PSFM_00003_96_1
		case "97":
			return models.PSFM_00003_97_1
		case "98":
			return models.PSFM_00003_98_1
		default:
			return ""
		}

	case models.PSF_ON_00003:
		switch inputMath {
		case "95":
			return models.PSFM_00004_95_1
		case "96":
			return models.PSFM_00004_96_1
		case "97":
			return models.PSFM_00004_97_1
		case "98":
			return models.PSFM_00004_98_1
		default:
			return ""
		}

	case models.PSF_ON_00004:
		switch inputMath {
		case "95":
			return models.PSFM_00005_95_1
		case "96":
			return models.PSFM_00005_96_1
		case "97":
			return models.PSFM_00005_97_1
		case "98":
			return models.PSFM_00005_98_1
		default:
			return ""
		}

	case models.PSF_ON_00005:
		switch inputMath {
		case "95":
			return models.PSFM_00006_95_1
		case "96":
			return models.PSFM_00006_96_1
		case "97":
			return models.PSFM_00006_97_1
		case "98":
			return models.PSFM_00006_98_1
		default:
			return ""
		}

	case models.PSF_ON_00007:
		switch inputMath {
		case "95":
			return models.PSFM_00008_95_1
		case "96":
			return models.PSFM_00008_96_1
		case "97":
			return models.PSFM_00008_97_1
		case "98":
			return models.PSFM_00008_98_1
		default:
			return ""
		}
	default:
		return " "
	}
}

func rngRedEnvelope() int32 {
	options := make([]rng.Option, 0, 5)

	for i := 0; i < 5; i++ {
		options = append(options, rng.Option{Weight: 1, Item: i})
	}

	return int32(rng.New(options).Item.(int))
}

func getTriggerEnvelope(gameId, mathModuleId string, bonusTypeId int, triggerIconId int32) uint64 {
	result := probability.Service.Calc(
		gameId,
		mathModuleId,
		strconv.Itoa(bonusTypeId),
		triggerIconId,
		-1,
		0,
		0,
		0,
	)

	return uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
}

func getRate(x, y uint64) string {
	result := float32(x) / float32(y) * 100
	return fmt.Sprintf("%.2f", result)
}

func machineGunShot(gameId, rtpId string, bullets int, machineGunType int32) (pay, redEnvelope uint64) {
	var fishId int32 = -1
	var newBullets = 0
	var checkMachineGunType int32 = 0

	switch gameId {
	case models.PSF_ON_00001, models.PSF_ON_00002:
		if machineGunType == 27 {
			checkMachineGunType = 26
		}
	default:
		if machineGunType == 201 {
			checkMachineGunType = 202
		} else {
			checkMachineGunType = 201
		}
	}

	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId(fishList[gameId])
			if fishId != checkMachineGunType {
				break
			}
		}

		_, result, _, _, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId, mathModuleId, "", "", "", BonusIconList[gameId][machineGunType],
			-1, fishId, machineGunType,
			1, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 1,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay * result.Multiplier)
		}

		switch gameId {
		case models.PSF_ON_00001, models.PSF_ON_00002:
			if result.TriggerIconId == 28 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}

			if result.TriggerIconId == 29 {
				pay += uint64(result.ExtraData[0].(int))
			}
		default:
			if result.TriggerIconId == 100 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}

			if result.TriggerIconId == 101 {
				pay += uint64(result.ExtraData[0].(int))
			}
		}

		if result.ExtraTriggerBonus != nil {
			if result.ExtraTriggerBonus[0].TriggerIconId == triggerIconList[gameId] {
				result := probability.Service.Calc(
					gameId,
					mathModuleId,
					strconv.Itoa(result.ExtraTriggerBonus[0].BonusTypeId),
					int32(triggerIconList[gameId]),
					machineGunType,
					0,
					0,
					0,
				)
				redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}
		} else {
			if result.TriggerIconId == triggerIconList[gameId] {
				result := probability.Service.Calc(
					gameId,
					mathModuleId,
					strconv.Itoa(result.BonusTypeId),
					int32(triggerIconList[gameId]),
					machineGunType,
					0,
					0,
					0,
				)

				if gameId != models.PSF_ON_00004 {
					redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				}
			}
		}

		if result.TriggerIconId == int(machineGunType) {
			newBullets += result.BonusPayload.(int)
			if newBullets > 999 {
				newBullets = 999
			}
		}
	}

	if newBullets > 0 {
		p, r := machineGunShot(gameId, rtpId, newBullets, machineGunType)
		pay += p
		redEnvelope += r
	}

	return pay, redEnvelope
}

func shotDrill(gameId, rtpId string, bullets int, bulletTypeId int32) (pay, drillBullet uint64) {
	var fishId int32 = -1
	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	for i := 0; i < bullets; i++ {
		// Random Fish ID
		for {
			fishId = rngFishId(fishList[gameId])

			if gameId == models.PSF_ON_00001 || gameId == models.PSF_ON_00002 {
				if fishId != 26 && fishId != 27 {
					break
				}
			}
			if gameId == models.PSF_ON_00003 {
				if fishId != 12 && fishId != 13 && fishId != 14 &&
					fishId != 15 && fishId != 16 && fishId != 17 {
					break
				}
			}
		}

		_, result, _, _, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId, mathModuleId, "", "", "", rtpId,
			-1, fishId, bulletTypeId,
			1, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 1,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay * result.Multiplier)

			switch gameId {
			case models.PSF_ON_00003:
				// Slot
				if fishId == 6 || fishId == 7 || fishId == 8 ||
					fishId == 9 || fishId == 10 || fishId == 11 {
					pay += uint64(result.ExtraData[0].(int))
				}
			}
		}

		if gameId == models.PSF_ON_00001 || gameId == models.PSF_ON_00002 {
			if result.TriggerIconId == 28 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}

			if result.TriggerIconId == 29 {
				pay += uint64(result.ExtraData[0].(int))
			}

			if result.ExtraTriggerBonus != nil {
				if result.ExtraTriggerBonus[0].TriggerIconId == 31 {
					result := probability.Service.Calc(
						gameId,
						mathModuleId,
						strconv.Itoa(result.ExtraTriggerBonus[0].BonusTypeId),
						int32(triggerIconList[gameId]),
						bulletTypeId,
						0,
						0,
						0,
					)

					pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				}
			} else {
				if result.TriggerIconId == 31 {
					result := probability.Service.Calc(
						gameId,
						mathModuleId,
						strconv.Itoa(result.BonusTypeId),
						int32(triggerIconList[gameId]),
						bulletTypeId,
						0,
						0,
						0,
					)

					pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				}
			}
		}

		if result.Bullet > 0 {
			drillBullet += uint64(result.Bullet)
		}
	}

	return pay, drillBullet
}
