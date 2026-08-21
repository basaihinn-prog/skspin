package PSF_ON_00007_1

import PSF_ON_00007_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00007-1/PSF-ON-00007-1-MONSTER"

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
		Fish0                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_0"`
		Fish1                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_1"`
		Fish2                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_2"`
		Fish3                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_3"`
		Fish4                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_4"`
		Fish5                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_5"`
		Fish6                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_6"`
		Fish7                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_7"`
		Fish8                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_8"`
		Fish9                 PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_9"`
		Fish10                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_10"`
		Fish11                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_11"`
		Fish12                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_12"`
		Fish13                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_13"`
		Fish14                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_14"`
		Fish15                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_15"`
		Fish16                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_16"`
		Fish17                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_17"`
		Fish18                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_18"`
		Fish19                PSF_ON_00007_1_MONSTER.BsFishMath                  `json:"Fish_19"`
		MagicShellFaFaFa      PSF_ON_00007_1_MONSTER.BsMagicShellFaFaFaMath      `json:"Magic Shell FAFAFA"`
		LuckSlot              PSF_ON_00007_1_MONSTER.BsLuckSlotMath              `json:"Luck Slot"`
		RandomTriggerEnvelope PSF_ON_00007_1_MONSTER.BsRandomTriggerEnvelopeMath `json:"Random Trigger Envelope"`
		MachineGun            PSF_ON_00007_1_MONSTER.BsMachineGunMath            `json:"Machine gun"`
		SuperMachineGun       PSF_ON_00007_1_MONSTER.BsMachineGunMath            `json:"Super Machine Gun"`
		PirateShip            PSF_ON_00007_1_MONSTER.BsPirateShip                `json:"Pirate Ship"`
		TreasureDish          PSF_ON_00007_1_MONSTER.BsPirateShip                `json:"Treasure Dish"`
		TreasureHermitCrab    PSF_ON_00007_1_MONSTER.BsTreasureHermitCrabMath    `json:"Treasure Hermit Crab"`
		OceanKing             PSF_ON_00007_1_MONSTER.BsOceanKingMath             `json:"Ocean King"`
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
		PSF_ON_00007_1_MONSTER.FsStripTableMath
	} `json:"StripTable"`
	Icons struct {
		Fish0                     PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_1"`
		Fish4                     PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_2"`
		Fish9                     PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_4"`
		Fish11                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_7"`
		Fish13                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_9"`
		Fish14                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_11"`
		Fish15                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_14"`
		Fish16                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_15"`
		Fish17                    PSF_ON_00007_1_MONSTER.FsFishMath `json:"Fish_17"`
		PirateShip                PSF_ON_00007_1_MONSTER.FsFishMath `json:"Pirate Ship"`
		TreasureDish              PSF_ON_00007_1_MONSTER.FsFishMath `json:"Treasure Dish"`
		OceanKing                 PSF_ON_00007_1_MONSTER.FsFishMath `json:"Ocean King"`
		MysterySymbolRandomChange PSF_ON_00007_1_MONSTER.FsFishMath `json:"mystery symbol random change"`
	} `json:"Icons"`
	ChangeLink1 PSF_ON_00007_1_MONSTER.FsChangeLink1 `json:"ChangeLink_1"`
	MixAward    []PSF_ON_00007_1_MONSTER.FsSlotMath  `json:"MixAward"`
}
