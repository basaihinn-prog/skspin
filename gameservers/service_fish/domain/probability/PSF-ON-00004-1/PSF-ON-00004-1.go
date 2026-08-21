package PSF_ON_00004_1

import PSF_ON_00004_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00004-1/PSF-ON-00004-1-MONSTER"

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
		Clam               PSF_ON_00004_1_MONSTER.BsFishMath
		Shrimp             PSF_ON_00004_1_MONSTER.BsFishMath
		Surmullet          PSF_ON_00004_1_MONSTER.BsFishMath
		Squid              PSF_ON_00004_1_MONSTER.BsFishMath
		FlyingFish         PSF_ON_00004_1_MONSTER.BsFishMath `json:"Flying Fish"`
		Halibut            PSF_ON_00004_1_MONSTER.BsFishMath
		ButterflyFish      PSF_ON_00004_1_MONSTER.BsFishMath `json:"Butterflyfish"`
		Oplegnathus        PSF_ON_00004_1_MONSTER.BsFishMath
		Snapper            PSF_ON_00004_1_MONSTER.BsFishMath
		MahiMahi           PSF_ON_00004_1_MONSTER.BsFishMath `json:"Mahi Mahi"`
		Stingray           PSF_ON_00004_1_MONSTER.BsFishMath
		LobsterWaiter      PSF_ON_00004_1_MONSTER.BsFishMath `json:"Lobster Waiter"`
		PenguinWaiter      PSF_ON_00004_1_MONSTER.BsFishMath `json:"Penguin Waiter"`
		PlatypusSeniorChef PSF_ON_00004_1_MONSTER.BsFishMath `json:"Platypus Senior Chef"`
		SeaLionChef        PSF_ON_00004_1_MONSTER.BsFishMath `json:"Sea Lion Chef"`
		HairtCrab          PSF_ON_00004_1_MONSTER.BsFishMath `json:"Hairt Crab"`
		SnaggletoothShark  PSF_ON_00004_1_MONSTER.BsFishMath `json:"Snaggletooth Shark"`
		SwordFish          PSF_ON_00004_1_MONSTER.BsFishMath
		WhaleShark         PSF_ON_00004_1_MONSTER.BsFishMath        `json:"Whale Shark"`
		GiantOarfish       PSF_ON_00004_1_MONSTER.BsFishMath        `json:"Giant Oarfish"`
		RedEnvelope        PSF_ON_00004_1_MONSTER.BsRedEnvelopeMath `json:"Red Envelope"`
		Slot               PSF_ON_00004_1_MONSTER.BsSlotMath
		MachineGun         PSF_ON_00004_1_MONSTER.BsMachineGunMath     `json:"Machine Gun"`
		SuperMachineGun    PSF_ON_00004_1_MONSTER.BsMachineGunMath     `json:"Super Machine Gun"`
		LobsterDash        PSF_ON_00004_1_MONSTER.BsDashMath           `json:"Lobster Dish"`
		FruitDash          PSF_ON_00004_1_MONSTER.BsDashMath           `json:"Fruit Dish"`
		XiaoLongBao        PSF_ON_00004_1_MONSTER.BsXiaoLongBaoMath    `json:"XIAO LONG BAO"`
		WhiteTigerChef     PSF_ON_00004_1_MONSTER.BsWhiteTigerChefMath `json:"WHITE TIGER CHEF"`
	}
}

type DrbMath struct {
	FirstMultiplier1 []int `json:"FirstMultiplier_1"`
	FirstMultiplier2 []int `json:"FirstMultiplier_2"`
	FirstMultiplier3 []int `json:"FirstMultiplier_3"`
	BetGroup         struct {
		Bet1  struct{ DrbBetGroupMath }
		Bet2  struct{ DrbBetGroupMath }
		Bet3  struct{ DrbBetGroupMath }
		Bet4  struct{ DrbBetGroupMath }
		Bet5  struct{ DrbBetGroupMath }
		Bet6  struct{ DrbBetGroupMath }
		Bet7  struct{ DrbBetGroupMath }
		Bet8  struct{ DrbBetGroupMath }
		Bet9  struct{ DrbBetGroupMath }
		Bet10 struct{ DrbBetGroupMath }
	}
	WinGroup struct {
		Win1  struct{ DrbWinGroupMath }
		Win2  struct{ DrbWinGroupMath }
		Win3  struct{ DrbWinGroupMath }
		Win4  struct{ DrbWinGroupMath }
		Win5  struct{ DrbWinGroupMath }
		Win6  struct{ DrbWinGroupMath }
		Win7  struct{ DrbWinGroupMath }
		Win8  struct{ DrbWinGroupMath }
		Win9  struct{ DrbWinGroupMath }
		Win10 struct{ DrbWinGroupMath }
	}
}

type DrbBetGroupMath struct {
	Turnover int
	Weight   int
}

type DrbWinGroupMath struct {
	Multiplier int
	Win        int
	Weight     int
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
		PSF_ON_00004_1_MONSTER.FsStripTableMath
	}

	Icons struct {
		Clam                      PSF_ON_00004_1_MONSTER.FsFishMath
		FlyingFish                PSF_ON_00004_1_MONSTER.FsFishMath `json:"Flying Fish"`
		MahiMahi                  PSF_ON_00004_1_MONSTER.FsFishMath `json:"Mahi Mahi"`
		LobsterWaiter             PSF_ON_00004_1_MONSTER.FsFishMath `json:"Lobster Waiter"`
		PlatypusSeniorChef        PSF_ON_00004_1_MONSTER.FsFishMath `json:"Platypus Senior Chef"`
		SeaLionChef               PSF_ON_00004_1_MONSTER.FsFishMath `json:"Sea Lion Chef"`
		HairtCrab                 PSF_ON_00004_1_MONSTER.FsFishMath `json:"Hairt Crab"`
		SnaggletoothShark         PSF_ON_00004_1_MONSTER.FsFishMath `json:"Snaggletooth Shark"`
		Swordfish                 PSF_ON_00004_1_MONSTER.FsFishMath
		LobsterDish               PSF_ON_00004_1_MONSTER.FsFishMath `json:"Lobster Dish"`
		FruitDish                 PSF_ON_00004_1_MONSTER.FsFishMath `json:"Fruit Dish"`
		WhiteTigerChef            PSF_ON_00004_1_MONSTER.FsFishMath `json:"WHITE TIGER CHEF"`
		MysterySymbolRandomChange PSF_ON_00004_1_MONSTER.FsFishMath `json:"mystery symbol random change"`
		BigSymbol                 PSF_ON_00004_1_MONSTER.FsFishMath `json:"big symbol"`
	}
	ChangeLink1 PSF_ON_00004_1_MONSTER.FsFaceOffChanger1Math `json:"ChangeLink_1"`
	ChangeLink2 PSF_ON_00004_1_MONSTER.FsFaceOffChanger2Math `json:"ChangeLink_2"`
	MixAward    []PSF_ON_00004_1_MONSTER.FsSlotMath
}
