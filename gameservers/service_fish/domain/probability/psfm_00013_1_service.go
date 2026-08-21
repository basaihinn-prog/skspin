package probability

import (
	"reflect"
	PSFM_00013_93_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-93-1"
	PSFM_00013_94_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-94-1"
	PSFM_00013_95_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-95-1"
	PSFM_00013_96_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-96-1"
	PSFM_00013_97_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-97-1"
	PSFM_00013_98_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-98-1"
	RKF_H5_00001_1 "serve/service_fish/domain/probability/RKF-H5-00001-1"
	RKF_H5_00001_1_MONSTER "serve/service_fish/domain/probability/RKF-H5-00001-1/RKF-H5-00001-1-MONSTER"
	RKF_H5_00001_1_RTP "serve/service_fish/domain/probability/RKF-H5-00001-1/RKF-H5-00001-1-RTP"
	"serve/service_fish/models"
	"strconv"
)

func chooseMath_psfm00013(mathModuleId string) (bsMath *RKF_H5_00001_1.BsMath, fsMath *RKF_H5_00001_1.FsMath, drbMath *RKF_H5_00001_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00013_93_1:
		return PSFM_00013_93_1.RTP93BS.RKF_H5_00001_1_BsMath, PSFM_00013_93_1.RTP93FS.RKF_H5_00001_1_FsMath, PSFM_00013_93_1.RTP93DRB.RKF_H5_00001_1_DrbMath
	case models.PSFM_00013_94_1:
		return PSFM_00013_94_1.RTP94BS.RKF_H5_00001_1_BsMath, PSFM_00013_94_1.RTP94FS.RKF_H5_00001_1_FsMath, PSFM_00013_94_1.RTP94DRB.RKF_H5_00001_1_DrbMath
	case models.PSFM_00013_95_1:
		return PSFM_00013_95_1.RTP95BS.RKF_H5_00001_1_BsMath, PSFM_00013_95_1.RTP95FS.RKF_H5_00001_1_FsMath, PSFM_00013_95_1.RTP95DRB.RKF_H5_00001_1_DrbMath
	case models.PSFM_00013_96_1:
		return PSFM_00013_96_1.RTP96BS.RKF_H5_00001_1_BsMath, PSFM_00013_96_1.RTP96FS.RKF_H5_00001_1_FsMath, PSFM_00013_96_1.RTP96DRB.RKF_H5_00001_1_DrbMath
	case models.PSFM_00013_97_1:
		return PSFM_00013_97_1.RTP97BS.RKF_H5_00001_1_BsMath, PSFM_00013_97_1.RTP97FS.RKF_H5_00001_1_FsMath, PSFM_00013_97_1.RTP97DRB.RKF_H5_00001_1_DrbMath
	case models.PSFM_00013_98_1:
		return PSFM_00013_98_1.RTP98BS.RKF_H5_00001_1_BsMath, PSFM_00013_98_1.RTP98FS.RKF_H5_00001_1_FsMath, PSFM_00013_98_1.RTP98DRB.RKF_H5_00001_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00013_1(fishId int32, rtpId string, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*RKF_H5_00001_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*RKF_H5_00001_1.FsMath)

	switch int(fishId) {
	case bsMath.Icons.Fish0.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish0.Hit(
			rtpId,
			&bsMath.Icons.Fish0,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish1.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish1.Hit(
			rtpId,
			&bsMath.Icons.Fish1,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish2.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish2.Hit(
			rtpId,
			&bsMath.Icons.Fish2,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish3.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish3.Hit(
			rtpId,
			&bsMath.Icons.Fish3,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish4.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish4.Hit(
			rtpId,
			&bsMath.Icons.Fish4,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish5.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish5.Hit(
			rtpId,
			&bsMath.Icons.Fish5,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish6.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish6.Hit(
			rtpId,
			&bsMath.Icons.Fish6,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish7.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish7.Hit(
			rtpId,
			&bsMath.Icons.Fish7,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish8.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish8.Hit(
			rtpId,
			&bsMath.Icons.Fish8,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish9.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish9.Hit(
			rtpId,
			&bsMath.Icons.Fish9,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish10.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish10.Hit(
			rtpId,
			&bsMath.Icons.Fish10,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish11.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish11.Hit(
			rtpId,
			&bsMath.Icons.Fish11,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish12.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish12.Hit(
			rtpId,
			&bsMath.Icons.Fish12,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish13.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish13.Hit(
			rtpId,
			&bsMath.Icons.Fish13,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish14.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish14.Hit(
			rtpId,
			&bsMath.Icons.Fish14,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish15.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish15.Hit(
			rtpId,
			&bsMath.Icons.Fish15,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish16.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish16.Hit(
			rtpId,
			&bsMath.Icons.Fish16,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish17.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish17.Hit(
			rtpId,
			&bsMath.Icons.Fish17,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish18.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish18.Hit(
			rtpId,
			&bsMath.Icons.Fish18,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Fish19.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Fish19.Hit(
			rtpId,
			&bsMath.Icons.Fish19,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Yggdrasil.IconID:
		_, _, v3 := RKF_H5_00001_1_MONSTER.Yggdrasil.Hit(
			rtpId,
			&bsMath.Icons.Yggdrasil,
		)

		if len(v3) > 0 {
			return newOptionBonusPay(int(fishId), 0, 100, -1, 1, v3)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.HalonSlot.IconID:
		v1, _, _, v4, v5 := RKF_H5_00001_1_MONSTER.HalonSlot.Hit(
			rtpId,
			&bsMath.Icons.HalonSlot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.ChangeLink2,
			&fsMath.Icons.Fish0,
			&fsMath.Icons.Fish4,
			&fsMath.Icons.Fish9,
			&fsMath.Icons.Fish11,
			&fsMath.Icons.Fish13,
			&fsMath.Icons.Fish14,
			&fsMath.Icons.Fish15,
			&fsMath.Icons.Fish16,
			&fsMath.Icons.Fish17,
			&fsMath.Icons.FaceOff1,
			&fsMath.Icons.FaceOff3,
			&fsMath.Icons.FaceOff5,
		)

		if v1 > 0 && len(v4) == 9 && len(v5) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), 0, 101, -1, 1, make([]interface{}, 2), v5)
			slot.ExtraData[0] = v1
			slot.ExtraData[1] = v4
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.RandomTriggerYggdrasil.IconID:
		atoiRtpId, err := strconv.Atoi(rtpId)
		if err != nil {
			atoiRtpId = 1
		}

		v1, v2 := RKF_H5_00001_1_MONSTER.RandomTriggerYggdrasil.Pick(
			atoiRtpId,
			&bsMath.Icons.RandomTriggerYggdrasil,
		)
		return newOptionBonusPay(int(fishId), 0, 102, v1, 1, v2)

	case bsMath.Icons.ThunderHammer.IconID:
		v1, v2, v3, v4 := RKF_H5_00001_1_MONSTER.ThunderHammer.Hit(
			rtpId,
			&bsMath.Icons.ThunderHammer,
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

	case bsMath.Icons.Stormbreaker.IconID:
		v1, v2, v3, v4 := RKF_H5_00001_1_MONSTER.StormBreaker.Hit(
			rtpId,
			&bsMath.Icons.Stormbreaker,
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

	case bsMath.Icons.DeepSeaTreasure.IconID:
		_, _, _, v4 := RKF_H5_00001_1_MONSTER.DeepSeaTreasure.Hit(
			rtpId,
			&bsMath.Icons.DeepSeaTreasure,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.CasketOfAncientWinters.IconID:
		_, _, _, v4 := RKF_H5_00001_1_MONSTER.CasketOfAncientWinters.Hit(
			rtpId,
			&bsMath.Icons.CasketOfAncientWinters,
		)

		if v4 > 0 {
			return newSimplePay(int(fishId), 1, v4)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.ChangeLink.ChangeIconId[0]:
		v1 := RKF_H5_00001_1_MONSTER.PoseidonChanger.Hit(&bsMath.Icons.ChangeLink)
		return newSimplePay(v1, 0, 1)

	case bsMath.Icons.Poseidon1.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Poseidon1.Hit(
			rtpId,
			&bsMath.Icons.Poseidon1,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Poseidon2.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Poseidon2.Hit(
			rtpId,
			&bsMath.Icons.Poseidon2,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Poseidon3.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Poseidon3.Hit(
			rtpId,
			&bsMath.Icons.Poseidon3,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Poseidon4.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Poseidon4.Hit(
			rtpId,
			&bsMath.Icons.Poseidon4,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	case bsMath.Icons.Poseidon5.IconID:
		v1, v2, v3 := RKF_H5_00001_1_MONSTER.Poseidon5.Hit(
			rtpId,
			&bsMath.Icons.Poseidon5,
		)
		return getProbabilityPay(int(fishId), v1, v2, v3)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00013_1_RngRtp(state int, drbMath interface{}) (rtpState int, rtpId string, bullet uint64) {
	return RKF_H5_00001_1_RTP.Service.RngRtp(state, drbMath)
}

func psfm_00013_1_useRtp(bulletType int32, rtpId string, bsMathI interface{}) string {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*RKF_H5_00001_1.BsMath)

	switch bulletType {
	case MACHINE_GUN_TYPE_ID: // Thunder Hammer
		return RKF_H5_00001_1_MONSTER.ThunderHammer.UseRtp(&bsMath.Icons.ThunderHammer)

	case SUPER_MACHINE_GUN_TYPE_ID: // Storm Breaker
		return RKF_H5_00001_1_MONSTER.StormBreaker.UseRtp(&bsMath.Icons.Stormbreaker)

	default:
		return rtpId
	}
}
