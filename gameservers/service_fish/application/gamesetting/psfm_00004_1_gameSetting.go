package gamesetting

import (
	"reflect"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	PSF_ON_00003 "serve/service_fish/domain/fish/PSF-ON-00003"
	PSF_ON_00003_1 "serve/service_fish/domain/probability/PSF-ON-00003-1"
	PSFM_00004_95_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-95-1"
	PSFM_00004_96_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-96-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00004_98_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-98-1"
	"serve/service_fish/models"
)

const (
	K_SYMBOL      = game_setting_proto.StripsRecall_SymbolDef_K_SYMBOL
	K_XBET        = game_setting_proto.StripsRecall_SymbolDef_K_xBET
	K_BONUS       = game_setting_proto.StripsRecall_SymbolDef_K_BONUS
	K_TOTAL_BONUS = game_setting_proto.StripsRecall_SymbolDef_K_xTOTALBET_BONUSTIMES
)

var math00004List = map[string]bool{
	models.PSFM_00004_95_1: true,
	models.PSFM_00004_96_1: true,
	models.PSFM_00004_97_1: true,
	models.PSFM_00004_98_1: true,
}

func psfm_00004_1_PayTable(sr *game_setting_proto.StripsRecall, mathModuleId string) bool {
	if !math00004List[mathModuleId] {
		return false
	}

	for _, fishId := range fishList00003 {
		sr.AllSymbolDef = append(sr.AllSymbolDef, psfm_00004_1_symbolDef(fishId, mathModuleId, fishSymbol00003[fishId].symbolType, fishSymbol00003[fishId].payType))
	}

	return true
}

func psfm_00004_1_symbolDef(
	fishId int,
	mathModuleId string,
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType,
	payType game_setting_proto.StripsRecall_SymbolDef_PayType,
) *game_setting_proto.StripsRecall_SymbolDef {
	sd := &game_setting_proto.StripsRecall_SymbolDef{
		SymbolId:   uint32(fishId),
		SymbolType: symbolType,
		LayerId:    psfm_00004_getObjectLayer(fishId),
		PayType:    payType,
	}

	switch mathModuleId {
	case models.PSFM_00004_95_1:
		psfm_00004_getPay(fishId, sd, PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath, PSFM_00004_95_1.RTP95FS.PSF_ON_00003_1_FsMath)
	case models.PSFM_00004_96_1:
		psfm_00004_getPay(fishId, sd, PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath, PSFM_00004_96_1.RTP96FS.PSF_ON_00003_1_FsMath)
	case models.PSFM_00004_97_1:
		psfm_00004_getPay(fishId, sd, PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath, PSFM_00004_97_1.RTP97FS.PSF_ON_00003_1_FsMath)
	case models.PSFM_00004_98_1:
		psfm_00004_getPay(fishId, sd, PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath, PSFM_00004_98_1.RTP98FS.PSF_ON_00003_1_FsMath)
	}
	return sd
}

func psfm_00004_getPay(fishId int, sd *game_setting_proto.StripsRecall_SymbolDef, bsMathI, fsMathI interface{}) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00003_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00003_1.FsMath)

	switch fishId {
	case 0:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie1.IconPays[0]))
	case 1:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie2.IconPays[0]))
	case 2:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie3.IconPays[0]))
	case 3:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie4.IconPays[0]))
	case 4:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie5.IconPays[0]))
	case 5:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie6.IconPays[0]))
	case 6:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie1Slot.IconPays[0]))
	case 7:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie2Slot.IconPays[0]))
	case 8:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie3Slot.IconPays[0]))
	case 9:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie4Slot.IconPays[0]))
	case 10:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie5Slot.IconPays[0]))
	case 11:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie6Slot.IconPays[0]))
	case 12:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie1Drill.IconPays[0]))
	case 13:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie2Drill.IconPays[0]))
	case 14:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie3Drill.IconPays[0]))
	case 15:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie4Drill.IconPays[0]))
	case 16:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie5Drill.IconPays[0]))
	case 17:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Zombie6Drill.IconPays[0]))
	case 203:
		sd.Pays = append(sd.Pays, 0)
	case 303:
		iconPays := bsMath.Icons.LittleZombie.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 401:
		iconPays := bsMath.Icons.Demon1.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 402:
		iconPays := bsMath.Icons.Demon2.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 403:
		iconPays := bsMath.Icons.Demon3.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	// Process Slot
	case 3031:
		iconPay := fsMath.MixAward[0].MixIconGroup[0].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 3032:
		iconPay := fsMath.MixAward[0].MixIconGroup[1].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 3033:
		iconPay := fsMath.MixAward[0].MixIconGroup[2].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 3034:
		iconPay := fsMath.MixAward[0].MixIconGroup[3].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 4031:
		iconPay := fsMath.MixAward[0].MixIconGroup[10].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 4032:
		iconPay := fsMath.MixAward[0].MixIconGroup[11].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))

	case 4033:
		iconPay := fsMath.MixAward[0].MixIconGroup[12].Quantity
		sd.Pays = append(sd.Pays, int64(iconPay))
	}
}

func psfm_00004_getObjectLayer(fishId int) uint32 {
	switch fishId {
	case 3031, 3032, 3033, 3034:
		fallthrough
	case 4031, 4032, 4033:
		return 0
	default:
		return uint32(PSF_ON_00003.Objects.Layer(fishId))
	}
}

var fishList00003 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 203, 303, 401, 402, 403, 3031, 3032, 3033, 3034, 4031, 4032, 4033}
var fishSymbol00003 = map[int]symbol{
	0: {symbolType: K_SYMBOL, payType: K_XBET}, 1: {symbolType: K_SYMBOL, payType: K_XBET}, 2: {symbolType: K_SYMBOL, payType: K_XBET},
	3: {symbolType: K_SYMBOL, payType: K_XBET}, 4: {symbolType: K_SYMBOL, payType: K_XBET}, 5: {symbolType: K_SYMBOL, payType: K_XBET},
	6: {symbolType: K_BONUS, payType: K_XBET}, 7: {symbolType: K_BONUS, payType: K_XBET}, 8: {symbolType: K_BONUS, payType: K_XBET},
	9: {symbolType: K_BONUS, payType: K_XBET}, 10: {symbolType: K_BONUS, payType: K_XBET}, 11: {symbolType: K_BONUS, payType: K_XBET},
	12: {symbolType: K_BONUS, payType: K_XBET}, 13: {symbolType: K_BONUS, payType: K_XBET}, 14: {symbolType: K_BONUS, payType: K_XBET},
	15: {symbolType: K_BONUS, payType: K_XBET}, 16: {symbolType: K_BONUS, payType: K_XBET}, 17: {symbolType: K_BONUS, payType: K_XBET},
	203: {symbolType: K_BONUS, payType: K_XBET}, 303: {symbolType: K_BONUS, payType: K_XBET},
	401: {symbolType: K_BONUS, payType: K_XBET}, 402: {symbolType: K_BONUS, payType: K_XBET}, 403: {symbolType: K_BONUS, payType: K_XBET},
	3031: {symbolType: K_SYMBOL, payType: K_XBET}, 3032: {symbolType: K_SYMBOL, payType: K_XBET}, 3033: {symbolType: K_SYMBOL, payType: K_XBET}, 3034: {symbolType: K_SYMBOL, payType: K_XBET},
	4031: {symbolType: K_SYMBOL, payType: K_XBET}, 4032: {symbolType: K_SYMBOL, payType: K_XBET}, 4033: {symbolType: K_SYMBOL, payType: K_XBET},
}
