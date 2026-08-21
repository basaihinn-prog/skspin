package PSF_ON_00005_1_MONSTER

import "serve/fish_comm/rng"

var GiantWhaleChanger = &giantWhaleChanger{}

type giantWhaleChanger struct {
}

type BsGiantWhaleChangerMath struct {
	ChangeIconId    []int `json:"ChangeIconId"`
	ChangeCandidate struct {
		ChangeCandidate1 struct{ HitBsGiantWhaleChangerMath } `json:"ChangeCandidate_1"`
		ChangeCandidate2 struct{ HitBsGiantWhaleChangerMath } `json:"ChangeCandidate_2"`
		ChangeCandidate3 struct{ HitBsGiantWhaleChangerMath } `json:"ChangeCandidate_3"`
		ChangeCandidate4 struct{ HitBsGiantWhaleChangerMath } `json:"ChangeCandidate_4"`
		ChangeCandidate5 struct{ HitBsGiantWhaleChangerMath } `json:"ChangeCandidate_5"`
	} `json:"ChangeCandidate"`
}

type HitBsGiantWhaleChangerMath struct {
	IconID int `json:"IconID"`
	Weight int `json:"Weight"`
}

type FsGiantWhaleChanger1Math struct {
	ChangeIconId          []int                            `json:"ChangeIconId"`
	ChangeCandidateNormal struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock  struct{ FsChangeCandidate1Math } `json:"ChangeCandidate_Block"`
}

type FsGiantWhaleChanger2Math struct {
	ChangeIconId          []int                            `json:"ChangeIconId"`
	ChangeCandidateNormal struct{ FsChangeCandidate2Math } `json:"ChangeCandidate_Normal"`
	ChangeCandidateBlock  struct{ FsChangeCandidate2Math } `json:"ChangeCandidate_Block"`
}

type FsChangeCandidate2Math struct {
	ChangeCandidate1 struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_1"`
	ChangeCandidate2 struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_2"`
	ChangeCandidate3 struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_3"`
}

type FsChangeCandidate1Math struct {
	ChangeCandidate1  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_1"`
	ChangeCandidate2  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_2"`
	ChangeCandidate3  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_3"`
	ChangeCandidate4  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_4"`
	ChangeCandidate5  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_5"`
	ChangeCandidate6  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_6"`
	ChangeCandidate7  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_7"`
	ChangeCandidate8  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_8"`
	ChangeCandidate9  struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_9"`
	ChangeCandidate10 struct{ ChangeFsGiantWhaleChangerMath } `json:"ChangeCandidate_10"`
}

type ChangeFsGiantWhaleChangerMath struct {
	IconID int `json:"IconID"`
	Weight int `json:"Weight"`
}

func (g *giantWhaleChanger) Hit(math *BsGiantWhaleChangerMath) (iconId int) {
	giantWhaleChange := make([]rng.Option, 0, 4)

	giantWhaleChange = append(giantWhaleChange, rng.Option{math.ChangeCandidate.ChangeCandidate1.Weight, math.ChangeCandidate.ChangeCandidate1.IconID})
	giantWhaleChange = append(giantWhaleChange, rng.Option{math.ChangeCandidate.ChangeCandidate2.Weight, math.ChangeCandidate.ChangeCandidate2.IconID})
	giantWhaleChange = append(giantWhaleChange, rng.Option{math.ChangeCandidate.ChangeCandidate3.Weight, math.ChangeCandidate.ChangeCandidate3.IconID})
	giantWhaleChange = append(giantWhaleChange, rng.Option{math.ChangeCandidate.ChangeCandidate4.Weight, math.ChangeCandidate.ChangeCandidate4.IconID})
	giantWhaleChange = append(giantWhaleChange, rng.Option{math.ChangeCandidate.ChangeCandidate5.Weight, math.ChangeCandidate.ChangeCandidate5.IconID})

	return MONSTER.rng(giantWhaleChange).(int)
}

func (g *giantWhaleChanger) rngCandidate1(math *struct{ FsChangeCandidate1Math }) (iconId int) {
	giantWhaleOff := make([]rng.Option, 0, 9)

	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate1.Weight, math.ChangeCandidate1.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate2.Weight, math.ChangeCandidate2.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate3.Weight, math.ChangeCandidate3.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate4.Weight, math.ChangeCandidate4.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate5.Weight, math.ChangeCandidate5.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate6.Weight, math.ChangeCandidate6.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate7.Weight, math.ChangeCandidate7.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate8.Weight, math.ChangeCandidate8.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate9.Weight, math.ChangeCandidate9.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate10.Weight, math.ChangeCandidate10.IconID})

	return MONSTER.rng(giantWhaleOff).(int)
}

func (g *giantWhaleChanger) rngCandidate2(math *struct{ FsChangeCandidate2Math }) (iconId int) {
	giantWhaleOff := make([]rng.Option, 0, 2)

	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate1.Weight, math.ChangeCandidate1.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate2.Weight, math.ChangeCandidate2.IconID})
	giantWhaleOff = append(giantWhaleOff, rng.Option{math.ChangeCandidate3.Weight, math.ChangeCandidate3.IconID})

	return MONSTER.rng(giantWhaleOff).(int)
}

func (g *giantWhaleChanger) HitFs(math *FsFishMath) (iconPays int) {
	return StarFish.HitFs(math)
}
