package rtp

import (
	"fmt"
	"serve/fish_comm/rng"
	"serve/service_fish/domain/probability"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	"serve/service_fish/models"
	"strconv"
	"testing"
)

const (
	max_machineGun_bullet_00005 = 999
	game_id_00005               = models.PSF_ON_00005
	secWebSocketKey_00005       = "jerry00005"
	bet_00005                   = 1
	trigger_icon_id             = 102
)

var mathModule_id_00005 string

func TestService_00005(t *testing.T) {
	var low_fish_0 uint64 = 0
	var low_hit_0 uint64 = 0
	var low_win_0 uint64 = 0
	var low_fish_1 uint64 = 0
	var low_hit_1 uint64 = 0
	var low_win_1 uint64 = 0
	var low_fish_2 uint64 = 0
	var low_hit_2 uint64 = 0
	var low_win_2 uint64 = 0
	var low_fish_3 uint64 = 0
	var low_hit_3 uint64 = 0
	var low_win_3 uint64 = 0
	var low_fish_4 uint64 = 0
	var low_hit_4 uint64 = 0
	var low_win_4 uint64 = 0
	var low_fish_5 uint64 = 0
	var low_hit_5 uint64 = 0
	var low_win_5 uint64 = 0
	var low_fish_6 uint64 = 0
	var low_hit_6 uint64 = 0
	var low_win_6 uint64 = 0
	var low_fish_7 uint64 = 0
	var low_hit_7 uint64 = 0
	var low_win_7 uint64 = 0
	var low_fish_8 uint64 = 0
	var low_hit_8 uint64 = 0
	var low_win_8 uint64 = 0
	var low_fish_9 uint64 = 0
	var low_hit_9 uint64 = 0
	var low_win_9 uint64 = 0
	var low_fish_10 uint64 = 0
	var low_hit_10 uint64 = 0
	var low_win_10 uint64 = 0
	var low_fish_11 uint64 = 0
	var low_hit_11 uint64 = 0
	var low_win_11 uint64 = 0
	var low_fish_12 uint64 = 0
	var low_hit_12 uint64 = 0
	var low_win_12 uint64 = 0
	var low_fish_13 uint64 = 0
	var low_hit_13 uint64 = 0
	var low_win_13 uint64 = 0
	var low_fish_14 uint64 = 0
	var low_hit_14 uint64 = 0
	var low_win_14 uint64 = 0
	var low_fish_15 uint64 = 0
	var low_hit_15 uint64 = 0
	var low_win_15 uint64 = 0
	var low_fish_16 uint64 = 0
	var low_hit_16 uint64 = 0
	var low_win_16 uint64 = 0
	var low_fish_17 uint64 = 0
	var low_hit_17 uint64 = 0
	var low_win_17 uint64 = 0
	var low_fish_18 uint64 = 0
	var low_hit_18 uint64 = 0
	var low_win_18 uint64 = 0
	var low_fish_19 uint64 = 0
	var low_hit_19 uint64 = 0
	var low_win_19 uint64 = 0
	var low_fish_100 uint64 = 0
	var low_hit_100 uint64 = 0
	var low_win_100 uint64 = 0
	var low_fish_101 uint64 = 0
	var low_hit_101 uint64 = 0
	var low_win_101 uint64 = 0
	var low_fish_201 uint64 = 0
	var low_hit_201 uint64 = 0
	var low_win_201 uint64 = 0
	var low_fish_202 uint64 = 0
	var low_hit_202 uint64 = 0
	var low_win_202 uint64 = 0
	var low_fish_300 uint64 = 0
	var low_hit_300 uint64 = 0
	var low_win_300 uint64 = 0
	var low_fish_301 uint64 = 0
	var low_hit_301 uint64 = 0
	var low_win_301 uint64 = 0
	var low_fish_501 uint64 = 0
	var low_hit_501 uint64 = 0
	var low_win_501 uint64 = 0
	var low_fish_502 uint64 = 0
	var low_hit_502 uint64 = 0
	var low_win_502 uint64 = 0
	var low_fish_503 uint64 = 0
	var low_hit_503 uint64 = 0
	var low_win_503 uint64 = 0
	var low_fish_504 uint64 = 0
	var low_hit_504 uint64 = 0
	var low_win_504 uint64 = 0
	var low_fish_505 uint64 = 0
	var low_hit_505 uint64 = 0
	var low_win_505 uint64 = 0

	var high_fish_0 uint64 = 0
	var high_hit_0 uint64 = 0
	var high_win_0 uint64 = 0
	var high_fish_1 uint64 = 0
	var high_hit_1 uint64 = 0
	var high_win_1 uint64 = 0
	var high_fish_2 uint64 = 0
	var high_hit_2 uint64 = 0
	var high_win_2 uint64 = 0
	var high_fish_3 uint64 = 0
	var high_hit_3 uint64 = 0
	var high_win_3 uint64 = 0
	var high_fish_4 uint64 = 0
	var high_hit_4 uint64 = 0
	var high_win_4 uint64 = 0
	var high_fish_5 uint64 = 0
	var high_hit_5 uint64 = 0
	var high_win_5 uint64 = 0
	var high_fish_6 uint64 = 0
	var high_hit_6 uint64 = 0
	var high_win_6 uint64 = 0
	var high_fish_7 uint64 = 0
	var high_hit_7 uint64 = 0
	var high_win_7 uint64 = 0
	var high_fish_8 uint64 = 0
	var high_hit_8 uint64 = 0
	var high_win_8 uint64 = 0
	var high_fish_9 uint64 = 0
	var high_hit_9 uint64 = 0
	var high_win_9 uint64 = 0
	var high_fish_10 uint64 = 0
	var high_hit_10 uint64 = 0
	var high_win_10 uint64 = 0
	var high_fish_11 uint64 = 0
	var high_hit_11 uint64 = 0
	var high_win_11 uint64 = 0
	var high_fish_12 uint64 = 0
	var high_hit_12 uint64 = 0
	var high_win_12 uint64 = 0
	var high_fish_13 uint64 = 0
	var high_hit_13 uint64 = 0
	var high_win_13 uint64 = 0
	var high_fish_14 uint64 = 0
	var high_hit_14 uint64 = 0
	var high_win_14 uint64 = 0
	var high_fish_15 uint64 = 0
	var high_hit_15 uint64 = 0
	var high_win_15 uint64 = 0
	var high_fish_16 uint64 = 0
	var high_hit_16 uint64 = 0
	var high_win_16 uint64 = 0
	var high_fish_17 uint64 = 0
	var high_hit_17 uint64 = 0
	var high_win_17 uint64 = 0
	var high_fish_18 uint64 = 0
	var high_hit_18 uint64 = 0
	var high_win_18 uint64 = 0
	var high_fish_19 uint64 = 0
	var high_hit_19 uint64 = 0
	var high_win_19 uint64 = 0
	var high_fish_100 uint64 = 0
	var high_hit_100 uint64 = 0
	var high_win_100 uint64 = 0
	var high_fish_101 uint64 = 0
	var high_hit_101 uint64 = 0
	var high_win_101 uint64 = 0
	var high_fish_201 uint64 = 0
	var high_hit_201 uint64 = 0
	var high_win_201 uint64 = 0
	var high_fish_202 uint64 = 0
	var high_hit_202 uint64 = 0
	var high_win_202 uint64 = 0
	var high_fish_300 uint64 = 0
	var high_hit_300 uint64 = 0
	var high_win_300 uint64 = 0
	var high_fish_301 uint64 = 0
	var high_hit_301 uint64 = 0
	var high_win_301 uint64 = 0
	var high_fish_501 uint64 = 0
	var high_hit_501 uint64 = 0
	var high_win_501 uint64 = 0
	var high_fish_502 uint64 = 0
	var high_hit_502 uint64 = 0
	var high_win_502 uint64 = 0
	var high_fish_503 uint64 = 0
	var high_hit_503 uint64 = 0
	var high_win_503 uint64 = 0
	var high_fish_504 uint64 = 0
	var high_hit_504 uint64 = 0
	var high_win_504 uint64 = 0
	var high_fish_505 uint64 = 0
	var high_hit_505 uint64 = 0
	var high_win_505 uint64 = 0

	var high_redenvelop_0 uint64 = 0
	var low_redenvelop_0 uint64 = 0
	var high_redenvelop_1 uint64 = 0
	var low_redenvelop_1 uint64 = 0
	var high_redenvelop_2 uint64 = 0
	var low_redenvelop_2 uint64 = 0
	var high_redenvelop_3 uint64 = 0
	var low_redenvelop_3 uint64 = 0
	var high_redenvelop_4 uint64 = 0
	var low_redenvelop_4 uint64 = 0
	var high_redenvelop_5 uint64 = 0
	var low_redenvelop_5 uint64 = 0
	var high_redenvelop_6 uint64 = 0
	var low_redenvelop_6 uint64 = 0
	var high_redenvelop_7 uint64 = 0
	var low_redenvelop_7 uint64 = 0
	var high_redenvelop_8 uint64 = 0
	var low_redenvelop_8 uint64 = 0
	var high_redenvelop_9 uint64 = 0
	var low_redenvelop_9 uint64 = 0
	var high_redenvelop_10 uint64 = 0
	var low_redenvelop_10 uint64 = 0
	var high_redenvelop_11 uint64 = 0
	var low_redenvelop_11 uint64 = 0
	var high_redenvelop_12 uint64 = 0
	var low_redenvelop_12 uint64 = 0
	var high_redenvelop_13 uint64 = 0
	var low_redenvelop_13 uint64 = 0
	var high_redenvelop_14 uint64 = 0
	var low_redenvelop_14 uint64 = 0
	var high_redenvelop_15 uint64 = 0
	var low_redenvelop_15 uint64 = 0
	var high_redenvelop_16 uint64 = 0
	var low_redenvelop_16 uint64 = 0
	var high_redenvelop_17 uint64 = 0
	var low_redenvelop_17 uint64 = 0
	var high_redenvelop_18 uint64 = 0
	var low_redenvelop_18 uint64 = 0
	var high_redenvelop_19 uint64 = 0
	var low_redenvelop_19 uint64 = 0
	var high_redenvelop_201 uint64 = 0
	var low_redenvelop_201 uint64 = 0
	var high_redenvelop_202 uint64 = 0
	var low_redenvelop_202 uint64 = 0
	var high_redenvelop_300 uint64 = 0
	var low_redenvelop_300 uint64 = 0
	var high_redenvelop_301 uint64 = 0
	var low_redenvelop_301 uint64 = 0
	var high_redenvelop_501 uint64 = 0
	var low_redenvelop_501 uint64 = 0
	var high_redenvelop_502 uint64 = 0
	var low_redenvelop_502 uint64 = 0
	var high_redenvelop_503 uint64 = 0
	var low_redenvelop_503 uint64 = 0
	var high_redenvelop_504 uint64 = 0
	var low_redenvelop_504 uint64 = 0
	var high_redenvelop_505 uint64 = 0
	var low_redenvelop_505 uint64 = 0

	fmt.Print("MathModuleId(95、96、97、98), RunTimes(x 萬), FishId -> ")
	var rtpId, run_times int
	var inputFishId string
	fmt.Scanf("%d, %d, %s", &rtpId, &run_times, &inputFishId)

	switch rtpId {
	case 95:
		mathModule_id_00005 = models.PSFM_00006_95_1
	case 96:
		mathModule_id_00005 = models.PSFM_00006_96_1
	case 97:
		mathModule_id_00005 = models.PSFM_00006_97_1
	case 98:
		mathModule_id_00005 = models.PSFM_00006_98_1
	}
	run_times_00005 := run_times * 10000

	for i := 0; i < run_times_00005; i++ {
		var fishId int32
		if inputFishId != "" {
			tempFish, _ := strconv.Atoi(inputFishId)
			fishId = int32(tempFish)
		} else {
			fishId = rngFishId_00005()
		}

		Service.Decrease(game_id_00005, 0, mathModule_id_00005, secWebSocketKey_00005, bet_00005)

		rtpState := Service.RtpState(secWebSocketKey_00005, game_id_00005, subgameId)
		rtpId := Service.RtpId(secWebSocketKey_00005, game_id_00005, 0)
		//rtpState := 0
		//rtpId := "300"

		result := probability.Service.Calc(
			game_id_00005,
			mathModule_id_00005,
			rtpId,
			fishId,
			-1,
			0,
			0,
		)

		switch fishId {
		case 0:
			high_fish_0, high_hit_0, high_win_0, low_fish_0, low_hit_0, low_win_0 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_0, high_hit_0, high_win_0, low_fish_0, low_hit_0, low_win_0,
			)
			high_redenvelop_0, low_redenvelop_0 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_0, low_redenvelop_0)
		case 1:
			high_fish_1, high_hit_1, high_win_1, low_fish_1, low_hit_1, low_win_1 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_1, high_hit_1, high_win_1, low_fish_1, low_hit_1, low_win_1,
			)
			high_redenvelop_1, low_redenvelop_1 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_1, low_redenvelop_1)
		case 2:
			high_fish_2, high_hit_2, high_win_2, low_fish_2, low_hit_2, low_win_2 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_2, high_hit_2, high_win_2, low_fish_2, low_hit_2, low_win_2,
			)
			high_redenvelop_2, low_redenvelop_2 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_2, low_redenvelop_2)
		case 3:
			high_fish_3, high_hit_3, high_win_3, low_fish_3, low_hit_3, low_win_3 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_3, high_hit_3, high_win_3, low_fish_3, low_hit_3, low_win_3,
			)
			high_redenvelop_3, low_redenvelop_3 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_3, low_redenvelop_3)
		case 4:
			high_fish_4, high_hit_4, high_win_4, low_fish_4, low_hit_4, low_win_4 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_4, high_hit_4, high_win_4, low_fish_4, low_hit_4, low_win_4,
			)
			high_redenvelop_4, low_redenvelop_4 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_4, low_redenvelop_4)
		case 5:
			high_fish_5, high_hit_5, high_win_5, low_fish_5, low_hit_5, low_win_5 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_5, high_hit_5, high_win_5, low_fish_5, low_hit_5, low_win_5,
			)
			high_redenvelop_5, low_redenvelop_5 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_5, low_redenvelop_5)
		case 6:
			high_fish_6, high_hit_6, high_win_6, low_fish_6, low_hit_6, low_win_6 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_6, high_hit_6, high_win_6, low_fish_6, low_hit_6, low_win_6,
			)
			high_redenvelop_6, low_redenvelop_6 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_6, low_redenvelop_6)
		case 7:
			high_fish_7, high_hit_7, high_win_7, low_fish_7, low_hit_7, low_win_7 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_7, high_hit_7, high_win_7, low_fish_7, low_hit_7, low_win_7,
			)
			high_redenvelop_7, low_redenvelop_7 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_7, low_redenvelop_7)
		case 8:
			high_fish_8, high_hit_8, high_win_8, low_fish_8, low_hit_8, low_win_8 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_8, high_hit_8, high_win_8, low_fish_8, low_hit_8, low_win_8,
			)
			high_redenvelop_8, low_redenvelop_8 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_8, low_redenvelop_8)
		case 9:
			high_fish_9, high_hit_9, high_win_9, low_fish_9, low_hit_9, low_win_9 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_9, high_hit_9, high_win_9, low_fish_9, low_hit_9, low_win_9,
			)
			high_redenvelop_9, low_redenvelop_9 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_9, low_redenvelop_9)
		case 10:
			high_fish_10, high_hit_10, high_win_10, low_fish_10, low_hit_10, low_win_10 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_10, high_hit_10, high_win_10, low_fish_10, low_hit_10, low_win_10,
			)
			high_redenvelop_10, low_redenvelop_10 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_10, low_redenvelop_10)
		case 11:
			high_fish_11, high_hit_11, high_win_11, low_fish_11, low_hit_11, low_win_11 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_11, high_hit_11, high_win_11, low_fish_11, low_hit_11, low_win_11,
			)
			high_redenvelop_11, low_redenvelop_11 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_11, low_redenvelop_11)
		case 12:
			high_fish_12, high_hit_12, high_win_12, low_fish_12, low_hit_12, low_win_12 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_12, high_hit_12, high_win_12, low_fish_12, low_hit_12, low_win_12,
			)
			high_redenvelop_12, low_redenvelop_12 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_12, low_redenvelop_12)
		case 13:
			high_fish_13, high_hit_13, high_win_13, low_fish_13, low_hit_13, low_win_13 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_13, high_hit_13, high_win_13, low_fish_13, low_hit_13, low_win_13,
			)
			high_redenvelop_13, low_redenvelop_13 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_13, low_redenvelop_13)
		case 14:
			high_fish_14, high_hit_14, high_win_14, low_fish_14, low_hit_14, low_win_14 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_14, high_hit_14, high_win_14, low_fish_14, low_hit_14, low_win_14,
			)
			high_redenvelop_14, low_redenvelop_14 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_14, low_redenvelop_14)
		case 15:
			high_fish_15, high_hit_15, high_win_15, low_fish_15, low_hit_15, low_win_15 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_15, high_hit_15, high_win_15, low_fish_15, low_hit_15, low_win_15,
			)
			high_redenvelop_15, low_redenvelop_15 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_15, low_redenvelop_15)
		case 16:
			high_fish_16, high_hit_16, high_win_16, low_fish_16, low_hit_16, low_win_16 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_16, high_hit_16, high_win_16, low_fish_16, low_hit_16, low_win_16,
			)
			high_redenvelop_16, low_redenvelop_16 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_16, low_redenvelop_16)
		case 17:
			high_fish_17, high_hit_17, high_win_17, low_fish_17, low_hit_17, low_win_17 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_17, high_hit_17, high_win_17, low_fish_17, low_hit_17, low_win_17,
			)
			high_redenvelop_17, low_redenvelop_17 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_17, low_redenvelop_17)
		case 18:
			high_fish_18, high_hit_18, high_win_18, low_fish_18, low_hit_18, low_win_18 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_18, high_hit_18, high_win_18, low_fish_18, low_hit_18, low_win_18,
			)
			high_redenvelop_18, low_redenvelop_18 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_18, low_redenvelop_18)
		case 19:
			high_fish_19, high_hit_19, high_win_19, low_fish_19, low_hit_19, low_win_19 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_19, high_hit_19, high_win_19, low_fish_19, low_hit_19, low_win_19,
			)
			high_redenvelop_19, low_redenvelop_19 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_19, low_redenvelop_19)
		case 100:
			switch rtpState {
			case high:
				if result.TriggerIconId == 100 {
					high_hit_100++
					high_win_100 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				}
				high_fish_100++
			case low:
				if result.TriggerIconId == 100 {
					low_hit_100++
					low_win_100 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				}
				low_fish_100++
			default:
				t.Fatal("RTP State Error")
			}

		case 101:
			switch rtpState {
			case high:
				if result.TriggerIconId == 101 {
					high_hit_101++
					high_win_101 += uint64(result.ExtraData[0].(int))
				}
				high_fish_101++
			case low:
				if result.TriggerIconId == 101 {
					low_hit_101++
					low_win_101 += uint64(result.ExtraData[0].(int))
				}
				low_fish_101++
			default:
				t.Fatal("RTP State Error")
			}

		case 201:
			high_fish_201, high_hit_201, high_win_201, low_fish_201, low_hit_201, low_win_201,
				high_redenvelop_201, low_redenvelop_201 = setMachineGun(
				rtpState, result.Pay, t, result.ExtraTriggerBonus, 201, result.BonusPayload.(int),
				high_fish_201, high_hit_201, high_win_201, low_fish_201, low_hit_201, low_win_201,
				high_redenvelop_201, low_redenvelop_201,
			)
			high_redenvelop_201, low_redenvelop_201 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_201, low_redenvelop_201)
		case 202:
			high_fish_202, high_hit_202, high_win_202, low_fish_202, low_hit_202, low_win_202,
				high_redenvelop_202, low_redenvelop_202 = setMachineGun(
				rtpState, result.Pay, t, result.ExtraTriggerBonus, 202, result.BonusPayload.(int),
				high_fish_202, high_hit_202, high_win_202, low_fish_202, low_hit_202, low_win_202,
				high_redenvelop_202, low_redenvelop_202,
			)
			high_redenvelop_202, low_redenvelop_202 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_202, low_redenvelop_202)
		case 300:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					high_win_300 += uint64(result.Pay * result.Multiplier)
					high_hit_300++
				}
				high_fish_300++
			case low:
				if result.Pay > 0 {
					low_win_300 += uint64(result.Pay * result.Multiplier)
					low_hit_300++
				}
				low_fish_300++
			default:
				t.Fatal("RTP State Error")
			}
			if result.ExtraTriggerBonus != nil {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(result.BonusTypeId),
					102,
					-1,
					0,
					0,
				)

				switch rtpState {
				case high:
					high_redenvelop_300 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				case low:
					low_redenvelop_300 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				default:
					t.Fatal("RTP State Error")
				}
			} else {
				high_redenvelop_300, low_redenvelop_300 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_300, low_redenvelop_300)
			}

		case 301:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					high_win_301 += uint64(result.Pay * result.Multiplier)
					high_hit_301++
				}
				high_fish_301++
			case low:
				if result.Pay > 0 {
					low_win_301 += uint64(result.Pay * result.Multiplier)
					low_hit_301++
				}
				low_fish_301++
			default:
				t.Fatal("RTP State Error")
			}
			if result.ExtraTriggerBonus != nil {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(result.BonusTypeId),
					102,
					-1,
					0,
					0,
				)

				switch rtpState {
				case high:
					high_redenvelop_301 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				case low:
					low_redenvelop_301 += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				default:
					t.Fatal("RTP State Error")
				}
			} else {
				high_redenvelop_301, low_redenvelop_301 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_301, low_redenvelop_301)
			}

		case 501:
			high_fish_501, high_hit_501, high_win_501, low_fish_501, low_hit_501, low_win_501 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_501, high_hit_501, high_win_501, low_fish_501, low_hit_501, low_win_501,
			)
			high_redenvelop_501, low_redenvelop_501 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_501, low_redenvelop_501)
		case 502:
			high_fish_502, high_hit_502, high_win_502, low_fish_502, low_hit_502, low_win_502 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_502, high_hit_502, high_win_502, low_fish_502, low_hit_502, low_win_502,
			)
			high_redenvelop_502, low_redenvelop_502 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_502, low_redenvelop_502)
		case 503:
			high_fish_503, high_hit_503, high_win_503, low_fish_503, low_hit_503, low_win_503 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_503, high_hit_503, high_win_503, low_fish_503, low_hit_503, low_win_503,
			)
			high_redenvelop_503, low_redenvelop_503 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_503, low_redenvelop_503)
		case 504:
			high_fish_504, high_hit_504, high_win_504, low_fish_504, low_hit_504, low_win_504 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_504, high_hit_504, high_win_504, low_fish_504, low_hit_504, low_win_504,
			)
			high_redenvelop_504, low_redenvelop_504 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_504, low_redenvelop_504)
		case 505:
			high_fish_505, high_hit_505, high_win_505, low_fish_505, low_hit_505, low_win_505 = setResult(
				rtpState, result.Pay*result.Multiplier, t,
				high_fish_505, high_hit_505, high_win_505, low_fish_505, low_hit_505, low_win_505,
			)
			high_redenvelop_505, low_redenvelop_505 = setEnvelope(rtpState, result.TriggerIconId, result.BonusTypeId, t, high_redenvelop_505, low_redenvelop_505)
		}
	}

	low_fish := low_fish_0 + low_fish_1 + low_fish_2 + low_fish_3 + low_fish_4 + low_fish_5 + low_fish_6 + low_fish_7 + low_fish_8 + low_fish_9 +
		low_fish_10 + low_fish_11 + low_fish_12 + low_fish_13 + low_fish_14 + low_fish_15 + low_fish_16 + low_fish_17 + low_fish_18 + low_fish_19 +
		low_fish_100 + low_fish_101 + low_fish_201 + low_fish_202 + low_fish_300 + low_fish_301 + low_fish_501 + low_fish_502 + low_fish_503 + low_fish_504 + low_fish_505

	low_hit := low_hit_0 + low_hit_1 + low_hit_2 + low_hit_3 + low_hit_4 + low_hit_5 + low_hit_6 + low_hit_7 + low_hit_8 + low_hit_9 +
		low_hit_10 + low_hit_11 + low_hit_12 + low_hit_13 + low_hit_14 + low_hit_15 + low_hit_16 + low_hit_17 + low_hit_18 + low_hit_19 +
		low_hit_100 + low_hit_101 + low_hit_201 + low_hit_202 + low_hit_300 + low_hit_301 + low_hit_501 + low_hit_502 + low_hit_503 + low_hit_504 + low_hit_505

	low_win := low_win_0 + low_win_1 + low_win_2 + low_win_3 + low_win_4 + low_win_5 + low_win_6 + low_win_7 + low_win_8 + low_win_9 +
		low_win_10 + low_win_11 + low_win_12 + low_win_13 + low_win_14 + low_win_15 + low_win_16 + low_win_17 + low_win_18 + low_win_19 +
		low_win_100 + low_win_101 + low_win_201 + low_win_202 + low_win_300 + low_win_301 + low_win_501 + low_win_502 + low_win_503 + low_win_504 + low_win_505

	high_fish := high_fish_0 + high_fish_1 + high_fish_2 + high_fish_3 + high_fish_4 + high_fish_5 + high_fish_6 + high_fish_7 + high_fish_8 + high_fish_9 +
		high_fish_10 + high_fish_11 + high_fish_12 + high_fish_13 + high_fish_14 + high_fish_15 + high_fish_16 + high_fish_17 + high_fish_18 + high_fish_19 +
		high_fish_100 + high_fish_101 + high_fish_201 + high_fish_202 + high_fish_300 + high_fish_301 + high_fish_501 + high_fish_502 + high_fish_503 + high_fish_504 + high_fish_505

	high_hit := high_hit_0 + high_hit_1 + high_hit_2 + high_hit_3 + high_hit_4 + high_hit_5 + high_hit_6 + high_hit_7 + high_hit_8 + high_hit_9 +
		high_hit_10 + high_hit_11 + high_hit_12 + high_hit_13 + high_hit_14 + high_hit_15 + high_hit_16 + high_hit_17 + high_hit_18 + high_hit_19 +
		high_hit_100 + high_hit_101 + high_hit_201 + high_hit_202 + high_hit_300 + high_hit_301 + high_hit_501 + high_hit_502 + high_hit_503 + high_hit_504 + high_hit_505

	high_win := high_win_0 + high_win_1 + high_win_2 + high_win_3 + high_win_4 + high_win_5 + high_win_6 + high_win_7 + high_win_8 + high_win_9 +
		high_win_10 + high_win_11 + high_win_12 + high_win_13 + high_win_14 + high_win_15 + high_win_16 + high_win_17 + high_win_18 + high_win_19 +
		high_win_100 + high_win_101 + high_win_201 + high_win_202 + high_win_300 + high_win_301 + high_win_501 + high_win_502 + high_win_503 + high_win_504 + high_win_505

	low_redenvelop := low_redenvelop_0 + low_redenvelop_1 + low_redenvelop_2 + low_redenvelop_3 + low_redenvelop_4 + low_redenvelop_5 + low_redenvelop_6 +
		low_redenvelop_7 + low_redenvelop_8 + low_redenvelop_9 + low_redenvelop_10 + low_redenvelop_11 + low_redenvelop_12 + low_redenvelop_13 +
		low_redenvelop_14 + low_redenvelop_15 + low_redenvelop_16 + low_redenvelop_17 + low_redenvelop_18 + low_redenvelop_19 + low_redenvelop_201 +
		low_redenvelop_202 + low_redenvelop_300 + low_redenvelop_301 + low_redenvelop_501 + low_redenvelop_502 + low_redenvelop_503 + low_redenvelop_504 + low_redenvelop_505

	high_redenvelop := high_redenvelop_0 + high_redenvelop_1 + high_redenvelop_2 + high_redenvelop_3 + high_redenvelop_4 + high_redenvelop_5 + high_redenvelop_6 +
		high_redenvelop_7 + high_redenvelop_8 + high_redenvelop_9 + high_redenvelop_10 + high_redenvelop_11 + high_redenvelop_12 + high_redenvelop_13 +
		high_redenvelop_14 + high_redenvelop_15 + high_redenvelop_16 + high_redenvelop_17 + high_redenvelop_18 + high_redenvelop_19 + high_redenvelop_201 +
		high_redenvelop_202 + high_redenvelop_300 + high_redenvelop_301 + high_redenvelop_501 + high_redenvelop_502 + high_redenvelop_503 + high_redenvelop_504 + high_redenvelop_505

	fmt.Println("FishId", ":", "LowTimes", "LowHit", "LowWin", "HighTimes", "HighHit", "HighWin", "LowRedenvelop", "HighRedenvelop",
		"LowHitRate", "LowPayRate", "HighHitRate", "HighPayRate", "TotalPayRate")
	fmt.Println(0, ":", low_fish_0, low_hit_0, low_win_0, high_fish_0, high_hit_0, high_win_0, low_redenvelop_0, high_redenvelop_0,
		getRate(low_hit_0, low_fish_0), getRate(low_win_0+low_redenvelop_0, low_fish_0),
		getRate(high_hit_0, high_fish_0), getRate(high_win_0+high_redenvelop_0, high_fish_0),
		getRate(low_win_0+high_win_0+low_redenvelop_0+high_redenvelop_0, low_fish_0+high_fish_0))
	fmt.Println(1, ":", low_fish_1, low_hit_1, low_win_1, high_fish_1, high_hit_1, high_win_1, low_redenvelop_1, high_redenvelop_1,
		getRate(low_hit_1, low_fish_1), getRate(low_win_1+low_redenvelop_1, low_fish_1),
		getRate(high_hit_1, high_fish_1), getRate(high_win_1+high_redenvelop_1, high_fish_1),
		getRate(low_win_1+high_win_1+low_redenvelop_1+high_redenvelop_1, low_fish_1+high_fish_1))
	fmt.Println(2, ":", low_fish_2, low_hit_2, low_win_2, high_fish_2, high_hit_2, high_win_2, low_redenvelop_2, high_redenvelop_2,
		getRate(low_hit_2, low_fish_2), getRate(low_win_2+low_redenvelop_2, low_fish_2),
		getRate(high_hit_2, high_fish_2), getRate(high_win_2+high_redenvelop_2, high_fish_2),
		getRate(low_win_2+high_win_2+low_redenvelop_2+high_redenvelop_2, low_fish_2+high_fish_2))
	fmt.Println(3, ":", low_fish_3, low_hit_3, low_win_3, high_fish_3, high_hit_3, high_win_3, low_redenvelop_3, high_redenvelop_3,
		getRate(low_hit_3, low_fish_3), getRate(low_win_3+low_redenvelop_3, low_fish_3),
		getRate(high_hit_3, high_fish_3), getRate(high_win_3+high_redenvelop_3, high_fish_3),
		getRate(low_win_3+high_win_3+low_redenvelop_3+high_redenvelop_3, low_fish_3+high_fish_3))
	fmt.Println(4, ":", low_fish_4, low_hit_4, low_win_4, high_fish_4, high_hit_4, high_win_4, low_redenvelop_4, high_redenvelop_4,
		getRate(low_hit_4, low_fish_4), getRate(low_win_4+low_redenvelop_4, low_fish_4),
		getRate(high_hit_4, high_fish_4), getRate(high_win_4+high_redenvelop_4, high_fish_4),
		getRate(low_win_4+high_win_4+low_redenvelop_4+high_redenvelop_4, low_fish_4+high_fish_4))
	fmt.Println(5, ":", low_fish_5, low_hit_5, low_win_5, high_fish_5, high_hit_5, high_win_5, low_redenvelop_5, high_redenvelop_5,
		getRate(low_hit_5, low_fish_5), getRate(low_win_5+low_redenvelop_5, low_fish_5),
		getRate(high_hit_5, high_fish_5), getRate(high_win_5+high_redenvelop_5, high_fish_5),
		getRate(low_win_5+high_win_5+low_redenvelop_5+high_redenvelop_5, low_fish_5+high_fish_5))
	fmt.Println(6, ":", low_fish_6, low_hit_6, low_win_6, high_fish_6, high_hit_6, high_win_6, low_redenvelop_6, high_redenvelop_6,
		getRate(low_hit_6, low_fish_6), getRate(low_win_6+low_redenvelop_6, low_fish_6),
		getRate(high_hit_6, high_fish_6), getRate(high_win_6+high_redenvelop_6, high_fish_6),
		getRate(low_win_6+high_win_6+low_redenvelop_6+high_redenvelop_6, low_fish_6+high_fish_6))
	fmt.Println(7, ":", low_fish_7, low_hit_7, low_win_7, high_fish_7, high_hit_7, high_win_7, low_redenvelop_7, high_redenvelop_7,
		getRate(low_hit_7, low_fish_7), getRate(low_win_7+low_redenvelop_7, low_fish_7),
		getRate(high_hit_7, high_fish_7), getRate(high_win_7+high_redenvelop_7, high_fish_7),
		getRate(low_win_7+high_win_7+low_redenvelop_7+high_redenvelop_7, low_fish_7+high_fish_7))
	fmt.Println(8, ":", low_fish_8, low_hit_8, low_win_8, high_fish_8, high_hit_8, high_win_8, low_redenvelop_8, high_redenvelop_8,
		getRate(low_hit_8, low_fish_8), getRate(low_win_8+low_redenvelop_8, low_fish_8),
		getRate(high_hit_8, high_fish_8), getRate(high_win_8+high_redenvelop_8, high_fish_8),
		getRate(low_win_8+high_win_8+low_redenvelop_8+high_redenvelop_8, low_fish_8+high_fish_8))
	fmt.Println(9, ":", low_fish_9, low_hit_9, low_win_9, high_fish_9, high_hit_9, high_win_9, low_redenvelop_9, high_redenvelop_9,
		getRate(low_hit_9, low_fish_9), getRate(low_win_9+low_redenvelop_9, low_fish_9),
		getRate(high_hit_9, high_fish_9), getRate(high_win_9+high_redenvelop_9, high_fish_9),
		getRate(low_win_9+high_win_9+low_redenvelop_9+high_redenvelop_9, low_fish_9+high_fish_9))
	fmt.Println(10, ":", low_fish_10, low_hit_10, low_win_10, high_fish_10, high_hit_10, high_win_10, low_redenvelop_10, high_redenvelop_10,
		getRate(low_hit_10, low_fish_10), getRate(low_win_10+low_redenvelop_10, low_fish_10),
		getRate(high_hit_10, high_fish_10), getRate(high_win_10+high_redenvelop_10, high_fish_10),
		getRate(low_win_10+high_win_10+low_redenvelop_10+high_redenvelop_10, low_fish_10+high_fish_10))
	fmt.Println(11, ":", low_fish_11, low_hit_11, low_win_11, high_fish_11, high_hit_11, high_win_11, low_redenvelop_11, high_redenvelop_11,
		getRate(low_hit_11, low_fish_11), getRate(low_win_11+low_redenvelop_11, low_fish_11),
		getRate(high_hit_11, high_fish_11), getRate(high_win_11+high_redenvelop_11, high_fish_11),
		getRate(low_win_11+high_win_11+low_redenvelop_11+high_redenvelop_11, low_fish_11+high_fish_11))
	fmt.Println(12, ":", low_fish_12, low_hit_12, low_win_12, high_fish_12, high_hit_12, high_win_12, low_redenvelop_12, high_redenvelop_12,
		getRate(low_hit_12, low_fish_12), getRate(low_win_12+low_redenvelop_12, low_fish_12),
		getRate(high_hit_12, high_fish_12), getRate(high_win_12+high_redenvelop_12, high_fish_12),
		getRate(low_win_12+high_win_12+low_redenvelop_12+high_redenvelop_12, low_fish_12+high_fish_12))
	fmt.Println(13, ":", low_fish_13, low_hit_13, low_win_13, high_fish_13, high_hit_13, high_win_13, low_redenvelop_13, high_redenvelop_13,
		getRate(low_hit_13, low_fish_13), getRate(low_win_13+low_redenvelop_13, low_fish_13),
		getRate(high_hit_13, high_fish_13), getRate(high_win_13+high_redenvelop_13, high_fish_13),
		getRate(low_win_13+high_win_13+low_redenvelop_13+high_redenvelop_13, low_fish_13+high_fish_13))
	fmt.Println(14, ":", low_fish_14, low_hit_14, low_win_14, high_fish_14, high_hit_14, high_win_14, low_redenvelop_14, high_redenvelop_14,
		getRate(low_hit_14, low_fish_14), getRate(low_win_14+low_redenvelop_14, low_fish_14),
		getRate(high_hit_14, high_fish_14), getRate(high_win_14+high_redenvelop_14, high_fish_14),
		getRate(low_win_14+high_win_14+low_redenvelop_14+high_redenvelop_14, low_fish_14+high_fish_14))
	fmt.Println(15, ":", low_fish_15, low_hit_15, low_win_15, high_fish_15, high_hit_15, high_win_15, low_redenvelop_15, high_redenvelop_15,
		getRate(low_hit_15, low_fish_15), getRate(low_win_15+low_redenvelop_15, low_fish_15),
		getRate(high_hit_15, high_fish_15), getRate(high_win_15+high_redenvelop_15, high_fish_15),
		getRate(low_win_15+high_win_15+low_redenvelop_15+high_redenvelop_15, low_fish_15+high_fish_15))
	fmt.Println(16, ":", low_fish_16, low_hit_16, low_win_16, high_fish_16, high_hit_16, high_win_16, low_redenvelop_16, high_redenvelop_16,
		getRate(low_hit_16, low_fish_16), getRate(low_win_16+low_redenvelop_16, low_fish_16),
		getRate(high_hit_16, high_fish_16), getRate(high_win_16+high_redenvelop_16, high_fish_16),
		getRate(low_win_16+high_win_16+low_redenvelop_16+high_redenvelop_16, low_fish_16+high_fish_16))
	fmt.Println(17, ":", low_fish_17, low_hit_17, low_win_17, high_fish_17, high_hit_17, high_win_17, low_redenvelop_17, high_redenvelop_17,
		getRate(low_hit_17, low_fish_17), getRate(low_win_17+low_redenvelop_17, low_fish_17),
		getRate(high_hit_17, high_fish_17), getRate(high_win_17+high_redenvelop_17, high_fish_17),
		getRate(low_win_17+high_win_17+low_redenvelop_17+high_redenvelop_17, low_fish_17+high_fish_17))
	fmt.Println(18, ":", low_fish_18, low_hit_18, low_win_18, high_fish_18, high_hit_18, high_win_18, low_redenvelop_18, high_redenvelop_18,
		getRate(low_hit_18, low_fish_18), getRate(low_win_18+low_redenvelop_18, low_fish_18),
		getRate(high_hit_18, high_fish_18), getRate(high_win_18+high_redenvelop_18, high_fish_18),
		getRate(low_win_18+high_win_18+low_redenvelop_18+high_redenvelop_18, low_fish_18+high_fish_18))
	fmt.Println(19, ":", low_fish_19, low_hit_19, low_win_19, high_fish_19, high_hit_19, high_win_19, low_redenvelop_19, high_redenvelop_19,
		getRate(low_hit_19, low_fish_19), getRate(low_win_19+low_redenvelop_19, low_fish_19),
		getRate(high_hit_19, high_fish_19), getRate(high_win_19+high_redenvelop_19, high_fish_19),
		getRate(low_win_19+high_win_19+low_redenvelop_19+high_redenvelop_19, low_fish_19+high_fish_19))
	fmt.Println(100, ":", low_fish_100, low_hit_100, low_win_100, high_fish_100, high_hit_100, high_win_100,
		getRate(low_hit_100, low_fish_100), getRate(low_win_100, low_fish_100), getRate(high_hit_100, high_fish_100), getRate(high_win_100, high_fish_100),
		getRate(low_win_100+high_win_100, low_fish_100+high_fish_100))
	fmt.Println(101, ":", low_fish_101, low_hit_101, low_win_101, high_fish_101, high_hit_101, high_win_101,
		getRate(low_hit_101, low_fish_101), getRate(low_win_101, low_fish_101), getRate(high_hit_101, high_fish_101), getRate(high_win_101, high_fish_101),
		getRate(low_win_101+high_win_101, low_fish_101+high_fish_101))
	fmt.Println(201, ":", low_fish_201, low_hit_201, low_win_201, high_fish_201, high_hit_201, high_win_201, low_redenvelop_201, high_redenvelop_201,
		getRate(low_hit_201, low_fish_201), getRate(low_win_201+low_redenvelop_201, low_fish_201),
		getRate(high_hit_201, high_fish_201), getRate(high_win_201+high_redenvelop_201, high_fish_201),
		getRate(low_win_201+high_win_201+low_redenvelop_201+high_redenvelop_201, low_fish_201+high_fish_201))
	fmt.Println(202, ":", low_fish_202, low_hit_202, low_win_202, high_fish_202, high_hit_202, high_win_202, low_redenvelop_202, high_redenvelop_202,
		getRate(low_hit_202, low_fish_202), getRate(low_win_202+low_redenvelop_202, low_fish_202),
		getRate(high_hit_202, high_fish_202), getRate(high_win_202+high_redenvelop_202, high_fish_202),
		getRate(low_win_202+high_win_202+low_redenvelop_202+high_redenvelop_202, low_fish_202+high_fish_202))
	fmt.Println(300, ":", low_fish_300, low_hit_300, low_win_300, high_fish_300, high_hit_300, high_win_300, low_redenvelop_300, high_redenvelop_300,
		getRate(low_hit_300, low_fish_300), getRate(low_win_300+low_redenvelop_300, low_fish_300),
		getRate(high_hit_300, high_fish_300), getRate(high_win_300+high_redenvelop_300, high_fish_300),
		getRate(low_win_300+high_win_300+low_redenvelop_300+high_redenvelop_300, low_fish_300+high_fish_300))
	fmt.Println(301, ":", low_fish_301, low_hit_301, low_win_301, high_fish_301, high_hit_301, high_win_301, low_redenvelop_301, high_redenvelop_301,
		getRate(low_hit_301, low_fish_301), getRate(low_win_301+low_redenvelop_301, low_fish_301),
		getRate(high_hit_301, high_fish_301), getRate(high_win_301+high_redenvelop_301, high_fish_301),
		getRate(low_win_301+high_win_301+low_redenvelop_301+high_redenvelop_301, low_fish_301+high_fish_301))
	fmt.Println(501, ":", low_fish_501, low_hit_501, low_win_501, high_fish_501, high_hit_501, high_win_501, low_redenvelop_501, high_redenvelop_501,
		getRate(low_hit_501, low_fish_501), getRate(low_win_501+low_redenvelop_501, low_fish_501),
		getRate(high_hit_501, high_fish_501), getRate(high_win_501+high_redenvelop_501, high_fish_501),
		getRate(low_win_501+high_win_501+low_redenvelop_501+high_redenvelop_501, low_fish_501+high_fish_501))
	fmt.Println(502, ":", low_fish_502, low_hit_502, low_win_502, high_fish_502, high_hit_502, high_win_502, low_redenvelop_502, high_redenvelop_502,
		getRate(low_hit_502, low_fish_502), getRate(low_win_502+low_redenvelop_502, low_fish_502),
		getRate(high_hit_502, high_fish_502), getRate(high_win_502+high_redenvelop_502, high_fish_502),
		getRate(low_win_502+high_win_502+low_redenvelop_502+high_redenvelop_502, low_fish_502+high_fish_502))
	fmt.Println(503, ":", low_fish_503, low_hit_503, low_win_503, high_fish_503, high_hit_503, high_win_503, low_redenvelop_503, high_redenvelop_503,
		getRate(low_hit_503, low_fish_503), getRate(low_win_503+low_redenvelop_503, low_fish_503),
		getRate(high_hit_503, high_fish_503), getRate(high_win_503+high_redenvelop_503, high_fish_503),
		getRate(low_win_503+high_win_503+low_redenvelop_503+high_redenvelop_503, low_fish_503+high_fish_503))
	fmt.Println(504, ":", low_fish_504, low_hit_504, low_win_504, high_fish_504, high_hit_504, high_win_504, low_redenvelop_504, high_redenvelop_504,
		getRate(low_hit_504, low_fish_504), getRate(low_win_504+low_redenvelop_504, low_fish_504),
		getRate(high_hit_504, high_fish_504), getRate(high_win_504+high_redenvelop_504, high_fish_504),
		getRate(low_win_504+high_win_504+low_redenvelop_504+high_redenvelop_504, low_fish_504+high_fish_504))
	fmt.Println(505, ":", low_fish_505, low_hit_505, low_win_505, high_fish_505, high_hit_505, high_win_505, low_redenvelop_505, high_redenvelop_505,
		getRate(low_hit_505, low_fish_505), getRate(low_win_505+low_redenvelop_505, low_fish_505),
		getRate(high_hit_505, high_fish_505), getRate(high_win_505+high_redenvelop_505, high_fish_505),
		getRate(low_win_505+high_win_505+low_redenvelop_505+high_redenvelop_505, low_fish_505+high_fish_505))
	fmt.Println("Total", low_fish, low_hit, low_win, high_fish, high_hit, high_win, low_redenvelop, high_redenvelop,
		getRate(low_win+high_win+low_redenvelop+high_redenvelop, low_fish+high_fish))
}

