package redenvelope

import (
	"serve/service_fish/domain/probability"
	"serve/service_fish/models"
	"testing"
)

const (
	TIMES uint64 = 10000000
)

func TestService_RngRedEnvelope0_4(t *testing.T) {
	var f int32 = 0

	for f = 0; f < 5; f++ {
		var i uint64 = 0
		var winCount uint64 = 0
		var bonusCount10 uint64 = 0
		var bonusCount12 uint64 = 0
		var bonusCount15 uint64 = 0
		var bonusCount18 uint64 = 0
		var bonusCount20 uint64 = 0
		var bonusCount22 uint64 = 0
		var bonusCount25 uint64 = 0
		var bonusCount28 uint64 = 0
		var bonusCount30 uint64 = 0
		var bonusCount50 uint64 = 0
		var bonusCount100 uint64 = 0
		var bonusCount200 uint64 = 0
		var bonusCount300 uint64 = 0
		var bonusCount500 uint64 = 0
		var bonusCount1000 uint64 = 0

		var bonusPay10 uint64 = 0
		var bonusPay12 uint64 = 0
		var bonusPay15 uint64 = 0
		var bonusPay18 uint64 = 0
		var bonusPay20 uint64 = 0
		var bonusPay22 uint64 = 0
		var bonusPay25 uint64 = 0
		var bonusPay28 uint64 = 0
		var bonusPay30 uint64 = 0
		var bonusPay50 uint64 = 0
		var bonusPay100 uint64 = 0
		var bonusPay200 uint64 = 0
		var bonusPay300 uint64 = 0
		var bonusPay500 uint64 = 0
		var bonusPay1000 uint64 = 0

		for i = 0; i < TIMES; i++ {
			result := probability.Service.Calc(models.PSF_ON_00001, models.PSFM_00001_98_1, "", 28, 0)

			if result.TriggerIconId == 28 {
				winCount++
				re := New(
					"",
					"",
					"",
					"",
					0,
					28,
					result.BonusPayload.([]int),
					0,
					0,
					0,
					0,
				)

				re.PlayerOptIndex = f

				switch Service.rngRedEnvelope(re).Pay {
				case 10:
					bonusPay10 = 10
					bonusCount10++
				case 12:
					bonusPay12 = 12
					bonusCount12++
				case 15:
					bonusPay15 = 15
					bonusCount15++
				case 18:
					bonusPay18 = 18
					bonusCount18++
				case 20:
					bonusPay20 = 20
					bonusCount20++
				case 22:
					bonusPay22 = 22
					bonusCount22++
				case 25:
					bonusPay25 = 25
					bonusCount25++
				case 28:
					bonusPay28 = 28
					bonusCount28++
				case 30:
					bonusPay30 = 30
					bonusCount30++
				case 50:
					bonusPay50 = 50
					bonusCount50++
				case 100:
					bonusPay100 = 100
					bonusCount100++
				case 200:
					bonusPay200 = 200
					bonusCount200++
				case 300:
					bonusPay300 = 300
					bonusCount300++
				case 500:
					bonusPay500 = 500
					bonusCount500++
				case 1000:
					bonusPay1000 = 1000
					bonusCount1000++
				}
			}
		}

		t.Log("index", f, " :", winCount,
			bonusCount10,
			bonusCount12,
			bonusCount15,
			bonusCount18,
			bonusCount20,
			bonusCount22,
			bonusCount25,
			bonusCount28,
			bonusCount30,
			bonusCount50,
			bonusCount100,
			bonusCount200,
			bonusCount300,
			bonusCount500,
			bonusCount1000,
		)

		t.Log("index", f, " :", winCount,
			bonusCount10*bonusPay10,
			bonusCount12*bonusPay12,
			bonusCount15*bonusPay15,
			bonusCount18*bonusPay18,
			bonusCount20*bonusPay20,
			bonusCount22*bonusPay22,
			bonusCount25*bonusPay25,
			bonusCount28*bonusPay28,
			bonusCount30*bonusPay30,
			bonusCount50*bonusPay50,
			bonusCount100*bonusPay100,
			bonusCount200*bonusPay200,
			bonusCount300*bonusPay300,
			bonusCount500*bonusPay500,
			bonusCount1000*bonusPay1000,
		)
	}
}
