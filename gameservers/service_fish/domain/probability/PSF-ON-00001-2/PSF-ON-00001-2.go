package PSF_ON_00001_2

import PSF_ON_00001_2_MONSTER "serve/service_fish/domain/probability/PSF-ON-00001-2/PSF-ON-00001-2-MONSTER"

type BsMath struct {
	ModuleType string
	BaseData   struct {
		ProjectID      string `json:"ProjectID"`
		ModuleID       string
		GameType       string
		StripType      string
		GameTypeDetail string
		BonusList      []string
	}
	Icons struct {
		StarFish            PSF_ON_00001_2_MONSTER.BsFishMath
		SeaHorse            PSF_ON_00001_2_MONSTER.BsFishMath
		Guppy               PSF_ON_00001_2_MONSTER.BsFishMath
		ClownFish           PSF_ON_00001_2_MONSTER.BsFishMath
		Dory                PSF_ON_00001_2_MONSTER.BsFishMath
		Grammidae           PSF_ON_00001_2_MONSTER.BsFishMath
		RajahCichlasoma     PSF_ON_00001_2_MONSTER.BsFishMath `json:"Rajah Cichlasoma"`
		PufferFish          PSF_ON_00001_2_MONSTER.BsFishMath `json:"Puffer fish"`
		LanternFish         PSF_ON_00001_2_MONSTER.BsFishMath `json:"Lantern fish"`
		LionFish            PSF_ON_00001_2_MONSTER.BsFishMath
		Turtle              PSF_ON_00001_2_MONSTER.BsFishMath
		Lobster             PSF_ON_00001_2_MONSTER.BsFishMath
		Penguin             PSF_ON_00001_2_MONSTER.BsFishMath
		Platypus            PSF_ON_00001_2_MONSTER.BsFishMath
		Manatee             PSF_ON_00001_2_MONSTER.BsFishMath
		Dolphin             PSF_ON_00001_2_MONSTER.BsFishMath
		Koi                 PSF_ON_00001_2_MONSTER.BsFishMath
		HermitCrab          PSF_ON_00001_2_MONSTER.BsFishMath `json:"Hermit crab"`
		Shark               PSF_ON_00001_2_MONSTER.BsFishMath
		AsianArowana        PSF_ON_00001_2_MONSTER.BsFishMath    `json:"Asian arowana"`
		HotPot              PSF_ON_00001_2_MONSTER.BsHotPotMath  `json:"Hot Pot"`
		FaceOff1            PSF_ON_00001_2_MONSTER.BsFaceOffMath `json:"Face Off 1"`
		FaceOff2            PSF_ON_00001_2_MONSTER.BsFaceOffMath `json:"Face Off 2"`
		FaceOff3            PSF_ON_00001_2_MONSTER.BsFaceOffMath `json:"Face Off 3"`
		FaceOff4            PSF_ON_00001_2_MONSTER.BsFaceOffMath `json:"Face Off 4"`
		FaceOff5            PSF_ON_00001_2_MONSTER.BsFaceOffMath `json:"Face Off 5"`
		Drill               PSF_ON_00001_2_MONSTER.BsDrillMath
		MachineGun          PSF_ON_00001_2_MONSTER.BsMachineGunMath  `json:"Machine gun"`
		RedEnvelope         PSF_ON_00001_2_MONSTER.BsRedEnvelopeMath `json:"Red envelope"`
		Slot                PSF_ON_00001_2_MONSTER.BsSlotMath        `json:"slot"`
		FaceOffRandomChange struct {
			IconID           int
			IsWildSubstitute bool
			BonusID          string
			BonusTimes       []int
			IconPays         []int
		} `json:"Face Off random change"`
		Envelope PSF_ON_00001_2_MONSTER.BsEnvelopeMath `json:"Envelope"`
	}
	ChangeLink PSF_ON_00001_2_MONSTER.BsFaceOffChangerMath
}

type DrbMath struct {
	InitLowGroupWeight []int `json:"InitLowGroupWeight"`
	RTPGroupID         []int `json:"RTPGroupID"`
	HighRTPGroupWeight []int `json:"HighRTPGroupWeight"`
	LowRTPGroupWeight  []int `json:"LowRTPGroupWeight"`
	RTPGroup           struct {
		RTP1 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP1"`
		RTP2 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP2"`
		RTP3 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP3"`
		RTP4 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP4"`
		RTP5 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP5"`
		RTP6 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP6"`
		RTP7 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP7"`
		RTP8 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP8"`
		RTP9 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP9"`
		RTP10 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP10"`
		RTP11 struct {
			ID        string `json:"ID"`
			MinBullet int    `json:"MinBullet"`
			MaxBullet int    `json:"MaxBullet"`
		} `json:"RTP11"`
	} `json:"RTPGroup"`
}

type FsMath struct {
	ModuleType int
	BaseData   struct {
		ProjectID      string `json:"ProjectID"`
		ModuleID       string
		GameType       int
		StripType      int
		GameTypeDetail int
		BonusList      []interface{}
	}
	StripTable struct {
		PSF_ON_00001_2_MONSTER.FsReelMath
	}
	Icons struct {
		StarFish            PSF_ON_00001_2_MONSTER.FsFishMath
		Dory                PSF_ON_00001_2_MONSTER.FsFishMath
		LionFish            PSF_ON_00001_2_MONSTER.FsFishMath
		Lobster             PSF_ON_00001_2_MONSTER.FsFishMath
		Platypus            PSF_ON_00001_2_MONSTER.FsFishMath
		Manatee             PSF_ON_00001_2_MONSTER.FsFishMath
		Dolphin             PSF_ON_00001_2_MONSTER.FsFishMath
		Koi                 PSF_ON_00001_2_MONSTER.FsFishMath
		HermitCrab          PSF_ON_00001_2_MONSTER.FsFishMath    `json:"Hermit crab"`
		FaceOff1            PSF_ON_00001_2_MONSTER.FsFaceOffMath `json:"Face Off 1"`
		FaceOff3            PSF_ON_00001_2_MONSTER.FsFaceOffMath `json:"Face Off 3"`
		FaceOff5            PSF_ON_00001_2_MONSTER.FsFaceOffMath `json:"Face Off 5"`
		FaceOffRandomChange struct {
			IconID   int
			IconPays []int
		} `json:"Face Off random change"`
	}
	ChangeLink PSF_ON_00001_2_MONSTER.FsFaceOffChangerMath
	MixAward   []PSF_ON_00001_2_MONSTER.FsSlotMath
}
