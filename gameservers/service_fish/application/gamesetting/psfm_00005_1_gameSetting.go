package gamesetting

import (
	"reflect"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	PSF_ON_00004 "serve/service_fish/domain/fish/PSF-ON-00004"
	PSF_ON_00004_1 "serve/service_fish/domain/probability/PSF-ON-00004-1"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00005_96_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-96-1"
	PSFM_00005_97_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-97-1"
	PSFM_00005_98_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-98-1"
	"serve/service_fish/models"
)

var math00005List = map[string]bool{
	models.PSFM_00005_95_1: true,
	models.PSFM_00005_96_1: true,
	models.PSFM_00005_97_1: true,
	models.PSFM_00005_98_1: true,
}

func psfm_00005_1_PayTable(sr *game_setting_proto.StripsRecall, mathModuleId string) bool {
	if !math00005List[mathModuleId] {
		return false
	}

	for _, fishId := range fishList00004 {
		sr.AllSymbolDef = append(sr.AllSymbolDef, psfm_00005_1_symbolDef(fishId, mathModuleId, fishSymbol00004[fishId].symbolType, fishSymbol00004[fishId].payType))
	}

	return true
}

func psfm_00005_1_symbolDef(
	fishId int,
	mathModuleId string,
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType,
	payType game_setting_proto.StripsRecall_SymbolDef_PayType,
) *game_setting_proto.StripsRecall_SymbolDef {
	sd := &game_setting_proto.StripsRecall_SymbolDef{
		SymbolId:   uint32(fishId),
		SymbolType: symbolType,
		LayerId:    psfm_00005_getObjectLayer(fishId),
		PayType:    payType,
	}

	switch mathModuleId {
	case models.PSFM_00005_95_1:
		psfm_00005_getPay(fishId, sd, PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath)
	case models.PSFM_00005_96_1:
		psfm_00005_getPay(fishId, sd, PSFM_00005_96_1.RTP96BS.PSF_ON_00004_1_BsMath)
	case models.PSFM_00005_97_1:
		psfm_00005_getPay(fishId, sd, PSFM_00005_97_1.RTP97BS.PSF_ON_00004_1_BsMath)
	case models.PSFM_00005_98_1:
		psfm_00005_getPay(fishId, sd, PSFM_00005_98_1.RTP98BS.PSF_ON_00004_1_BsMath)
	}
	return sd
}

func psfm_00005_getPay(fishId int, sd *game_setting_proto.StripsRecall_SymbolDef, bsMathI interface{}) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00004_1.BsMath)

	switch fishId {
	case 0:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Clam.IconPays[0]))
	case 1:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Shrimp.IconPays[0]))
	case 2:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Surmullet.IconPays[0]))
	case 3:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Squid.IconPays[0]))
	case 4:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FlyingFish.IconPays[0]))
	case 5:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Halibut.IconPays[0]))
	case 6:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.ButterflyFish.IconPays[0]))
	case 7:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Oplegnathus.IconPays[0]))
	case 8:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Snapper.IconPays[0]))
	case 9:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.MahiMahi.IconPays[0]))
	case 10:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Stingray.IconPays[0]))
	case 11:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.LobsterWaiter.IconPays[0]))
	case 12:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.PenguinWaiter.IconPays[0]))
	case 13:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.PlatypusSeniorChef.IconPays[0]))
	case 14:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SeaLionChef.IconPays[0]))
	case 15:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.HairtCrab.IconPays[0]))
	case 16:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SnaggletoothShark.IconPays[0]))
	case 17:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SwordFish.IconPays[0]))
	case 18:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.WhaleShark.IconPays[0]))
	case 19:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.GiantOarfish.IconPays[0]))
	case 100:
		sd.Pays = append(sd.Pays, 0)

		lowIconPays := bsMath.Icons.RedEnvelope.RTP1.LowIconPays
		highIconPays := bsMath.Icons.RedEnvelope.RTP1.HighIconPays

		for _, v := range lowIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

		for _, v := range highIconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 101:
		// slot do nothing

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
		iconPays := bsMath.Icons.LobsterDash.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 301:
		iconPays := bsMath.Icons.FruitDash.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))

	case 302:
		sd.Pays = append(sd.Pays, 0)
		iconPays := bsMath.Icons.XiaoLongBao.RTP1.Gain1.IconPays

		for _, v := range iconPays {
			sd.Rounds = append(sd.Rounds, int32(v))
		}

	case 400:
		iconPays := bsMath.Icons.WhiteTigerChef.RTP1.IconPays
		sd.Pays = append(sd.Pays, int64(iconPays[len(iconPays)-1]))
	}
}

func psfm_00005_getObjectLayer(fishId int) uint32 {
	return uint32(PSF_ON_00004.Objects.Layer(fishId))
}

var fishList00004 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 100, 101, 201, 202, 300, 301, 302, 400}
var fishSymbol00004 = map[int]symbol{
	0: {symbolType: K_SYMBOL, payType: K_XBET}, 1: {symbolType: K_SYMBOL, payType: K_XBET}, 2: {symbolType: K_SYMBOL, payType: K_XBET},
	3: {symbolType: K_SYMBOL, payType: K_XBET}, 4: {symbolType: K_SYMBOL, payType: K_XBET}, 5: {symbolType: K_SYMBOL, payType: K_XBET},
	6: {symbolType: K_SYMBOL, payType: K_XBET}, 7: {symbolType: K_SYMBOL, payType: K_XBET}, 8: {symbolType: K_SYMBOL, payType: K_XBET},
	9: {symbolType: K_SYMBOL, payType: K_XBET}, 10: {symbolType: K_SYMBOL, payType: K_XBET}, 11: {symbolType: K_SYMBOL, payType: K_XBET},
	12: {symbolType: K_SYMBOL, payType: K_XBET}, 13: {symbolType: K_SYMBOL, payType: K_XBET}, 14: {symbolType: K_SYMBOL, payType: K_XBET},
	15: {symbolType: K_SYMBOL, payType: K_XBET}, 16: {symbolType: K_SYMBOL, payType: K_XBET}, 17: {symbolType: K_SYMBOL, payType: K_XBET},
	18: {symbolType: K_SYMBOL, payType: K_XBET}, 19: {symbolType: K_SYMBOL, payType: K_XBET},
	100: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 101: {symbolType: K_BONUS, payType: K_XBET},
	201: {symbolType: K_BONUS, payType: K_TOTAL_BONUS}, 202: {symbolType: K_BONUS, payType: K_TOTAL_BONUS},
	300: {symbolType: K_SYMBOL, payType: K_XBET}, 301: {symbolType: K_SYMBOL, payType: K_XBET}, 302: {symbolType: K_BONUS, payType: K_TOTAL_BONUS},
	400: {symbolType: K_BONUS, payType: K_XBET},
}
