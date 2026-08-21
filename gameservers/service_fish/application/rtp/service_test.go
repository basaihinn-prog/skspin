package rtp

import (
	"fmt"
	"serve/service_fish/application/lottery"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	PSFM_00007_97_1 "serve/service_fish/domain/probability/PSFM-00007-1/PSFM-00007-97-1"
	PSFM_00008_97_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-97-1"
	PSFM_00013_93_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-93-1"
	"serve/service_fish/models"
	"strconv"

	"serve/fish_comm/rng"
)

const (
	high = 1
	low  = 0

	//TIMES           = 100 * 100 * 100 * 50
	//GAME_ID         = models.PSF_ON_00001
	//MATH_MODULE_ID  = models.PSFM_00002_97_1
	BET             = 1
	secWebSocketKey = "Levi"

	subgameId = 0
)

func shot_26(gameId, mathModuleId string, rtpId string, bullets int) (pay uint64) {
	var fishId int32 = -1

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId()

			if fishId != 26 && fishId != 27 {
				break
			}
		}

		result := probability.Service.Calc(
			gameId,
			mathModuleId,
			rtpId,
			fishId,
			26,
			Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
			0,
		)

		if result.Pay > 0 {
			pay += (uint64(result.Pay) * uint64(result.Multiplier))
		}

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
					31,
					26,
					Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
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
					31,
					26,
					Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
					0,
				)

				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}
		}
	}
	return pay
}

func shot_27(gameId, mathModuleId string, rtpId string, bullets int) (pay uint64) {
	var fishId int32 = -1
	var newBullets = 0

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId()

			if fishId != 26 {
				break
			}
		}

		result := probability.Service.Calc(
			gameId,
			mathModuleId,
			rtpId,
			fishId,
			27,
			Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
			0,
		)

		if result.Pay > 0 {
			pay += (uint64(result.Pay) * uint64(result.Multiplier))
		}

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
					31,
					27,
					Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
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
					31,
					27,
					Service.RtpBudget(secWebSocketKey, gameId, -1, mathModuleId),
					0,
				)

				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}
		}
		if result.TriggerIconId == 27 {
			newBullets += result.BonusPayload.(int)

			if newBullets > 999 {
				newBullets = 999
			}
		}
	}

	if newBullets > 0 {
		pay += shot_27(gameId, mathModuleId, rtpId, newBullets)
	}

	return pay
}

func rngFishId() int32 {
	options := make([]rng.Option, 0, 29)

	for i := 0; i < 30; i++ {
		options = append(options, rng.Option{1, i})
	}

	return int32(rng.New(options).Item.(int))
}

func rngFishIdNew(fishList []int32) int32 {
	options := make([]rng.Option, 0)

	for _, v := range fishList {
		options = append(options, rng.Option{Weight: 1, Item: v})
	}

	return rng.New(options).Item.(int32)
}

func rngRedEnvelope() int32 {
	options := make([]rng.Option, 0, 5)

	for i := 0; i < 5; i++ {
		options = append(options, rng.Option{1, i})
	}

	return int32(rng.New(options).Item.(int))
}

func getRate(x, y uint64) string {
	result := float32(x) / float32(y) * 100
	return fmt.Sprintf("%.2f", result)
}

// --------------------------------------------------------------------------------------------------------------------------------
var fishList = map[string][]int32{
	models.PSF_ON_00001: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
	models.PSF_ON_00002: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
	models.PSF_ON_00003: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 303},
	models.PSF_ON_00004: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400},
	models.PSF_ON_00005: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 501, 502, 503, 504, 505},
	models.PSF_ON_00006: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400},
	models.PSF_ON_00007: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400},
	models.PSF_ON_20002: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
	models.RKF_H5_00001: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 501, 502, 503, 504, 505},
}
var triggerIconList = map[string]int{
	models.PSF_ON_00001: 31, models.PSF_ON_00002: 31, models.PSF_ON_00003: -1,
	models.PSF_ON_00004: -1, models.PSF_ON_00005: 102, models.PSF_ON_00006: 102, models.PSF_ON_00007: 102,
	models.PSF_ON_20002: 31, models.RKF_H5_00001: 102,
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
	models.PSF_ON_00006: {
		201: PSFM_00007_97_1.RTP97BS.PSF_ON_00006_1_BsMath.Icons.MachineGun.UseRTP,
		202: PSFM_00007_97_1.RTP97BS.PSF_ON_00006_1_BsMath.Icons.SuperMachineGun.UseRTP,
	},
	models.PSF_ON_00007: {
		201: PSFM_00008_97_1.RTP97BS.PSF_ON_00007_1_BsMath.Icons.MachineGun.UseRTP,
		202: PSFM_00008_97_1.RTP97BS.PSF_ON_00007_1_BsMath.Icons.SuperMachineGun.UseRTP,
	},
	models.PSF_ON_20002: {
		26: PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath.Icons.Drill.RTP8.ID,
		27: PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID,
	},
	models.RKF_H5_00001: {
		201: PSFM_00013_93_1.RTP93BS.RKF_H5_00001_1_BsMath.Icons.ThunderHammer.UseRTP,
		202: PSFM_00013_93_1.RTP93BS.RKF_H5_00001_1_BsMath.Icons.Stormbreaker.UseRTP,
	},
}

