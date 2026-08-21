package probability

import (
	"reflect"
	PSF_ON_00004_1 "serve/service_fish/domain/probability/PSF-ON-00004-1"
	PSF_ON_00004_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00004-1/PSF-ON-00004-1-MONSTER"
	PSF_ON_00004_1_RTP "serve/service_fish/domain/probability/PSF-ON-00004-1/PSF-ON-00004-1-RTP"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00005_96_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-96-1"
	PSFM_00005_97_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-97-1"
	PSFM_00005_98_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-98-1"
	"serve/service_fish/models"
	"strconv"
)

func chooseMath_psfm00005(mathModuleId string) (bsMath *PSF_ON_00004_1.BsMath, fsMath *PSF_ON_00004_1.FsMath, drbMath *PSF_ON_00004_1.DrbMath) {
	switch mathModuleId {
	case models.PSFM_00005_95_1:
		return PSFM_00005_95_1.RTP95BS.PSF_ON_00004_1_BsMath, PSFM_00005_95_1.RTP95FS.PSF_ON_00004_1_FsMath, PSFM_00005_95_1.RTP95DRB.PSF_ON_00004_1_DrbMath
	case models.PSFM_00005_96_1:
		return PSFM_00005_96_1.RTP96BS.PSF_ON_00004_1_BsMath, PSFM_00005_96_1.RTP96FS.PSF_ON_00004_1_FsMath, PSFM_00005_96_1.RTP96DRB.PSF_ON_00004_1_DrbMath
	case models.PSFM_00005_97_1:
		return PSFM_00005_97_1.RTP97BS.PSF_ON_00004_1_BsMath, PSFM_00005_97_1.RTP97FS.PSF_ON_00004_1_FsMath, PSFM_00005_97_1.RTP97DRB.PSF_ON_00004_1_DrbMath
	case models.PSFM_00005_98_1:
		return PSFM_00005_98_1.RTP98BS.PSF_ON_00004_1_BsMath, PSFM_00005_98_1.RTP98FS.PSF_ON_00004_1_FsMath, PSFM_00005_98_1.RTP98DRB.PSF_ON_00004_1_DrbMath
	}

	return nil, nil, nil
}

