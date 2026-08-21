package rtp

import (
	"fmt"
	"serve/service_fish/domain/probability"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"
	"serve/service_fish/models"
	"strconv"
	"testing"
)

const gameId_00001 = models.PSF_ON_00001

var mathModuleId_00001 string

func TestService_00001_Calc(t *testing.T) {
	var rtpLow_fish_0 uint64 = 0
	var rtpLow_fish_1 uint64 = 0
	var rtpLow_fish_2 uint64 = 0
	var rtpLow_fish_3 uint64 = 0
	var rtpLow_fish_4 uint64 = 0
	var rtpLow_fish_5 uint64 = 0
	var rtpLow_fish_6 uint64 = 0
	var rtpLow_fish_7 uint64 = 0
	var rtpLow_fish_8 uint64 = 0
	var rtpLow_fish_9 uint64 = 0
	var rtpLow_fish_10 uint64 = 0
	var rtpLow_fish_11 uint64 = 0
	var rtpLow_fish_12 uint64 = 0
	var rtpLow_fish_13 uint64 = 0
	var rtpLow_fish_14 uint64 = 0
	var rtpLow_fish_15 uint64 = 0
	var rtpLow_fish_16 uint64 = 0
	var rtpLow_fish_17 uint64 = 0
	var rtpLow_fish_18 uint64 = 0
	var rtpLow_fish_19 uint64 = 0
	var rtpLow_fish_20 uint64 = 0
	var rtpLow_fish_21 uint64 = 0
	var rtpLow_fish_22 uint64 = 0
	var rtpLow_fish_23 uint64 = 0
	var rtpLow_fish_24 uint64 = 0
	var rtpLow_fish_25 uint64 = 0
	var rtpLow_fish_26 uint64 = 0
	var rtpLow_fish_27 uint64 = 0
	var rtpLow_fish_28 uint64 = 0
	var rtpLow_fish_29 uint64 = 0

	var rtpHigh_fish_0 uint64 = 0
	var rtpHigh_fish_1 uint64 = 0
	var rtpHigh_fish_2 uint64 = 0
	var rtpHigh_fish_3 uint64 = 0
	var rtpHigh_fish_4 uint64 = 0
	var rtpHigh_fish_5 uint64 = 0
	var rtpHigh_fish_6 uint64 = 0
	var rtpHigh_fish_7 uint64 = 0
	var rtpHigh_fish_8 uint64 = 0
	var rtpHigh_fish_9 uint64 = 0
	var rtpHigh_fish_10 uint64 = 0
	var rtpHigh_fish_11 uint64 = 0
	var rtpHigh_fish_12 uint64 = 0
	var rtpHigh_fish_13 uint64 = 0
	var rtpHigh_fish_14 uint64 = 0
	var rtpHigh_fish_15 uint64 = 0
	var rtpHigh_fish_16 uint64 = 0
	var rtpHigh_fish_17 uint64 = 0
	var rtpHigh_fish_18 uint64 = 0
	var rtpHigh_fish_19 uint64 = 0
	var rtpHigh_fish_20 uint64 = 0
	var rtpHigh_fish_21 uint64 = 0
	var rtpHigh_fish_22 uint64 = 0
	var rtpHigh_fish_23 uint64 = 0
	var rtpHigh_fish_24 uint64 = 0
	var rtpHigh_fish_25 uint64 = 0
	var rtpHigh_fish_26 uint64 = 0
	var rtpHigh_fish_27 uint64 = 0
	var rtpHigh_fish_28 uint64 = 0
	var rtpHigh_fish_29 uint64 = 0

	var rtpLow_win_0 uint64 = 0
	var rtpLow_win_1 uint64 = 0
	var rtpLow_win_2 uint64 = 0
	var rtpLow_win_3 uint64 = 0
	var rtpLow_win_4 uint64 = 0
	var rtpLow_win_5 uint64 = 0
	var rtpLow_win_6 uint64 = 0
	var rtpLow_win_7 uint64 = 0
	var rtpLow_win_8 uint64 = 0
	var rtpLow_win_9 uint64 = 0
	var rtpLow_win_10 uint64 = 0
	var rtpLow_win_11 uint64 = 0
	var rtpLow_win_12 uint64 = 0
	var rtpLow_win_13 uint64 = 0
	var rtpLow_win_14 uint64 = 0
	var rtpLow_win_15 uint64 = 0
	var rtpLow_win_16 uint64 = 0
	var rtpLow_win_17 uint64 = 0
	var rtpLow_win_18 uint64 = 0
	var rtpLow_win_19 uint64 = 0
	var rtpLow_win_20 uint64 = 0
	var rtpLow_win_21 uint64 = 0
	var rtpLow_win_22 uint64 = 0
	var rtpLow_win_23 uint64 = 0
	var rtpLow_win_24 uint64 = 0
	var rtpLow_win_25 uint64 = 0
	var rtpLow_win_26 uint64 = 0
	var rtpLow_win_27 uint64 = 0
	var rtpLow_win_28 uint64 = 0
	var rtpLow_win_29 uint64 = 0

	var rtpHigh_win_0 uint64 = 0
	var rtpHigh_win_1 uint64 = 0
	var rtpHigh_win_2 uint64 = 0
	var rtpHigh_win_3 uint64 = 0
	var rtpHigh_win_4 uint64 = 0
	var rtpHigh_win_5 uint64 = 0
	var rtpHigh_win_6 uint64 = 0
	var rtpHigh_win_7 uint64 = 0
	var rtpHigh_win_8 uint64 = 0
	var rtpHigh_win_9 uint64 = 0
	var rtpHigh_win_10 uint64 = 0
	var rtpHigh_win_11 uint64 = 0
	var rtpHigh_win_12 uint64 = 0
	var rtpHigh_win_13 uint64 = 0
	var rtpHigh_win_14 uint64 = 0
	var rtpHigh_win_15 uint64 = 0
	var rtpHigh_win_16 uint64 = 0
	var rtpHigh_win_17 uint64 = 0
	var rtpHigh_win_18 uint64 = 0
	var rtpHigh_win_19 uint64 = 0
	var rtpHigh_win_20 uint64 = 0
	var rtpHigh_win_21 uint64 = 0
	var rtpHigh_win_22 uint64 = 0
	var rtpHigh_win_23 uint64 = 0
	var rtpHigh_win_24 uint64 = 0
	var rtpHigh_win_25 uint64 = 0
	var rtpHigh_win_26 uint64 = 0
	var rtpHigh_win_27 uint64 = 0
	var rtpHigh_win_28 uint64 = 0
	var rtpHigh_win_29 uint64 = 0

	var rtpLow_hit_0 uint64 = 0
	var rtpLow_hit_1 uint64 = 0
	var rtpLow_hit_2 uint64 = 0
	var rtpLow_hit_3 uint64 = 0
	var rtpLow_hit_4 uint64 = 0
	var rtpLow_hit_5 uint64 = 0
	var rtpLow_hit_6 uint64 = 0
	var rtpLow_hit_7 uint64 = 0
	var rtpLow_hit_8 uint64 = 0
	var rtpLow_hit_9 uint64 = 0
	var rtpLow_hit_10 uint64 = 0
	var rtpLow_hit_11 uint64 = 0
	var rtpLow_hit_12 uint64 = 0
	var rtpLow_hit_13 uint64 = 0
	var rtpLow_hit_14 uint64 = 0
	var rtpLow_hit_15 uint64 = 0
	var rtpLow_hit_16 uint64 = 0
	var rtpLow_hit_17 uint64 = 0
	var rtpLow_hit_18 uint64 = 0
	var rtpLow_hit_19 uint64 = 0
	var rtpLow_hit_20 uint64 = 0
	var rtpLow_hit_21 uint64 = 0
	var rtpLow_hit_22 uint64 = 0
	var rtpLow_hit_23 uint64 = 0
	var rtpLow_hit_24 uint64 = 0
	var rtpLow_hit_25 uint64 = 0
	var rtpLow_hit_26 uint64 = 0
	var rtpLow_hit_27 uint64 = 0
	var rtpLow_hit_28 uint64 = 0
	var rtpLow_hit_29 uint64 = 0

	var rtpHigh_hit_0 uint64 = 0
	var rtpHigh_hit_1 uint64 = 0
	var rtpHigh_hit_2 uint64 = 0
	var rtpHigh_hit_3 uint64 = 0
	var rtpHigh_hit_4 uint64 = 0
	var rtpHigh_hit_5 uint64 = 0
	var rtpHigh_hit_6 uint64 = 0
	var rtpHigh_hit_7 uint64 = 0
	var rtpHigh_hit_8 uint64 = 0
	var rtpHigh_hit_9 uint64 = 0
	var rtpHigh_hit_10 uint64 = 0
	var rtpHigh_hit_11 uint64 = 0
	var rtpHigh_hit_12 uint64 = 0
	var rtpHigh_hit_13 uint64 = 0
	var rtpHigh_hit_14 uint64 = 0
	var rtpHigh_hit_15 uint64 = 0
	var rtpHigh_hit_16 uint64 = 0
	var rtpHigh_hit_17 uint64 = 0
	var rtpHigh_hit_18 uint64 = 0
	var rtpHigh_hit_19 uint64 = 0
	var rtpHigh_hit_20 uint64 = 0
	var rtpHigh_hit_21 uint64 = 0
	var rtpHigh_hit_22 uint64 = 0
	var rtpHigh_hit_23 uint64 = 0
	var rtpHigh_hit_24 uint64 = 0
	var rtpHigh_hit_25 uint64 = 0
	var rtpHigh_hit_26 uint64 = 0
	var rtpHigh_hit_27 uint64 = 0
	var rtpHigh_hit_28 uint64 = 0
	var rtpHigh_hit_29 uint64 = 0

	var rtpLow_redenvelop_0 uint64 = 0
	var rtpLow_redenvelop_1 uint64 = 0
	var rtpLow_redenvelop_2 uint64 = 0
	var rtpLow_redenvelop_3 uint64 = 0
	var rtpLow_redenvelop_4 uint64 = 0
	var rtpLow_redenvelop_5 uint64 = 0
	var rtpLow_redenvelop_6 uint64 = 0
	var rtpLow_redenvelop_7 uint64 = 0
	var rtpLow_redenvelop_8 uint64 = 0
	var rtpLow_redenvelop_9 uint64 = 0
	var rtpLow_redenvelop_10 uint64 = 0
	var rtpLow_redenvelop_11 uint64 = 0
	var rtpLow_redenvelop_12 uint64 = 0
	var rtpLow_redenvelop_13 uint64 = 0
	var rtpLow_redenvelop_14 uint64 = 0
	var rtpLow_redenvelop_15 uint64 = 0
	var rtpLow_redenvelop_16 uint64 = 0
	var rtpLow_redenvelop_17 uint64 = 0
	var rtpLow_redenvelop_18 uint64 = 0
	var rtpLow_redenvelop_19 uint64 = 0

	var rtpLow_redenvelop_21 uint64 = 0
	var rtpLow_redenvelop_22 uint64 = 0
	var rtpLow_redenvelop_23 uint64 = 0
	var rtpLow_redenvelop_24 uint64 = 0
	var rtpLow_redenvelop_25 uint64 = 0

	var rtpHigh_redenvelop_0 uint64 = 0
	var rtpHigh_redenvelop_1 uint64 = 0
	var rtpHigh_redenvelop_2 uint64 = 0
	var rtpHigh_redenvelop_3 uint64 = 0
	var rtpHigh_redenvelop_4 uint64 = 0
	var rtpHigh_redenvelop_5 uint64 = 0
	var rtpHigh_redenvelop_6 uint64 = 0
	var rtpHigh_redenvelop_7 uint64 = 0
	var rtpHigh_redenvelop_8 uint64 = 0
	var rtpHigh_redenvelop_9 uint64 = 0
	var rtpHigh_redenvelop_10 uint64 = 0
	var rtpHigh_redenvelop_11 uint64 = 0
	var rtpHigh_redenvelop_12 uint64 = 0
	var rtpHigh_redenvelop_13 uint64 = 0
	var rtpHigh_redenvelop_14 uint64 = 0
	var rtpHigh_redenvelop_15 uint64 = 0
	var rtpHigh_redenvelop_16 uint64 = 0
	var rtpHigh_redenvelop_17 uint64 = 0
	var rtpHigh_redenvelop_18 uint64 = 0
	var rtpHigh_redenvelop_19 uint64 = 0

	var rtpHigh_redenvelop_21 uint64 = 0
	var rtpHigh_redenvelop_22 uint64 = 0
	var rtpHigh_redenvelop_23 uint64 = 0
	var rtpHigh_redenvelop_24 uint64 = 0
	var rtpHigh_redenvelop_25 uint64 = 0

	fmt.Print("MathModuleId(95、96、97、98), RunTimes(x 萬), FishId -> ")
	var rtpId, run_times int
	var inputFishId string
	fmt.Scanf("%d, %d, %s", &rtpId, &run_times, &inputFishId)

	switch rtpId {
	case 95:
		mathModuleId_00001 = models.PSFM_00002_95_1
	case 96:
		mathModuleId_00001 = models.PSFM_00002_96_1
	case 97:
		mathModuleId_00001 = models.PSFM_00002_97_1
	case 98:
		mathModuleId_00001 = models.PSFM_00002_98_1
	}
	run_times = run_times * 10000

	for i := 0; i < run_times; i++ {
		rtpId := Service.RtpId(secWebSocketKey, gameId_00001, -1)

		Service.Decrease(gameId_00001, -1, mathModuleId_00001, secWebSocketKey, BET)

		rtpState := Service.RtpState(secWebSocketKey, gameId_00001, subgameId)

		var fishId int32
		if inputFishId != "" {
			tempFish, _ := strconv.Atoi(inputFishId)
			fishId = int32(tempFish)
		} else {
			fishId = rngFishId()
		}

		result := probability.Service.Calc(
			gameId_00001,
			mathModuleId_00001,
			rtpId,
			fishId,
			-1,
			Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
			0,
		)

		switch fishId {
		case 0:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_0 += uint64(result.Pay)
					rtpHigh_win_0++
				}

				rtpHigh_hit_0++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_0 += uint64(result.Pay)
					rtpLow_win_0++
				}

				rtpLow_hit_0++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_0 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_0 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 1:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_1 += uint64(result.Pay)
					rtpHigh_win_1++
				}

				rtpHigh_hit_1++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_1 += uint64(result.Pay)
					rtpLow_win_1++
				}

				rtpLow_hit_1++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_1 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_1 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 2:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_2 += uint64(result.Pay)
					rtpHigh_win_2++
				}

				rtpHigh_hit_2++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_2 += uint64(result.Pay)
					rtpLow_win_2++
				}

				rtpLow_hit_2++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_2 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_2 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 3:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_3 += uint64(result.Pay)
					rtpHigh_win_3++
				}

				rtpHigh_hit_3++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_3 += uint64(result.Pay)
					rtpLow_win_3++
				}

				rtpLow_hit_3++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_3 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_3 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 4:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_4 += uint64(result.Pay)
					rtpHigh_win_4++
				}

				rtpHigh_hit_4++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_4 += uint64(result.Pay)
					rtpLow_win_4++
				}

				rtpLow_hit_4++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_4 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_4 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 5:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_5 += uint64(result.Pay)
					rtpHigh_win_5++
				}

				rtpHigh_hit_5++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_5 += uint64(result.Pay)
					rtpLow_win_5++
				}

				rtpLow_hit_5++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_5 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_5 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 6:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_6 += uint64(result.Pay)
					rtpHigh_win_6++
				}

				rtpHigh_hit_6++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_6 += uint64(result.Pay)
					rtpLow_win_6++
				}

				rtpLow_hit_6++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_6 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_6 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 7:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_7 += uint64(result.Pay)
					rtpHigh_win_7++
				}

				rtpHigh_hit_7++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_7 += uint64(result.Pay)
					rtpLow_win_7++
				}

				rtpLow_hit_7++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_7 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_7 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 8:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_8 += uint64(result.Pay)
					rtpHigh_win_8++
				}

				rtpHigh_hit_8++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_8 += uint64(result.Pay)
					rtpLow_win_8++
				}

				rtpLow_hit_8++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_8 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_8 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 9:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_9 += uint64(result.Pay)
					rtpHigh_win_9++
				}

				rtpHigh_hit_9++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_9 += uint64(result.Pay)
					rtpLow_win_9++
				}

				rtpLow_hit_9++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_9 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_9 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 10:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_10 += uint64(result.Pay)
					rtpHigh_win_10++
				}

				rtpHigh_hit_10++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_10 += uint64(result.Pay)
					rtpLow_win_10++
				}

				rtpLow_hit_10++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_10 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_10 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 11:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_11 += uint64(result.Pay)
					rtpHigh_win_11++
				}

				rtpHigh_hit_11++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_11 += uint64(result.Pay)
					rtpLow_win_11++
				}

				rtpLow_hit_11++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_11 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_11 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 12:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_12 += uint64(result.Pay)
					rtpHigh_win_12++
				}

				rtpHigh_hit_12++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_12 += uint64(result.Pay)
					rtpLow_win_12++
				}

				rtpLow_hit_12++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_12 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_12 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 13:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_13 += uint64(result.Pay)
					rtpHigh_win_13++
				}

				rtpHigh_hit_13++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_13 += uint64(result.Pay)
					rtpLow_win_13++
				}

				rtpLow_hit_13++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_13 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_13 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 14:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_14 += uint64(result.Pay)
					rtpHigh_win_14++
				}

				rtpHigh_hit_14++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_14 += uint64(result.Pay)
					rtpLow_win_14++
				}

				rtpLow_hit_14++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_14 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_14 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 15:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_15 += uint64(result.Pay)
					rtpHigh_win_15++
				}

				rtpHigh_hit_15++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_15 += uint64(result.Pay)
					rtpLow_win_15++
				}

				rtpLow_hit_15++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_15 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_15 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 16:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_16 += uint64(result.Pay)
					rtpHigh_win_16++
				}

				rtpHigh_hit_16++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_16 += uint64(result.Pay)
					rtpLow_win_16++
				}

				rtpLow_hit_16++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_16 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_16 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 17:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_17 += uint64(result.Pay)
					rtpHigh_win_17++
				}

				rtpHigh_hit_17++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_17 += uint64(result.Pay)
					rtpLow_win_17++
				}

				rtpLow_hit_17++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_17 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_17 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 18:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_18 += uint64(result.Pay)
					rtpHigh_win_18++
				}

				rtpHigh_hit_18++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_18 += uint64(result.Pay)
					rtpLow_win_18++
				}

				rtpLow_hit_18++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_18 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_18 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 19:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_19 += uint64(result.Pay)
					rtpHigh_win_19++
				}

				rtpHigh_hit_19++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_19 += uint64(result.Pay)
					rtpLow_win_19++
				}

				rtpLow_hit_19++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_19 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_19 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 20:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_20 += uint64(result.Pay) * uint64(result.Multiplier)
					rtpHigh_win_20++
				}

				rtpHigh_hit_20++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_20 += uint64(result.Pay) * uint64(result.Multiplier)
					rtpLow_win_20++
				}

				rtpLow_hit_20++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				t.Fatal("cannot trigger 31")
			}

		case 21:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_21 += uint64(result.Pay)
					rtpHigh_win_21++
				}

				rtpHigh_hit_21++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_21 += uint64(result.Pay)
					rtpLow_win_21++
				}

				rtpLow_hit_21++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_21 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_21 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 22:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_22 += uint64(result.Pay)
					rtpHigh_win_22++
				}

				rtpHigh_hit_22++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_22 += uint64(result.Pay)
					rtpLow_win_22++
				}

				rtpLow_hit_22++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_22 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_22 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 23:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_23 += uint64(result.Pay)
					rtpHigh_win_23++
				}

				rtpHigh_hit_23++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_23 += uint64(result.Pay)
					rtpLow_win_23++
				}

				rtpLow_hit_23++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_23 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_23 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 24:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_24 += uint64(result.Pay)
					rtpHigh_win_24++
				}

				rtpHigh_hit_24++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_24 += uint64(result.Pay)
					rtpLow_win_24++
				}

				rtpLow_hit_24++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_24 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_24 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 25:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_25 += uint64(result.Pay)
					rtpHigh_win_25++
				}

				rtpHigh_hit_25++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_25 += uint64(result.Pay)
					rtpLow_win_25++
				}

				rtpLow_hit_25++

			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				result := probability.Service.Calc(
					gameId_00001,
					mathModuleId_00001,
					strconv.Itoa(result.BonusTypeId),
					31,
					-1,
					Service.RtpBudget(secWebSocketKey, gameId_00001, -1, mathModuleId_00001),
					0,
				)

				switch rtpState {
				case high:
					rtpHigh_redenvelop_25 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				case low:
					rtpLow_redenvelop_25 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
				default:
					t.Fatal("RTP State Error")
				}
			}

		case 26:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_26 += uint64(result.Pay)
					rtpHigh_fish_26 += shot_26(gameId_00001, mathModuleId_00001, PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.Drill.RTP8.ID, result.BonusPayload.(int))
					rtpHigh_win_26++
				}

				rtpHigh_hit_26++

			case low:
				if result.Pay > 0 {
					rtpLow_fish_26 += uint64(result.Pay)
					rtpLow_fish_26 += shot_26(gameId_00001, mathModuleId_00001, PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.Drill.RTP8.ID, result.BonusPayload.(int))
					rtpLow_win_26++
				}

				rtpLow_hit_26++
			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				t.Fatal("cannot trigger 31")
			}

		case 27:
			switch rtpState {
			case high:
				if result.Pay > 0 {
					rtpHigh_fish_27 += uint64(result.Pay)
					rtpHigh_fish_27 += shot_27(gameId_00001, mathModuleId_00001, PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID, result.BonusPayload.(int))
					rtpHigh_win_27++
				}

				rtpHigh_hit_27++
			case low:
				if result.Pay > 0 {
					rtpLow_fish_27 += uint64(result.Pay)
					rtpLow_fish_27 += shot_27(gameId_00001, mathModuleId_00001, PSFM_00002_98_1.RTP98BS.PSF_ON_00001_2_BsMath.Icons.MachineGun.RTP8.ID, result.BonusPayload.(int))
					rtpLow_win_27++
				}

				rtpLow_hit_27++
			default:
				t.Fatal("RTP State Error")
			}

			if result.TriggerIconId == 31 {
				t.Fatal("cannot trigger 31")
			}

		case 28:
			switch rtpState {
			case high:
				if result.TriggerIconId == 28 {
					rtpHigh_fish_28 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
					rtpHigh_win_28++
				}
				rtpHigh_hit_28++
			case low:
				if result.TriggerIconId == 28 {
					rtpLow_fish_28 += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
					rtpLow_win_28++
				}
				rtpLow_hit_28++
			default:
				t.Fatal("RTP State Error")
			}

		case 29:

			switch rtpState {
			case high:
				if result.TriggerIconId == 29 {
					rtpHigh_fish_29 += uint64(result.ExtraData[0].(int))
					rtpHigh_win_29++
				}
				rtpHigh_hit_29++
			case low:
				if result.TriggerIconId == 29 {
					rtpLow_fish_29 += uint64(result.ExtraData[0].(int))
					rtpLow_win_29++
				}
				rtpLow_hit_29++
			default:
				t.Fatal("RTP State Error")
			}

		default:
			t.Fatal("Fish Id Error")
		}
	}

	rtpLow_win := rtpLow_win_0 + rtpLow_win_1 + rtpLow_win_2 + rtpLow_win_3 + rtpLow_win_4 + rtpLow_win_5 + rtpLow_win_6 + rtpLow_win_7 + rtpLow_win_8 +
		rtpLow_win_9 + rtpLow_win_10 + rtpLow_win_11 + rtpLow_win_12 + rtpLow_win_13 + rtpLow_win_14 + rtpLow_win_15 + rtpLow_win_16 + rtpLow_win_17 +
		rtpLow_win_18 + rtpLow_win_19 + rtpLow_win_20 + rtpLow_win_21 + rtpLow_win_22 + rtpLow_win_23 + rtpLow_win_24 + rtpLow_win_25 + rtpLow_win_26 +
		rtpLow_win_27 + rtpLow_win_28 + rtpLow_win_29

	rtpLow_hit := rtpLow_hit_0 + rtpLow_hit_1 + rtpLow_hit_2 + rtpLow_hit_3 + rtpLow_hit_4 + rtpLow_hit_5 + rtpLow_hit_6 + rtpLow_hit_7 + rtpLow_hit_8 +
		rtpLow_hit_9 + rtpLow_hit_10 + rtpLow_hit_11 + rtpLow_hit_12 + rtpLow_hit_13 + rtpLow_hit_14 + rtpLow_hit_15 + rtpLow_hit_16 + rtpLow_hit_17 +
		rtpLow_hit_18 + rtpLow_hit_19 + rtpLow_hit_20 + rtpLow_hit_21 + rtpLow_hit_22 + rtpLow_hit_23 + rtpLow_hit_24 + rtpLow_hit_25 + rtpLow_hit_26 +
		rtpLow_hit_27 + rtpLow_hit_28 + rtpLow_hit_29

	rtpHigh_win := rtpHigh_win_0 + rtpHigh_win_1 + rtpHigh_win_2 + rtpHigh_win_3 + rtpHigh_win_4 + rtpHigh_win_5 + rtpHigh_win_6 + rtpHigh_win_7 + rtpHigh_win_8 +
		rtpHigh_win_9 + rtpHigh_win_10 + rtpHigh_win_11 + rtpHigh_win_12 + rtpHigh_win_13 + rtpHigh_win_14 + rtpHigh_win_15 + rtpHigh_win_16 + rtpHigh_win_17 +
		rtpHigh_win_18 + rtpHigh_win_19 + rtpHigh_win_20 + rtpHigh_win_21 + rtpHigh_win_22 + rtpHigh_win_23 + rtpHigh_win_24 + rtpHigh_win_25 + rtpHigh_win_26 +
		rtpHigh_win_27 + rtpHigh_win_28 + rtpHigh_win_29

	rtpHigh_hit := rtpHigh_hit_0 + rtpHigh_hit_1 + rtpHigh_hit_2 + rtpHigh_hit_3 + rtpHigh_hit_4 + rtpHigh_hit_5 + rtpHigh_hit_6 + rtpHigh_hit_7 + rtpHigh_hit_8 +
		rtpHigh_hit_9 + rtpHigh_hit_10 + rtpHigh_hit_11 + rtpHigh_hit_12 + rtpHigh_hit_13 + rtpHigh_hit_14 + rtpHigh_hit_15 + rtpHigh_hit_16 + rtpHigh_hit_17 +
		rtpHigh_hit_18 + rtpHigh_hit_19 + rtpHigh_hit_20 + rtpHigh_hit_21 + rtpHigh_hit_22 + rtpHigh_hit_23 + rtpHigh_hit_24 + rtpHigh_hit_25 + rtpHigh_hit_26 +
		rtpHigh_hit_27 + rtpHigh_hit_28 + rtpHigh_hit_29

	rtpLow_fish := rtpLow_fish_0 + rtpLow_fish_1 + rtpLow_fish_2 + rtpLow_fish_3 + rtpLow_fish_4 + rtpLow_fish_5 + rtpLow_fish_6 + rtpLow_fish_7 + rtpLow_fish_8 +
		rtpLow_fish_9 + rtpLow_fish_10 + rtpLow_fish_11 + rtpLow_fish_12 + rtpLow_fish_13 + rtpLow_fish_14 + rtpLow_fish_15 + rtpLow_fish_16 + rtpLow_fish_17 +
		rtpLow_fish_18 + rtpLow_fish_19 + rtpLow_fish_20 + rtpLow_fish_21 + rtpLow_fish_22 + rtpLow_fish_23 + rtpLow_fish_24 + rtpLow_fish_25 + rtpLow_fish_26 +
		rtpLow_fish_27 + rtpLow_fish_28 + rtpLow_fish_29

	rtpHigh_fish := rtpHigh_fish_0 + rtpHigh_fish_1 + rtpHigh_fish_2 + rtpHigh_fish_3 + rtpHigh_fish_4 + rtpHigh_fish_5 + rtpHigh_fish_6 + rtpHigh_fish_7 + rtpHigh_fish_8 +
		rtpHigh_fish_9 + rtpHigh_fish_10 + rtpHigh_fish_11 + rtpHigh_fish_12 + rtpHigh_fish_13 + rtpHigh_fish_14 + rtpHigh_fish_15 + rtpHigh_fish_16 + rtpHigh_fish_17 +
		rtpHigh_fish_18 + rtpHigh_fish_19 + rtpHigh_fish_20 + rtpHigh_fish_21 + rtpHigh_fish_22 + rtpHigh_fish_23 + rtpHigh_fish_24 + rtpHigh_fish_25 + rtpHigh_fish_26 +
		rtpHigh_fish_27 + rtpHigh_fish_28 + rtpHigh_fish_29

	rtpLow_redenvelop := rtpLow_redenvelop_0 + rtpLow_redenvelop_1 + rtpLow_redenvelop_2 + rtpLow_redenvelop_3 + rtpLow_redenvelop_4 + rtpLow_redenvelop_5 +
		rtpLow_redenvelop_6 + rtpLow_redenvelop_7 + rtpLow_redenvelop_8 + rtpLow_redenvelop_9 + rtpLow_redenvelop_10 + rtpLow_redenvelop_11 + rtpLow_redenvelop_12 +
		rtpLow_redenvelop_13 + rtpLow_redenvelop_14 + rtpLow_redenvelop_15 + rtpLow_redenvelop_16 + rtpLow_redenvelop_17 + rtpLow_redenvelop_18 + rtpLow_redenvelop_19 +
		rtpLow_redenvelop_21 + rtpLow_redenvelop_22 + rtpLow_redenvelop_23 + rtpLow_redenvelop_24 + rtpLow_redenvelop_25

	rtpHigh_redenvelop := rtpHigh_redenvelop_0 + rtpHigh_redenvelop_1 + rtpHigh_redenvelop_2 + rtpHigh_redenvelop_3 + rtpHigh_redenvelop_4 + rtpHigh_redenvelop_5 +
		rtpHigh_redenvelop_6 + rtpHigh_redenvelop_7 + rtpHigh_redenvelop_8 + rtpHigh_redenvelop_9 + rtpHigh_redenvelop_10 + rtpHigh_redenvelop_11 + rtpHigh_redenvelop_12 +
		rtpHigh_redenvelop_13 + rtpHigh_redenvelop_14 + rtpHigh_redenvelop_15 + rtpHigh_redenvelop_16 + rtpHigh_redenvelop_17 + rtpHigh_redenvelop_18 + rtpHigh_redenvelop_19 +
		rtpHigh_redenvelop_21 + rtpHigh_redenvelop_22 + rtpHigh_redenvelop_23 + rtpHigh_redenvelop_24 + rtpHigh_redenvelop_25

	fmt.Println("FishId", ":", "LowHit", "LowTimes", "HighHit", "HighTimes", "LowWin", "HighWin", "LowRedenvelop", "HighRedenvelop",
		"LowHitRate", "LowPayRate", "HighHitRate", "HighPayRate", "TotalPayRate")
	fmt.Println(0, ":", rtpLow_win_0, rtpLow_hit_0, rtpHigh_win_0, rtpHigh_hit_0, rtpLow_fish_0, rtpHigh_fish_0, rtpLow_redenvelop_0, rtpHigh_redenvelop_0,
		getRate(rtpLow_win_0, rtpLow_hit_0), getRate(rtpLow_fish_0+rtpLow_redenvelop_0, rtpLow_hit_0),
		getRate(rtpHigh_win_0, rtpHigh_hit_0), getRate(rtpHigh_fish_0+rtpHigh_redenvelop_0, rtpHigh_hit_0),
		getRate(rtpLow_fish_0+rtpHigh_fish_0+rtpLow_redenvelop_0+rtpHigh_redenvelop_0, rtpLow_hit_0+rtpHigh_hit_0))
	fmt.Println(1, ":", rtpLow_win_1, rtpLow_hit_1, rtpHigh_win_1, rtpHigh_hit_1, rtpLow_fish_1, rtpHigh_fish_1, rtpLow_redenvelop_1, rtpHigh_redenvelop_1,
		getRate(rtpLow_win_1, rtpLow_hit_1), getRate(rtpLow_fish_1+rtpLow_redenvelop_1, rtpLow_hit_1),
		getRate(rtpHigh_win_1, rtpHigh_hit_1), getRate(rtpHigh_fish_1+rtpHigh_redenvelop_0, rtpHigh_hit_1),
		getRate(rtpLow_fish_1+rtpHigh_fish_1+rtpLow_redenvelop_1+rtpHigh_redenvelop_1, rtpLow_hit_1+rtpHigh_hit_1))
	fmt.Println(2, ":", rtpLow_win_2, rtpLow_hit_2, rtpHigh_win_2, rtpHigh_hit_2, rtpLow_fish_2, rtpHigh_fish_2, rtpLow_redenvelop_2, rtpHigh_redenvelop_2,
		getRate(rtpLow_win_2, rtpLow_hit_2), getRate(rtpLow_fish_2+rtpLow_redenvelop_2, rtpLow_hit_2),
		getRate(rtpHigh_win_2, rtpHigh_hit_2), getRate(rtpHigh_fish_2+rtpHigh_redenvelop_2, rtpHigh_hit_2),
		getRate(rtpLow_fish_2+rtpHigh_fish_2+rtpLow_redenvelop_2+rtpHigh_redenvelop_2, rtpLow_hit_2+rtpHigh_hit_2))
	fmt.Println(3, ":", rtpLow_win_3, rtpLow_hit_3, rtpHigh_win_3, rtpHigh_hit_3, rtpLow_fish_3, rtpHigh_fish_3, rtpLow_redenvelop_3, rtpHigh_redenvelop_3,
		getRate(rtpLow_win_3, rtpLow_hit_3), getRate(rtpLow_fish_3+rtpLow_redenvelop_3, rtpLow_hit_3),
		getRate(rtpHigh_win_3, rtpHigh_hit_3), getRate(rtpHigh_fish_3+rtpHigh_redenvelop_3, rtpHigh_hit_3),
		getRate(rtpLow_fish_3+rtpHigh_fish_3+rtpLow_redenvelop_3+rtpHigh_redenvelop_3, rtpLow_hit_3+rtpHigh_hit_3))
	fmt.Println(4, ":", rtpLow_win_4, rtpLow_hit_4, rtpHigh_win_4, rtpHigh_hit_4, rtpLow_fish_4, rtpHigh_fish_4, rtpLow_redenvelop_4, rtpHigh_redenvelop_4,
		getRate(rtpLow_win_4, rtpLow_hit_4), getRate(rtpLow_fish_4+rtpLow_redenvelop_4, rtpLow_hit_4),
		getRate(rtpHigh_win_4, rtpHigh_hit_4), getRate(rtpHigh_fish_4+rtpHigh_redenvelop_4, rtpHigh_hit_4),
		getRate(rtpLow_fish_4+rtpHigh_fish_4+rtpLow_redenvelop_4+rtpHigh_redenvelop_4, rtpLow_hit_4+rtpHigh_hit_4))
	fmt.Println(5, ":", rtpLow_win_5, rtpLow_hit_5, rtpHigh_win_5, rtpHigh_hit_5, rtpLow_fish_5, rtpHigh_fish_5, rtpLow_redenvelop_5, rtpHigh_redenvelop_5,
		getRate(rtpLow_win_5, rtpLow_hit_5), getRate(rtpLow_fish_5+rtpLow_redenvelop_5, rtpLow_hit_5),
		getRate(rtpHigh_win_5, rtpHigh_hit_5), getRate(rtpHigh_fish_5+rtpHigh_redenvelop_5, rtpHigh_hit_5),
		getRate(rtpLow_fish_5+rtpHigh_fish_5+rtpLow_redenvelop_5+rtpHigh_redenvelop_5, rtpLow_hit_5+rtpHigh_hit_5))
	fmt.Println(6, ":", rtpLow_win_6, rtpLow_hit_6, rtpHigh_win_6, rtpHigh_hit_6, rtpLow_fish_6, rtpHigh_fish_6, rtpLow_redenvelop_6, rtpHigh_redenvelop_6,
		getRate(rtpLow_win_6, rtpLow_hit_6), getRate(rtpLow_fish_6+rtpLow_redenvelop_6, rtpLow_hit_6),
		getRate(rtpHigh_win_6, rtpHigh_hit_6), getRate(rtpHigh_fish_6+rtpHigh_redenvelop_6, rtpHigh_hit_6),
		getRate(rtpLow_fish_6+rtpHigh_fish_6+rtpLow_redenvelop_6+rtpHigh_redenvelop_6, rtpLow_hit_6+rtpHigh_hit_6))
	fmt.Println(7, ":", rtpLow_win_7, rtpLow_hit_7, rtpHigh_win_7, rtpHigh_hit_7, rtpLow_fish_7, rtpHigh_fish_7, rtpLow_redenvelop_7, rtpHigh_redenvelop_7,
		getRate(rtpLow_win_7, rtpLow_hit_7), getRate(rtpLow_fish_7+rtpLow_redenvelop_7, rtpLow_hit_7),
		getRate(rtpHigh_win_7, rtpHigh_hit_7), getRate(rtpHigh_fish_7+rtpHigh_redenvelop_7, rtpHigh_hit_7),
		getRate(rtpLow_fish_7+rtpHigh_fish_7+rtpLow_redenvelop_7+rtpHigh_redenvelop_7, rtpLow_hit_7+rtpHigh_hit_7))
	fmt.Println(8, ":", rtpLow_win_8, rtpLow_hit_8, rtpHigh_win_8, rtpHigh_hit_8, rtpLow_fish_8, rtpHigh_fish_8, rtpLow_redenvelop_8, rtpHigh_redenvelop_8,
		getRate(rtpLow_win_8, rtpLow_hit_8), getRate(rtpLow_fish_8+rtpLow_redenvelop_8, rtpLow_hit_8),
		getRate(rtpHigh_win_8, rtpHigh_hit_8), getRate(rtpHigh_fish_8+rtpHigh_redenvelop_8, rtpHigh_hit_8),
		getRate(rtpLow_fish_8+rtpHigh_fish_8+rtpLow_redenvelop_8+rtpHigh_redenvelop_8, rtpLow_hit_8+rtpHigh_hit_8))
	fmt.Println(9, ":", rtpLow_win_9, rtpLow_hit_9, rtpHigh_win_9, rtpHigh_hit_9, rtpLow_fish_9, rtpHigh_fish_9, rtpLow_redenvelop_9, rtpHigh_redenvelop_9,
		getRate(rtpLow_win_9, rtpLow_hit_9), getRate(rtpLow_fish_9+rtpLow_redenvelop_9, rtpLow_hit_9),
		getRate(rtpHigh_win_9, rtpHigh_hit_9), getRate(rtpHigh_fish_9+rtpHigh_redenvelop_9, rtpHigh_hit_9),
		getRate(rtpLow_fish_9+rtpHigh_fish_9+rtpLow_redenvelop_9+rtpHigh_redenvelop_9, rtpLow_hit_9+rtpHigh_hit_9))

	fmt.Println(10, ":", rtpLow_win_10, rtpLow_hit_10, rtpHigh_win_10, rtpHigh_hit_10, rtpLow_fish_10, rtpHigh_fish_10, rtpLow_redenvelop_10, rtpHigh_redenvelop_10,
		getRate(rtpLow_win_10, rtpLow_hit_10), getRate(rtpLow_fish_10+rtpLow_redenvelop_10, rtpLow_hit_10),
		getRate(rtpHigh_win_10, rtpHigh_hit_10), getRate(rtpHigh_fish_10+rtpHigh_redenvelop_10, rtpHigh_hit_10),
		getRate(rtpLow_fish_10+rtpHigh_fish_10+rtpLow_redenvelop_10+rtpHigh_redenvelop_10, rtpLow_hit_10+rtpHigh_hit_10))
	fmt.Println(11, ":", rtpLow_win_11, rtpLow_hit_11, rtpHigh_win_11, rtpHigh_hit_11, rtpLow_fish_11, rtpHigh_fish_11, rtpLow_redenvelop_11, rtpHigh_redenvelop_11,
		getRate(rtpLow_win_11, rtpLow_hit_11), getRate(rtpLow_fish_11+rtpLow_redenvelop_11, rtpLow_hit_11),
		getRate(rtpHigh_win_11, rtpHigh_hit_11), getRate(rtpHigh_fish_11+rtpHigh_redenvelop_11, rtpHigh_hit_11),
		getRate(rtpLow_fish_11+rtpHigh_fish_11+rtpLow_redenvelop_11+rtpHigh_redenvelop_11, rtpLow_hit_11+rtpHigh_hit_11))
	fmt.Println(12, ":", rtpLow_win_12, rtpLow_hit_12, rtpHigh_win_12, rtpHigh_hit_12, rtpLow_fish_12, rtpHigh_fish_12, rtpLow_redenvelop_12, rtpHigh_redenvelop_12,
		getRate(rtpLow_win_12, rtpLow_hit_12), getRate(rtpLow_fish_12+rtpLow_redenvelop_12, rtpLow_hit_12),
		getRate(rtpHigh_win_12, rtpHigh_hit_12), getRate(rtpHigh_fish_12+rtpHigh_redenvelop_12, rtpHigh_hit_12),
		getRate(rtpLow_fish_12+rtpHigh_fish_12+rtpLow_redenvelop_12+rtpHigh_redenvelop_12, rtpLow_hit_12+rtpHigh_hit_12))
	fmt.Println(13, ":", rtpLow_win_13, rtpLow_hit_13, rtpHigh_win_13, rtpHigh_hit_13, rtpLow_fish_13, rtpHigh_fish_13, rtpLow_redenvelop_13, rtpHigh_redenvelop_13,
		getRate(rtpLow_win_13, rtpLow_hit_13), getRate(rtpLow_fish_13+rtpLow_redenvelop_13, rtpLow_hit_13),
		getRate(rtpHigh_win_13, rtpHigh_hit_13), getRate(rtpHigh_fish_13+rtpHigh_redenvelop_13, rtpHigh_hit_13),
		getRate(rtpLow_fish_13+rtpHigh_fish_13+rtpLow_redenvelop_13+rtpHigh_redenvelop_13, rtpLow_hit_13+rtpHigh_hit_13))
	fmt.Println(14, ":", rtpLow_win_14, rtpLow_hit_14, rtpHigh_win_14, rtpHigh_hit_14, rtpLow_fish_14, rtpHigh_fish_14, rtpLow_redenvelop_14, rtpHigh_redenvelop_14,
		getRate(rtpLow_win_14, rtpLow_hit_14), getRate(rtpLow_fish_14+rtpLow_redenvelop_14, rtpLow_hit_14),
		getRate(rtpHigh_win_14, rtpHigh_hit_14), getRate(rtpHigh_fish_14+rtpHigh_redenvelop_14, rtpHigh_hit_14),
		getRate(rtpLow_fish_14+rtpHigh_fish_14+rtpLow_redenvelop_14+rtpHigh_redenvelop_14, rtpLow_hit_14+rtpHigh_hit_14))
	fmt.Println(15, ":", rtpLow_win_15, rtpLow_hit_15, rtpHigh_win_15, rtpHigh_hit_15, rtpLow_fish_15, rtpHigh_fish_15, rtpLow_redenvelop_15, rtpHigh_redenvelop_15,
		getRate(rtpLow_win_15, rtpLow_hit_15), getRate(rtpLow_fish_15+rtpLow_redenvelop_15, rtpLow_hit_15),
		getRate(rtpHigh_win_15, rtpHigh_hit_15), getRate(rtpHigh_fish_15+rtpHigh_redenvelop_15, rtpHigh_hit_15),
		getRate(rtpLow_fish_15+rtpHigh_fish_15+rtpLow_redenvelop_15+rtpHigh_redenvelop_15, rtpLow_hit_15+rtpHigh_hit_15))
	fmt.Println(16, ":", rtpLow_win_16, rtpLow_hit_16, rtpHigh_win_16, rtpHigh_hit_16, rtpLow_fish_16, rtpHigh_fish_16, rtpLow_redenvelop_16, rtpHigh_redenvelop_16,
		getRate(rtpLow_win_16, rtpLow_hit_16), getRate(rtpLow_fish_16+rtpLow_redenvelop_16, rtpLow_hit_16),
		getRate(rtpHigh_win_16, rtpHigh_hit_16), getRate(rtpHigh_fish_16+rtpHigh_redenvelop_16, rtpHigh_hit_16),
		getRate(rtpLow_fish_16+rtpHigh_fish_16+rtpLow_redenvelop_16+rtpHigh_redenvelop_16, rtpLow_hit_16+rtpHigh_hit_16))
	fmt.Println(17, ":", rtpLow_win_17, rtpLow_hit_17, rtpHigh_win_17, rtpHigh_hit_17, rtpLow_fish_17, rtpHigh_fish_17, rtpLow_redenvelop_17, rtpHigh_redenvelop_17,
		getRate(rtpLow_win_17, rtpLow_hit_17), getRate(rtpLow_fish_17+rtpLow_redenvelop_17, rtpLow_hit_17),
		getRate(rtpHigh_win_17, rtpHigh_hit_17), getRate(rtpHigh_fish_17+rtpHigh_redenvelop_17, rtpHigh_hit_17),
		getRate(rtpLow_fish_17+rtpHigh_fish_17+rtpLow_redenvelop_17+rtpHigh_redenvelop_17, rtpLow_hit_17+rtpHigh_hit_17))
	fmt.Println(18, ":", rtpLow_win_18, rtpLow_hit_18, rtpHigh_win_18, rtpHigh_hit_18, rtpLow_fish_18, rtpHigh_fish_18, rtpLow_redenvelop_18, rtpHigh_redenvelop_18,
		getRate(rtpLow_win_18, rtpLow_hit_18), getRate(rtpLow_fish_18+rtpLow_redenvelop_18, rtpLow_hit_18),
		getRate(rtpHigh_win_18, rtpHigh_hit_18), getRate(rtpHigh_fish_18+rtpHigh_redenvelop_18, rtpHigh_hit_18),
		getRate(rtpLow_fish_18+rtpHigh_fish_18+rtpLow_redenvelop_18+rtpHigh_redenvelop_18, rtpLow_hit_18+rtpHigh_hit_18))
	fmt.Println(19, ":", rtpLow_win_19, rtpLow_hit_19, rtpHigh_win_19, rtpHigh_hit_19, rtpLow_fish_19, rtpHigh_fish_19, rtpLow_redenvelop_19, rtpHigh_redenvelop_19,
		getRate(rtpLow_win_19, rtpLow_hit_19), getRate(rtpLow_fish_19+rtpLow_redenvelop_19, rtpLow_hit_19),
		getRate(rtpHigh_win_19, rtpHigh_hit_19), getRate(rtpHigh_fish_19+rtpHigh_redenvelop_19, rtpHigh_hit_19),
		getRate(rtpLow_fish_19+rtpHigh_fish_19+rtpLow_redenvelop_19+rtpHigh_redenvelop_19, rtpLow_hit_19+rtpHigh_hit_19))

	fmt.Println(20, ":", rtpLow_win_20, rtpLow_hit_20, rtpHigh_win_20, rtpHigh_hit_20, rtpLow_fish_20, rtpHigh_fish_20,
		getRate(rtpLow_win_20, rtpLow_hit_20), getRate(rtpLow_fish_20, rtpLow_hit_20),
		getRate(rtpHigh_win_20, rtpHigh_hit_20), getRate(rtpHigh_fish_20, rtpHigh_hit_20),
		getRate(rtpLow_fish_20+rtpHigh_fish_20, rtpLow_hit_20+rtpHigh_hit_20))
	fmt.Println(21, ":", rtpLow_win_21, rtpLow_hit_21, rtpHigh_win_21, rtpHigh_hit_21, rtpLow_fish_21, rtpHigh_fish_21, rtpLow_redenvelop_21, rtpHigh_redenvelop_21,
		getRate(rtpLow_win_21, rtpLow_hit_21), getRate(rtpLow_fish_21+rtpLow_redenvelop_21, rtpLow_hit_21),
		getRate(rtpHigh_win_21, rtpHigh_hit_21), getRate(rtpHigh_fish_21+rtpHigh_redenvelop_21, rtpHigh_hit_21),
		getRate(rtpLow_fish_21+rtpHigh_fish_21+rtpLow_redenvelop_21+rtpHigh_redenvelop_21, rtpLow_hit_21+rtpHigh_hit_21))
	fmt.Println(22, ":", rtpLow_win_22, rtpLow_hit_22, rtpHigh_win_22, rtpHigh_hit_22, rtpLow_fish_22, rtpHigh_fish_22, rtpLow_redenvelop_22, rtpHigh_redenvelop_22,
		getRate(rtpLow_win_22, rtpLow_hit_22), getRate(rtpLow_fish_22+rtpLow_redenvelop_22, rtpLow_hit_22),
		getRate(rtpHigh_win_22, rtpHigh_hit_22), getRate(rtpHigh_fish_22+rtpHigh_redenvelop_22, rtpHigh_hit_22),
		getRate(rtpLow_fish_22+rtpHigh_fish_22+rtpLow_redenvelop_22+rtpHigh_redenvelop_22, rtpLow_hit_22+rtpHigh_hit_22))
	fmt.Println(23, ":", rtpLow_win_23, rtpLow_hit_23, rtpHigh_win_23, rtpHigh_hit_23, rtpLow_fish_23, rtpHigh_fish_23, rtpLow_redenvelop_23, rtpHigh_redenvelop_23,
		getRate(rtpLow_win_23, rtpLow_hit_23), getRate(rtpLow_fish_23+rtpLow_redenvelop_23, rtpLow_hit_23),
		getRate(rtpHigh_win_23, rtpHigh_hit_23), getRate(rtpHigh_fish_23+rtpHigh_redenvelop_23, rtpHigh_hit_23),
		getRate(rtpLow_fish_23+rtpHigh_fish_23+rtpLow_redenvelop_23+rtpHigh_redenvelop_23, rtpLow_hit_23+rtpHigh_hit_23))
	fmt.Println(24, ":", rtpLow_win_24, rtpLow_hit_24, rtpHigh_win_24, rtpHigh_hit_24, rtpLow_fish_24, rtpHigh_fish_24, rtpLow_redenvelop_24, rtpHigh_redenvelop_24,
		getRate(rtpLow_win_24, rtpLow_hit_24), getRate(rtpLow_fish_24+rtpLow_redenvelop_24, rtpLow_hit_24),
		getRate(rtpHigh_win_24, rtpHigh_hit_24), getRate(rtpHigh_fish_24+rtpHigh_redenvelop_24, rtpHigh_hit_24),
		getRate(rtpLow_fish_24+rtpHigh_fish_24+rtpLow_redenvelop_24+rtpHigh_redenvelop_24, rtpLow_hit_24+rtpHigh_hit_24))
	fmt.Println(25, ":", rtpLow_win_25, rtpLow_hit_25, rtpHigh_win_25, rtpHigh_hit_25, rtpLow_fish_25, rtpHigh_fish_25, rtpLow_redenvelop_25, rtpHigh_redenvelop_25,
		getRate(rtpLow_win_25, rtpLow_hit_25), getRate(rtpLow_fish_25+rtpLow_redenvelop_25, rtpLow_hit_25),
		getRate(rtpHigh_win_25, rtpHigh_hit_25), getRate(rtpHigh_fish_25+rtpHigh_redenvelop_25, rtpHigh_hit_25),
		getRate(rtpLow_fish_25+rtpHigh_fish_25+rtpLow_redenvelop_25+rtpHigh_redenvelop_25, rtpLow_hit_25+rtpHigh_hit_25))
	fmt.Println(26, ":", rtpLow_win_26, rtpLow_hit_26, rtpHigh_win_26, rtpHigh_hit_26, rtpLow_fish_26, rtpHigh_fish_26,
		getRate(rtpLow_win_26, rtpLow_hit_26), getRate(rtpLow_fish_26, rtpLow_hit_26),
		getRate(rtpHigh_win_26, rtpHigh_hit_26), getRate(rtpHigh_fish_26, rtpHigh_hit_26),
		getRate(rtpLow_fish_26+rtpHigh_fish_26, rtpLow_hit_26+rtpHigh_hit_26))
	fmt.Println(27, ":", rtpLow_win_27, rtpLow_hit_27, rtpHigh_win_27, rtpHigh_hit_27, rtpLow_fish_27, rtpHigh_fish_27,
		getRate(rtpLow_win_27, rtpLow_hit_27), getRate(rtpLow_fish_27, rtpLow_hit_27),
		getRate(rtpHigh_win_27, rtpHigh_hit_27), getRate(rtpHigh_fish_27, rtpHigh_hit_27),
		getRate(rtpLow_fish_27+rtpHigh_fish_27, rtpLow_hit_27+rtpHigh_hit_27))
	fmt.Println(28, ":", rtpLow_win_28, rtpLow_hit_28, rtpHigh_win_28, rtpHigh_hit_28, rtpLow_fish_28, rtpHigh_fish_28,
		getRate(rtpLow_win_28, rtpLow_hit_28), getRate(rtpLow_fish_28, rtpLow_hit_28),
		getRate(rtpHigh_win_28, rtpHigh_hit_28), getRate(rtpHigh_fish_28, rtpHigh_hit_28),
		getRate(rtpLow_fish_28+rtpHigh_fish_28, rtpLow_hit_28+rtpHigh_hit_28))
	fmt.Println(29, ":", rtpLow_win_29, rtpLow_hit_29, rtpHigh_win_29, rtpHigh_hit_29, rtpLow_fish_29, rtpHigh_fish_29,
		getRate(rtpLow_win_29, rtpLow_hit_29), getRate(rtpLow_fish_29, rtpLow_hit_29),
		getRate(rtpHigh_win_29, rtpHigh_hit_29), getRate(rtpHigh_fish_29, rtpHigh_hit_29),
		getRate(rtpLow_fish_29+rtpHigh_fish_29, rtpLow_hit_29+rtpHigh_hit_29))
	fmt.Println("Total", rtpLow_win, rtpLow_hit, rtpHigh_win, rtpHigh_hit, rtpLow_fish, rtpHigh_fish, rtpLow_redenvelop, rtpHigh_redenvelop,
		getRate(rtpLow_fish+rtpHigh_fish+rtpLow_redenvelop+rtpHigh_redenvelop, rtpLow_hit+rtpHigh_hit))
}