const (
	defaultBet  = 1
	defaultRate = 1
)

var randomBet = []int{1, 3, 5}

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

	case models.PSF_ON_00006:
		switch inputMath {
		case "95":
			return models.PSFM_00007_95_1
		case "96":
			return models.PSFM_00007_96_1
		case "97":
			return models.PSFM_00007_97_1
		case "98":
			return models.PSFM_00007_98_1
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

	case models.PSF_ON_20002:
		switch inputMath {
		case "93":
			return models.PSFM_00003_93_1
		case "94":
			return models.PSFM_00003_94_1
		default:
			return ""
		}

	case models.RKF_H5_00001:
		switch inputMath {
		case "93":
			return models.PSFM_00013_93_1
		case "94":
			return models.PSFM_00013_94_1
		case "95":
			return models.PSFM_00013_95_1
		case "96":
			return models.PSFM_00013_96_1
		case "97":
			return models.PSFM_00013_97_1
		case "98":
			return models.PSFM_00013_98_1
		default:
			return ""
		}

	default:
		return " "
	}
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

func resultProcess(gameId, mathModuleId string, fishId int32, result *probability.Probability, bet uint64, weaponFish bool) (
	hits, wins, envelopes, bullet,
	drillTimes, drillHits, drillBullet, drillFreeBulletTimes, drillWin uint64,
	machineGunTimes, machineGunHits, superMachineGunTimes, superMachineGunHits, drillTimesMap, drillHitsMap, drillWinsMap map[int32]uint64) {
	hits = 0
	wins = 0
	envelopes = 0
	bullet = 0
	drillTimes = 0
	drillHits = 0
	drillBullet = 0
	drillFreeBulletTimes = 0
	drillWin = 0

	machineGunTimes = make(map[int32]uint64, 0)
	machineGunHits = make(map[int32]uint64, 0)
	superMachineGunTimes = make(map[int32]uint64, 0)
	superMachineGunHits = make(map[int32]uint64, 0)
	drillTimesMap = make(map[int32]uint64, 0)
	drillHitsMap = make(map[int32]uint64, 0)
	drillWinsMap = make(map[int32]uint64, 0)

	if result.Pay > 0 {
		hits++
		wins += uint64(result.Pay*result.Multiplier) * bet
	}

	// Mercenary Bullet
	if result.Bullet > 0 {
		bullet = uint64(result.Bullet)
	}

	if gameId != models.PSF_ON_00004 && gameId != models.PSF_ON_00003 {
		if result.TriggerIconId == triggerIconList[gameId] {
			envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId])) * bet
		}
	}

	switch gameId {
	case models.PSF_ON_00001, models.PSF_ON_00002, models.PSF_ON_20002:
		// Drill
		if fishId == 26 && result.Pay > 0 {
			p, _, _, _, _, envelope, drillTimeMap, drillHitMap, _ := shotDrill(gameId, mathModuleId, BonusIconList[gameId][fishId], result.BonusPayload.(int), fishId, bet)
			wins += p
			envelopes += envelope
			for k, v := range drillTimeMap {
				drillTimesMap[k] += v
			}
			for k, v := range drillHitMap {
				drillHitsMap[k] += v
			}
		}

		// MachineGun
		if fishId == 27 && result.Pay > 0 {
			p, r, t, h := machineGunShot(gameId, mathModuleId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId, bet, weaponFish,
			)

			wins += p
			envelopes += r

			for k, v := range t {
				machineGunTimes[k] += v
			}
			for k, v := range h {
				machineGunHits[k] += v
			}
		}

		// Red Envelope
		if fishId == 28 && result.TriggerIconId == 28 {
			hits++
			wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
		}

		// Slot
		if fishId == 29 && result.TriggerIconId == 29 {
			hits++
			wins += uint64(result.ExtraData[0].(int)) * bet
		}

	default:
		// Red Envelope
		if fishId == 100 && result.TriggerIconId == 100 {
			hits++
			wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
		}

		// Slot
		if gameId == models.PSF_ON_00003 {
			if (fishId == 6 || fishId == 7 || fishId == 8 ||
				fishId == 9 || fishId == 10 || fishId == 11) && result.TriggerIconId == 101 {
				wins += uint64(result.ExtraData[0].(int)) * bet
			}
		} else {
			if fishId == 101 && result.TriggerIconId == 101 {
				hits++
				wins += uint64(result.ExtraData[0].(int)) * bet
			}
		}

		// MachineGun
		if fishId == 201 && result.Pay > 0 {
			p, r, t, h := machineGunShot(gameId, mathModuleId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId, bet, weaponFish,
			)

			wins += p
			envelopes += r

			for k, v := range t {
				machineGunTimes[k] += v
			}
			for k, v := range h {
				machineGunHits[k] += v
			}
		}

		// SuperMachineGun
		if fishId == 202 && result.Pay > 0 {
			p, r, t, h := machineGunShot(gameId, mathModuleId, BonusIconList[gameId][fishId],
				result.BonusPayload.(int), fishId, bet, weaponFish,
			)

			wins += p
			envelopes += r

			for k, v := range t {
				superMachineGunTimes[k] += v
			}
			for k, v := range h {
				superMachineGunHits[k] += v
			}
		}

		// Drill
		if gameId == models.PSF_ON_00003 {
			if fishId == 12 || fishId == 13 || fishId == 14 ||
				fishId == 15 || fishId == 16 || fishId == 17 {
				p, b, dt, dh, fbt, _, drillTimeMap, drillHitMap, drillWinMap :=
					shotDrill(gameId, mathModuleId, BonusIconList[gameId][fishId], result.BonusPayload.(int), fishId, bet)
				wins += p
				drillTimes = dt
				drillHits = dh
				drillBullet = b
				drillFreeBulletTimes = fbt
				drillWin = p

				drillTimesMap = insertMapData(drillTimeMap, drillTimesMap)
				drillHitsMap = insertMapData(drillHitMap, drillHitsMap)
				drillWinsMap = insertMapData(drillWinMap, drillWinsMap)
			}
		} else {
			// TODO
		}
	}

	// Fruit Dish
	if fishId == 300 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId])) * bet
	}

	// A pack of beer
	if fishId == 301 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId])) * bet
	}

	// XIAO LONG BAO
	if fishId == 302 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope(gameId, mathModuleId, result.BonusTypeId, int32(triggerIconList[gameId])) * bet
	}

	return hits, wins, envelopes, bullet, drillTimes, drillHits, drillBullet, drillFreeBulletTimes, drillWin,
		machineGunTimes, machineGunHits, superMachineGunTimes, superMachineGunHits, drillTimesMap, drillHitsMap, drillWinsMap
}

