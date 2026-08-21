package rtp

import (
	"fmt"
	"os"
	"serve/service_fish/domain/probability"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	"serve/service_fish/models"
	"strconv"
	"testing"
)

const (
	game_id_20002         = models.PSF_ON_20002
	secWebSocketKey_20002 = "jerry"
	bet_20002             = 1
	trigger_icon_id_20002 = 31

	//mathModule_id_20002 = models.PSFM_00003_93_1
	//run_times_20002 = 100 * 100 * 100 * 100
)

var fishList20002 = []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
var rtpList_20002 = []string{"0", "20", "40", "60", "80", "901", "902", "100", "150", "200", "300"}

var mathModule_id_20002 string

func TestService_20002(t *testing.T) {
	times := make(map[string]uint64, 0)
	hits := make(map[string]uint64, 0)
	wins := make(map[string]uint64, 0)
	envelopes := make(map[string]uint64, 0)

	//lowTimes := make(map[int32]uint64, 0)
	//lowHits := make(map[int32]uint64, 0)
	//lowWins := make(map[int32]uint64, 0)
	//lowEnvelopes := make(map[int32]uint64, 0)
	//
	//highTimes := make(map[int32]uint64, 0)
	//highHits := make(map[int32]uint64, 0)
	//highWins := make(map[int32]uint64, 0)
	//highEnvelopes := make(map[int32]uint64, 0)

	runTimes20002 := 0
	var fishId int32 = -1
	randomFishId := true
	if len(os.Args) >= 3 {
		switch os.Args[1] {
		case "93":
			mathModule_id_20002 = models.PSFM_00003_93_1
		case "94":
			mathModule_id_20002 = models.PSFM_00003_94_1
		default:
			fmt.Println("MathModule is Error")
			return
		}

		inputTimes, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("RunTimes is Error")
			return
		}
		runTimes20002 = inputTimes * 10000

		if len(os.Args) == 4 {
			tempFishId, _ := strconv.Atoi(os.Args[3])
			fishId = int32(tempFishId)
			randomFishId = false
		}
	} else {
		fmt.Println("Format Error. Format : MathModule Times FishId(Optional)")
		return
	}

	for i := 0; i < runTimes20002; i++ {
		if randomFishId {
			fishId = rngFishIdNew(fishList20002)
		}

		//fishId := int32(0)

		Service.Decrease(game_id_20002, 0, mathModule_id_20002, secWebSocketKey_20002, bet_20002)
		//rtpState := Service.RtpState(secWebSocketKey_20002, game_id_20002, 0)
		rtpId := Service.RtpId(secWebSocketKey_20002, game_id_20002, 0)

		result := probability.Service.Calc(
			game_id_20002,
			mathModule_id_20002,
			rtpId,
			fishId,
			-1,
			0,
			0,
			0,
		)
		hit, win, envelope := processResult20002(fishId, result)

		mapKey := string(fishId) + rtpId
		times[mapKey]++
		hits[mapKey] += hit
		wins[mapKey] += win
		envelopes[mapKey] += envelope

		//switch rtpState {
		//case low:
		//	lowTimes[fishId]++
		//	lowHits[fishId] += hit
		//	lowWins[fishId] += win
		//	lowEnvelopes[fishId] += envelope
		//
		//case high:
		//	highTimes[fishId]++
		//	highHits[fishId] += hit
		//	highWins[fishId] += win
		//	highEnvelopes[fishId] += envelope
		//}
	}

	//for _, v := range fishList20002 {
	//	if lowTimes[v] != 0 || highTimes[v] != 0 {
	//		t.Log(v, ":", lowTimes[v], lowHits[v], lowWins[v], lowEnvelopes[v], highTimes[v], highHits[v], highWins[v], highEnvelopes[v],
	//			getRate(lowHits[v], lowTimes[v]), getRate(lowWins[v]+lowEnvelopes[v], lowTimes[v]),
	//			getRate(highHits[v], highTimes[v]), getRate(highWins[v]+highEnvelopes[v], highTimes[v]),
	//			getRate(lowWins[v]+lowEnvelopes[v]+highWins[v]+highEnvelopes[v], lowTimes[v]+highTimes[v]),
	//		)
	//	}
	//}

	title := fmt.Sprint("FishId", " :")
	for _, rtp := range rtpList_20002 {
		title = fmt.Sprint(title, " RTP"+rtp+"Times", " RTP"+rtp+"Hits", " RTP"+rtp+"Wins", " RTP"+rtp+"Envelope")
	}
	title = fmt.Sprint(title, " HitRate", " RtpRate")
	fmt.Println(title)

	var totalTimes uint64 = 0
	var totalWins uint64 = 0
	var totalEnvelope uint64 = 0

	for _, fish := range fishList20002 {
		resultStr := fmt.Sprint(fish) + " : "
		var fishTimes uint64 = 0
		var fishHits uint64 = 0
		var fishWins uint64 = 0
		var fishEnvelope uint64 = 0

		for _, rtp := range rtpList_20002 {
			mapKey := string(fish) + rtp
			resultStr += fmt.Sprint(times[mapKey], hits[mapKey], wins[mapKey], envelopes[mapKey], " ")

			fishTimes += times[mapKey]
			fishHits += hits[mapKey]
			fishWins += wins[mapKey]
			fishEnvelope += envelopes[mapKey]
		}

		totalTimes += fishTimes
		totalWins += fishWins
		totalEnvelope += fishEnvelope
		fmt.Println(resultStr, getRate(fishHits, fishTimes), getRate(fishWins+fishEnvelope, fishTimes))
	}
	fmt.Println("Total RTP", ":", getRate(totalWins+totalEnvelope, totalTimes))
}

func processResult20002(fishId int32, result *probability.Probability) (hits, wins, envelopes uint64) {
	hits = 0
	wins = 0
	envelopes = 0

	if result.Pay > 0 {
		hits++
		wins += uint64(result.Pay * result.Multiplier)
	}

	if result.TriggerIconId == trigger_icon_id_20002 {
		envelopes += bonusEnvelope20002(result.BonusTypeId, trigger_icon_id_20002)
	}

	// Drill
	if fishId == 26 && result.Pay > 0 {
		wins += shot_26(game_id_20002, mathModule_id_20002, PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath.Icons.Drill.RTP8.ID, result.BonusPayload.(int))
	}

	// MachineGun
	if fishId == 27 && result.Pay > 0 {
		wins += shot_27(game_id_20002, mathModule_id_20002, PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath.Icons.MachineGun.RTP8.ID, result.BonusPayload.(int))
	}

	// RedEnvelope
	if fishId == 28 && result.TriggerIconId == 28 {
		hits++
		wins += uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
	}

	// Slot
	if fishId == 29 && result.TriggerIconId == 29 {
		hits++
		wins += uint64(result.ExtraData[0].(int))
	}

	return hits, wins, envelopes
}

func bonusEnvelope20002(bonusTypeId int, triggerIconId int32) uint64 {
	result := probability.Service.Calc(
		game_id_20002,
		mathModule_id_20002,
		strconv.Itoa(bonusTypeId),
		triggerIconId,
		-1,
		0,
		0,
		0,
	)

	return uint64(result.BonusPayload.([]int)[rngRedEnvelope()])
}
