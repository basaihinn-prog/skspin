package probability

import (
	"serve/service_fish/models"
	"testing"
)

const (
	TIMES          uint64 = 10000000
	GAME_ID               = models.PSF_ON_00001
	MATH_MODULE_ID        = models.PSFM_00001_98_1
	RTP_ID                = ""
)

func TestService_Calc_0_19(t *testing.T) {
	var f int32 = 0

	for f = 0; f < 20; f++ {
		var winPay uint64 = 0
		var winCount uint64 = 0
		var bonusCount uint64 = 0
		var i uint64 = 0

		for i = 0; i < TIMES; i++ {
			result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, f)

			if result.Pay > 0 {
				winPay = uint64(result.Pay)
				winCount++
			}

			if result.TriggerIconId == 31 {
				bonusCount++
			}
		}
		t.Log(f, ":", winCount, winCount*winPay, bonusCount)
	}
}

func TestService_Calc_20(t *testing.T) {
	var winCount uint64 = 0
	var bonusCount50 uint64 = 0
	var bonusCount60 uint64 = 0
	var bonusCount70 uint64 = 0
	var bonusCount80 uint64 = 0
	var bonusCount90 uint64 = 0
	var bonusCount100 uint64 = 0
	var bonusCount110 uint64 = 0
	var bonusCount120 uint64 = 0
	var bonusCount130 uint64 = 0
	var bonusCount150 uint64 = 0

	var bonusPay50 uint64 = 0
	var bonusPay60 uint64 = 0
	var bonusPay70 uint64 = 0
	var bonusPay80 uint64 = 0
	var bonusPay90 uint64 = 0
	var bonusPay100 uint64 = 0
	var bonusPay110 uint64 = 0
	var bonusPay120 uint64 = 0
	var bonusPay130 uint64 = 0
	var bonusPay150 uint64 = 0

	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 20)

		if result.Pay > 0 {
			winCount++
		}

		switch result.Pay {
		case 50:
			bonusPay50 = uint64(result.Pay)
			bonusCount50++
		case 60:
			bonusPay60 = uint64(result.Pay)
			bonusCount60++
		case 70:
			bonusPay70 = uint64(result.Pay)
			bonusCount70++
		case 80:
			bonusPay80 = uint64(result.Pay)
			bonusCount80++
		case 90:
			bonusPay90 = uint64(result.Pay)
			bonusCount90++
		case 100:
			bonusPay100 = uint64(result.Pay)
			bonusCount100++
		case 110:
			bonusPay110 = uint64(result.Pay)
			bonusCount110++
		case 120:
			bonusPay120 = uint64(result.Pay)
			bonusCount120++
		case 130:
			bonusPay130 = uint64(result.Pay)
			bonusCount130++
		case 150:
			bonusPay150 = uint64(result.Pay)
			bonusCount150++
		}

	}
	t.Log(20, ":", winCount,
		bonusCount50,
		bonusCount60,
		bonusCount70,
		bonusCount80,
		bonusCount90,
		bonusCount100,
		bonusCount110,
		bonusCount120,
		bonusCount130,
		bonusCount150,
	)
	t.Log(20, ":", winCount,
		bonusCount50*bonusPay50,
		bonusCount60*bonusPay60,
		bonusCount70*bonusPay70,
		bonusCount80*bonusPay80,
		bonusCount90*bonusPay90,
		bonusCount100*bonusPay100,
		bonusCount110*bonusPay110,
		bonusCount120*bonusPay120,
		bonusCount130*bonusPay130,
		bonusCount150*bonusPay150,
	)
}

func TestService_Calc_21_25(t *testing.T) {
	var winCount1 uint64 = 0
	var winCount2 uint64 = 0
	var winCount3 uint64 = 0
	var winCount4 uint64 = 0
	var winCount5 uint64 = 0

	var winPay1 uint64 = 0
	var winPay2 uint64 = 0
	var winPay3 uint64 = 0
	var winPay4 uint64 = 0
	var winPay5 uint64 = 0

	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		p1 := uint64(Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 21).Pay)
		if p1 > 0 {
			winPay1 = p1
			winCount1++
		}

		p2 := uint64(Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 22).Pay)
		if p2 > 0 {
			winPay2 = p2
			winCount2++
		}

		p3 := uint64(Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 23).Pay)
		if p3 > 0 {
			winPay3 = p3
			winCount3++
		}

		p4 := uint64(Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 24).Pay)
		if p4 > 0 {
			winPay4 = p4
			winCount4++
		}

		p5 := uint64(Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 25).Pay)
		if p5 > 0 {
			winPay5 = p5
			winCount5++
		}

	}
	t.Log(21, ":", winCount1, winCount1*winPay1)
	t.Log(22, ":", winCount2, winCount2*winPay2)
	t.Log(23, ":", winCount3, winCount3*winPay3)
	t.Log(24, ":", winCount4, winCount4*winPay4)
	t.Log(25, ":", winCount5, winCount5*winPay5)
}