func mercenaryBulletProcessNew(gameId, mathModuleId string, bullets uint64, selectFishList []int32, bet uint64) (
	bulletWin, leftBullet, drillTimes, drillHits, mercenaryTimes, mercenaryHits, mercenaryBullet, mercenaryDrillBullet,
	mercenaryFreeBulletTimes, mercenaryDrillFreeBulletTimes, drillWin uint64,
	drillTimesMap, drillHitsMap, drillWinsMap map[int32]uint64) {
	bulletWin = 0
	leftBullet = 0
	drillTimes = 0
	drillHits = 0
	mercenaryTimes = 0
	mercenaryHits = 0
	mercenaryBullet = 0
	mercenaryDrillBullet = 0
	mercenaryFreeBulletTimes = 0
	mercenaryDrillFreeBulletTimes = 0
	drillWin = 0

	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}

	for {
		times := bullets / 60
		leftBullet = bullets - (times * 60)

		for i := 0; uint64(i) < (times * 60); i++ {
			mercenaryTimes++
			fishId := rngFishIdNew(selectFishList)

			_, result, _, _, _, _, _, _, _ := lottery.Service.MathProcess(
				secWebSocketKey, gameId, mathModuleId, "", "", "", BonusIconList[gameId][31],
				1, fishId, -1,
				bet, 1, 0,
				&hitBullet, &hitFish,
				lottery.PLAYER, "",
				0, 0,
				true, 0, 0,
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
					pay, drillBullet, _, dt, dh, fbt, drillTimeMap, drillHitMap, drillWinMap := shotDrill(gameId, mathModuleId, BonusIconList[gameId][fishId], result.BonusPayload.(int), fishId, bet)
					bulletWin += pay
					leftBullet += drillBullet
					drillTimes = dt
					drillHits = dh
					mercenaryDrillBullet = drillBullet
					mercenaryDrillFreeBulletTimes += fbt
					drillWin += pay
					drillWinsMap = insertMapData(drillWinMap, drillWinsMap)
					drillHitsMap = insertMapData(drillHitMap, drillHitsMap)
					drillTimesMap = insertMapData(drillTimeMap, drillTimesMap)
				}

				mercenaryHits++
			}

			if result.Bullet > 0 {
				leftBullet += uint64(result.Bullet)
				//mercenaryHits++
				mercenaryBullet += uint64(result.Bullet)
				mercenaryFreeBulletTimes++
			}
		}

		if leftBullet < 60 {
			break
		}
	}

	return bulletWin, leftBullet, drillTimes, drillHits, mercenaryTimes, mercenaryHits, mercenaryBullet, mercenaryDrillBullet,
		mercenaryFreeBulletTimes, mercenaryDrillFreeBulletTimes, drillWin, drillTimesMap, drillHitsMap, drillWinsMap
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
	)

	return uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
}