func setResult(rtpState, iconPay int, t *testing.T,
	highFish, highHit, highWin uint64,
	lowFish, lowHit, lowWin uint64,
) (result_high_Fish, result_high_Hit, result_high_Win uint64,
	result_low_Fish, result_low_Hit, result_low_Win uint64,
) {
	switch rtpState {
	case high:
		result_high_Fish, result_high_Hit, result_high_Win = setCalc(iconPay, highFish, highHit, highWin)
		result_low_Fish = lowFish
		result_low_Hit = lowHit
		result_low_Win = lowWin

	case low:
		result_low_Fish, result_low_Hit, result_low_Win = setCalc(iconPay, lowFish, lowHit, lowWin)
		result_high_Fish = highFish
		result_high_Hit = highHit
		result_high_Win = highWin

	default:
		t.Fatal("RTP State Error")
	}

	return result_high_Fish, result_high_Hit, result_high_Win, result_low_Fish, result_low_Hit, result_low_Win
}

func setCalc(iconPay int,
	fish, hit, win uint64,
) (resultFish, resultHit, resultWin uint64) {
	if iconPay > 0 {
		resultWin = win + uint64(iconPay)
		resultHit = hit + 1
	} else {
		resultWin = win
		resultHit = hit
	}

	resultFish = fish + 1

	return resultFish, resultHit, resultWin
}

