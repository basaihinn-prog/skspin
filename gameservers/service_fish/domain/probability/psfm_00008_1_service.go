package probability

import (
	"reflect"
	PSF_ON_00007_1 "serve/service_fish/domain/probability/PSF-ON-00007-1"
	PSF_ON_00007_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00007-1/PSF-ON-00007-1-MONSTER"
	PSF_ON_00007_1_RTP "serve/service_fish/domain/probability/PSF-ON-00007-1/PSF-ON-00007-1-RTP"
	PSFM_00008_95_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-95-1"
	PSFM_00008_96_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-96-1"
	PSFM_00008_97_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-97-1"
	PSFM_00008_98_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-98-1"
	"serve/service_fish/models"
	"strconv"
)

func chooseMath_psfm00008(mathModuleId string) (bsMath *PSF_ON_00007_1.BsMath, fsMath *PSF_ON_00007_1.FsMath, drbMath *PSF_ON_00007_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00008_95_1:
		return PSFM_00008_95_1.RTP95BS.PSF_ON_00007_1_BsMath, PSFM_00008_95_1.RTP95FS.PSF_ON_00007_1_FsMath, PSFM_00008_95_1.RTP95DRB.PSF_ON_00007_1_DrbMath
	case models.PSFM_00008_96_1:
		return PSFM_00008_96_1.RTP96BS.PSF_ON_00007_1_BsMath, PSFM_00008_96_1.RTP96FS.PSF_ON_00007_1_FsMath, PSFM_00008_96_1.RTP96DRB.PSF_ON_00007_1_DrbMath
	case models.PSFM_00008_97_1:
		return PSFM_00008_97_1.RTP97BS.PSF_ON_00007_1_BsMath, PSFM_00008_97_1.RTP97FS.PSF_ON_00007_1_FsMath, PSFM_00008_97_1.RTP97DRB.PSF_ON_00007_1_DrbMath
	case models.PSFM_00008_98_1:
		return PSFM_00008_98_1.RTP98BS.PSF_ON_00007_1_BsMath, PSFM_00008_98_1.RTP98FS.PSF_ON_00007_1_FsMath, PSFM_00008_98_1.RTP98DRB.PSF_ON_00007_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00008_1(fishId int32, rtpId string, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00007_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00007_1.FsMath)

	switch int(fishId) {
	case bsMath.Icons.Fish0.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish0.Hit(
			rtpId,
			&bsMath.Icons.Fish0,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish1.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish1.Hit(
			rtpId,
			&bsMath.Icons.Fish1,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish2.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish2.Hit(
			rtpId,
			&bsMath.Icons.Fish2,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish3.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish3.Hit(
			rtpId,
			&bsMath.Icons.Fish3,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish4.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish4.Hit(
			rtpId,
			&bsMath.Icons.Fish4,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish5.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish5.Hit(
			rtpId,
			&bsMath.Icons.Fish5,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish6.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish6.Hit(
			rtpId,
			&bsMath.Icons.Fish6,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish7.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish7.Hit(
			rtpId,
			&bsMath.Icons.Fish7,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish8.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish8.Hit(
			rtpId,
			&bsMath.Icons.Fish8,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish9.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish9.Hit(
			rtpId,
			&bsMath.Icons.Fish9,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish10.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish10.Hit(
			rtpId,
			&bsMath.Icons.Fish10,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish11.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish11.Hit(
			rtpId,
			&bsMath.Icons.Fish11,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish12.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish12.Hit(
			rtpId,
			&bsMath.Icons.Fish12,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish13.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish13.Hit(
			rtpId,
			&bsMath.Icons.Fish13,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish14.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish14.Hit(
			rtpId,
			&bsMath.Icons.Fish14,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish15.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish15.Hit(
			rtpId,
			&bsMath.Icons.Fish15,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish16.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish16.Hit(
			rtpId,
			&bsMath.Icons.Fish16,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish17.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish17.Hit(
			rtpId,
			&bsMath.Icons.Fish17,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish18.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish18.Hit(
			rtpId,
			&bsMath.Icons.Fish18,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish19.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.Fish19.Hit(
			rtpId,
			&bsMath.Icons.Fish19,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.MagicShellFaFaFa.IconID:
		_, _, v3 := PSF_ON_00007_1_MONSTER.MagicShellFaFaFa.Hit(
			rtpId,
			&bsMath.Icons.MagicShellFaFaFa,
		)

		if len(v3) > 0 {
			return newOptionBonusPay(int(fishId), 0, 100, -1, 1, v3)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.LuckSlot.IconID:
		v1, _, _, v4, v5 := PSF_ON_00007_1_MONSTER.LuckSlot.Hit(
			rtpId,
			&bsMath.Icons.LuckSlot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.Fish0,
			&fsMath.Icons.Fish4,
			&fsMath.Icons.Fish9,
			&fsMath.Icons.Fish11,
			&fsMath.Icons.Fish13,
			&fsMath.Icons.Fish14,
			&fsMath.Icons.Fish15,
			&fsMath.Icons.Fish16,
			&fsMath.Icons.Fish17,
			&fsMath.Icons.PirateShip,
			&fsMath.Icons.TreasureDish,
			&fsMath.Icons.OceanKing,
		)

		if v1 > 0 && len(v4) == 12 && len(v5) == 4 {
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

		v1, v2 := PSF_ON_00007_1_MONSTER.RandomTriggerEnvelope.Pick(
			atoiRtpId,
			&bsMath.Icons.RandomTriggerEnvelope,
		)
		return newOptionBonusPay(int(fishId), 0, 102, v1, 1, v2)

	case bsMath.Icons.MachineGun.IconID:
		v1, v2, v3, v4 := PSF_ON_00007_1_MONSTER.MachineGun.Hit(
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
		v1, v2, v3, v4 := PSF_ON_00007_1_MONSTER.SuperMachineGun.Hit(
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

	case bsMath.Icons.PirateShip.IconID:
		_, _, _, v4 := PSF_ON_00007_1_MONSTER.PirateShip.Hit(
			rtpId,
			&bsMath.Icons.PirateShip,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.TreasureDish.IconID:
		_, _, _, v4 := PSF_ON_00007_1_MONSTER.TreasureDish.Hit(
			rtpId,
			&bsMath.Icons.TreasureDish,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.TreasureHermitCrab.IconID:
		v1, v2, v3, _, _ := PSF_ON_00007_1_MONSTER.TreasureHermitCrab.Hit(
			rtpId,
			&bsMath.Icons.TreasureHermitCrab,
		)

		if v1 > 0 {
			xiaoLongBao := newExtraOptionBonusPay(int(fishId), v1, 302, -1, 1, make([]interface{}, 2), v2)
			xiaoLongBao.ExtraData[0] = v2
			xiaoLongBao.ExtraData[1] = v3
			return xiaoLongBao
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.OceanKing.IconID:
		v1, v2, v3 := PSF_ON_00007_1_MONSTER.OceanKing.Hit(
			rtpId,
			&bsMath.Icons.OceanKing,
		)

		return getProbabilityPay(int(fishId), v1, v2, v3)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00008_1_RngRtp(state int, drbMath interface{}) (rtpState int, rtpId string, bullet uint64) {
	return PSF_ON_00007_1_RTP.Service.RngRtp(state, drbMath)
}

func psfm_00008_1_useRtp(bulletType int32, rtpId string, bsMathI interface{}) string {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00007_1.BsMath)

	switch bulletType {
	case MACHINE_GUN_TYPE_ID:
		return PSF_ON_00007_1_MONSTER.MachineGun.UseRtp(&bsMath.Icons.MachineGun)

	case SUPER_MACHINE_GUN_TYPE_ID:
		return PSF_ON_00007_1_MONSTER.SuperMachineGun.UseRtp(&bsMath.Icons.SuperMachineGun)

	default:
		return rtpId
	}
}
