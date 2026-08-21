package rtp

import (
	"fmt"
	"serve/fish_comm/rng"
	"serve/service_fish/domain/probability"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00005_96_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-96-1"
	PSFM_00005_97_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-97-1"
	PSFM_00005_98_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-98-1"
	"serve/service_fish/models"
	"strconv"
	"testing"
)

const (
	gameId = models.PSF_ON_00004
	//mathModuleId            = models.PSFM_00005_95_1
	secWebSocketKey_00004 = "jerry"
	//run_times               = 100 * 100 * 100 * 10
	//run_times               = 100 * 100 *100 * 1
	//run_times               = 5000
	max_machineGun_bullet = 999
	bet                   = 1
	rate                  = 100
	min_bet               = 10
)

var mathModuleId string

func TestService_00004_Random(t *testing.T) {
	// Clam
	var fish_0 uint64 = 0
	var hit_0 uint64 = 0
	var win_0 uint64 = 0

	// Shrimp
	var fish_1 uint64 = 0
	var hit_1 uint64 = 0
	var win_1 uint64 = 0

	// Surmullet
	var fish_2 uint64 = 0
	var hit_2 uint64 = 0
	var win_2 uint64 = 0

	// Squid
	var fish_3 uint64 = 0
	var hit_3 uint64 = 0
	var win_3 uint64 = 0

	// Flying Fish
	var fish_4 uint64 = 0
	var hit_4 uint64 = 0
	var win_4 uint64 = 0

	// Halibut
	var fish_5 uint64 = 0
	var hit_5 uint64 = 0
	var win_5 uint64 = 0

	// Butterfly Fish
	var fish_6 uint64 = 0
	var hit_6 uint64 = 0
	var win_6 uint64 = 0

	// Oplegnathus
	var fish_7 uint64 = 0
	var hit_7 uint64 = 0
	var win_7 uint64 = 0

	// Snapper
	var fish_8 uint64 = 0
	var hit_8 uint64 = 0
	var win_8 uint64 = 0

	// Mahi Mahi
	var fish_9 uint64 = 0
	var hit_9 uint64 = 0
	var win_9 uint64 = 0

	// Stingray
	var fish_10 uint64 = 0
	var hit_10 uint64 = 0
	var win_10 uint64 = 0

	// Lobster Waiter
	var fish_11 uint64 = 0
	var hit_11 uint64 = 0
	var win_11 uint64 = 0

	// Penguin Waiter
	var fish_12 uint64 = 0
	var hit_12 uint64 = 0
	var win_12 uint64 = 0

	// Platypus Senior Chef
	var fish_13 uint64 = 0
	var hit_13 uint64 = 0
	var win_13 uint64 = 0

	// Sea Lion Chef
	var fish_14 uint64 = 0
	var hit_14 uint64 = 0
	var win_14 uint64 = 0

	// HairtCrab
	var fish_15 uint64 = 0
	var hit_15 uint64 = 0
	var win_15 uint64 = 0

	// Snaggletooth Shark
	var fish_16 uint64 = 0
	var hit_16 uint64 = 0
	var win_16 uint64 = 0

	// Sword Fish
	var fish_17 uint64 = 0
	var hit_17 uint64 = 0
	var win_17 uint64 = 0

	// Whale Shark
	var fish_18 uint64 = 0
	var hit_18 uint64 = 0
	var win_18 uint64 = 0

	// GiantOarfish
	var fish_19 uint64 = 0
	var hit_19 uint64 = 0
	var win_19 uint64 = 0

	// Red Envelope
	var fish_100 uint64 = 0
	var hit_100 uint64 = 0
	var win_100 uint64 = 0

	// Slot
	var fish_101 uint64 = 0
	var hit_101 uint64 = 0
	var win_101 uint64 = 0

	// MachineGun
	var fish_201 uint64 = 0
	var hit_201 uint64 = 0
	var win_201 uint64 = 0

	// SuperMachineGun
	var fish_202 uint64 = 0
	var hit_202 uint64 = 0
	var win_202 uint64 = 0

	// Lobster Dash
	var fish_300 uint64 = 0
	var hit_300 uint64 = 0
	var win_300 uint64 = 0

	// Fruit Dash
	var fish_301 uint64 = 0
	var hit_301 uint64 = 0
	var win_301 uint64 = 0

	// XiaoLongBao
	var fish_302 uint64 = 0
	var hit_302 uint64 = 0
	var win_302 uint64 = 0

	// White Tiger Chef
	var fish_400 uint64 = 0
	var hit_400 uint64 = 0
	var win_400 uint64 = 0

	fmt.Print("MathModuleId(95、96、97、98), RunTimes(x 萬), FishId -> ")
	var rtpId, run_times int
	var inputFishId string
	fmt.Scanf("%d, %d, %s", &rtpId, &run_times, &inputFishId)

	switch rtpId {
	case 95:
		mathModuleId = models.PSFM_00005_95_1
	case 96:
		mathModuleId = models.PSFM_00005_96_1
	case 97:
		mathModuleId = models.PSFM_00005_97_1
	case 98:
		mathModuleId = models.PSFM_00005_98_1
	}
	run_times = run_times * 10000

	for i := 0; i < run_times; i++ {
		var fishId int32
		if inputFishId != "" {
			tempFish, _ := strconv.Atoi(inputFishId)
			fishId = int32(tempFish)
		} else {
			fishId = rngFishId_00004()
		}

		Service.DecreaseDenominator(gameId, subgameId, mathModuleId, secWebSocketKey_00004, fishId, "", bet*rate, min_bet)
		rtpId := Service.RtpId(secWebSocketKey_00004, gameId, subgameId)

		budget := Service.RtpBudget(secWebSocketKey_00004, gameId, subgameId, mathModuleId)

		result := probability.Service.Calc(
			gameId,
			mathModuleId,
			rtpId,
			fishId,
			-1,
			budget,
			bet*rate,
		)

		Service.DecreaseMolecular(gameId, subgameId, mathModuleId, secWebSocketKey_00004,
			fishId, result.Pay*result.Multiplier, bet*rate, result.TriggerIconId)

		switch fishId {
		case 0:
			if result.Pay > 0 {
				fish_0 += uint64(result.Pay * result.Multiplier)
				win_0++
			}
			hit_0++

		case 1:
			if result.Pay > 0 {
				fish_1 += uint64(result.Pay * result.Multiplier)
				win_1++
			}
			hit_1++

		case 2:
			if result.Pay > 0 {
				fish_2 += uint64(result.Pay * result.Multiplier)
				win_2++
			}
			hit_2++

		case 3:
			if result.Pay > 0 {
				fish_3 += uint64(result.Pay * result.Multiplier)
				win_3++
			}
			hit_3++

		case 4:
			if result.Pay > 0 {
				fish_4 += uint64(result.Pay * result.Multiplier)
				win_4++
			}
			hit_4++

		case 5:
			if result.Pay > 0 {
				fish_5 += uint64(result.Pay * result.Multiplier)
				win_5++
			}
			hit_5++

		case 6:
			if result.Pay > 0 {
				fish_6 += uint64(result.Pay * result.Multiplier)
				win_6++
			}
			hit_6++

		case 7:
			if result.Pay > 0 {
				fish_7 += uint64(result.Pay * result.Multiplier)
				win_7++
			}
			hit_7++

		case 8:
			if result.Pay > 0 {
				fish_8 += uint64(result.Pay * result.Multiplier)
				win_8++
			}
			hit_8++

		case 9:
			if result.Pay > 0 {
				fish_9 += uint64(result.Pay * result.Multiplier)
				win_9++
			}
			hit_9++

		case 10:
			if result.Pay > 0 {
				fish_10 += uint64(result.Pay * result.Multiplier)
				win_10++
			}
			hit_10++

		case 11:
			if result.Pay > 0 {
				fish_11 += uint64(result.Pay * result.Multiplier)
				win_11++
			}
			hit_11++

		case 12:
			if result.Pay > 0 {
				fish_12 += uint64(result.Pay * result.Multiplier)
				win_12++
			}
			hit_12++

		case 13:
			if result.Pay > 0 {
				fish_13 += uint64(result.Pay * result.Multiplier)
				win_13++
			}
			hit_13++

		case 14:
			if result.Pay > 0 {
				fish_14 += uint64(result.Pay * result.Multiplier)
				win_14++
			}
			hit_14++

		case 15:
			if result.Pay > 0 {
				fish_15 += uint64(result.Pay * result.Multiplier)
				win_15++
			}
			hit_15++

		case 16:
			if result.Pay > 0 {
				fish_16 += uint64(result.Pay * result.Multiplier)
				win_16++
			}
			hit_16++

		case 17:
			if result.Pay > 0 {
				fish_17 += uint64(result.Pay * result.Multiplier)
				win_17++
			}
			hit_17++

		case 18:
			if result.Pay > 0 {
				fish_18 += uint64(result.Pay * result.Multiplier)
				win_18++
			}
			hit_18++

		case 19:
			if result.Pay > 0 {
				fish_19 += uint64(result.Pay * result.Multiplier)
				win_19++
			}
			hit_19++

		case 100:
			if result.TriggerIconId == 100 {
				fish_100 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				win_100++
			}
			hit_100++

		case 101:
			if result.TriggerIconId == 101 {
				fish_101 += uint64(result.ExtraData[0].(int))
				win_101++
			}
			hit_101++

		case 201:
			if result.Pay > 0 {
				fish_201 += uint64(result.Pay)
				switch mathModuleId {
				case models.PSFM_00005_95_1:
					fish_201 += shot_201(PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath.Icons.MachineGun.RTP2.ID, result.BonusPayload.(int))
				case models.PSFM_00005_96_1:
					fish_201 += shot_201(PSFM_00005_96_1.RTP96BS.PSF_ON_00004_1_BsMath.Icons.MachineGun.RTP2.ID, result.BonusPayload.(int))
				case models.PSFM_00005_97_1:
					fish_201 += shot_201(PSFM_00005_97_1.RTP97BS.PSF_ON_00004_1_BsMath.Icons.MachineGun.RTP2.ID, result.BonusPayload.(int))
				case models.PSFM_00005_98_1:
					fish_201 += shot_201(PSFM_00005_98_1.RTP98BS.PSF_ON_00004_1_BsMath.Icons.MachineGun.RTP2.ID, result.BonusPayload.(int))
				}
				win_201++
			}
			hit_201++

		case 202:
			if result.Pay > 0 {
				fish_202 += uint64(result.Pay)
				switch mathModuleId {
				case models.PSFM_00005_95_1:
					fish_202 += shot_202(PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath.Icons.SuperMachineGun.RTP4.ID, result.BonusPayload.(int))
				case models.PSFM_00005_96_1:
					fish_202 += shot_202(PSFM_00005_96_1.RTP96BS.PSF_ON_00004_1_BsMath.Icons.SuperMachineGun.RTP4.ID, result.BonusPayload.(int))
				case models.PSFM_00005_97_1:
					fish_202 += shot_202(PSFM_00005_97_1.RTP97BS.PSF_ON_00004_1_BsMath.Icons.SuperMachineGun.RTP4.ID, result.BonusPayload.(int))
				case models.PSFM_00005_98_1:
					fish_202 += shot_202(PSFM_00005_98_1.RTP98BS.PSF_ON_00004_1_BsMath.Icons.SuperMachineGun.RTP4.ID, result.BonusPayload.(int))
				}
				win_202++
			}
			hit_202++

		case 300:
			if result.Pay > 0 {
				fish_300 += uint64(result.Pay * result.Multiplier)
				win_300++
			}
			hit_300++

		case 301:
			if result.Pay > 0 {
				fish_301 += uint64(result.Pay * result.Multiplier)
				win_301++
			}
			hit_301++

		case 302:
			if result.TriggerIconId == 302 {
				fish_302 += uint64(result.Pay * result.Multiplier)
				win_302++
			}
			hit_302++

		case 400:
			if result.Pay > 0 {
				fish_400 += uint64(result.Pay * result.Multiplier)
				win_400++
			}
			hit_400++

		default:
			t.Fatal("Fish Id Error")
		}
	}

	totalHit := hit_0 + hit_1 + hit_2 + hit_3 + hit_4 + hit_5 + hit_6 + hit_7 + hit_8 + hit_9 +
		hit_10 + hit_11 + hit_12 + hit_13 + hit_14 + hit_15 + hit_16 + hit_17 + hit_18 + hit_19 +
		hit_100 + hit_101 + hit_201 + hit_202 + hit_300 + hit_301 + hit_302 + hit_400

	totalWin := win_0 + win_1 + win_2 + win_3 + win_4 + win_5 + win_6 + win_7 + win_8 + win_9 +
		win_10 + win_11 + win_12 + win_13 + win_14 + win_15 + win_16 + win_17 + win_18 + win_19 +
		win_100 + win_101 + win_201 + win_202 + win_300 + win_301 + win_302 + win_400

	totalFish := fish_0 + fish_1 + fish_2 + fish_3 + fish_4 + fish_5 + fish_6 + fish_7 + fish_8 + fish_9 +
		fish_10 + fish_11 + fish_12 + fish_13 + fish_14 + fish_15 + fish_16 + fish_17 + fish_18 + fish_19 +
		fish_100 + fish_101 + fish_201 + fish_202 + fish_300 + fish_301 + fish_302 + fish_400

	fmt.Println("FishId", "Hit", "Win", "Pay", "Win Rate", "Pay Rate")
	fmt.Println("000", hit_0, win_0, fish_0, getRate(win_0, hit_0), getRate(fish_0, hit_0))
	fmt.Println("001", hit_1, win_1, fish_1, getRate(win_1, hit_1), getRate(fish_1, hit_1))
	fmt.Println("002", hit_2, win_2, fish_2, getRate(win_2, hit_2), getRate(fish_2, hit_2))
	fmt.Println("003", hit_3, win_3, fish_3, getRate(win_3, hit_3), getRate(fish_3, hit_3))
	fmt.Println("004", hit_4, win_4, fish_4, getRate(win_4, hit_4), getRate(fish_4, hit_4))
	fmt.Println("005", hit_5, win_5, fish_5, getRate(win_5, hit_5), getRate(fish_5, hit_5))
	fmt.Println("006", hit_6, win_6, fish_6, getRate(win_6, hit_6), getRate(fish_6, hit_6))
	fmt.Println("007", hit_7, win_7, fish_7, getRate(win_7, hit_7), getRate(fish_7, hit_7))
	fmt.Println("008", hit_8, win_8, fish_8, getRate(win_8, hit_8), getRate(fish_8, hit_8))
	fmt.Println("009", hit_9, win_9, fish_9, getRate(win_9, hit_9), getRate(fish_9, hit_9))

	fmt.Println("010", hit_10, win_10, fish_10, getRate(win_10, hit_10), getRate(fish_10, hit_10))
	fmt.Println("011", hit_11, win_11, fish_11, getRate(win_11, hit_11), getRate(fish_11, hit_11))
	fmt.Println("012", hit_12, win_12, fish_12, getRate(win_12, hit_12), getRate(fish_12, hit_12))
	fmt.Println("013", hit_13, win_13, fish_13, getRate(win_13, hit_13), getRate(fish_13, hit_13))
	fmt.Println("014", hit_14, win_14, fish_14, getRate(win_14, hit_14), getRate(fish_14, hit_14))
	fmt.Println("015", hit_15, win_15, fish_15, getRate(win_15, hit_15), getRate(fish_15, hit_15))
	fmt.Println("016", hit_16, win_16, fish_16, getRate(win_16, hit_16), getRate(fish_16, hit_16))
	fmt.Println("017", hit_17, win_17, fish_17, getRate(win_17, hit_17), getRate(fish_17, hit_17))
	fmt.Println("018", hit_18, win_18, fish_18, getRate(win_18, hit_18), getRate(fish_18, hit_18))
	fmt.Println("019", hit_19, win_19, fish_19, getRate(win_19, hit_19), getRate(fish_19, hit_19))

	fmt.Println("100", hit_100, win_100, fish_100, getRate(win_100, hit_100), getRate(fish_100, hit_100))
	fmt.Println("101", hit_101, win_101, fish_101, getRate(win_101, hit_101), getRate(fish_101, hit_101))
	fmt.Println("201", hit_201, win_201, fish_201, getRate(win_201, hit_201), getRate(fish_201, hit_201))
	fmt.Println("202", hit_202, win_202, fish_202, getRate(win_202, hit_202), getRate(fish_202, hit_202))
	fmt.Println("300", hit_300, win_300, fish_300, getRate(win_300, hit_300), getRate(fish_300, hit_300))
	fmt.Println("301", hit_301, win_301, fish_301, getRate(win_301, hit_301), getRate(fish_301, hit_301))
	fmt.Println("302", hit_302, win_302, fish_302, getRate(win_302, hit_302), getRate(fish_302, hit_302))
	fmt.Println("400", hit_400, win_400, fish_400, getRate(win_400, hit_400), getRate(fish_400, hit_400))
	fmt.Println("Total", totalHit, totalWin, totalFish, getRate(totalWin, totalHit), getRate(totalFish, totalHit))
}