func setEnvelope(rtpState, triggerIconId, bonusTypeId int, t *testing.T,
	high_envelop, low_envelop uint64,
) (highEnvelop, lowEnvelop uint64) {
	if triggerIconId == trigger_icon_id {
		result := probability.Service.Calc(
			game_id_00005,
			mathModule_id_00005,
			strconv.Itoa(bonusTypeId),
			trigger_icon_id,
			-1,
			0,
			0,
		)

		switch rtpState {
		case high:
			highEnvelop = high_envelop + uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
			lowEnvelop = low_envelop
		case low:
			highEnvelop = high_envelop
			lowEnvelop = low_envelop + uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
		default:
			t.Fatal("RTP State Error")
		}
	} else {
		highEnvelop = high_envelop
		lowEnvelop = low_envelop
	}

	return highEnvelop, lowEnvelop
}

func rngFishId_00005() int32 {
	options := make([]rng.Option, 0)

	for i := 0; i < 20; i++ {
		options = append(options, rng.Option{1, i})
	}
	options = append(options, rng.Option{1, 100})
	options = append(options, rng.Option{1, 101})
	options = append(options, rng.Option{1, 201})
	options = append(options, rng.Option{1, 202})
	options = append(options, rng.Option{1, 300})
	options = append(options, rng.Option{1, 301})
	options = append(options, rng.Option{1, 501})
	options = append(options, rng.Option{1, 502})
	options = append(options, rng.Option{1, 503})
	options = append(options, rng.Option{1, 504})
	options = append(options, rng.Option{1, 505})

	return int32(rng.New(options).Item.(int))
}

