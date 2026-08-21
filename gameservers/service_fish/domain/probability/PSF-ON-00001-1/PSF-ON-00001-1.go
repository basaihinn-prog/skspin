package PSF_ON_00001_1

import (
	PSF_ON_00001_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00001-1/PSF-ON-00001-1-MONSTER"
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
		StarFish            PSF_ON_00001_1_MONSTER.BsFishMath
		SeaHorse            PSF_ON_00001_1_MONSTER.BsFishMath
		Guppy               PSF_ON_00001_1_MONSTER.BsFishMath
		ClownFish           PSF_ON_00001_1_MONSTER.BsFishMath
		Dory                PSF_ON_00001_1_MONSTER.BsFishMath
		Grammidae           PSF_ON_00001_1_MONSTER.BsFishMath
		RajahCichlasoma     PSF_ON_00001_1_MONSTER.BsFishMath `json:"Rajah Cichlasoma"`
		PufferFish          PSF_ON_00001_1_MONSTER.BsFishMath `json:"Puffer fish"`
		LanternFish         PSF_ON_00001_1_MONSTER.BsFishMath `json:"Lantern fish"`
		LionFish            PSF_ON_00001_1_MONSTER.BsFishMath
		Turtle              PSF_ON_00001_1_MONSTER.BsFishMath
		Lobster             PSF_ON_00001_1_MONSTER.BsFishMath
		Penguin             PSF_ON_00001_1_MONSTER.BsFishMath
		Platypus            PSF_ON_00001_1_MONSTER.BsFishMath
		Manatee             PSF_ON_00001_1_MONSTER.BsFishMath
		Dolphin             PSF_ON_00001_1_MONSTER.BsFishMath
		Koi                 PSF_ON_00001_1_MONSTER.BsFishMath
		HermitCrab          PSF_ON_00001_1_MONSTER.BsFishMath `json:"Hermit crab"`
		Shark               PSF_ON_00001_1_MONSTER.BsFishMath
		AsianArowana        PSF_ON_00001_1_MONSTER.BsFishMath    `json:"Asian arowana"`
		HotPot              PSF_ON_00001_1_MONSTER.BsHotPotMath  `json:"Hot Pot"`
		FaceOff1            PSF_ON_00001_1_MONSTER.BsFaceOffMath `json:"Face Off 1"`
		FaceOff2            PSF_ON_00001_1_MONSTER.BsFaceOffMath `json:"Face Off 2"`
		FaceOff3            PSF_ON_00001_1_MONSTER.BsFaceOffMath `json:"Face Off 3"`
		FaceOff4            PSF_ON_00001_1_MONSTER.BsFaceOffMath `json:"Face Off 4"`
		FaceOff5            PSF_ON_00001_1_MONSTER.BsFaceOffMath `json:"Face Off 5"`
		Drill               PSF_ON_00001_1_MONSTER.BsDrillMath
		MachineGun          PSF_ON_00001_1_MONSTER.BsMachineGunMath  `json:"Machine gun"`
		RedEnvelope         PSF_ON_00001_1_MONSTER.BsRedEnvelopeMath `json:"Red envelope"`
		Slot                PSF_ON_00001_1_MONSTER.BsSlotMath        `json:"slot"`
		FaceOffRandomChange struct {
			IconID           int
			IsWildSubstitute bool
			BonusID          string
			BonusTimes       []int
			IconPays         []int
		} `json:"Face Off random change"`
		Envelope PSF_ON_00001_1_MONSTER.BsRedEnvelopeMath `json:"Envelope"`
	}
	ChangeLink PSF_ON_00001_1_MONSTER.BsFaceOffChangerMath
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
		PSF_ON_00001_1_MONSTER.FsReelMath
	}
	Icons struct {
		StarFish            PSF_ON_00001_1_MONSTER.FsFishMath
		Dory                PSF_ON_00001_1_MONSTER.FsFishMath
		LionFish            PSF_ON_00001_1_MONSTER.FsFishMath
		Lobster             PSF_ON_00001_1_MONSTER.FsFishMath
		Platypus            PSF_ON_00001_1_MONSTER.FsFishMath
		Manatee             PSF_ON_00001_1_MONSTER.FsFishMath
		Dolphin             PSF_ON_00001_1_MONSTER.FsFishMath
		Koi                 PSF_ON_00001_1_MONSTER.FsFishMath
		HermitCrab          PSF_ON_00001_1_MONSTER.FsFishMath    `json:"Hermit crab"`
		FaceOff1            PSF_ON_00001_1_MONSTER.FsFaceOffMath `json:"Face Off 1"`
		FaceOff3            PSF_ON_00001_1_MONSTER.FsFaceOffMath `json:"Face Off 3"`
		FaceOff5            PSF_ON_00001_1_MONSTER.FsFaceOffMath `json:"Face Off 5"`
		FaceOffRandomChange struct {
			IconID   int
			IconPays []int
		} `json:"Face Off random change"`
	}
	ChangeLink PSF_ON_00001_1_MONSTER.FsFaceOffChangerMath
	MixAward   []PSF_ON_00001_1_MONSTER.FsSlotMath
}