func psfm_00005_1(fishId, bulletId int32, rtpId string, budget, cent uint64, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00004_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00004_1.FsMath)

	switch int(fishId) {
	case bsMath.Icons.Clam.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Clam.Hit(
			rtpId,
			&bsMath.Icons.Clam,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Shrimp.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Shrimp.Hit(
			rtpId,
			&bsMath.Icons.Shrimp,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Surmullet.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Surmullet.Hit(
			rtpId,
			&bsMath.Icons.Surmullet,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Squid.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Squid.Hit(
			rtpId,
			&bsMath.Icons.Squid,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.FlyingFish.IconID:
		v1 := PSF_ON_00004_1_MONSTER.FlyingFish.Hit(
			rtpId,
			&bsMath.Icons.FlyingFish,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Halibut.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Halibut.Hit(
			rtpId,
			&bsMath.Icons.Halibut,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.ButterflyFish.IconID:
		v1 := PSF_ON_00004_1_MONSTER.ButterflyFish.Hit(
			rtpId,
			&bsMath.Icons.ButterflyFish,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Oplegnathus.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Oplegnathus.Hit(
			rtpId,
			&bsMath.Icons.Oplegnathus,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Snapper.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Snapper.Hit(
			rtpId,
			&bsMath.Icons.Snapper,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.MahiMahi.IconID:
		v1 := PSF_ON_00004_1_MONSTER.MahiMahi.Hit(
			rtpId,
			&bsMath.Icons.MahiMahi,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Stingray.IconID:
		v1 := PSF_ON_00004_1_MONSTER.Stingray.Hit(
			rtpId,
			&bsMath.Icons.Stingray,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.LobsterWaiter.IconID:
		v1 := PSF_ON_00004_1_MONSTER.LobsterWaiter.Hit(
			rtpId,
			&bsMath.Icons.LobsterWaiter,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.PenguinWaiter.IconID:
		v1 := PSF_ON_00004_1_MONSTER.PenguinWaiter.Hit(
			rtpId,
			&bsMath.Icons.PenguinWaiter,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.PlatypusSeniorChef.IconID:
		v1 := PSF_ON_00004_1_MONSTER.PlatypusSeniorChef.Hit(
			rtpId,
			&bsMath.Icons.PlatypusSeniorChef,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.SeaLionChef.IconID:
		v1 := PSF_ON_00004_1_MONSTER.SeaLionChef.Hit(
			rtpId,
			&bsMath.Icons.SeaLionChef,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.HairtCrab.IconID:
		v1 := PSF_ON_00004_1_MONSTER.HairtCrab.Hit(
			rtpId,
			&bsMath.Icons.HairtCrab,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.SnaggletoothShark.IconID:
		v1 := PSF_ON_00004_1_MONSTER.SnaggletoothShark.Hit(
			rtpId,
			&bsMath.Icons.SnaggletoothShark,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.SwordFish.IconID:
		v1 := PSF_ON_00004_1_MONSTER.SwordFish.Hit(
			rtpId,
			&bsMath.Icons.SwordFish,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.WhaleShark.IconID:
		v1 := PSF_ON_00004_1_MONSTER.WhaleShark.Hit(
			rtpId,
			&bsMath.Icons.WhaleShark,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.GiantOarfish.IconID:
		v1 := PSF_ON_00004_1_MONSTER.GiantOarfish.Hit(
			rtpId,
			&bsMath.Icons.GiantOarfish,
		)
		if checkBudget(budget, v1*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.RedEnvelope.IconID:
		v1, v2 := PSF_ON_00004_1_MONSTER.RedEnvelope.Hit(
			rtpId,
			&bsMath.Icons.RedEnvelope,
		)

		if len(v1) > 0 && checkBudget(budget, v2*int(cent), bulletId) {
			return newOptionBonusPay(int(fishId), 0, 100, -1, 1, v1)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.Slot.IconID:
		v1, v2, v3, v4 := PSF_ON_00004_1_MONSTER.Slot.Hit(
			rtpId,
			&bsMath.Icons.Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsStripTableMath,
			&fsMath.ChangeLink1,
			&fsMath.ChangeLink2,
			&fsMath.Icons.Clam,
			&fsMath.Icons.FlyingFish,
			&fsMath.Icons.MahiMahi,
			&fsMath.Icons.LobsterWaiter,
			&fsMath.Icons.PlatypusSeniorChef,
			&fsMath.Icons.SeaLionChef,
			&fsMath.Icons.HairtCrab,
			&fsMath.Icons.SnaggletoothShark,
			&fsMath.Icons.Swordfish,
			&fsMath.Icons.LobsterDish,
			&fsMath.Icons.FruitDish,
			&fsMath.Icons.WhiteTigerChef,
		)

		if v1 > 0 && len(v2) == 9 && len(v3) == 3 && checkBudget(budget, v4*int(cent), bulletId) {
			// client don't want to get pay when hit slot
			slot := newExtraOptionBonusPay(int(fishId), 0, 101, -1, 1, make([]interface{}, 2), v3)
			slot.ExtraData[0] = v1
			slot.ExtraData[1] = v2
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.MachineGun.IconID:
		v1, v2, v3 := PSF_ON_00004_1_MONSTER.MachineGun.Hit(
			rtpId,
			&bsMath.Icons.MachineGun,
		)

		if v1 > 0 && checkBudget(budget, v3*int(cent), bulletId) {
			return newOptionBonusPay(int(fishId), v1, 201, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.SuperMachineGun.IconID:
		v1, v2, v3 := PSF_ON_00004_1_MONSTER.SuperMachineGun.Hit(
			rtpId,
			&bsMath.Icons.SuperMachineGun,
		)

		if v1 > 0 && checkBudget(budget, v3*int(cent), bulletId) {
			return newOptionBonusPay(int(fishId), v1, 202, -1, 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.LobsterDash.IconID:
		v1, v2 := PSF_ON_00004_1_MONSTER.LobsterDash.Hit(
			rtpId,
			&bsMath.Icons.LobsterDash,
		)
		if checkBudget(budget, v1*int(cent), bulletId) && v2 != 0 {
			return newSimplePay(int(fishId), 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.FruitDash.IconID:
		v1, v2 := PSF_ON_00004_1_MONSTER.FruitDash.Hit(
			rtpId,
			&bsMath.Icons.FruitDash,
		)
		if checkBudget(budget, v1*int(cent), bulletId) && v2 != 0 {
			return newSimplePay(int(fishId), 1, v2)
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.XiaoLongBao.IconID:
		v1, v2, v3, v4 := PSF_ON_00004_1_MONSTER.XiaoLongBao.Hit(
			rtpId,
			&bsMath.Icons.XiaoLongBao,
		)

		if v1 > 0 && checkBudget(budget, v4*int(cent), bulletId) {
			xiaoLongBao := newExtraOptionBonusPay(int(fishId), v1, 302, -1, 1, make([]interface{}, 2), v2)
			xiaoLongBao.ExtraData[0] = v2
			xiaoLongBao.ExtraData[1] = v3
			return xiaoLongBao
		}
		return newSimplePay(int(fishId), 0, 1)

	case bsMath.Icons.WhiteTigerChef.IconID:
		v1, v2 := PSF_ON_00004_1_MONSTER.WhiteTigerChef.Hit(
			rtpId,
			&bsMath.Icons.WhiteTigerChef,
		)
		if checkBudget(budget, v2*int(cent), bulletId) {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newSimplePay(int(fishId), 0, 1)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}

func psfm_00005_1_RngFraction(drbMath interface{}) (denominator, molecular, multiplier int) {
	return PSF_ON_00004_1_RTP.Service.RngFraction(drbMath)
}

func psfm_00005_1_RngRtp(fishId int32, rtpOrder string, bsMathI interface{}) (rtpId string) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00004_1.BsMath)

	var order = 1

	if rtpOrder == "" {
		order = PSF_ON_00004_1_RTP.Service.RngRtp()
	} else {
		order, _ = strconv.Atoi(rtpOrder)
	}

	switch fishId {
	case 0:
		return PSF_ON_00004_1_MONSTER.Clam.Rtp(order, &bsMath.Icons.Clam)

	case 1:
		return PSF_ON_00004_1_MONSTER.Shrimp.Rtp(order, &bsMath.Icons.Shrimp)

	case 2:
		return PSF_ON_00004_1_MONSTER.Surmullet.Rtp(order, &bsMath.Icons.Surmullet)

	case 3:
		return PSF_ON_00004_1_MONSTER.Squid.Rtp(order, &bsMath.Icons.Squid)

	case 4:
		return PSF_ON_00004_1_MONSTER.FlyingFish.Rtp(order, &bsMath.Icons.FlyingFish)

	case 5:
		return PSF_ON_00004_1_MONSTER.Halibut.Rtp(order, &bsMath.Icons.Halibut)

	case 6:
		return PSF_ON_00004_1_MONSTER.ButterflyFish.Rtp(order, &bsMath.Icons.ButterflyFish)

	case 7:
		return PSF_ON_00004_1_MONSTER.Oplegnathus.Rtp(order, &bsMath.Icons.Oplegnathus)

	case 8:
		return PSF_ON_00004_1_MONSTER.Snapper.Rtp(order, &bsMath.Icons.Snapper)

	case 9:
		return PSF_ON_00004_1_MONSTER.MahiMahi.Rtp(order, &bsMath.Icons.MahiMahi)

	case 10:
		return PSF_ON_00004_1_MONSTER.Stingray.Rtp(order, &bsMath.Icons.Stingray)

	case 11:
		return PSF_ON_00004_1_MONSTER.LobsterWaiter.Rtp(order, &bsMath.Icons.LobsterWaiter)

	case 12:
		return PSF_ON_00004_1_MONSTER.PenguinWaiter.Rtp(order, &bsMath.Icons.PenguinWaiter)

	case 13:
		return PSF_ON_00004_1_MONSTER.PlatypusSeniorChef.Rtp(order, &bsMath.Icons.PlatypusSeniorChef)

	case 14:
		return PSF_ON_00004_1_MONSTER.SeaLionChef.Rtp(order, &bsMath.Icons.SeaLionChef)

	case 15:
		return PSF_ON_00004_1_MONSTER.HairtCrab.Rtp(order, &bsMath.Icons.HairtCrab)

	case 16:
		return PSF_ON_00004_1_MONSTER.SnaggletoothShark.Rtp(order, &bsMath.Icons.SnaggletoothShark)

	case 17:
		return PSF_ON_00004_1_MONSTER.SwordFish.Rtp(order, &bsMath.Icons.SwordFish)

	case 18:
		return PSF_ON_00004_1_MONSTER.WhaleShark.Rtp(order, &bsMath.Icons.WhaleShark)

	case 19:
		return PSF_ON_00004_1_MONSTER.GiantOarfish.Rtp(order, &bsMath.Icons.GiantOarfish)

	case 100:
		return PSF_ON_00004_1_MONSTER.RedEnvelope.Rtp(order, &bsMath.Icons.RedEnvelope)

	case 101:
		return PSF_ON_00004_1_MONSTER.Slot.Rtp(order, &bsMath.Icons.Slot)

	case 201:
		return PSF_ON_00004_1_MONSTER.MachineGun.Rtp(order, &bsMath.Icons.MachineGun)

	case 202:
		return PSF_ON_00004_1_MONSTER.SuperMachineGun.Rtp(order, &bsMath.Icons.SuperMachineGun)

	case 300:
		return PSF_ON_00004_1_MONSTER.LobsterDash.Rtp(order, &bsMath.Icons.LobsterDash)

	case 301:
		return PSF_ON_00004_1_MONSTER.FruitDash.Rtp(order, &bsMath.Icons.FruitDash)

	case 302:
		return PSF_ON_00004_1_MONSTER.XiaoLongBao.Rtp(order, &bsMath.Icons.XiaoLongBao)

	case 400:
		return PSF_ON_00004_1_MONSTER.WhiteTigerChef.Rtp(order, &bsMath.Icons.WhiteTigerChef)

	default:
		return ""
	}
}

func psfm_00005_1_useRtp(bulletType int32, rtpId string, bsMathI interface{}) string {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00004_1.BsMath)

	switch bulletType {
	case 201:
		return PSF_ON_00004_1_MONSTER.MachineGun.UseRtp(&bsMath.Icons.MachineGun)

	case 202:
		return PSF_ON_00004_1_MONSTER.SuperMachineGun.UseRtp(&bsMath.Icons.SuperMachineGun)

	default:
		return rtpId
	}
}

func psfm_00005_1_avgPay(fishId int, bsMathI interface{}) (haveAvgPay bool, avgPay int) {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00004_1.BsMath)

	haveAvgPay = false
	avgPay = 0

	switch fishId {
	case 100:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.RedEnvelope.AvgPay(&bsMath.Icons.RedEnvelope)
	case 101:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.Slot.AvgPay(&bsMath.Icons.Slot)
	case 201:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.MachineGun.AvgPay(&bsMath.Icons.MachineGun)
	case 202:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.SuperMachineGun.AvgPay(&bsMath.Icons.SuperMachineGun)
	case 300:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.LobsterDash.AvgPay(&bsMath.Icons.LobsterDash)
	case 301:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.FruitDash.AvgPay(&bsMath.Icons.FruitDash)
	case 302:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.XiaoLongBao.AvgPay(&bsMath.Icons.XiaoLongBao)
	case 400:
		haveAvgPay = true
		avgPay = PSF_ON_00004_1_MONSTER.WhiteTigerChef.AvgPay(&bsMath.Icons.WhiteTigerChef)
	default:
		return false, 0
	}

	return haveAvgPay, avgPay
}

func psfm_00005_1_firstMultiplier(subgameId int, drbMath interface{}) int {
	return PSF_ON_00004_1_RTP.Service.GetFirstMultiplier(subgameId, drbMath)
}
func checkBudget(budget uint64, iconPay int, bulletId int32) bool {
	switch bulletId {
	case 201:
		return true
	case 202:
		return true
	default:
		if iconPay > int(budget) {
			return false
		} else {
			return true
		}
	}
}
