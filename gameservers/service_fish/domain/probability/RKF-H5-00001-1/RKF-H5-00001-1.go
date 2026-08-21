package RKF_H5_00001_1

import (
	RKF_H5_00001_1_MONSTER "serve/service_fish/domain/probability/RKF-H5-00001-1/RKF-H5-00001-1-MONSTER"
)

type BsMath struct {
	ModuleType string
	BaseData   struct {
		ProjectID      string   `json:"ProjectID"`
		ModuleID       string   `json:"ModuleId"`
		GameType       string   `json:"GameType"`
		StripType      string   `json:"StripType"`
		GameTypeDetail string   `json:"GameTypeDetail"`
		BonusList      []string `json:"BonusList"`
	}
	Icons struct {
		Fish0                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_0"`
		Fish1                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_1"`
		Fish2                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_2"`
		Fish3                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_3"`
		Fish4                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_4"`
		Fish5                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_5"`
		Fish6                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_6"`
		Fish7                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_7"`
		Fish8                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_8"`
		Fish9                  RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_9"`
		Fish10                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_10"`
		Fish11                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_11"`
		Fish12                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_12"`
		Fish13                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_13"`
		Fish14                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_14"`
		Fish15                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_15"`
		Fish16                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_16"`
		Fish17                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_17"`
		Fish18                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_18"`
		Fish19                 RKF_H5_00001_1_MONSTER.BsFishMath                   `json:"Fish_19"`
		Yggdrasil              RKF_H5_00001_1_MONSTER.BsYggdrasilMath              `json:"YGGDRASIL"`
		HalonSlot              RKF_H5_00001_1_MONSTER.BsHalonSlotMath              `json:"HALON SLOT"`
		RandomTriggerYggdrasil RKF_H5_00001_1_MONSTER.BsRandomTriggerYggdrasilMath `json:"Random Trigger YGGDRASIL"`
		ThunderHammer          RKF_H5_00001_1_MONSTER.BsThunderHammerMath          `json:"THUNDER HAMMER"`
		Stormbreaker           RKF_H5_00001_1_MONSTER.BsThunderHammerMath          `json:"STORMBREAKER"`
		DeepSeaTreasure        RKF_H5_00001_1_MONSTER.BsDeepSeaTreasureMath        `json:"DEEP SEA TREASURE"`
		CasketOfAncientWinters RKF_H5_00001_1_MONSTER.BsDeepSeaTreasureMath        `json:"Casket of Ancient Winters"`
		PoseidonRandomChange   struct {
			IconID           int
			IsWildSubstitute bool
			BonusID          string
			BonusTimes       []int
			IconPays         []int
		} `json:"POSEIDON random change"`
		ChangeLink RKF_H5_00001_1_MONSTER.BsChangeLinkMath `json:"ChangeLink"`
		Poseidon1  RKF_H5_00001_1_MONSTER.BsPoseidonMath   `json:"POSEIDON 1"`
		Poseidon2  RKF_H5_00001_1_MONSTER.BsPoseidonMath   `json:"POSEIDON 2"`
		Poseidon3  RKF_H5_00001_1_MONSTER.BsPoseidonMath   `json:"POSEIDON 3"`
		Poseidon4  RKF_H5_00001_1_MONSTER.BsPoseidonMath   `json:"POSEIDON 4"`
		Poseidon5  RKF_H5_00001_1_MONSTER.BsPoseidonMath   `json:"POSEIDON 5"`
	}
}

type DrbMath struct {
	InitLowGroupWeight []int `json:"InitLowGroupWeight"`
	RTPGroupID         []int `json:"RTPGroupID"`
	HighRTPGroupWeight []int `json:"HighRTPGroupWeight"`
	LowRTPGroupWeight  []int `json:"LowRTPGroupWeight"`
	RTPGroup           struct {
		RTP1 struct{ HitDrbMath } `json:"RTP1"`
		RTP2 struct{ HitDrbMath } `json:"RTP2"`
		RTP3 struct{ HitDrbMath } `json:"RTP3"`
		RTP4 struct{ HitDrbMath } `json:"RTP4"`
		RTP5 struct{ HitDrbMath } `json:"RTP5"`
		RTP6 struct{ HitDrbMath } `json:"RTP6"`
		RTP7 struct{ HitDrbMath } `json:"RTP7"`
		RTP8 struct{ HitDrbMath } `json:"RTP8"`
		RTP9 struct{ HitDrbMath } `json:"RTP9"`
	} `json:"RTPGroup"`
}

type HitDrbMath struct {
	ID        string `json:"ID"`
	MinBullet int    `json:"MinBullet"`
	MaxBullet int    `json:"MaxBullet"`
}

type FsMath struct {
	ModuleType int `json:"ModuleType"`
	BaseData   struct {
		ProjectID      string        `json:"ProjectID"`
		ModuleID       string        `json:"ModuleID"`
		GameType       int           `json:"GameType"`
		StripType      int           `json:"StripType"`
		GameTypeDetail int           `json:"GameTypeDetail"`
		BonusList      []interface{} `json:"BonusList"`
	} `json:"BaseData"`
	StripTable struct {
		RKF_H5_00001_1_MONSTER.FsStripTableMath
	} `json:"StripTable"`
	Icons struct {
		Fish0                     RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_0"`
		Fish4                     RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_4"`
		Fish9                     RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_9"`
		Fish11                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_11"`
		Fish13                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_13"`
		Fish14                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_14"`
		Fish15                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_15"`
		Fish16                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_16"`
		Fish17                    RKF_H5_00001_1_MONSTER.FsFishMath `json:"Fish_17"`
		FaceOff1                  RKF_H5_00001_1_MONSTER.FsFishMath `json:"Face Off 1"`
		FaceOff3                  RKF_H5_00001_1_MONSTER.FsFishMath `json:"Face Off 3"`
		FaceOff5                  RKF_H5_00001_1_MONSTER.FsFishMath `json:"Face Off 5"`
		MysterySymbolRandomChange RKF_H5_00001_1_MONSTER.FsFishMath `json:"mystery symbol random change"`
		BigSymbol                 RKF_H5_00001_1_MONSTER.FsFishMath `json:"big symbol"`
	} `json:"Icons"`

	ChangeLink1 RKF_H5_00001_1_MONSTER.FsChangeLink1Math `json:"ChangeLink_1"`
	ChangeLink2 RKF_H5_00001_1_MONSTER.FsChangeLink2Math `json:"ChangeLink_2"`
	MixAward    []RKF_H5_00001_1_MONSTER.FsSlotMath      `json:"MixAward"`
}