func rngRedEnvelope_00005() int32 {
	options := make([]rng.Option, 0, 5)

	for i := 0; i < 5; i++ {
		options = append(options, rng.Option{1, i})
	}

	return int32(rng.New(options).Item.(int))
}

func setMachineGun(rtpState, iconPay int, t *testing.T, extraTriggerBonus []*probability.Probability, machineGun_type int32, bullets int,
	high_fish, high_hit, high_win, low_fish, low_hit, low_win, high_envelope, low_envelope uint64,
) (result_high_fish, result_high_hit, result_high_win,
	result_low_fish, result_low_hit, result_low_win,
	result_high_envelope, result_low_envelope uint64) {
	switch rtpState {
	case high:
		if iconPay > 0 {
			result_high_win = high_win + uint64(iconPay)
			result_low_win = low_win

			if extraTriggerBonus != nil {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(extraTriggerBonus[0].BonusTypeId),
					102,
					-1,
					0,
					0,
				)
				switch rtpState {
				case high:
					result_high_envelope = high_envelope + uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
					result_low_envelope = low_envelope
				case low:
					result_high_envelope = high_envelope
					result_low_envelope = low_envelope + uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				default:
					t.Fatal("RTP State Error")
				}
			} else {
				result_high_envelope = high_envelope
				result_low_envelope = low_envelope
			}

			if machineGun_type == 201 {
				result_high_win += shot_machineGun(PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.MachineGun.RTP6.ID,
					bullets,
					machineGun_type,
				)
			} else {
				result_high_win += shot_machineGun(PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.SuperMachineGun.RTP8.ID,
					bullets,
					machineGun_type,
				)
			}
			result_high_hit = high_hit + 1
			result_low_hit = low_hit
		} else {
			result_high_hit = high_hit
			result_low_hit = low_hit
			result_high_win = high_win
			result_low_win = low_win
			result_high_envelope = high_envelope
			result_low_envelope = low_envelope
		}
		result_high_fish = high_fish + 1
		result_low_fish = low_fish

	case low:
		if iconPay > 0 {
			result_low_win = low_win + uint64(iconPay)
			result_high_win = high_win

			if extraTriggerBonus != nil {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(extraTriggerBonus[0].BonusTypeId),
					102,
					-1,
					0,
					0,
				)

				switch rtpState {
				case high:
					result_high_envelope = high_envelope + uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
					result_low_envelope = low_envelope
				case low:
					result_high_envelope = high_envelope
					result_low_envelope = low_envelope + uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
				default:
					t.Fatal("RTP State Error")
				}
			} else {
				result_high_envelope = high_envelope
				result_low_envelope = low_envelope
			}

			if machineGun_type == 201 {
				result_low_win += shot_machineGun(PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.MachineGun.RTP6.ID,
					bullets,
					machineGun_type,
				)
			} else {
				result_low_win += shot_machineGun(PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath.Icons.SuperMachineGun.RTP8.ID,
					bullets,
					machineGun_type,
				)
			}
			result_low_hit = low_hit + 1
			result_high_hit = high_hit
		} else {
			result_high_hit = high_hit
			result_low_hit = low_hit
			result_high_win = high_win
			result_low_win = low_win
			result_high_envelope = high_envelope
			result_low_envelope = low_envelope
		}
		result_low_fish = low_fish + 1
		result_high_fish = high_fish

	default:
		t.Fatal("RTP State Error")
	}

	return result_high_fish, result_high_hit, result_high_win,
		result_low_fish, result_low_hit, result_low_win,
		result_high_envelope, result_low_envelope
}

