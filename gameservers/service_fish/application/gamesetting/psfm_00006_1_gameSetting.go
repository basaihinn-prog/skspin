package gamesetting

import (
	"reflect"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	PSF_ON_00005 "serve/service_fish/domain/fish/PSF-ON-00005"
	PSF_ON_00005_1 "serve/service_fish/domain/probability/PSF-ON-00005-1"
	PSFM_00006_95_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-95-1"
	PSFM_00006_96_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-96-1"
	PSFM_00006_97_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-97-1"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	"serve/service_fish/models"
)

var math00006List = map[string]bool{
	models.PSFM_00006_95_1: true,
	models.PSFM_00006_96_1: true,
	models.PSFM_00006_97_1: true,
	models.PSFM_00006_98_1: true,
}

func psfm_00006_1_PayTable(sr *game_setting_proto.StripsRecall, mathModuleId string) bool {
	if !math00006List[mathModuleId] {
		return false
	}

	for _, fishId := range fishList00005 {
		sr.AllSymbolDef = append(sr.AllSymbolDef, psfm_00006_1_symbolDef(fishId, mathModuleId, fishSymbol00005[fishId].symbolType, fishSymbol00005[fishId].payType))
	}

	return true
}

func psfm_00006_1_symbolDef(
	fishId int,
	mathModuleId string,
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType,
	payType game_setting_proto.StripsRecall_SymbolDef_PayType,
) *game_setting_proto.StripsRecall_SymbolDef {
	sd := &game_setting_proto.StripsRecall_SymbolDef{
		SymbolId:   uint32(fishId),
		SymbolType: symbolType,
		LayerId:    psfm_00006_getObjectLayer(fishId),
		PayType:    payType,
	}

	switch mathModuleId {
	case models.PSFM_00006_95_1:
		psfm_00006_getPay(fishId, sd, PSFM_00006_95_1.RTP95BS.PSF_ON_00005_1_BsMath, PSFM_00006_95_1.RTP95FS.PSF_ON_00005_1_FsMath)
	case models.PSFM_00006_96_1:
		psfm_00006_getPay(fishId, sd, PSFM_00006_96_1.RTP96BS.PSF_ON_00005_1_BsMath, PSFM_00006_96_1.RTP96FS.PSF_ON_00005_1_FsMath)
	case models.PSFM_00006_97_1:
		psfm_00006_getPay(fishId, sd, PSFM_00006_97_1.RTP97BS.PSF_ON_00005_1_BsMath, PSFM_00006_97_1.RTP97FS.PSF_ON_00005_1_FsMath)
	case models.PSFM_00006_98_1:
		psfm_00006_getPay(fishId, sd, PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath, PSFM_00006_98_1.RTP98FS.PSF_ON_00005_1_FsMath)
	}
	return sd
}

func psfm_00006_getPay(fishId int, sd *game_setting_proto.StripsRecall_SymbolDef, bsMathI, fsMathI interface{}) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00005_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00005_1.FsMath)

	switch fishId {
	case 0:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.StarFish.IconPays[0]))
	case 1:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.JellyFish.IconPays[0]))
	case 2:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.RiverShrimo.IconPays[0]))
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
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SiameseFightingFish.IconPays[0]))
	case 9:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.LionFish.IconPays[0]))
	case 10:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SiameseTigerFish.IconPays[0]))
	case 11:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.DurianFish.IconPays[0]))
	case 12:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.MangoFish.IconPays[0]))
	case 13:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.MuayThaiLobster.IconPays[0]))
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
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Crocodile.IconPays[0]))
	case 100:
		sd.Pays = append(sd.Pays, 0)

		lowIconPays := bsMath.Icons.RedEnvelopeTree.LowIconPays
		highIconPays := bsMath.Icons.RedEnvelopeTree.HighIconPays

		for _, v := range lowIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

		for _, v := range highIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 101:
		// slot do nothing
	case 102:
		// Envelope do nothing
	case 201:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.MachineGun.IconPays[0]))
		bonusTimes := bsMath.Icons.MachineGun.BonusTimes

		for _, v := range bonusTimes {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 202:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SuperMachineGun.IconPays[0]))
		bonusTimes := bsMath.Icons.SuperMachineGun.BonusTimes

		for _, v := range bonusTimes {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 300:
		iconPays := bsMath.Icons.FruitDish.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 301:
		iconPays := bsMath.Icons.APackOfBeer.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 500:
		iconPays := fsMath.MixAward[0].MixIconGroup
		iconPay := 0
		for i := 0; i < len(iconPays); i++ {
			if iconPays[i].IconId == 500 {
				iconPay = iconPays[i].Quantity
			}
		}
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 501:
		iconPays := bsMath.Icons.GiantWhale1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 502:
		iconPays := bsMath.Icons.GiantWhale2.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 503:
		iconPays := bsMath.Icons.GiantWhale3.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 504:
		iconPays := bsMath.Icons.GiantWhale4.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 505:
		iconPays := bsMath.Icons.GiantWhale5.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))
	}
}

func psfm_00006_getObjectLayer(fishId int) uint32 {
	switch fishId {
	case 501, 502, 503, 504, 505:
		return uint32(PSF_ON_00005.Objects.Layer(500))
	default:
		return uint32(PSF_ON_00005.Objects.Layer(fishId))
	}
}

var fishList00005 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 500, 501, 502, 503, 504, 505}
var fishSymbol00005 = map[int]symbol{
	0: {symbolType: K_SYMBOL, payType: K_XBET}, 1: {symbolType: K_SYMBOL, payType: K_XBET}, 2: {symbolType: K_SYMBOL, payType: K_XBET},
	3: {symbolType: K_SYMBOL, payType: K_XBET}, 4: {symbolType: K_SYMBOL, payType: K_XBET}, 5: {symbolType: K_SYMBOL, payType: K_XBET},
	6: {symbolType: K_SYMBOL, payType: K_XBET}, 7: {symbolType: K_SYMBOL, payType: K_XBET}, 8: {symbolType: K_SYMBOL, payType: K_XBET},
	9: {symbolType: K_SYMBOL, payType: K_XBET}, 10: {symbolType: K_SYMBOL, payType: K_XBET}, 11: {symbolType: K_SYMBOL, payType: K_XBET},
	12: {symbolType: K_SYMBOL, payType: K_XBET}, 13: {symbolType: K_SYMBOL, payType: K_XBET}, 14: {symbolType: K_SYMBOL, payType: K_XBET},
	15: {symbolType: K_SYMBOL, payType: K_XBET}, 16: {symbolType: K_SYMBOL, payType: K_XBET}, 17: {symbolType: K_SYMBOL, payType: K_XBET},
	18: {symbolType: K_SYMBOL, payType: K_XBET}, 19: {symbolType: K_SYMBOL, payType: K_XBET},
	100: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 101: {symbolType: K_BONUS, payType: K_XBET},
	201: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 202: {symbolType: K_BONUS, payType: K_TOTAL_BONUS},
	300: {symbolType: K_SYMBOL, payType: K_XBET}, 301: {symbolType: K_SYMBOL, payType: K_XBET},
	500: {symbolType: K_SYMBOL, payType: K_XBET}, 501: {symbolType: K_SYMBOL, payType: K_XBET}, 502: {symbolType: K_SYMBOL, payType: K_XBET},
	503: {symbolType: K_SYMBOL, payType: K_XBET}, 504: {symbolType: K_SYMBOL, payType: K_XBET}, 505: {symbolType: K_SYMBOL, payType: K_XBET},
}
