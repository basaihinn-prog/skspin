package PSF_ON_00003_1

import (
	PSF_ON_00003_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00003-1/PSF-ON-00003-1-MONSTER"
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
		LittleZombie PSF_ON_00003_1_MONSTER.BsLittleZombieMath `json:"little zombie"`
		Zombie1      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie1"`
		Zombie2      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie2"`
		Zombie3      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie3"`
		Zombie4      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie4"`
		Zombie5      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie5"`
		Zombie6      PSF_ON_00003_1_MONSTER.BsZombieMath       `json:"zombie6"`
		Zombie1Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie1_slot"`
		Zombie2Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie2_slot"`
		Zombie3Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie3_slot"`
		Zombie4Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie4_slot"`
		Zombie5Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie5_slot"`
		Zombie6Slot  PSF_ON_00003_1_MONSTER.BsZombieSlotMath   `json:"zombie6_slot"`
		Zombie1Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie1_drill"`
		Zombie2Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie2_drill"`
		Zombie3Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie3_drill"`
		Zombie4Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie4_drill"`
		Zombie5Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie5_drill"`
		Zombie6Drill PSF_ON_00003_1_MONSTER.BsZombieDrillMath  `json:"zombie6_drill"`
		BonusZombie  PSF_ON_00003_1_MONSTER.BsBonusZombieMath  `json:"bonus zombie"`
		Demon1       PSF_ON_00003_1_MONSTER.BsDemonMath        `json:"demon1"`
		Demon2       PSF_ON_00003_1_MONSTER.BsDemonMath        `json:"demon2"`
		Demon3       PSF_ON_00003_1_MONSTER.BsDemonMath        `json:"demon3"`
		Bullets      PSF_ON_00003_1_MONSTER.BsBullet
	}
}

type DrbMath struct {
	NetWinGroup struct {
		InitialWeight []int `json:"InitialWeight"`
		NetWin        []int `json:"NetWin"`
		NetWinWeight  []int `json:"NetWInWeight"`
		NetLoss       []int `json:"NetLoss"`
		NetLossWeight []int `json:"NetLossWeight"`

		NetWinRTP  struct{ DrbWinMath }  `json:"NetWinRTP"`
		NetLossRTP struct{ DrbLossMath } `json:"NetLossRTP"`
	} `json:"NetWin_Group"`
}

type DrbWinMath struct {
	NetWin1 struct{ DrbNetWin } `json:"NetWin1"`
	NetWin2 struct{ DrbNetWin } `json:"NetWin2"`
	NetWin3 struct{ DrbNetWin } `json:"NetWin3"`
	NetWin4 struct{ DrbNetWin } `json:"NetWin4"`
	NetWin5 struct{ DrbNetWin } `json:"NetWin5"`
}

type DrbLossMath struct {
	NetLoss1 struct{ DrbNetLoss } `json:"NetLoss1"`
	NetLoss2 struct{ DrbNetLoss } `json:"NetLoss2"`
	NetLoss3 struct{ DrbNetLoss } `json:"NetLoss3"`
	NetLoss4 struct{ DrbNetLoss } `json:"NetLoss4"`
	NetLoss5 struct{ DrbNetLoss } `json:"NetLoss5"`
}

type DrbNetWin struct {
	Win      []int            `json:"Win"`
	RTPGroup []int            `json:"RTPGroup"`
	Weight   []int            `json:"Weight"`
	RTP1     struct{ DrbRTP } `json:"RTP1"`
	RTP2     struct{ DrbRTP } `json:"RTP2"`
	RTP3     struct{ DrbRTP } `json:"RTP3"`
	RTP4     struct{ DrbRTP } `json:"RTP4"`
	RTP5     struct{ DrbRTP } `json:"RTP5"`
}

type DrbNetLoss struct {
	Loss     []int            `json:"Loss"`
	RTPGroup []int            `json:"RTPGroup"`
	Weight   []int            `json:"Weight"`
	RTP1     struct{ DrbRTP } `json:"RTP1"`
	RTP2     struct{ DrbRTP } `json:"RTP2"`
	RTP3     struct{ DrbRTP } `json:"RTP3"`
	RTP4     struct{ DrbRTP } `json:"RTP4"`
	RTP5     struct{ DrbRTP } `json:"RTP5"`
}

type DrbRTP struct {
	ID     string `json:"ID"`
	Bullet []int  `json:"Bullet"`
}

type FsMath struct {
	ModuleType string
	BaseData   struct {
		ProjectID      string
		ModuleID       string
		GameType       string
		StripType      string
		GameTypeDetail string
		BonusList      []string
	}
	StripTable struct {
		PSF_ON_00003_1_MONSTER.FsStripTableMath
	}
	Icons struct {
		LittleZombie1             PSF_ON_00003_1_MONSTER.FsZombieMath `json:"little zombie 1"`
		LittleZombie2             PSF_ON_00003_1_MONSTER.FsZombieMath `json:"little zombie 2"`
		LittleZombie3             PSF_ON_00003_1_MONSTER.FsZombieMath `json:"little zombie 3"`
		LittleZombie4             PSF_ON_00003_1_MONSTER.FsZombieMath `json:"little zombie 4"`
		Zombie1                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 1"`
		Zombie2                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 2"`
		Zombie3                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 3"`
		Zombie4                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 4"`
		Zombie5                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 5"`
		Zombie6                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"zombie 6"`
		Demon31                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"demon 3_1"`
		Demon32                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"demon 3_2"`
		Demon33                   PSF_ON_00003_1_MONSTER.FsZombieMath `json:"demon 3_3"`
		MysterySymbolRandomChange PSF_ON_00003_1_MONSTER.FsZombieMath `json:"mystery symbol random change"`
	}
	ChangeLink1 PSF_ON_00003_1_MONSTER.FsSlotChanger1Math `json:"ChangeLink_1"`
	MixAward    []PSF_ON_00003_1_MONSTER.FsSlotMath
}
