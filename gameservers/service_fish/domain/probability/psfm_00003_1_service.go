package probability

import (
	"reflect"
	PSF_ON_00002_1 "serve/service_fish/domain/probability/PSF-ON-00002-1"
	PSF_ON_00002_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00002-1/PSF-ON-00002-1-MONSTER"
	PSF_ON_00002_1_RTP "serve/service_fish/domain/probability/PSF-ON-00002-1/PSF-ON-00002-1-RTP"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	PSFM_00003_94_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-94-1"
	PSFM_00003_95_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-95-1"
	PSFM_00003_96_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-96-1"
	PSFM_00003_97_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-97-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	"serve/service_fish/models"
	"strconv"
)

func chooseMath_psfm00003(mathModuleId string) (bsMath *PSF_ON_00002_1.BsMath, fsMath *PSF_ON_00002_1.FsMath, drbMath *PSF_ON_00002_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00003_93_1:
		return PSFM_00003_93_1.RTP93BS.PSF_ON_00002_1_BsMath, PSFM_00003_93_1.RTP93FS.PSF_ON_00002_1_FsMath, PSFM_00003_93_1.RTP93DRB.PSF_ON_00002_1_DrbMath
	case models.PSFM_00003_94_1:
		return PSFM_00003_94_1.RTP94BS.PSF_ON_00002_1_BsMath, PSFM_00003_94_1.RTP94FS.PSF_ON_00002_1_FsMath, PSFM_00003_94_1.RTP94DRB.PSF_ON_00002_1_DrbMath
	case models.PSFM_00003_95_1:
		return PSFM_00003_95_1.RTP95BS.PSF_ON_00002_1_BsMath, PSFM_00003_95_1.RTP95FS.PSF_ON_00002_1_FsMath, PSFM_00003_95_1.RTP95DRB.PSF_ON_00002_1_DrbMath
	case models.PSFM_00003_96_1:
		return PSFM_00003_96_1.RTP96BS.PSF_ON_00002_1_BsMath, PSFM_00003_96_1.RTP96FS.PSF_ON_00002_1_FsMath, PSFM_00003_96_1.RTP96DRB.PSF_ON_00002_1_DrbMath
	case models.PSFM_00003_97_1:
		return PSFM_00003_97_1.RTP97BS.PSF_ON_00002_1_BsMath, PSFM_00003_97_1.RTP97FS.PSF_ON_00002_1_FsMath, PSFM_00003_97_1.RTP97DRB.PSF_ON_00002_1_DrbMath
	case models.PSFM_00003_98_1:
		return PSFM_00003_98_1.RTP98BS.PSF_ON_00002_1_BsMath, PSFM_00003_98_1.RTP98FS.PSF_ON_00002_1_FsMath, PSFM_00003_98_1.RTP98DRB.PSF_ON_00002_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00003_1(fishId int32, rtpId string, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00002_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00002_1.FsMath)

	switch fishId {
	case 0:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.StarFish.Hit(
			rtpId,
			&bsMath.Icons.StarFish,
		)
		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 1:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.SeaHorse.Hit(
			rtpId,
			&bsMath.Icons.SeaHorse,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 2:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.SmallJellyfish.Hit(
			rtpId,
			&bsMath.Icons.SmallJellyFish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 3:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Guppy.Hit(
			rtpId,
			&bsMath.Icons.Guppy,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 4:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Clownfish.Hit(
			rtpId,
			&bsMath.Icons.ClownFish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 5:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Dory.Hit(
			rtpId,
			&bsMath.Icons.Dory,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 6:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Grammidae.Hit(
			rtpId,
			&bsMath.Icons.Grammidae,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 7:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.RajahCichlasoma.Hit(
			rtpId,
			&bsMath.Icons.RajahCichlasoma,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 8:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.PufferFish.Hit(
			rtpId,
			&bsMath.Icons.PufferFish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 9:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.LionFish.Hit(
			rtpId,
			&bsMath.Icons.LionFish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 10:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Lanternfish.Hit(
			rtpId,
			&bsMath.Icons.Lanternfish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 11:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.GiantJellyfish.Hit(
			rtpId,
			&bsMath.Icons.GiantJellyfish,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 12:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Turtle.Hit(
			rtpId,
			&bsMath.Icons.Turtle,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 13:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.Lobster.Hit(
			rtpId,
			&bsMath.Icons.Lobster,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 14:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.HermitCrab.Hit(
			rtpId,
			&bsMath.Icons.HermitCrab,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 15:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.GoldGemTurtle.Hit(
			rtpId,
			&bsMath.Icons.GoldGemTurtle,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 16:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.GoldCoinToad.Hit(
			rtpId,
			&bsMath.Icons.GoldCoinToad,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 17:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.GoldCrab.Hit(
			rtpId,
			&bsMath.Icons.GoldCrab,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 18:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.GoldShark.Hit(
			rtpId,
			&bsMath.Icons.GoldShark,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 19:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.AsianArowana.Hit(
			rtpId,
			&bsMath.Icons.AsianArowana,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 20:
		v1, v2, v3, v4 := PSF_ON_00002_1_MONSTER.HotPot.Hit(
			rtpId,
			&bsMath.Icons.HotPot,
		)

		if v2 == -1 || v3 == -1 {
			if v1 > 0 {
				return newSimplePay(int(fishId), v1, v4)
			}
			return newSimplePay(int(fishId), v1, 1)
		}

		if v1 > 0 {
			p := newSimplePay(int(fishId), v1, v4)
			p.ExtraTriggerBonus = append(p.ExtraTriggerBonus, newTriggerBonusPay(int(fishId), 0, v2, v3, 1))
			return p
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 21:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.FaceOff1.Hit(
			rtpId,
			&bsMath.Icons.FaceOff1,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 22:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.FaceOff2.Hit(
			rtpId,
			&bsMath.Icons.FaceOff2,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 23:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.FaceOff3.Hit(
			rtpId,
			&bsMath.Icons.FaceOff3,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 24:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.FaceOff4.Hit(
			rtpId,
			&bsMath.Icons.FaceOff4,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 25:
		v1, v2, v3 := PSF_ON_00002_1_MONSTER.FaceOff5.Hit(
			rtpId,
			&bsMath.Icons.FaceOff5,
		)

		if v2 == -1 || v3 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 26:
		v1, v2, v3, v4 := PSF_ON_00002_1_MONSTER.Drill.Hit(
			rtpId,
			&bsMath.Icons.Drill,
		)
		// v2:triggerIconId, v3:bonusTypeId,

		if v2 == -1 || v3 == -1 {
			if v1 > 0 {
				return newOptionBonusPay(int(fishId), v1, 26, -1, 1, v4)
			}
			return newSimplePay(int(fishId), 0, 1)
		}

		if v1 > 0 {
			p := newOptionBonusPay(int(fishId), v1, 26, v3, 1, v4)
			p.ExtraTriggerBonus = append(p.ExtraTriggerBonus, newTriggerBonusPay(int(fishId), 0, v2, v3, 1))
			return p
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 27:
		v1, v2, v3, v4 := PSF_ON_00002_1_MONSTER.MachineGun.Hit(
			rtpId,
			&bsMath.Icons.MachineGun,
		)

		if v2 == -1 || v3 == -1 {
			if v1 > 0 {
				return newOptionBonusPay(int(fishId), v1, 27, -1, 1, v4)
			}
			return newSimplePay(int(fishId), 0, 1)
		}

		if v1 > 0 {
			p := newOptionBonusPay(int(fishId), v1, 27, v3, 1, v4)
			p.ExtraTriggerBonus = append(p.ExtraTriggerBonus, newTriggerBonusPay(int(fishId), 0, 31, v3, 1))
			return p
		}
		return newTriggerBonusPay(int(fishId), v1, v2, v3, 1)

	case 28:
		_, _, v3 := PSF_ON_00002_1_MONSTER.RedEnvelope.Hit(
			rtpId,
			&bsMath.Icons.RedEnvelope,
		)

		if len(v3) > 0 {
			return newOptionBonusPay(int(fishId), 0, 28, -1, 1, v3)
		}
		return newSimplePay(int(fishId), 0, 1)

	case 29:
		v1, _, _, v4, v5 := PSF_ON_00002_1_MONSTER.Slot.Hit(
			rtpId,
			&bsMath.Icons.Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsReelMath,
			&fsMath.Changelink1,
			&fsMath.Changelink2,
			&fsMath.Icons.StarFish,
			&fsMath.Icons.Clownfish,
			&fsMath.Icons.LionFish,
			&fsMath.Icons.GiantJellyfish,
			&fsMath.Icons.Lobster,
			&fsMath.Icons.HermitCrab,
			&fsMath.Icons.GoldGemTurtle,
			&fsMath.Icons.GoldCoinToad,
			&fsMath.Icons.GoldCrab,
			&fsMath.Icons.FaceOff1,
			&fsMath.Icons.FaceOff3,
			&fsMath.Icons.FaceOff5,
		)

		if v1 > 0 && len(v4) == 9 && len(v5) == 3 {
			// client don't want to get pay when hit slot
			slot := newExtraOptionBonusPay(int(fishId), 0, 29, -1, 1, make([]interface{}, 2), v5)
			slot.ExtraData[0] = v1
			slot.ExtraData[1] = v4
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case 30:
		v1 := PSF_ON_00002_1_MONSTER.FaceOffChanger.Hit(&bsMath.ChangeLink)
		return newSimplePay(v1, 0, 1)

	case 31:
		atoiRtpId, err := strconv.Atoi(rtpId)

		if err != nil {
			atoiRtpId = 1
		}

		v1, v2 := PSF_ON_00002_1_MONSTER.Envelope.Pick(
			atoiRtpId,
			&bsMath.Icons.Envelope,
		)
		return newOptionBonusPay(int(fishId), 0, 31, v1, 1, v2)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00003_1_RngRtp(state int, drbMath interface{}) (rtpState int, rtpId string, rtpBullets uint64) {
	return PSF_ON_00002_1_RTP.Service.RngRtp(state, drbMath)
}
