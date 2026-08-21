package PSF_ON_00003_1_MONSTER

import "serve/fish_comm/rng"

var SlotChanger = &slotChanger{}

type slotChanger struct {
}

type FsSlotChanger1Math struct {
	ChangeIconId          []int
	ChangeCandidateNormal struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock  struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Block"`
}

type FsChangeCandidate1Math struct {
	ChangeCandidate1  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_1"`
	ChangeCandidate2  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_2"`
	ChangeCandidate3  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_3"`
	ChangeCandidate4  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_4"`
	ChangeCandidate5  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_5"`
	ChangeCandidate6  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_6"`
	ChangeCandidate7  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_7"`
	ChangeCandidate8  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_8"`
	ChangeCandidate9  struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_9"`
	ChangeCandidate10 struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_10"`
	ChangeCandidate11 struct{ HitFsSlotChangerMath } `json:"ChangeCandidate_11"`
}

type HitFsSlotChangerMath struct {
	IconID int
	Weight int
}

func (s *slotChanger) rngCandidate1(math *struct{ FsChangeCandidate1Math }) (iconId int) {
	faceOff := make([]rng.Option, 0, 11)

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
	faceOff = append(faceOff, rng.Option{math.ChangeCandidate11.Weight, math.ChangeCandidate11.IconID})

	return MONSTER.rng(faceOff).(int)
}
