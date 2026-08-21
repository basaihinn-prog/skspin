package probability

import (
	"reflect"
	PSF_ON_00005_1 "serve/service_fish/domain/probability/PSF-ON-00005-1"
	PSF_ON_00005_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00005-1/PSF-ON-00005-1-MONSTER"
	PSF_ON_00005_1_RTP "serve/service_fish/domain/probability/PSF-ON-00005-1/PSF-ON-00005-1-RTP"
	PSFM_00006_95_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-95-1"
	PSFM_00006_96_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-96-1"
	PSFM_00006_97_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-97-1"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	"serve/service_fish/models"
	"strconv"
)

const (
	NO_TRIGGERID              = -1
	NO_BONUSTYPE              = -1
	MACHINE_GUN_TYPE_ID       = 201
	SUPER_MACHINE_GUN_TYPE_ID = 202
)

func chooseMath_psfm00006(mathModuleId string) (bsMath *PSF_ON_00005_1.BsMath, fsMath *PSF_ON_00005_1.FsMath, drbMath *PSF_ON_00005_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00006_95_1:
		return PSFM_00006_95_1.RTP95BS.PSF_ON_00005_1_BsMath, PSFM_00006_95_1.RTP95FS.PSF_ON_00005_1_FsMath, PSFM_00006_95_1.RTP95DRB.PSF_ON_00005_1_DrbMath
	case models.PSFM_00006_96_1:
		return PSFM_00006_96_1.RTP96BS.PSF_ON_00005_1_BsMath, PSFM_00006_96_1.RTP96FS.PSF_ON_00005_1_FsMath, PSFM_00006_96_1.RTP96DRB.PSF_ON_00005_1_DrbMath
	case models.PSFM_00006_97_1:
		return PSFM_00006_97_1.RTP97BS.PSF_ON_00005_1_BsMath, PSFM_00006_97_1.RTP97FS.PSF_ON_00005_1_FsMath, PSFM_00006_97_1.RTP97DRB.PSF_ON_00005_1_DrbMath
	case models.PSFM_00006_98_1:
		return PSFM_00006_98_1.RTP98BS.PSF_ON_00005_1_BsMath, PSFM_00006_98_1.RTP98FS.PSF_ON_00005_1_FsMath, PSFM_00006_98_1.RTP98DRB.PSF_ON_00005_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00006_1(fishId int32, rtpId string, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00005_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00005_1.FsMath)

	switch int(fishId) {
	case bsMath.Icons.StarFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.StarFish.Hit(
			rtpId,
			&bsMath.Icons.StarFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.JellyFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.JellyFish.Hit(
			rtpId,
			&bsMath.Icons.JellyFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.RiverShrimo.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.RiverShrimo.Hit(
			rtpId,
			&bsMath.Icons.RiverShrimo,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Guppy.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.Guppy.Hit(
			rtpId,
			&bsMath.Icons.Guppy,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.ClownFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.ClownFish.Hit(
			rtpId,
			&bsMath.Icons.ClownFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Dory.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.Dory.Hit(
			rtpId,
			&bsMath.Icons.Dory,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Grammidae.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.Grammidae.Hit(
			rtpId,
			&bsMath.Icons.Grammidae,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.RajahCichlasoma.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.RajahCichlasoma.Hit(
			rtpId,
			&bsMath.Icons.RajahCichlasoma,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.SiameseFightingFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.SiameseFightingFish.Hit(
			rtpId,
			&bsMath.Icons.SiameseFightingFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.LionFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.LionFish.Hit(
			rtpId,
			&bsMath.Icons.LionFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.SiameseTigerFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.SiameseTigerFish.Hit(
			rtpId,
			&bsMath.Icons.SiameseTigerFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.DurianFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.DurianFish.Hit(
			rtpId,
			&bsMath.Icons.DurianFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.MangoFish.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.MangoFish.Hit(
			rtpId,
			&bsMath.Icons.MangoFish,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.MuayThaiLobster.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.MuayThaiLobster.Hit(
			rtpId,
			&bsMath.Icons.MuayThaiLobster,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.HermitCrab.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.HermitCrab.Hit(
			rtpId,
			&bsMath.Icons.HermitCrab,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GoldGemTurtle.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GoldGemTurtle.Hit(
			rtpId,
			&bsMath.Icons.GoldGemTurtle,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GoldCoinToad.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GoldCoinToad.Hit(
			rtpId,
			&bsMath.Icons.GoldCoinToad,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GoldCrab.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GoldCrab.Hit(
			rtpId,
			&bsMath.Icons.GoldCrab,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GoldShark.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GoldShark.Hit(
			rtpId,
			&bsMath.Icons.GoldShark,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Crocodile.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.Crocodile.Hit(
			rtpId,
			&bsMath.Icons.Crocodile,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.RedEnvelopeTree.IconID:
		_, _, v3 := PSF_ON_00005_1_MONSTER.RedEnvelopeTree.Hit(
			rtpId,
			&bsMath.Icons.RedEnvelopeTree,
		)

		if len(v3) > 0 {
			return newOptionBonusPay(int(fishId), 0, 100, -1, 1, v3)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Slot.IconID:
		v1, _, _, v4, v5 := PSF_ON_00005_1_MONSTER.Slot.Hit(
			rtpId,
			&bsMath.Icons.Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.ChangeLink2,
			&fsMath.Icons.StarFish,
			&fsMath.Icons.ClownFish,
			&fsMath.Icons.LionFish,
			&fsMath.Icons.GiantJellyFish,
			&fsMath.Icons.Lobster,
			&fsMath.Icons.HermitCrab,
			&fsMath.Icons.GoldGemTurtle,
			&fsMath.Icons.GoldCoinToad,
			&fsMath.Icons.GoldCrab,
			&fsMath.Icons.FruitDish,
			&fsMath.Icons.APackOfBeer,
			&fsMath.Icons.GiantWhale,
		)

		if v1 > 0 && len(v4) == 9 && len(v5) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), 0, 101, -1, 1, make([]interface{}, 2), v5)
			slot.ExtraData[0] = v1
			slot.ExtraData[1] = v4
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.RandomTriggerEnvelope.IconID:
		atoiRtpId, err := strconv.Atoi(rtpId)
		if err != nil {
			atoiRtpId = 1
		}

		v1, v2 := PSF_ON_00005_1_MONSTER.RandomTriggerEnvelope.Pick(
			atoiRtpId,
			&bsMath.Icons.RandomTriggerEnvelope,
		)
		return newOptionBonusPay(int(fishId), 0, 102, v1, 1, v2)

	case bsMath.Icons.MachineGun.IconID:
		v1, v2, v3, v4 := PSF_ON_00005_1_MONSTER.MachineGun.Hit(
			rtpId,
			&bsMath.Icons.MachineGun,
		)

		if v2 == -1 || v3 == -1 {
			if v1 > 0 {
				return newOptionBonusPay(int(fishId), v1, 201, -1, 1, v4)
			}
			return newSimplePay(int(fishId), 0, 1)
		}

		if v1 > 0 {
			p := newOptionBonusPay(int(fishId), v1, 201, v3, 1, v4)
			p.ExtraTriggerBonus = append(p.ExtraTriggerBonus, newTriggerBonusPay(int(fishId), 0, 102, v3, 1))
			return p
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case bsMath.Icons.SuperMachineGun.IconID:
		v1, v2, v3, v4 := PSF_ON_00005_1_MONSTER.SuperMachineGun.Hit(
			rtpId,
			&bsMath.Icons.SuperMachineGun,
		)

		if v2 == -1 || v3 == -1 {
			if v1 > 0 {
				return newOptionBonusPay(int(fishId), v1, 202, -1, 1, v4)
			}
			return newSimplePay(int(fishId), 0, 1)
		}

		if v1 > 0 {
			p := newOptionBonusPay(int(fishId), v1, 202, v3, 1, v4)
			p.ExtraTriggerBonus = append(p.ExtraTriggerBonus, newTriggerBonusPay(int(fishId), 0, 102, v3, 1))
			return p
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case bsMath.Icons.FruitDish.IconID:
		_, _, _, v4 := PSF_ON_00005_1_MONSTER.FruitDish.Hit(
			rtpId,
			&bsMath.Icons.FruitDish,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.APackOfBeer.IconID:
		_, _, _, v4 := PSF_ON_00005_1_MONSTER.APackOfBeer.Hit(
			rtpId,
			&bsMath.Icons.APackOfBeer,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.ChangeLink.ChangeIconId[0]:
		v1 := PSF_ON_00005_1_MONSTER.GiantWhaleChanger.Hit(&bsMath.Icons.ChangeLink)
		return newSimplePay(v1, 0, 1)

	case bsMath.Icons.GiantWhale1.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GiantWhale1.Hit(
			rtpId,
			&bsMath.Icons.GiantWhale1,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GiantWhale2.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GiantWhale2.Hit(
			rtpId,
			&bsMath.Icons.GiantWhale2,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GiantWhale3.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GiantWhale3.Hit(
			rtpId,
			&bsMath.Icons.GiantWhale3,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GiantWhale4.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GiantWhale4.Hit(
			rtpId,
			&bsMath.Icons.GiantWhale4,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.GiantWhale5.IconID:
		v1, v2, v3 := PSF_ON_00005_1_MONSTER.GiantWhale5.Hit(
			rtpId,
			&bsMath.Icons.GiantWhale5,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00006_1_RngRtp(state int, drbMath interface{}) (rtpState int, rtpId string, bullet uint64) {
	return PSF_ON_00005_1_RTP.Service.RngRtp(state, drbMath)
}

func psfm_00006_1_useRtp(bulletType int32, rtpId string, bsMathI interface{}) string {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00005_1.BsMath)

	switch bulletType {
	case MACHINE_GUN_TYPE_ID:
		return PSF_ON_00005_1_MONSTER.MachineGun.UseRtp(&bsMath.Icons.MachineGun)

	case SUPER_MACHINE_GUN_TYPE_ID:
		return PSF_ON_00005_1_MONSTER.SuperMachineGun.UseRtp(&bsMath.Icons.SuperMachineGun)

	default:
		return rtpId
	}
}

func getProbabilityPay(fishId, v1, v2, v3 int) *Probability {
	if v2 == NO_TRIGGERID || v3 == NO_BONUSTYPE {
		return newSimplePay(fishId, v1, 1)
	}
	return newTriggerBonusPay(fishId, v1, v2, v3, 1)
}
