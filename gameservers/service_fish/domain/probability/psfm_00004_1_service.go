package probability

import (
	"reflect"
	PSF_ON_00003_1 "serve/service_fish/domain/probability/PSF-ON-00003-1"
	PSF_ON_00003_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00003-1/PSF-ON-00003-1-MONSTER"
	PSF_ON_00003_1_RTP "serve/service_fish/domain/probability/PSF-ON-00003-1/PSF-ON-00003-1-RTP"
	PSFM_00004_95_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-95-1"
	PSFM_00004_96_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-96-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00004_98_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-98-1"
	"serve/service_fish/models"
)

func chooseMath_psfm00004(mathModuleId string) (bsMath *PSF_ON_00003_1.BsMath, fsMath *PSF_ON_00003_1.FsMath, drbMath *PSF_ON_00003_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00004_95_1:
		return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath, PSFM_00004_95_1.RTP95FS.PSF_ON_00003_1_FsMath, PSFM_00004_95_1.RTP95DRB.PSF_ON_00003_1_DrbMath
	case models.PSFM_00004_96_1:
		return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath, PSFM_00004_96_1.RTP96FS.PSF_ON_00003_1_FsMath, PSFM_00004_96_1.RTP96DRB.PSF_ON_00003_1_DrbMath
	case models.PSFM_00004_97_1:
		return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath, PSFM_00004_97_1.RTP97FS.PSF_ON_00003_1_FsMath, PSFM_00004_97_1.RTP97DRB.PSF_ON_00003_1_DrbMath
	case models.PSFM_00004_98_1:
		return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath, PSFM_00004_98_1.RTP98FS.PSF_ON_00003_1_FsMath, PSFM_00004_98_1.RTP98DRB.PSF_ON_00003_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00004_1(fishId int32, rtpId string, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00003_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00003_1.FsMath)

	switch int(fishId) {
	case bsMath.Icons.Zombie1.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie1.Hit(
			rtpId,
			&bsMath.Icons.Zombie1,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie2.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie2.Hit(
			rtpId,
			&bsMath.Icons.Zombie2,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie3.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie3.Hit(
			rtpId,
			&bsMath.Icons.Zombie3,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie4.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie4.Hit(
			rtpId,
			&bsMath.Icons.Zombie4,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie5.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie5.Hit(
			rtpId,
			&bsMath.Icons.Zombie5,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie6.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.Zombie6.Hit(
			rtpId,
			&bsMath.Icons.Zombie6,
			&bsMath.Icons.Bullets,
		)
		if v2 > 0 {
			return newSimpleBullet(int(fishId), v1, v2, 31, -1, 1)
		}
		return newSimplePay(int(fishId), v1, 1)

	case bsMath.Icons.Zombie1Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie1Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie1Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie2Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie2Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie2Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie3Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie3Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie3Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie4Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie4Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie4Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie5Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie5Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie5Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie6Slot.IconID:
		v1, v2, v3, v4, v5 := PSF_ON_00003_1_MONSTER.Zombie6Slot.Hit(
			rtpId,
			&bsMath.Icons.Zombie6Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.Icons.LittleZombie1,
			&fsMath.Icons.LittleZombie2,
			&fsMath.Icons.LittleZombie3,
			&fsMath.Icons.LittleZombie4,
			&fsMath.Icons.Zombie1,
			&fsMath.Icons.Zombie2,
			&fsMath.Icons.Zombie3,
			&fsMath.Icons.Zombie4,
			&fsMath.Icons.Zombie5,
			&fsMath.Icons.Zombie6,
			&fsMath.Icons.Demon31,
			&fsMath.Icons.Demon32,
			&fsMath.Icons.Demon33,
			&bsMath.Icons.Bullets,
		)
		if v5 > 0 {
			if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
				slot := newExtraOptionBonusBullet(int(fishId), v1, 101, -1, 1, v5, make([]interface{}, 2), v4)
				slot.ExtraData[0] = v2
				slot.ExtraData[1] = v3
				return slot
			} else {
				return newSimpleBullet(int(fishId), 0, v5, 31, -1, 1)
			}
		}
		if v2 > 0 && len(v3) == 9 && len(v4) == 3 {
			slot := newExtraOptionBonusPay(int(fishId), v1, 101, -1, 1, make([]interface{}, 2), v4)
			slot.ExtraData[0] = v2
			slot.ExtraData[1] = v3
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie1Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie1Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie1Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie2Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie2Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie2Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie3Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie3Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie3Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie4Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie4Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie4Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie5Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie5Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie5Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Zombie6Drill.IconID:
		v1, v2, v3 := PSF_ON_00003_1_MONSTER.Zombie6Drill.Hit(
			rtpId,
			&bsMath.Icons.Zombie6Drill,
			&bsMath.Icons.Bullets,
		)
		if v3 > 0 {
			if v1 > 0 {
				return newOptionBonusBullet(int(fishId), v1, 200, -1, 1, v3, v2)
			}
			return newOptionBonusBullet(int(fishId), v1, -1, -1, 1, v3, v2)
		}
		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 200, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.BonusZombie.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.BonusZombie.Hit(
			rtpId,
			&bsMath.Icons.BonusZombie,
		)

		if v1 > 0 {
			return newSimpleBullet(int(fishId), 0, v1, v2, -1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.LittleZombie.IconID:
		v1, v2 := PSF_ON_00003_1_MONSTER.LittleZombie.Hit(
			rtpId,
			&bsMath.Icons.LittleZombie,
		)
		if v1 > 0 && v2 <= 0 {
			return newSimplePay(int(fishId), v1, 1)
		}
		if v1 <= 0 && v2 > 0 {
			return newSimpleBullet(int(fishId), 0, v2, -1, -1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Demon1.IconID:
		v1 := PSF_ON_00003_1_MONSTER.Demon1.Hit(
			rtpId,
			&bsMath.Icons.Demon1,
		)
		if v1 > 0 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Demon2.IconID:
		v1 := PSF_ON_00003_1_MONSTER.Demon2.Hit(
			rtpId,
			&bsMath.Icons.Demon2,
		)
		if v1 > 0 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Demon3.IconID:
		v1 := PSF_ON_00003_1_MONSTER.Demon3.Hit(
			rtpId,
			&bsMath.Icons.Demon3,
		)
		if v1 > 0 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00004_1_RngRtp(state int, drbMath interface{}) (rtpState int, rtpId string, bullet uint64, netWinGroup int) {
	return PSF_ON_00003_1_RTP.Service.RngRtp(state, drbMath)
}

func psfm_00004_1_useRtp(bulletType, mercenaryType int32, rtpId string, bsMathI interface{}) string {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00003_1.BsMath)

	switch {
	case bulletType == 200:
		return PSF_ON_00003_1_MONSTER.Zombie1Drill.UseRTP(&bsMath.Icons.Zombie1Drill)
	case mercenaryType != -1:
		return PSF_ON_00003_1_MONSTER.Bullet.UseRTP(&bsMath.Icons.Bullets)
	default:
		return rtpId
	}
}
