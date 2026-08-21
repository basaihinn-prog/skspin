package gamesetting

import (
	"reflect"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	PSF_ON_00001 "serve/service_fish/domain/fish/PSF-ON-00001"
	PSF_ON_00001_1 "serve/service_fish/domain/probability/PSF-ON-00001-1"
	PSFM_00001_98_1 "serve/service_fish/domain/probability/PSFM-00001-1/PSFM-00001-98-1"
	"serve/service_fish/models"
)

var math00001List = map[string]bool{
	models.PSFM_00001_98_1: true,
}

func psfm_00001_1_PayTable(sr *game_setting_proto.StripsRecall, mathModuleId string) bool {
	if !math00001List[mathModuleId] {
		return false
	}

	for _, fishId := range fishList00001 {
		sr.AllSymbolDef = append(sr.AllSymbolDef, psfm_00001_1_symbolDef(fishId, mathModuleId, fishSymbol00001[fishId].symbolType, fishSymbol00001[fishId].payType))
	}

	return true
}

func psfm_00001_1_symbolDef(
	fishId int,
	mathModuleId string,
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType,
	payType game_setting_proto.StripsRecall_SymbolDef_PayType,
) *game_setting_proto.StripsRecall_SymbolDef {
	sd := &game_setting_proto.StripsRecall_SymbolDef{
		SymbolId:   uint32(fishId),
		SymbolType: symbolType,
		LayerId:    psfm_00001_getObjectLayer(fishId),
		PayType:    payType,
	}

	switch mathModuleId {
	case models.PSFM_00001_98_1:
		psfm_00001_getPay(fishId, sd, PSFM_00001_98_1.RTP98BS.PSF_ON_00001_1_BsMath)
	}
	return sd
}

func psfm_00001_getPay(fishId int, sd *game_setting_proto.StripsRecall_SymbolDef, bsMathI interface{}) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00001_1.BsMath)

	switch fishId {
	case 0:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.StarFish.IconPays[0]))
	case 1:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.SeaHorse.IconPays[0]))
	case 2:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Guppy.IconPays[0]))
	case 3:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.ClownFish.IconPays[0]))
	case 4:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Dory.IconPays[0]))
	case 5:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Grammidae.IconPays[0]))
	case 6:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.RajahCichlasoma.IconPays[0]))
	case 7:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.PufferFish.IconPays[0]))
	case 8:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.LanternFish.IconPays[0]))
	case 9:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.LionFish.IconPays[0]))
	case 10:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Turtle.IconPays[0]))
	case 11:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Lobster.IconPays[0]))
	case 12:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Penguin.IconPays[0]))
	case 13:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Platypus.IconPays[0]))
	case 14:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Manatee.IconPays[0]))
	case 15:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Dolphin.IconPays[0]))
	case 16:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Koi.IconPays[0]))
	case 17:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.HermitCrab.IconPays[0]))
	case 18:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.Shark.IconPays[0]))
	case 19:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.AsianArowana.IconPays[0]))
	case 20:
		iconPays := bsMath.Icons.HotPot.IconPays

		for v := range iconPays {
			sd.Pays = append(sd.Pays, int64(v))
		}

	case 21:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff1.IconPays[0]))
	case 22:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff2.IconPays[0]))
	case 23:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff3.IconPays[0]))
	case 24:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff4.IconPays[0]))
	case 25:
		sd.Pays = append(sd.Pays, int64(bsMath.Icons.FaceOff5.IconPays[0]))
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

func psfm_00001_getObjectLayer(fishId int) uint32 {
	switch fishId {
	case 21, 22, 23, 24, 25, 30:
		return uint32(PSF_ON_00001.Objects.Layer(30))
	default:
		return uint32(PSF_ON_00001.Objects.Layer(fishId))
	}
}
