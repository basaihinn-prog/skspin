package PSF_ON_00007_1_MONSTER

import "serve/fish_comm/rng"

var RandomTriggerEnvelope = &randomTriggerEnvelope{}

type randomTriggerEnvelope struct {
}

type BsRandomTriggerEnvelopeMath struct {
	IconID     int                                      `json:"IconID"`
	BonusID    string                                   `json:"BonusID"`
	BonusTimes []int                                    `json:"BonusTimes"`
	HitWeight  []int                                    `json:"HitWeight"`
	Type1      struct{ HitBsRandomTriggerEnvelopeMath } `json:"Type1"`
	Type2      struct{ HitBsRandomTriggerEnvelopeMath } `json:"Type2"`
}

type HitBsRandomTriggerEnvelopeMath struct {
	LowIconPays        []int `json:"LowIconPays"`
	LowIconPaysWeight  []int `json:"LowIconPaysWeight"`
	HighIconPays       []int `json:"HighIconPays"`
	HighIconPaysWeight []int `json:"HighIconPaysWeight"`
}

func (r *randomTriggerEnvelope) Pick(rtpId int, math *BsRandomTriggerEnvelopeMath) (bonusTypeId int, iconPays []int) {
	lowPayOptions := make([]rng.Option, 0, 5)
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
		panic("Error:PSFM_00006_MONSTER:pick:iconPays")
	}

	return bonusTypeId, iconPays
}