func TestService_Calc_26(t *testing.T) {
	var winCount uint64 = 0
	var winPay uint64 = 0

	var bonusCount40 uint64 = 0
	var bonusCount50 uint64 = 0
	var bonusCount60 uint64 = 0
	var bonusCount70 uint64 = 0
	var bonusCount80 uint64 = 0
	var bonusCount90 uint64 = 0
	var bonusCount100 uint64 = 0

	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 26)

		if result.Pay > 0 {
			winPay = uint64(result.Pay)
			winCount++
		}

		switch result.BonusPayload.(int) {
		case 40:
			bonusCount40++
		case 50:
			bonusCount50++
		case 60:
			bonusCount60++
		case 70:
			bonusCount70++
		case 80:
			bonusCount80++
		case 90:
			bonusCount90++
		case 100:
			bonusCount100++
		}

	}
	t.Log(26, ":", winCount, winCount*winPay,
		bonusCount40,
		bonusCount50,
		bonusCount60,
		bonusCount70,
		bonusCount80,
		bonusCount90,
		bonusCount100,
	)
}

func TestService_Calc_27(t *testing.T) {
	var winCount uint64 = 0
	var winPay uint64 = 0

	var bonusCount20 uint64 = 0
	var bonusCount25 uint64 = 0
	var bonusCount30 uint64 = 0
	var bonusCount35 uint64 = 0
	var bonusCount40 uint64 = 0
	var bonusCount45 uint64 = 0
	var bonusCount50 uint64 = 0
	var bonusCount60 uint64 = 0
	var bonusCount80 uint64 = 0
	var bonusCount100 uint64 = 0
	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 27)

		if result.Pay > 0 {
			winPay = uint64(result.Pay)
			winCount++
		}

		switch result.BonusPayload.(int) {
		case 20:
			bonusCount20++
		case 25:
			bonusCount25++
		case 30:
			bonusCount30++
		case 35:
			bonusCount35++
		case 40:
			bonusCount40++
		case 45:
			bonusCount45++
		case 50:
			bonusCount50++
		case 60:
			bonusCount60++
		case 80:
			bonusCount80++
		case 100:
			bonusCount100++
		}

	}
	t.Log(20, ":", winCount, winCount*winPay,
		bonusCount20,
		bonusCount25,
		bonusCount30,
		bonusCount35,
		bonusCount40,
		bonusCount45,
		bonusCount50,
		bonusCount60,
		bonusCount80,
		bonusCount100,
	)
}

func TestService_Calc_28(t *testing.T) {
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

	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 28)

		if result.TriggerIconId == 28 {
			winCount++

			for _, v := range result.BonusPayload.([]int) {
				switch v {
				case 10:
					bonusPay10 = uint64(v)
					bonusCount10++
				case 12:
					bonusPay12 = uint64(v)
					bonusCount12++
				case 15:
					bonusPay15 = uint64(v)
					bonusCount15++
				case 18:
					bonusPay18 = uint64(v)
					bonusCount18++
				case 20:
					bonusPay20 = uint64(v)
					bonusCount20++
				case 22:
					bonusPay22 = uint64(v)
					bonusCount22++
				case 25:
					bonusPay25 = uint64(v)
					bonusCount25++
				case 28:
					bonusPay28 = uint64(v)
					bonusCount28++
				case 30:
					bonusPay30 = uint64(v)
					bonusCount30++
				case 50:
					bonusPay50 = uint64(v)
					bonusCount50++
				case 100:
					bonusPay100 = uint64(v)
					bonusCount100++
				case 200:
					bonusPay200 = uint64(v)
					bonusCount200++
				case 300:
					bonusPay300 = uint64(v)
					bonusCount300++
				case 500:
					bonusPay500 = uint64(v)
					bonusCount500++
				case 1000:
					bonusPay1000 = uint64(v)
					bonusCount1000++
				}
			}
		}
	}
	t.Log(28, ":", winCount,
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

	t.Log(28, ":", winCount,
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

func TestService_Calc_29(t *testing.T) {
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
		result := Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 29)

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

func TestService_Calc_30(t *testing.T) {
	var faceCount21 uint64 = 0
	var faceCount22 uint64 = 0
	var faceCount23 uint64 = 0
	var faceCount24 uint64 = 0
	var faceCount25 uint64 = 0
	var i uint64 = 0

	for i = 0; i < TIMES; i++ {
		switch Service.Calc(GAME_ID, MATH_MODULE_ID, RTP_ID, 30).FishTypeId {
		case 21:
			faceCount21++
		case 22:
			faceCount22++
		case 23:
			faceCount23++
		case 24:
			faceCount24++
		case 25:
			faceCount25++
		}
	}
	t.Log(30, ":",
		faceCount21,
		faceCount22,
		faceCount23,
		faceCount24,
		faceCount25,
	)
}
