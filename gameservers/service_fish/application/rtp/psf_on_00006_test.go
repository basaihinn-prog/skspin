package rtp

import (
	"fmt"
	"os"
	"serve/service_fish/domain/probability"
	PSFM_00007_97_1 "serve/service_fish/domain/probability/PSFM-00007-1/PSFM-00007-97-1"
	"serve/service_fish/models"
	"strconv"
	"testing"

	"serve/fish_comm/rng"
)

const (
	game_id_00006         = models.PSF_ON_00006
	secWebSocketKey_00006 = "jerry"
	bet_00006             = 1
	trigger_icon_id_00006 = 102
)

var fishList00006 = []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400}
var mathModule_id_00006 string

func TestService_00006(t *testing.T) {
	openTest := false
	if len(os.Args) > 1 {
		if os.Args[1] == "test" {
			openTest = true
		}
	}

	lowTimes := make(map[int32]uint64, 0)
	lowHits := make(map[int32]uint64, 0)
	lowWins := make(map[int32]uint64, 0)
	lowEnvelope := make(map[int32]uint64, 0)

	highTimes := make(map[int32]uint64, 0)
	highHits := make(map[int32]uint64, 0)
	highWins := make(map[int32]uint64, 0)
	highEnvelope := make(map[int32]uint64, 0)

	var rtpId, run_times int
	var inputFishId string
	fmt.Print("MathModuleId(95、96、97、98), RunTimes(x 萬), FishId -> ")
	fmt.Scanf("%d, %d, %s", &rtpId, &run_times, &inputFishId)

	switch rtpId {
	case 95:
		mathModule_id_00006 = models.PSFM_00007_95_1
	case 96:
		mathModule_id_00006 = models.PSFM_00007_96_1
	case 97:
		mathModule_id_00006 = models.PSFM_00007_97_1
	case 98:
		mathModule_id_00006 = models.PSFM_00007_98_1
	}

	run_times_00006 := run_times * 10000
	if openTest {
		run_times_00006 = run_times
	}

	for i := 0; i < run_times_00006; i++ {
		var fishId int32
		if inputFishId != "" {
			tempFish, _ := strconv.Atoi(inputFishId)
			fishId = int32(tempFish)
		} else {
			fishId = rngFishId_00006()
		}

		//hitFish := fish.Fish{}
		//hitBullet := bullet.Bullet{}
		//
		//_, result, rtpState := lottery.Service.MathProcess(
		//	secWebSocketKey_00006, game_id_00006, mathModule_id_00006, "", "", "", "",
		//	-1, fishId, -1,
		//	1, 1, subgameId,
		//	&hitBullet, &hitFish,
		//	lottery.PLAYER, "",
		//	0, 0,
		//	true,
		//)

		Service.Decrease(game_id_00006, 0, mathModule_id_00006, secWebSocketKey_00006, bet_00006)
		rtpState := Service.RtpState(secWebSocketKey_00006, game_id_00006, 0)
		rtpId := Service.RtpId(secWebSocketKey_00006, game_id_00006, 0)

		result := probability.Service.Calc(
			game_id_00006,
			mathModule_id_00006,
			rtpId,
			fishId,
			-1,
			0,
			0,
		)
		hit, win, envelope := processResult_2(fishId, result)

		switch rtpState {
		case low:
			lowTimes[fishId]++
			lowHits[fishId] += hit
			lowWins[fishId] += win
			lowEnvelope[fishId] += envelope

		case high:
			highTimes[fishId]++
			highHits[fishId] += hit
			highWins[fishId] += win
			highEnvelope[fishId] += envelope
		}
	}

	var lowTotalWin, lowTotalTimes, lowTotalEnvelope uint64
	var highTotalWin, highTotalTimes, highTotalEnvelope uint64

	fmt.Println("FishID", ":", "LowTimes", "LowHits", "LowWins", "LowEnvelope", "HighTimes", "HighHits", "HighWins", "HighEnvelope",
		"LowHitRate", "LowRtpRate", "HighHitRate", "HighRtpRate", "TotalRtpRate")
	for _, v := range fishList00006 {
		fmt.Println(v, ":", lowTimes[v], lowHits[v], lowWins[v], lowEnvelope[v], highTimes[v], highHits[v], highWins[v], highEnvelope[v],
			getRate(lowHits[v], lowTimes[v]), getRate(lowWins[v]+lowEnvelope[v], lowTimes[v]),
			getRate(highHits[v], highTimes[v]), getRate(highWins[v]+highEnvelope[v], highTimes[v]),
			getRate(lowWins[v]+lowEnvelope[v]+highWins[v]+highEnvelope[v], lowTimes[v]+highTimes[v]),
		)

		lowTotalTimes += lowTimes[v]
		lowTotalWin += lowWins[v]
		lowTotalEnvelope += lowEnvelope[v]

		highTotalTimes += highTimes[v]
		highTotalWin += highWins[v]
		highTotalEnvelope += highEnvelope[v]
	}
	fmt.Println("Total RTP", ":", getRate(lowTotalWin+lowTotalEnvelope+highTotalWin+highTotalEnvelope, lowTotalTimes+highTotalTimes))
}

