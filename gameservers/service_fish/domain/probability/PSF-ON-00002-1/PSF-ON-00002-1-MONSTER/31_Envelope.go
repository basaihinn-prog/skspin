package PSF_ON_00002_1_MONSTER

import "serve/fish_comm/rng"

var Envelope = &envelope{}

type envelope struct{}

type BsEnvelopeMath struct {
	IconID     int          `json:"IconID"`
	BonusID    string       `json:"BonusID"`
	BonusTimes []int        `json:"BonusTimes"`
	HitWeight  []int        `json:"HitWeight"`
	Type1      EnvelopeType `json:"Type1"`
	Type2      EnvelopeType `json:"Type2"`
}

type EnvelopeType struct {
	LowIconPays        []int `json:"LowIconPays"`
	LowIconPaysWeight  []int `json:"LowIconPaysWeight"`
	HighIconPays       []int `json:"HighIconPays"`
	HighIconPaysWeight []int `json:"HighIconPaysWeight"`
}

func (e *envelope) Pick(rtpId int, math *BsEnvelopeMath) (bonusTypeId int, iconPays []int) {
	lowPayOptions := make([]rng.Option, 0, 6)
	highPayOptions := make([]rng.Option, 0, 5)

	switch rtpId {
	case 2:
		bonusTypeId = 2

		for i := 0; i < len(math.Type2.LowIconPaysWeight); i++ {
			lowPayOptions = append(lowPayOptions, rng.Option{
				math.Type2.LowIconPaysWeight[i],
				math.Type2.LowIconPays[i],
			})
		}

		for i := 0; i < len(math.Type2.HighIconPaysWeight); i++ {
			highPayOptions = append(highPayOptions, rng.Option{
				math.Type2.HighIconPaysWeight[i],
				math.Type2.HighIconPays[i],
			})
		}

	case 1:
		fallthrough
	default:
		bonusTypeId = 1

		for i := 0; i < len(math.Type1.LowIconPaysWeight); i++ {
			lowPayOptions = append(lowPayOptions, rng.Option{
				math.Type1.LowIconPaysWeight[i],
				math.Type1.LowIconPays[i],
			})
		}

		for i := 0; i < len(math.Type1.HighIconPaysWeight); i++ {
			highPayOptions = append(highPayOptions, rng.Option{
				math.Type1.HighIconPaysWeight[i],
				math.Type1.HighIconPays[i],
			})
		}
	}

	iconPays = []int{
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(lowPayOptions).(int),
		MONSTER.rng(highPayOptions).(int),
	}

	if len(iconPays) != 5 {
		panic("Error:PSFM_00001_MONSTER:pick:iconPays")
	}

	return bonusTypeId, iconPays
}