func machineGunShot(gameId, mathModuleId, rtpId string, bullets int, machineGunType int32, bet uint64, weaponFish bool) (pay, redEnvelope uint64, times, hits map[int32]uint64) {
	var fishId int32 = -1
	var newBullets = 0
	var checkMachineGunType int32 = 0

	times = make(map[int32]uint64, 0)
	hits = make(map[int32]uint64, 0)

	switch gameId {
	case models.PSF_ON_00001, models.PSF_ON_00002, models.PSF_ON_20002:
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
			if weaponFish {
				fishId = machineGunType
			} else {
				fishId = rngFishIdNew(fishList[gameId])
			}

			if fishId != checkMachineGunType {
				break
			}
		}

		_, result, _, _, _, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId, mathModuleId, "", strconv.Itoa(int(machineGunType)), "", BonusIconList[gameId][machineGunType],
			-1, fishId, machineGunType,
			bet, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0, 0,
		)

		times[fishId]++

		if result.Pay > 0 {
			pay += uint64(result.Pay*result.Multiplier) * bet
			hits[fishId]++
		}

		switch gameId {
		case models.PSF_ON_00001, models.PSF_ON_00002, models.PSF_ON_20002:
			if result.TriggerIconId == 28 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
				hits[fishId]++
			}

			if result.TriggerIconId == 29 {
				pay += uint64(result.ExtraData[0].(int)) * bet
				hits[fishId]++
			}
		default:
			if result.TriggerIconId == 100 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
				hits[fishId]++
			}

			if result.TriggerIconId == 101 {
				pay += uint64(result.ExtraData[0].(int)) * bet
				hits[fishId]++
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
				)
				redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
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
				)

				if gameId != models.PSF_ON_00004 {
					redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
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
		p, r, nt, nh := machineGunShot(gameId, mathModuleId, rtpId, newBullets, machineGunType, bet, weaponFish)
		pay += p
		redEnvelope += r

		for k, v := range nt {
			times[k] += v
		}
		for k, v := range nh {
			hits[k] += v
		}
	}

	return pay, redEnvelope, times, hits
}

