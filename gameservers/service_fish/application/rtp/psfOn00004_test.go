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

const gameId00004 = models.PSF_ON_00004

var mathModuleId00004 string

// 分母 = 0, 分子 = 1
var rtpList00004 = map[int][]int{
	0: {20, 30, 40, 50, 60},
	1: {2, 4, 5, 6, 8, 10, 15, 20, 25, 30},
}
var randomBet00004 = []int{1, 2, 5}

func TestService00004(t *testing.T) {
	// Process Input Value
	if len(os.Args) < 5 || len(os.Args) >= 7 {
		fmt.Println("參數說明:數學代號 總次數 指定魚種 隨機下注金額(true or false) 武器魚種驗證(true or false)")
		return
	}

	// 數學代號
	mathModuleId00004 = chooseMath(gameId00004, os.Args[1])
	if mathModuleId00004 == "" {
		fmt.Println("數學代號不存在")
		return
	}
	if mathModuleId00004 == " " {
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
		selectFishList = fishList[gameId00004]
	} else {
		inputFishIdList := strings.Split(os.Args[3], ",")
		for _, inputFishId := range inputFishIdList {
			if !isContain(fishList[gameId00004], inputFishId) {
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

	var totalTimes = make(map[int32]uint64, 0)
	var totalHits = make(map[int32]uint64, 0)
	var totalWin = make(map[int32]uint64, 0)
	var machineGunTimes = make(map[int32]uint64, 0)
	var machineGunHits = make(map[int32]uint64, 0)
	var superMachineGunTimes = make(map[int32]uint64, 0)
	var superMachineGunHits = make(map[int32]uint64, 0)
	var denominatorTimes = make(map[int]uint64, 0) // 分母
	var molecularTimes = make(map[int]uint64, 0)   // 分子
	var betTimes = make(map[int32]map[int]uint64, 0)

	var totalMolecular uint64 = 0

	for i := 1; i <= runTimes*10000; i++ {
		// Process Random Bet
		var bet = defaultBet
		if checkRandomBet {
			bet = randomBet00004[(i-1)%3]
		}

		fishId := rngFishIdNew(selectFishList)

		_, result, _, _, _, _, denominator, molecular, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId00004, mathModuleId00004, "", "", "", "",
			-1, fishId, -1,
			uint64(bet), defaultRate, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0, 0,
		)

		hit, win, _, _, _, _, _, _, _,
			machineGunTime, machineGunHit,
			superMachineGunTime, superMachineGunHit, _, _, _ := resultProcess(gameId00004, mathModuleId00004, fishId, result, uint64(bet), weaponFish)

		totalTimes[fishId]++
		totalHits[fishId] += hit
		totalWin[fishId] += win

		// MachineGun
		for k, v := range machineGunTime {
			machineGunTimes[k] += v
		}
		for k, v := range machineGunHit {
			machineGunHits[k] += v
		}

		// SuperMachineGun
		for k, v := range superMachineGunTime {
			superMachineGunTimes[k] += v
		}
		for k, v := range superMachineGunHit {
			superMachineGunHits[k] += v
		}

		// Bet Change Times
		betTimes = verifySamipleMap(betTimes, fishId)
		betTimes[fishId][bet]++

		if denominator != 0 && molecular != 0 {
			// Denominator
			denominatorTimes[denominator]++
			// Molecular
			molecularTimes[molecular]++
		}

		// Total Molecular (Denominator * Molecular * Bet)
		totalMolecular += uint64(denominator * molecular * bet)
	}

	var printData = make(map[int32]string, 0)
	showString := ""

	for _, oneFish := range fishList[gameId00004] {
		showString = fmt.Sprint(oneFish)

		showString = fmt.Sprint(showString, " ", totalTimes[oneFish])
		showString = fmt.Sprint(showString, " ", totalHits[oneFish])
		showString = fmt.Sprint(showString, " ", totalWin[oneFish])

		if checkRandomBet {
			for _, index := range randomBet00004 {
				showString = fmt.Sprint(showString, " ", betTimes[oneFish][index])
			}
		}

		showString = fmt.Sprint(showString, " ", machineGunTimes[oneFish])
		showString = fmt.Sprint(showString, " ", superMachineGunTimes[oneFish])
		showString = fmt.Sprint(showString, " ", machineGunHits[oneFish])
		showString = fmt.Sprint(showString, " ", superMachineGunHits[oneFish])

		printData[oneFish] += showString
	}

	for _, oneFish := range fishList[gameId00004] {
		fmt.Println(printData[oneFish])
	}

	showString = ""
	for _, denominatorIndex := range rtpList00004[0] {
		showString = fmt.Sprint(showString, denominatorTimes[denominatorIndex], " ")
	}

	for _, molecularIndex := range rtpList00004[1] {
		showString = fmt.Sprint(showString, molecularTimes[molecularIndex], " ")
	}
	showString = fmt.Sprint(showString, totalMolecular, " ")

	fmt.Println(showString)
	//fmt.Println("-----------")
	//fmt.Println(machineGunTimes)
	//fmt.Println("-----------")
	//fmt.Println(machineGunHits)
}
