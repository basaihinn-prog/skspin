package PSF_ON_00004_1_MONSTER

import "serve/fish_comm/rng"

var FaceOffChanger = &faceOffChanger{}

type faceOffChanger struct {
}

type FsFaceOffChanger1Math struct {
	ChangeIconId          []int
	ChangeCandidateNormal struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock  struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Block"`
}

type FsFaceOffChanger2Math struct {
	ChangeIconId          []int
	ChangeCandidateNormal struct{ FsChangeCandidate2Math } `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock  struct{ FsChangeCandidate2Math } `json:"ChangeCandidate_Block"`
}

type FsChangeCandidate1Math struct {
	ChangeCandidate1  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_1"`
	ChangeCandidate2  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_2"`
	ChangeCandidate3  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_3"`
	ChangeCandidate4  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_4"`
	ChangeCandidate5  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_5"`
	ChangeCandidate6  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_6"`
	ChangeCandidate7  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_7"`
	ChangeCandidate8  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_8"`
	ChangeCandidate9  struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_9"`
	ChangeCandidate10 struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_10"`
}

type FsChangeCandidate2Math struct {
	ChangeCandidate1 struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_1"`
	ChangeCandidate2 struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_2"`
	ChangeCandidate3 struct{ HitFsFaceOffChangerMath } `json:"ChangeCandidate_3"`
}

type HitFsFaceOffChangerMath struct {
	IconID int
	Weight int
}

func (f *faceOffChanger) rngCandidate1(math *struct{ FsChangeCandidate1Math }) (iconId int) {
	faceOff := make([]rng.Option, 0, 9)

	faceOff = append(faceOff, rng.Option{math.ChangeCandidate1.Weight, math.ChangeCandidate1.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate2.Weight, math.ChangeCandidate2.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate3.Weight, math.ChangeCandidate3.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate4.Weight, math.ChangeCandidate4.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate5.Weight, math.ChangeCandidate5.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate6.Weight, math.ChangeCandidate6.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate7.Weight, math.ChangeCandidate7.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate8.Weight, math.ChangeCandidate8.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate9.Weight, math.ChangeCandidate9.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate10.Weight, math.ChangeCandidate10.IconID})

	return MONSTER.rng(faceOff).(int)
}

func (f *faceOffChanger) rngCandidate2(math *struct{ FsChangeCandidate2Math }) (iconId int) {
	faceOff := make([]rng.Option, 0, 2)

	faceOff = append(faceOff, rng.Option{math.ChangeCandidate1.Weight, math.ChangeCandidate1.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate2.Weight, math.ChangeCandidate2.IconID})
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate3.Weight, math.ChangeCandidate3.IconID})

	return MONSTER.rng(faceOff).(int)
}