func shotDrill(gameId, mathModuleId, rtpId string, bullets int, bulletTypeId int32, bet uint64) (pay, drillBullet, drillTimes, drillHits, freeBulletTimes, envelope uint64,
	drillTimesMap, drillHitsMap, drillWinsMap map[int32]uint64) {
	var fishId int32 = -1
	hitFish := fish.Fish{}
	hitBullet := bullet.Bullet{}
	uBullets, _ := strconv.ParseUint(strconv.Itoa(bullets), 10, 64)
	drillTimes = uBullets
	drillHits = 0
	freeBulletTimes = 0
	envelope = 0

	drillTimesMap = make(map[int32]uint64, 0)
	drillHitsMap = make(map[int32]uint64, 0)
	drillWinsMap = make(map[int32]uint64, 0)

	for i := 0; i < bullets; i++ {
		// Random Fish ID
		for {
			fishId = rngFishIdNew(fishList[gameId])

			if gameId == models.PSF_ON_00001 || gameId == models.PSF_ON_00002 || gameId == models.PSF_ON_20002 {
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

		_, result, _, _, _, _, _, _, _ := lottery.Service.MathProcess(
			secWebSocketKey, gameId, mathModuleId, "", strconv.Itoa(int(bulletTypeId)), "", rtpId,
			-1, fishId, bulletTypeId,
			bet, 1, 0,
			&hitBullet, &hitFish,
			lottery.PLAYER, "",
			0, 0,
			true, 0, 0,
		)

		drillTimesMap[fishId]++

		if result.Pay > 0 {
			pay += uint64(result.Pay*result.Multiplier) * bet
			drillWinsMap[fishId] += uint64(result.Pay*result.Multiplier) * bet

			switch gameId {
			case models.PSF_ON_00003:
				// Slot
				if fishId == 6 || fishId == 7 || fishId == 8 ||
					fishId == 9 || fishId == 10 || fishId == 11 {
					pay += uint64(result.ExtraData[0].(int))
					drillWinsMap[fishId] += uint64(result.ExtraData[0].(int))
				}
			}

			drillHits++
			drillHitsMap[fishId]++
		}

		if gameId == models.PSF_ON_00001 || gameId == models.PSF_ON_00002 || gameId == models.PSF_ON_20002 {
			if result.TriggerIconId == 28 {
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
				drillHits++
				drillHitsMap[fishId]++
			}

			if result.TriggerIconId == 29 {
				pay += uint64(result.ExtraData[0].(int)) * bet
				drillWinsMap[fishId] += uint64(result.ExtraData[0].(int)) * bet
				drillHits++
				drillHitsMap[fishId]++
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
					)

					envelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
					drillWinsMap[fishId] += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
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
					)

					envelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
					drillWinsMap[fishId] += uint64(result.BonusPayload.([]int)[rngRedEnvelope()]) * bet
				}
			}
		}

		if result.Bullet > 0 {
			drillBullet += uint64(result.Bullet)
			// TODO Check Get Free Bullet is Hit?
			freeBulletTimes++
		}
	}

	return pay, drillBullet, drillTimes, drillHits, freeBulletTimes, envelope, drillTimesMap, drillHitsMap, drillWinsMap
}

func verifySamipleMap(inputMap map[int32]map[int]uint64, fishId int32) map[int32]map[int]uint64 {
	if inputMap[fishId] == nil {
		inputMap[fishId] = make(map[int]uint64, 0)
	}

	return inputMap
}

func verifyMap(inputMap map[int32]map[int]map[string]uint64, fishId int32, rtpState int) map[int32]map[int]map[string]uint64 {
	if inputMap[fishId] == nil {
		inputMap[fishId] = make(map[int]map[string]uint64, 0)
	}
	if inputMap[fishId][rtpState] == nil {
		inputMap[fishId][rtpState] = make(map[string]uint64, 0)
	}

	return inputMap
}

func verifySecondMap(inputMap map[int32]map[int]map[int]map[string]uint64, fishId int32, rtpState int, netWinGroup int) map[int32]map[int]map[int]map[string]uint64 {
	if inputMap[fishId] == nil {
		inputMap[fishId] = make(map[int]map[int]map[string]uint64, 0)
	}

	if inputMap[fishId][rtpState] == nil {
		inputMap[fishId][rtpState] = make(map[int]map[string]uint64, 0)
	}

	if inputMap[fishId][rtpState][netWinGroup] == nil {
		inputMap[fishId][rtpState][netWinGroup] = make(map[string]uint64, 0)
	}

	return inputMap
}

func getShowData(showString string, fishId int32, resultMap map[int32]map[int]map[string]uint64, rtpList map[int][]string) string {
	for rtpState := 0; rtpState < 2; rtpState++ {
		for _, rtp := range rtpList[rtpState] {
			showString = fmt.Sprint(showString, " ", resultMap[fishId][rtpState][rtp])
		}
	}

	return showString
}

func insertMapData(insertMap map[int32]uint64, targetMap map[int32]uint64) map[int32]uint64 {
	for k, v := range insertMap {
		targetMap[k] += v
	}

	return targetMap
}
