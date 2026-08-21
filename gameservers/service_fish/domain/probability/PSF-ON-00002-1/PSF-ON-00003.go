package PSF_ON_00002_1

import (
	PSF_ON_00002_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00002-1/PSF-ON-00002-1-MONSTER"
)

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
		StarFish            PSF_ON_00002_1_MONSTER.BsFishMath
		SeaHorse            PSF_ON_00002_1_MONSTER.BsFishMath
		SmallJellyFish      PSF_ON_00002_1_MONSTER.BsFishMath `json:"Small Jellyfish"`
		Guppy               PSF_ON_00002_1_MONSTER.BsFishMath
		ClownFish           PSF_ON_00002_1_MONSTER.BsFishMath
		Dory                PSF_ON_00002_1_MONSTER.BsFishMath
		Grammidae           PSF_ON_00002_1_MONSTER.BsFishMath `json:"Grammidae"`
		RajahCichlasoma     PSF_ON_00002_1_MONSTER.BsFishMath `json:"rajah cichlasoma"`
		PufferFish          PSF_ON_00002_1_MONSTER.BsFishMath `json:"Puffer fish"`
		LionFish            PSF_ON_00002_1_MONSTER.BsFishMath
		Lanternfish         PSF_ON_00002_1_MONSTER.BsFishMath `json:"Lantern fish"`
		GiantJellyfish      PSF_ON_00002_1_MONSTER.BsFishMath `json:"giant jellyfish"`
		Turtle              PSF_ON_00002_1_MONSTER.BsFishMath
		Lobster             PSF_ON_00002_1_MONSTER.BsFishMath
		HermitCrab          PSF_ON_00002_1_MONSTER.BsFishMath    `json:"hermit crab"`
		GoldGemTurtle       PSF_ON_00002_1_MONSTER.BsFishMath    `json:"gold gem turtle"`
		GoldCoinToad        PSF_ON_00002_1_MONSTER.BsFishMath    `json:"gold coin toad"`
		GoldCrab            PSF_ON_00002_1_MONSTER.BsFishMath    `json:"gold crab"`
		GoldShark           PSF_ON_00002_1_MONSTER.BsFishMath    `json:"gold shark"`
		AsianArowana        PSF_ON_00002_1_MONSTER.BsFishMath    `json:"Asian arowana"`
		HotPot              PSF_ON_00002_1_MONSTER.BsHotPotMath  `json:"Hot Pot"`
		FaceOff1            PSF_ON_00002_1_MONSTER.BsFaceOffMath `json:"Face Off 1"`
		FaceOff2            PSF_ON_00002_1_MONSTER.BsFaceOffMath `json:"Face Off 2"`
		FaceOff3            PSF_ON_00002_1_MONSTER.BsFaceOffMath `json:"Face Off 3"`
		FaceOff4            PSF_ON_00002_1_MONSTER.BsFaceOffMath `json:"Face Off 4"`
		FaceOff5            PSF_ON_00002_1_MONSTER.BsFaceOffMath `json:"Face Off 5"`
		Drill               PSF_ON_00002_1_MONSTER.BsDrillMath
		MachineGun          PSF_ON_00002_1_MONSTER.BsMachineGunMath  `json:"Machine gun"`
		RedEnvelope         PSF_ON_00002_1_MONSTER.BsRedEnvelopeMath `json:"Red envelope"`
		Slot                PSF_ON_00002_1_MONSTER.BsSlotMath        `json:"slot"`
		FaceOffRandomChange struct {
			IconID           int
			IsWildSubstitute bool
			BonusID          string
			BonusTimes       []int
			IconPays         []int
		} `json:"Face Off random change"`
		Envelope PSF_ON_00002_1_MONSTER.BsEnvelopeMath `json:"Envelope"`
	}
	ChangeLink PSF_ON_00002_1_MONSTER.BsFaceOffChangerMath
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
		PSF_ON_00002_1_MONSTER.FsReelMath
	}
	Icons struct {
		StarFish            PSF_ON_00002_1_MONSTER.FsFishMath
		Clownfish           PSF_ON_00002_1_MONSTER.FsFishMath `json:"Clownfish"`
		LionFish            PSF_ON_00002_1_MONSTER.FsFishMath
		GiantJellyfish      PSF_ON_00002_1_MONSTER.FsFishMath `json:"giant jellyfish"`
		Lobster             PSF_ON_00002_1_MONSTER.FsFishMath
		HermitCrab          PSF_ON_00002_1_MONSTER.FsFishMath    `json:"hermit crab"`
		GoldGemTurtle       PSF_ON_00002_1_MONSTER.FsFishMath    `json:"gold gem turtle"`
		GoldCoinToad        PSF_ON_00002_1_MONSTER.FsFishMath    `json:"gold coin toad"`
		GoldCrab            PSF_ON_00002_1_MONSTER.FsFishMath    `json:"gold crab"`
		FaceOff1            PSF_ON_00002_1_MONSTER.FsFaceOffMath `json:"Face Off 1"`
		FaceOff3            PSF_ON_00002_1_MONSTER.FsFaceOffMath `json:"Face Off 3"`
		FaceOff5            PSF_ON_00002_1_MONSTER.FsFaceOffMath `json:"Face Off 5"`
		FaceOffRandomChange struct {
			IconID   int
			IconPays []int
		} `json:"Face Off random change"`
		MysterySymbolRandomChange struct {
			IconID   int
			IconPays []int
		} `json:"mystery symbol random change"`
	}
	Changelink1 PSF_ON_00002_1_MONSTER.FsFaceOffChangerMath_31 `json:"ChangeLink_1"`
	Changelink2 PSF_ON_00002_1_MONSTER.FsFaceOffChangerMath    `json:"ChangeLink_2"`
	MixAward    []PSF_ON_00002_1_MONSTER.FsSlotMath
}
