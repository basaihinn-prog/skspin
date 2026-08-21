package jilicomm

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"serve/comm/redisx"
)

type JILINextPlayRespStruct struct {
	Player             JILIPlayer
	IsBuy              bool
	App                *redisx.AppStore
	Bet                int64
	SelfPoolGold       int64
	HitBigAwardPercent int
}

func JIILNextPlayRespFunc(JILINextPlayRespS JILINextPlayRespStruct, JILINextPlayRespP any) (playResp any, hitBigAward, forcedKill, buyKill bool, err error) {
	v := reflect.ValueOf(&JILINextPlayRespP)

	if JILINextPlayRespS.IsBuy {
		minAwardPercent := JILINextPlayRespS.Player.BuyMinAwardPercent
		if minAwardPercent == 0 {
			minAwardPercent = JILINextPlayRespS.App.BuyMinAwardPercent
		}
		if rand.IntN(1000) < minAwardPercent {

			method := v.MethodByName("GetBuyMinData")
			result := method.Call(nil)

			playResp = result[0].Interface()
			buyKill = true

		} else {

			param1 := reflect.ValueOf(fmt.Sprintf("%d", JILINextPlayRespS.Bet))
			params := []reflect.Value{param1}

			method := v.MethodByName("NextBuy")
			result := method.Call(params)
			playResp = result[0].Interface()

		}
		return
	}

	nextData := redisx.GetPlayerNextData(JILINextPlayRespS.Player.PID)
	if nextData != nil {

		param1 := reflect.ValueOf(nextData)
		param2 := reflect.ValueOf(JILINextPlayRespS.App.GamePatten)
		params := []reflect.Value{param1, param2}

		method := v.MethodByName("ControlNextData")
		result := method.Call(params)
		playResp = result[0].Interface()

		return
	}

	noAwardPercent := JILINextPlayRespS.Player.NoAwardPercent
	if noAwardPercent == 0 {
		noAwardPercent = redisx.LoadNoAwardPercent(JILINextPlayRespS.Player.AppID)
	}
	if JILINextPlayRespS.SelfPoolGold < 0 && rand.IntN(1000) < noAwardPercent {

		param1 := reflect.ValueOf(JILINextPlayRespS.App.GamePatten)
		params := []reflect.Value{param1}

		method := v.MethodByName("SampleForceSimulate")
		result := method.Call(params)
		playResp = result[0].Interface()

		forcedKill = true
		return
	}

	if rand.IntN(1000) < JILINextPlayRespS.HitBigAwardPercent {
		if JILINextPlayRespS.Bet*10 < JILINextPlayRespS.SelfPoolGold {
			avg := JILINextPlayRespS.Player.CaclAvgBet()
			if avg*10 < JILINextPlayRespS.SelfPoolGold {

				param1 := reflect.ValueOf(JILINextPlayRespS.App.GamePatten)
				params := []reflect.Value{param1}

				method := v.MethodByName("GetBigReward")
				result := method.Call(params)
				playResp = result[0].Interface()

				hitBigAward = true
				return

			}
		}
	}

	if JILINextPlayRespS.Player.PersonWinMaxMult != 0 || JILINextPlayRespS.Player.PersonWinMaxScore != 0 {
		for {
			param1 := reflect.ValueOf(fmt.Sprintf("%d", JILINextPlayRespS.Bet))
			param2 := reflect.ValueOf(JILINextPlayRespS.App.GamePatten)
			params := []reflect.Value{param1, param2}

			method := v.MethodByName("Next")
			result := method.Call(params)
			playResp = result[0].Interface()

			Times := result[0].FieldByName("Times").Float()
			if Times < float64(JILINextPlayRespS.Player.PersonWinMaxMult) && Times*float64(JILINextPlayRespS.Bet) <= float64(JILINextPlayRespS.Player.PersonWinMaxScore)*10000 {
				return
			}
		}

	}

	return
}