func rngFishId_00004() int32 {
	options := make([]rng.Option, 0, 27)

	for i := 0; i < 20; i++ {
		options = append(options, rng.Option{Weight: 1, Item: i})
	}
	options = append(options, rng.Option{Weight: 1, Item: 100})
	options = append(options, rng.Option{Weight: 1, Item: 101})
	options = append(options, rng.Option{Weight: 1, Item: 201})
	options = append(options, rng.Option{Weight: 1, Item: 202})
	options = append(options, rng.Option{Weight: 1, Item: 300})
	options = append(options, rng.Option{Weight: 1, Item: 301})
	options = append(options, rng.Option{Weight: 1, Item: 302})
	options = append(options, rng.Option{Weight: 1, Item: 400})

	return int32(rng.New(options).Item.(int))
}

// 機關槍
func shot_201(rtpId string, bullets int) (pay uint64) {
	var fishId int32 = -1
	var newBullets = 0

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId_00004()
			if fishId != 202 {
				break
			}
		}

		result := probability.Service.Calc(
			gameId,
			mathModuleId,
			rtpId,
			fishId,
			201,
			0,
			0,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay) * uint64(result.Multiplier)
		}

		// 紅包
		if result.TriggerIconId == 100 {
			pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		}

		// 老虎機
		if result.TriggerIconId == 101 {
			pay += uint64(result.ExtraData[0].(int))
		}

		// 補上子彈數
		if result.TriggerIconId == 201 {
			newBullets += result.BonusPayload.(int)

			if newBullets > max_machineGun_bullet {
				newBullets = max_machineGun_bullet
			}
		}
	}

	if newBullets > 0 {
		pay += shot_201(rtpId, newBullets)
	}

	return pay
}

// 超級機關槍
func shot_202(rtpId string, bullets int) (pay uint64) {
	var fishId int32 = -1
	var newBullets = 0

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId_00004()
			if fishId != 201 {
				break
			}
		}

		result := probability.Service.Calc(
			gameId,
			mathModuleId,
			rtpId,
			fishId,
			202,
			0,
			0,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay) * uint64(result.Multiplier)
		}

		// 紅包
		if result.TriggerIconId == 100 {
			pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		}

		// 老虎機
		if result.TriggerIconId == 101 {
			pay += uint64(result.ExtraData[0].(int))
		}

		// 補上子彈數
		if result.TriggerIconId == 202 {
			newBullets += result.BonusPayload.(int)

			if newBullets > max_machineGun_bullet {
				newBullets = max_machineGun_bullet
			}
		}
	}

	if newBullets > 0 {
		pay += shot_202(rtpId, newBullets)
	}

	return pay
}
