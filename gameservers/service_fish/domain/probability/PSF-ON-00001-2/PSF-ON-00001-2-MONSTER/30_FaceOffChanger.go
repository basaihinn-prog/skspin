package PSF_ON_00001_2_MONSTER

import "serve/fish_comm/rng"

var FaceOffChanger = &faceOffChanger{}

type faceOffChanger struct{}

type BsFaceOffChangerMath struct {
	ChangeIconId    []int
	ChangeCandidate struct {
		ChangeCandidate1 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_1"`
		ChangeCandidate2 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_2"`
		ChangeCandidate3 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_3"`
		ChangeCandidate4 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_4"`
		ChangeCandidate5 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_5"`
	}
}

type FsFaceOffChangerMath struct {
	ChangeIconId          []int
	ChangeCandidateNormal struct {
		ChangeCandidate1 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_1"`
		ChangeCandidate2 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_2"`
		ChangeCandidate3 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_3"`
	} `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock struct {
		ChangeCandidate1 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_1"`
		ChangeCandidate2 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_2"`
		ChangeCandidate3 struct {
			IconID int
			Weight int
		} `json:"ChangeCandidate_3"`
	} `json:"ChangeCandidate_Block"`
}

func (f *faceOffChanger) Hit(math *BsFaceOffChangerMath) (iconId int) {
	faceOff := make([]rng.Option, 0, 4)

	faceOff = append(faceOff, rng.Option{math.ChangeCandidate.ChangeCandidate1.Weight, math.ChangeCandidate.ChangeCandidate1.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate.ChangeCandidate2.Weight, math.ChangeCandidate.ChangeCandidate2.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate.ChangeCandidate3.Weight, math.ChangeCandidate.ChangeCandidate3.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate.ChangeCandidate4.Weight, math.ChangeCandidate.ChangeCandidate4.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate.ChangeCandidate5.Weight, math.ChangeCandidate.ChangeCandidate5.IconID})

	return MONSTER.rng(faceOff).(int)
}

func (f *faceOffChanger) rngNormalReel(math *FsFaceOffChangerMath) (iconId int) {
	faceOff := make([]rng.Option, 0, 2)

	faceOff = append(faceOff, rng.Option{math.ChangeCandidateNormal.ChangeCandidate1.Weight, math.ChangeCandidateNormal.ChangeCandidate1.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidateNormal.ChangeCandidate2.Weight, math.ChangeCandidateNormal.ChangeCandidate2.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidateNormal.ChangeCandidate3.Weight, math.ChangeCandidateNormal.ChangeCandidate3.IconID})

	return MONSTER.rng(faceOff).(int)
}

func (f *faceOffChanger) rngBlockReel(math *FsFaceOffChangerMath) (iconId int) {
	faceOff := make([]rng.Option, 0, 2)

	faceOff = append(faceOff, rng.Option{math.ChangeCandidateBlock.ChangeCandidate1.Weight, math.ChangeCandidateBlock.ChangeCandidate1.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidateBlock.ChangeCandidate2.Weight, math.ChangeCandidateBlock.ChangeCandidate2.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidateBlock.ChangeCandidate3.Weight, math.ChangeCandidateBlock.ChangeCandidate3.IconID})

	return MONSTER.rng(faceOff).(int)
}
