package slot

import (
	"serve/service_fish/domain/probability"
	"serve/service_fish/models"
	"testing"
)

const (
	TIMES uint64 = 10000000
)

func TestService_Slot(t *testing.T) {
	var winCount uint64 = 0
	var winPay uint64 = 0

	var bonusCount_0_5 uint64 = 0
	var bonusCount_6_10 uint64 = 0
	var bonusCount_11_50 uint64 = 0
	var bonusCount_51_100 uint64 = 0
	var bonusCount_101_500 uint64 = 0
	var bonusCount_501_1000 uint64 = 0
	var bonusCount_1001_up uint64 = 0

	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		result := probability.Service.Calc(models.PSF_ON_00001, models.PSFM_00001_98_1, "", 29, 0)

		if result.TriggerIconId == 29 {
			winCount++
			pay := uint64(result.ExtraData[0].(int))
			winPay += pay

			switch {
			case pay >= 0 && pay <= 5:
				bonusCount_0_5++

			case pay >= 6 && pay <= 10:
				bonusCount_6_10++

			case pay >= 11 && pay <= 50:
				bonusCount_11_50++

			case pay >= 51 && pay <= 100:
				bonusCount_51_100++

			case pay >= 101 && pay <= 500:
				bonusCount_101_500++

			case pay >= 501 && pay <= 1000:
				bonusCount_501_1000++

			case pay >= 1001:
				bonusCount_1001_up++

			default:
				t.Fatal("29 slot error pay")
			}
		}
	}

	t.Log(29, ":",
		winPay,
		winCount,
		bonusCount_0_5,
		bonusCount_6_10,
		bonusCount_11_50,
		bonusCount_51_100,
		bonusCount_101_500,
		bonusCount_501_1000,
		bonusCount_1001_up,
	)
}