func shot_machineGun(rtpId string, bullets int, machineGun_type int32) (pay uint64) {
	var fishId int32 = -1
	var newBullets = 0

	var check_machineGun_type int32 = 0
	if machineGun_type == 201 {
		check_machineGun_type = 202
	} else {
		check_machineGun_type = 201
	}

	for i := 0; i < bullets; i++ {
		for {
			fishId = rngFishId_00005()
			if fishId != check_machineGun_type {
				break
			}
		}

		result := probability.Service.Calc(
			game_id_00005,
			mathModule_id_00005,
			rtpId,
			fishId,
			machineGun_type,
			0,
			0,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay * result.Multiplier)
		}

		if result.TriggerIconId == 100 {
			pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
		}

		if result.TriggerIconId == 101 {
			pay += uint64(result.ExtraData[0].(int))
		}

		if result.ExtraTriggerBonus != nil {
			if result.ExtraTriggerBonus[0].TriggerIconId == 102 {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(result.ExtraTriggerBonus[0].BonusTypeId),
					102,
					machineGun_type,
					0,
					0,
				)
				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
			}
		} else {
			if result.TriggerIconId == 102 {
				result := probability.Service.Calc(
					game_id_00005,
					mathModule_id_00005,
					strconv.Itoa(result.BonusTypeId),
					102,
					machineGun_type,
					0,
					0,
				)

				pay += uint64(result.BonusPayload.([]int)[rngRedEnvelope_00005()])
			}
		}

		if result.TriggerIconId == int(machineGun_type) {
			newBullets += result.BonusPayload.(int)

			if newBullets > max_machineGun_bullet_00005 {
				newBullets = max_machineGun_bullet_00005
			}
		}
	}

	if newBullets > 0 {
		pay += shot_machineGun(rtpId, newBullets, machineGun_type)
	}

	return pay
}
