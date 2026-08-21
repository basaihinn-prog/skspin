package rtp

import (
	"fmt"
	"os"
	"serve/service_fish/application/lottery"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/models"
	"strconv"
	"strings"
	"testing"
)

const gameId00002 = models.PSF_ON_00002

var mathModuleId00002 string
var rtpList00002 = map[int][]string{
	0: {"0", "20", "40", "60"},
	1: {"80", "100", "150", "200", "300"},
}

func TestService00002(t *testing.T) {
	// Process Input Value
	if len(os.Args) < 5 || len(os.Args) >= 7 {
		fmt.Println("參數說明:數學代號 總次數 指定魚種 隨機下注金額(true or false) 武器魚種驗證(true or false)")
		return
	}

	// 數學代號
	mathModuleId00002 = chooseMath(gameId00002, os.Args[1])
	if mathModuleId00002 == "" {
		fmt.Println("數學代號不存在")
		return
	}
	if mathModuleId00002 == " " {
		fmt.Println("遊戲代號不存在")
		return
	}

	// 總次數
	runTimes, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("總次數輸入錯誤")
	}

	// 指定魚種
	var selectFishList []int32
	if os.Args[3] == "-1" {
		selectFishList = fishList[gameId00002]
	} else {
		inputFishIdList := strings.Split(os.Args[3], ",")
		for _, inputFishId := range inputFishIdList {
			if !isContain(fishList[gameId00002], inputFishId) {
				fmt.Println("魚種不存在")
				return
			}
		}

		for _, inputFishStr := range inputFishIdList {
			fishId, _ := strconv.Atoi(inputFishStr)
			selectFishList = append(selectFishList, int32(fishId))
		}
	}

	// 隨機下注金額
	var checkRandomBet = false
	if os.Args[4] == "true" {
		checkRandomBet = true
	}

	// 武器驗證
	var weaponFish = false
	if len(os.Args) == 6 {
		if os.Args[5] == "true" {
			if len(selectFishList) != 1 {
				fmt.Println("武器魚驗證只能指定一種魚種")
				return
			}

			weaponFish = true
		}
	}

	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	var checkRtpState = -1
	var lastRtpBudget uint64 = 0
	var rtpStateTimes = make(map[int32]map[int]map[string]uint64, 0)
	var rtpBudgets = make(map[int32]map[int]map[string]uint64, 0)
	var resultTimes = make(map[int32]map[int]map[string]uint64, 0)
	var resultHits = make(map[int32]map[int]map[string]uint64, 0)
	var resultWins = make(map[int32]map[int]uint64, 0)
	var resultHitEnvelope = make(map[int32]map[int]map[string]uint64, 0)
	var resultEnvelope = make(map[int32]map[int]uint64, 0)
	var machineGunTimes = make(map[int32]uint64, 0)
	var machineGunHits = make(map[int32]uint64, 0)
	var drillTimes = make(map[int32]uint64, 0)
	var drillHits = make(map[int32]uint64, 0)

	var betTimes = make(map[int32]map[int]uint64, 0)

	for i := 1; i <= runTimes*10000; i++ {
		// Process Random Bet
		var bet = defaultBet
		if checkRandomBet {
			bet = randomBet[(i-1)%3]
		}

		fishId := rngFishIdNew(selectFishList)

		_, result, rtpState, rtpId, _, rtpBudget, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId00002, mathModuleId00002, "", "", "", "",
			-1, fishId, -1,
			uint64(bet), defaultRate, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0, 0,
		)

		hit, win, envelope, _, _, _, _, _, _,
			machineGunTime, machineGunHit,
			_, _,
			drillTime, drillHit, _ := resultProcess(gameId00002, mathModuleId00002, fishId, result, uint64(bet), weaponFish)

		// Process 水位轉換次數
		rtpStateTimes = verifyMap(rtpStateTimes, fishId, rtpState)
		if checkRtpState == -1 || checkRtpState != rtpState {
			rtpStateTimes[fishId][rtpState][rtpId]++
			checkRtpState = rtpState
		}

		// Process 水庫子彈數量
		rtpBudgets = verifyMap(rtpBudgets, fishId, rtpState)
		if lastRtpBudget == 0 || lastRtpBudget < rtpBudget {
			rtpBudgets[fishId][rtpState][rtpId] += rtpBudget + (uint64(bet) * defaultRate)
		}
		lastRtpBudget = rtpBudget

		// Process 打擊次數
		resultTimes = verifyMap(resultTimes, fishId, rtpState)
		resultTimes[fishId][rtpState][rtpId]++

		// Process 打中的次數
		resultHits = verifyMap(resultHits, fishId, rtpState)
		resultHits[fishId][rtpState][rtpId] += hit

		// Process 贏分
		resultWins = verifySamipleMap(resultWins, fishId)
		resultWins[fishId][rtpState] += win

		// Process Bonus紅包次數
		resultHitEnvelope = verifyMap(resultHitEnvelope, fishId, rtpState)
		if envelope > 0 {
			if result.TriggerIconId == 31 {
				resultHitEnvelope[fishId][rtpState][rtpId]++
			}

			if result.TriggerIconId == 26 || result.TriggerIconId == 27 {
				if result.ExtraTriggerBonus != nil {
					if result.ExtraTriggerBonus[0].TriggerIconId == 31 {
						resultHitEnvelope[fishId][rtpState][rtpId]++
					}
				}
			}
		}

		// Process Bonus紅包贏分
		resultEnvelope = verifySamipleMap(resultEnvelope, fishId)
		resultEnvelope[fishId][rtpState] += envelope

		// Process 機關槍次數
		for k, v := range machineGunTime {
			machineGunTimes[k] += v
		}

		// Process 鑽頭砲次數
		for k, v := range drillTime {
			drillTimes[k] += v
		}

		// Process 機關槍擊中數
		for k, v := range machineGunHit {
			machineGunHits[k] += v
		}

		// Process 鑽頭砲擊中數
		for k, v := range drillHit {
			drillHits[k] += v
		}

		// Process BetHit Times
		betTimes = verifySamipleMap(betTimes, fishId)
		betTimes[fishId][bet]++
	}

	var printData = make(map[int32]string, 0)
	showString := ""

	for _, oneFish := range fishList[gameId00002] {
		showString = fmt.Sprint(oneFish)

		// Rtp State Times
		showString = getShowData(showString, oneFish, rtpStateTimes, rtpList00002)

		// Rtp Bullet Count
		showString = getShowData(showString, oneFish, rtpBudgets, rtpList00002)

		// Fish Times
		showString = getShowData(showString, oneFish, resultTimes, rtpList00002)

		// Hit Times
		showString = getShowData(showString, oneFish, resultHits, rtpList00002)

		// Wins
		for rtpState := 0; rtpState < 2; rtpState++ {
			showString = fmt.Sprint(showString, " ", resultWins[oneFish][rtpState])
		}

		// RedEnvelop Hits
		showString = getShowData(showString, oneFish, resultHitEnvelope, rtpList00002)

		// RedEnvelop Wins
		for rtpState := 0; rtpState < 2; rtpState++ {
			showString = fmt.Sprint(showString, " ", resultEnvelope[oneFish][rtpState])
		}

		// MachineGun Times
		showString = fmt.Sprint(showString, " ", machineGunTimes[oneFish])
		// Drill Times
		showString = fmt.Sprint(showString, " ", drillTimes[oneFish])
		// MachineGun Hits
		showString = fmt.Sprint(showString, " ", machineGunHits[oneFish])
		// Drill Hits
		showString = fmt.Sprint(showString, " ", drillHits[oneFish])

		if checkRandomBet {
			// BetHit Times
			for _, oneBet := range randomBet {
				showString = fmt.Sprint(showString, " ", betTimes[oneFish][oneBet])
			}
		}

		printData[oneFish] += showString
	}

	for _, oneFish := range fishList[gameId00002] {
		fmt.Println(printData[oneFish])
	}
}
