package PSF_ON_00005_1

import PSF_ON_00005_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00005-1/PSF-ON-00005-1-MONSTER"

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
		StarFish               PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"starfish"`
		JellyFish              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Jellyfish"`
		RiverShrimo            PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"River Shrimo"`
		Guppy                  PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Guppy"`
		ClownFish              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Clownfish"`
		Dory                   PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Dory"`
		Grammidae              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Grammidae"`
		RajahCichlasoma        PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Rajah Cichlasoma"`
		SiameseFightingFish    PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Siamese Fighting Fish"`
		LionFish               PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Lionfish"`
		SiameseTigerFish       PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Siamese Tigerfish"`
		DurianFish             PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Durian Fish"`
		MangoFish              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Mango Fish"`
		MuayThaiLobster        PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Muay Thai Lobster"`
		HermitCrab             PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Hermit Crab"`
		GoldGemTurtle          PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Gold Gem Turtle"`
		GoldCoinToad           PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Gold Coin Toad"`
		GoldCrab               PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Gold Crab"`
		GoldShark              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Gold Shark"`
		Crocodile              PSF_ON_00005_1_MONSTER.BsFishMath                  `json:"Crocodile"`
		RedEnvelopeTree        PSF_ON_00005_1_MONSTER.BsRedEnvelopeTreeMath       `json:"Red envelope Tree"`
		Slot                   PSF_ON_00005_1_MONSTER.BsSlotMath                  `json:"slot"`
		RandomTriggerEnvelope  PSF_ON_00005_1_MONSTER.BsRandomTriggerEnvelopeMath `json:"Random Trigger Envelope"`
		MachineGun             PSF_ON_00005_1_MONSTER.BsMachineGunMath            `json:"Machine gun"`
		SuperMachineGun        PSF_ON_00005_1_MONSTER.BsMachineGunMath            `json:"Super Machine Gun"`
		FruitDish              PSF_ON_00005_1_MONSTER.BsFruitDish                 `json:"Fruit Dish"`
		APackOfBeer            PSF_ON_00005_1_MONSTER.BsFruitDish                 `json:"A pack of beer"`
		GiantWhaleRandomChange struct {
			IconID int `json:"IconID"`

			BonusID    string `json:"BonusID"`
			BonusTimes []int  `json:"BonusTimes"`
			IconPays   []int  `json:"IconPays"`
		} `json:"Giant Whale random change"`
		ChangeLink  PSF_ON_00005_1_MONSTER.BsGiantWhaleChangerMath `json:"ChangeLink"`
		GiantWhale1 PSF_ON_00005_1_MONSTER.BsGiantWhaleMath        `json:"Giant Whale 1"`
		GiantWhale2 PSF_ON_00005_1_MONSTER.BsGiantWhaleMath        `json:"Giant Whale 2"`
		GiantWhale3 PSF_ON_00005_1_MONSTER.BsGiantWhaleMath        `json:"Giant Whale 3"`
		GiantWhale4 PSF_ON_00005_1_MONSTER.BsGiantWhaleMath        `json:"Giant Whale 4"`
		GiantWhale5 PSF_ON_00005_1_MONSTER.BsGiantWhaleMath        `json:"Giant Whale 5"`
	}
}

type DrbMath struct {
	InitLowGroupWeight []int `json:"InitLowGroupWeight"`
	RTPGroupID         []int `json:"RTPGroupID"`
	HighRTPGroupWeight []int `json:"HighRTPGroupWeight"`
	LowRTPGroupWeight  []int `json:"LowRTPGroupWeight"`
	RTPGroup           struct {
		RTP1  struct{ HitDrbMath } `json:"RTP1"`
		RTP2  struct{ HitDrbMath } `json:"RTP2"`
		RTP3  struct{ HitDrbMath } `json:"RTP3"`
		RTP4  struct{ HitDrbMath } `json:"RTP4"`
		RTP5  struct{ HitDrbMath } `json:"RTP5"`
		RTP6  struct{ HitDrbMath } `json:"RTP6"`
		RTP7  struct{ HitDrbMath } `json:"RTP7"`
		RTP8  struct{ HitDrbMath } `json:"RTP8"`
		RTP9  struct{ HitDrbMath } `json:"RTP9"`
		RTP10 struct{ HitDrbMath } `json:"RTP10"`
		RTP11 struct{ HitDrbMath } `json:"RTP11"`
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
		PSF_ON_00005_1_MONSTER.FsStripTableMath
	} `json:"StripTable"`
	Icons struct {
		StarFish                  PSF_ON_00005_1_MONSTER.FsFishMath `json:"starfish"`
		ClownFish                 PSF_ON_00005_1_MONSTER.FsFishMath `json:"Clownfish"`
		LionFish                  PSF_ON_00005_1_MONSTER.FsFishMath `json:"Lionfish"`
		GiantJellyFish            PSF_ON_00005_1_MONSTER.FsFishMath `json:"giant jellyfish"`
		Lobster                   PSF_ON_00005_1_MONSTER.FsFishMath `json:"lobster"`
		HermitCrab                PSF_ON_00005_1_MONSTER.FsFishMath `json:"hermit crab"`
		GoldGemTurtle             PSF_ON_00005_1_MONSTER.FsFishMath `json:"gold gem turtle"`
		GoldCoinToad              PSF_ON_00005_1_MONSTER.FsFishMath `json:"gold coin toad"`
		GoldCrab                  PSF_ON_00005_1_MONSTER.FsFishMath `json:"gold crab"`
		FruitDish                 PSF_ON_00005_1_MONSTER.FsFishMath `json:"Face Off 1"`
		APackOfBeer               PSF_ON_00005_1_MONSTER.FsFishMath `json:"Face Off 3"`
		GiantWhale                PSF_ON_00005_1_MONSTER.FsFishMath `json:"Face Off 5"`
		FaceOffRandomChange       PSF_ON_00005_1_MONSTER.FsFishMath `json:"Face Off random change"`
		MysterySymbolRandomChange PSF_ON_00005_1_MONSTER.FsFishMath `json:"mystery symbol random change"`
	} `json:"Icons"`

	ChangeLink1 PSF_ON_00005_1_MONSTER.FsGiantWhaleChanger1Math `json:"ChangeLink_1"`
	ChangeLink2 PSF_ON_00005_1_MONSTER.FsGiantWhaleChanger2Math `json:"ChangeLink_2"`
	MixAward    []PSF_ON_00005_1_MONSTER.FsSlotMath             `json:"MixAward"`
}
