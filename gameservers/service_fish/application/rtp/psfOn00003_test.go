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

const (
	gameId00003 = models.PSF_ON_00003
)

var rtpList00003 = map[int][]string{
	0: {"120", "150", "180", "200", "250"},
	1: {"30", "50", "60", "70", "80"},
}
var netWinGroupList = []int{50, 80, 100, 150, 200}
var mathModuleId00003 string

func TestService00003(t *testing.T) {
	// Process Input Value
	if len(os.Args) != 4 {
		fmt.Println("參數說明:數學代號 總次數 指定魚種")
		return
	}

	// 數學代號
	mathModuleId00003 = chooseMath(gameId00003, os.Args[1])
	if mathModuleId00003 == "" {
		fmt.Println("數學代號不存在")
		return
	}
	if mathModuleId00003 == " " {
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
		selectFishList = fishList[gameId00003]
	} else {
		inputFishIdList := strings.Split(os.Args[3], ",")
		for _, inputFishId := range inputFishIdList {
			if !isContain(fishList[gameId00003], inputFishId) {
				fmt.Println("魚種不存在")
				return
			}
		}

		for _, inputFishStr := range inputFishIdList {
			fishId, _ := strconv.Atoi(inputFishStr)
			selectFishList = append(selectFishList, int32(fishId))
		}
	}

	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	lastRtpState := ""
	var netWinTimes = make(map[int32]map[int]map[int]map[string]uint64, 0)
	var fishTimes = make(map[int32]map[int]map[string]uint64, 0)
	var hitTimes = make(map[int32]map[int]map[string]uint64, 0)
	var drillTimes = make(map[int32]uint64, 0)
	var mercenaryTimes = make(map[int32]uint64, 0)
	var drillHits = make(map[int32]uint64, 0)
	var mercenaryHits = make(map[int32]uint64, 0)
	var freeBullet = make(map[int32]map[int]map[string]uint64, 0)
	var freeBulletDrill = make(map[int32]uint64, 0)
	var mercenaryFreeBullet = make(map[int32]uint64, 0)
	var freeBulletTimes = make(map[int32]uint64, 0)
	var wins = make(map[int32]map[int]uint64, 0)
	var drillWins = make(map[int32]uint64, 0)
	var mercenaryWins = make(map[int32]uint64, 0)

	var bullets uint64 = 0

	for i := 1; i <= runTimes*10000; i++ {
		fishId := rngFishIdNew(selectFishList)

		_, result, rtpState, rtpId, netWinGroup, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId00003, mathModuleId00003, "", "", "", "",
			-1, fishId, -1,
			1, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0, 0,
		)

		hit, win, _, resultBullet, _, _, drillBullet, drillFreeBulletTimes, drillWin, _, _, _, _, drillTimesMap, drillHitsMap, drillWinMap :=
			resultProcess(gameId00003, mathModuleId00003, fishId, result, 1, false)
		bullets = bullets + resultBullet + drillBullet

		mercenaryWin, leftBullet, _, _, mercenaryTime, mercenaryHit, mercenaryBullet, mercenaryDrillBullet,
			mercenaryfbt, mercenaryDrillfbt, mercenaryDrillWin, //mercenaryDrillWin
			bulletDrillTimesMap, bulletDrillHitsMap, bulletDrillWinMap := mercenaryBulletProcessNew(gameId00003, mathModuleId00003, bullets, selectFishList, 1)
		bullets = leftBullet

		// Process 水位轉換次數
		netWinTimes = verifySecondMap(netWinTimes, fishId, rtpState, netWinGroup)
		key := fmt.Sprint(rtpState, netWinGroup, rtpId)
		if lastRtpState == "" || lastRtpState != key {
			netWinTimes[fishId][rtpState][netWinGroup][rtpId]++
			lastRtpState = key
		}

		// Process 魚種次數
		fishTimes = verifyMap(fishTimes, fishId, rtpState)
		fishTimes[fishId][rtpState][rtpId]++

		// Process 擊中次數
		hitTimes = verifyMap(hitTimes, fishId, rtpState)
		hitTimes[fishId][rtpState][rtpId] += hit

		// Process 鑽頭砲
		//drillTimes[fishId] = drillTimes[fishId] + drillTime + bulletDrillTime
		//drillHits[fishId] = drillHits[fishId] + drillHit + bulletDrillHit
		drillTimes = insertMapData(drillTimesMap, drillTimes)
		drillTimes = insertMapData(bulletDrillTimesMap, drillTimes)

		drillHits = insertMapData(drillHitsMap, drillHits)
		drillHits = insertMapData(bulletDrillHitsMap, drillHits)

		// Process 傭兵
		mercenaryTimes[fishId] += mercenaryTime
		mercenaryHits[fishId] += mercenaryHit

		// Process FreeBullet Times
		freeBullet = verifyMap(freeBullet, fishId, rtpState)
		freeBullet[fishId][rtpState][rtpId] += resultBullet

		// Process Drill's FreeBullet Times
		freeBulletDrill[fishId] = freeBulletDrill[fishId] + drillBullet + mercenaryDrillBullet

		// Process Mercenary's FreeBullet Times
		mercenaryFreeBullet[fishId] += mercenaryBullet

		// Process FreeBullet Times ( Normal, Drill, Mercenary )
		if resultBullet > 0 {
			freeBulletTimes[fishId]++
		}
		freeBulletTimes[fishId] += drillFreeBulletTimes + mercenaryDrillfbt
		freeBulletTimes[fishId] += mercenaryfbt

		// Process Wins
		wins = verifySamipleMap(wins, fishId)
		wins[fishId][rtpState] = wins[fishId][rtpState] + (win - drillWin)

		// Process Drill Win
		//drillWins[fishId] = drillWins[fishId] + drillWin + mercenaryDrillWin
		drillWins = insertMapData(drillWinMap, drillWins)
		drillWins = insertMapData(bulletDrillWinMap, drillWins)

		// Process Mercenary Win
		mercenaryWins[fishId] += mercenaryWin - mercenaryDrillWin
	}

	var printData = make(map[int32]string, 0)
	showString := ""
	for _, fishId := range fishList[gameId00003] {
		showString = fmt.Sprint(fishId)
		// Rtp State Times
		for rtpState := 0; rtpState < 2; rtpState++ {
			for _, netWinGroup := range netWinGroupList {
				for _, rtpId := range rtpList00003[rtpState] {
					showString = fmt.Sprint(showString, " ", netWinTimes[fishId][rtpState][netWinGroup][rtpId])
				}
			}
		}

		// Times
		for rtpState := 1; rtpState >= 0; rtpState-- {
			for _, rtpId := range rtpList00003[rtpState] {
				showString = fmt.Sprint(showString, " ", fishTimes[fishId][rtpState][rtpId])
			}
		}

		// Hit Times
		for rtpState := 1; rtpState >= 0; rtpState-- {
			for _, rtpId := range rtpList00003[rtpState] {
				showString = fmt.Sprint(showString, " ", hitTimes[fishId][rtpState][rtpId])
			}
		}

		// 鑽頭砲次數
		showString = fmt.Sprint(showString, " ", drillTimes[fishId])
		// 傭兵次數
		showString = fmt.Sprint(showString, " ", mercenaryTimes[fishId])
		// 鑽頭砲Hit
		showString = fmt.Sprint(showString, " ", drillHits[fishId])
		// 傭兵Hit
		showString = fmt.Sprint(showString, " ", mercenaryHits[fishId])

		// 免費子彈獲得數量
		for rtpState := 1; rtpState >= 0; rtpState-- {
			for _, rtpId := range rtpList00003[rtpState] {
				showString = fmt.Sprint(showString, " ", freeBullet[fishId][rtpState][rtpId])
			}
		}

		// 鑽頭砲獲得免費子彈數量
		showString = fmt.Sprint(showString, " ", freeBulletDrill[fishId])
		// 傭兵獲得免費子彈數量
		showString = fmt.Sprint(showString, " ", mercenaryFreeBullet[fishId])
		// 免費子彈次數
		showString = fmt.Sprint(showString, " ", freeBulletTimes[fishId])

		// Wins
		for rtpState := 1; rtpState >= 0; rtpState-- {
			showString = fmt.Sprint(showString, " ", wins[fishId][rtpState])
		}

		// Drill Win
		showString = fmt.Sprint(showString, " ", drillWins[fishId])
		// Mercenary Win
		showString = fmt.Sprint(showString, " ", mercenaryWins[fishId])

		printData[fishId] += showString
	}

	for _, fish := range fishList[gameId00003] {
		fmt.Println(printData[fish])
	}
}
