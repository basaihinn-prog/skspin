package gamesetting

import (
	"reflect"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	PSF_ON_00002 "serve/service_fish/domain/fish/PSF-ON-00002"
	PSF_ON_00002_1 "serve/service_fish/domain/probability/PSF-ON-00002-1"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	PSFM_00003_94_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-94-1"
	PSFM_00003_95_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-95-1"
	PSFM_00003_96_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-96-1"
	PSFM_00003_97_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-97-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	"serve/service_fish/models"
)

var math00003List = map[string]bool{
	models.PSFM_00003_93_1: true,
	models.PSFM_00003_94_1: true,
	models.PSFM_00003_95_1: true,
	models.PSFM_00003_96_1: true,
	models.PSFM_00003_97_1: true,
	models.PSFM_00003_98_1: true,
}

func psfm_00003_1_PayTable(sr *game_setting_proto.StripsRecall, mathModuleId string) bool {
	if !math00003List[mathModuleId] {
		return false
	}

	for _, fishId := range fishList00002 {
		sr.AllSymbolDef = append(sr.AllSymbolDef, psfm_00003_1_symbolDef(fishId, mathModuleId, fishSymbol00002[fishId].symbolType, fishSymbol00002[fishId].payType))
	}

	return true
}

func psfm_00003_1_symbolDef(
	fishId int,
	mathModuleId string,
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType,
	payType game_setting_proto.StripsRecall_SymbolDef_PayType,
) *game_setting_proto.StripsRecall_SymbolDef {
	sd := &game_setting_proto.StripsRecall_SymbolDef{
		SymbolId:   uint32(fishId),
		SymbolType: symbolType,
		LayerId:    psfm_00003_getObjectLayer(fishId),
		PayType:    payType,
	}

	switch mathModuleId {
	case models.PSFM_00003_93_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath)
	case models.PSFM_00003_94_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_94_1.RTP94BS.PSF_ON_00002_1_BsMath)
	case models.PSFM_00003_95_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_95_1.RTP95BS.PSF_ON_00002_1_BsMath)
	case models.PSFM_00003_96_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_96_1.RTP96BS.PSF_ON_00002_1_BsMath)
	case models.PSFM_00003_97_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_97_1.RTP97BS.PSF_ON_00002_1_BsMath)
	case models.PSFM_00003_98_1:
		psfm_00003_getPay(fishId, sd, PSFM_00003_98_1.RTP98BS.PSF_ON_00002_1_BsMath)
	}
	return sd
}

func psfm_00003_getPay(fishId int, sd *game_setting_proto.StripsRecall_SymbolDef, bsMathI interface{}) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00002_1.BsMath)

	switch fishId {
	case 0:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.StarFish.IconPays[0]))
	case 1:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SeaHorse.IconPays[0]))
	case 2:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SmallJellyFish.IconPays[0]))
	case 3:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Guppy.IconPays[0]))
	case 4:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.ClownFish.IconPays[0]))
	case 5:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Dory.IconPays[0]))
	case 6:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Grammidae.IconPays[0]))
	case 7:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.RajahCichlasoma.IconPays[0]))
	case 8:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.PufferFish.IconPays[0]))
	case 9:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.LionFish.IconPays[0]))
	case 10:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Lanternfish.IconPays[0]))
	case 11:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GiantJellyfish.IconPays[0]))
	case 12:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Turtle.IconPays[0]))
	case 13:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Lobster.IconPays[0]))
	case 14:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.HermitCrab.IconPays[0]))
	case 15:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GoldGemTurtle.IconPays[0]))
	case 16:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GoldCoinToad.IconPays[0]))
	case 17:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GoldCrab.IconPays[0]))
	case 18:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GoldShark.IconPays[0]))
	case 19:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.AsianArowana.IconPays[0]))
	case 20:
		iconPays := bsMath.Icons.HotPot.IconPays

		for v := range iconPays {
			sd.Pays = append(sd.Pays, int64(v))
		}

	case 21:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff1.IconPays[1]))
	case 22:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff2.IconPays[1]))
	case 23:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff3.IconPays[1]))
	case 24:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff4.IconPays[1]))
	case 25:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff5.IconPays[1]))
	case 26:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Drill.IconPays[0]))

		bonusTimes := bsMath.Icons.Drill.BonusTimes

		for v := range bonusTimes {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 27:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.MachineGun.IconPays[0]))

		bonusTimes := bsMath.Icons.MachineGun.BonusTimes

		for v := range bonusTimes {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 28:
		sd.Pays = append(sd.Pays, 0)

		lowIconPays := bsMath.Icons.RedEnvelope.LowIconPays
		highIconPays := bsMath.Icons.RedEnvelope.HighIconPays

		for v := range lowIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

		for v := range highIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 29:
	// slot do nothing

	case 30:
		sd.Pays = append(sd.Pays, 0)
	}
}

func psfm_00003_getObjectLayer(fishId int) uint32 {
	switch fishId {
	case 21, 22, 23, 24, 25, 30:
		return uint32(PSF_ON_00002.Objects.Layer(30))
	default:
		return uint32(PSF_ON_00002.Objects.Layer(fishId))
	}
}

var fishList00002 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
var fishSymbol00002 = map[int]symbol{
	0: {symbolType: K_SYMBOL, payType: K_XBET}, 1: {symbolType: K_SYMBOL, payType: K_XBET}, 2: {symbolType: K_SYMBOL, payType: K_XBET},
	3: {symbolType: K_SYMBOL, payType: K_XBET}, 4: {symbolType: K_SYMBOL, payType: K_XBET}, 5: {symbolType: K_SYMBOL, payType: K_XBET},
	6: {symbolType: K_SYMBOL, payType: K_XBET}, 7: {symbolType: K_SYMBOL, payType: K_XBET}, 8: {symbolType: K_SYMBOL, payType: K_XBET},
	9: {symbolType: K_SYMBOL, payType: K_XBET}, 10: {symbolType: K_SYMBOL, payType: K_XBET}, 11: {symbolType: K_SYMBOL, payType: K_XBET},
	12: {symbolType: K_SYMBOL, payType: K_XBET}, 13: {symbolType: K_SYMBOL, payType: K_XBET}, 14: {symbolType: K_SYMBOL, payType: K_XBET},
	15: {symbolType: K_SYMBOL, payType: K_XBET}, 16: {symbolType: K_SYMBOL, payType: K_XBET}, 17: {symbolType: K_SYMBOL, payType: K_XBET},
	18: {symbolType: K_SYMBOL, payType: K_XBET}, 19: {symbolType: K_SYMBOL, payType: K_XBET}, 20: {symbolType: K_SYMBOL, payType: K_XBET},
	21: {symbolType: K_SYMBOL, payType: K_XBET}, 22: {symbolType: K_SYMBOL, payType: K_XBET}, 23: {symbolType: K_SYMBOL, payType: K_XBET},
	24: {symbolType: K_SYMBOL, payType: K_XBET}, 25: {symbolType: K_SYMBOL, payType: K_XBET}, 26: {symbolType: K_BONUS, payType: K_TOTAL_BONUS},
	27: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 28: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 29: {symbolType: K_BONUS, payType: K_XBET},
	30: {symbolType: K_SYMBOL, payType: K_XBET},
}