func processResult_2(fishId int32, result *probability.Probability) (hits, wins, envelopes uint64) {
	hits = 0
	wins = 0
	envelopes = 0

	if result.Pay > 0 {
		hits++
		wins += uint64(result.Pay * result.Multiplier)
	}

	if result.TriggerIconId == trigger_icon_id_00006 {
		envelopes += getTriggerEnvelope00006(result.BonusTypeId, trigger_icon_id_00006)
	}

	// Red Envelope
	if fishId == 100 && result.TriggerIconId == 100 {
		hits++
		wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
	}

	// Slot
	if fishId == 101 && result.TriggerIconId == 101 {
		hits++
		wins += uint64(result.ExtraData[0].(int))
	}

	// MachineGun
	if fishId == 201 && result.Pay > 0 {
		p, r := machineGunShot00006(PSFM_00007_97_1.RTP97BS.PSF_ON_00006_1_BsMath.Icons.MachineGun.RTP6.ID, result.BonusPayload.(int), 201)
		wins += p
		envelopes += r
	}

	// SuperMachineGun
	if fishId == 202 && result.Pay > 0 {
		p, r := machineGunShot00006(PSFM_00007_97_1.RTP97BS.PSF_ON_00006_1_BsMath.Icons.SuperMachineGun.RTP8.ID, result.BonusPayload.(int), 202)
		wins += p
		envelopes += r
	}

	// Fruit Dish
	if fishId == 300 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope00006(result.BonusTypeId, trigger_icon_id_00006)
	}

	// A pack of beer
	if fishId == 301 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope00006(result.BonusTypeId, trigger_icon_id_00006)
	}

	// XIAO LONG BAO
	if fishId == 302 && result.ExtraTriggerBonus != nil {
		envelopes += getTriggerEnvelope00006(result.BonusTypeId, trigger_icon_id_00006)
	}

	return hits, wins, envelopes
}

func rngFishId_00006() int32 {
	options := make([]rng.Option, 0)

	for _, v := range fishList00006 {
		options = append(options, rng.Option{1, v})
	}

	return rng.New(options).Item.(int32)
}

func getTriggerEnvelope00006(bonusTypeId int, triggerIconId int32) uint64 {
	result := probability.Service.Calc(
		game_id_00006,
		mathModule_id_00006,
		strconv.Itoa(bonusTypeId),
		triggerIconId,
		-1,
		0,
		0,
	)

	return uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
}

func machineGunShot00006(rtpId string, bullets int, machineGunType int32) (pay, redEnvelope uint64) {
	var fishId int32 = -1
	var newBullets = 0
	var check_machineGun_type int32 = 0
	if machineGunType == 201 {
		check_machineGun_type = 202
	} else {
		check_machineGun_type = 201
	}

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId_00006()
			if fishId != check_machineGun_type {
				break
			}
		}

		result := probability.Service.Calc(
			game_id_00006,
			mathModule_id_00006,
			rtpId,
			fishId,
			machineGunType,
			0,
			0,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay * result.Multiplier)
		}

		if result.TriggerIconId == 100 {
			pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		}

		if result.TriggerIconId == 101 {
			pay += uint64(result.ExtraData[0].(int))
		}

		if result.ExtraTriggerBonus != nil {
			if result.ExtraTriggerBonus[0].TriggerIconId == 102 {
				result := probability.Service.Calc(
					game_id_00006,
					mathModule_id_00006,
					strconv.Itoa(result.ExtraTriggerBonus[0].BonusTypeId),
					102,
					machineGunType,
					0,
					0,
				)
				redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			}
		} else {
			if result.TriggerIconId == 102 {
				result := probability.Service.Calc(
					game_id_00006,
					mathModule_id_00006,
					strconv.Itoa(result.BonusTypeId),
					102,
					machineGunType,
					0,
					0,
				)
				redEnvelope += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
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
		p, r := machineGunShot00006(rtpId, newBullets, machineGunType)
		pay += p
		redEnvelope += r
	}

	return pay, redEnvelope
}
