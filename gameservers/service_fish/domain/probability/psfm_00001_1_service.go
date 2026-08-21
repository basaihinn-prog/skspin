package probability

import (
	"reflect"
	PSF_ON_00001_1 "serve/service_fish/domain/probability/PSF-ON-00001-1"
	PSF_ON_00001_1_MONSTER "serve/service_fish/domain/probability/PSF-ON-00001-1/PSF-ON-00001-1-MONSTER"
	PSFM_00001_98_1 "serve/service_fish/domain/probability/PSFM-00001-1/PSFM-00001-98-1"
	"serve/service_fish/models"
)

func chooseMath_psfm00001(mathModuleId string) (bsMath *PSF_ON_00001_1.BsMath, fsMath *PSF_ON_00001_1.FsMath) {
	switch mathModuleId {
	case models.PSFM_00001_98_1:
		return PSFM_00001_98_1.RTP98BS.PSF_ON_00001_1_BsMath, PSFM_00001_98_1.RTP98FS.PSF_ON_00001_1_FsMath
	}

	return nil, nil
}

func psfm_00001_1(fishId int32, bsMathI, fsMathI interface{}) *Probability {
	bsMath := reflect.ValueOf(bsMathI).Interface().(*PSF_ON_00001_1.BsMath)
	fsMath := reflect.ValueOf(fsMathI).Interface().(*PSF_ON_00001_1.FsMath)

	switch fishId {
	case 0:
		v1, v2 := PSF_ON_00001_1_MONSTER.StarFish.Hit(&bsMath.Icons.StarFish)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 1:
		v1, v2 := PSF_ON_00001_1_MONSTER.SeaHorse.Hit(&bsMath.Icons.SeaHorse)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 2:
		v1, v2 := PSF_ON_00001_1_MONSTER.Guppy.Hit(&bsMath.Icons.Guppy)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 3:
		v1, v2 := PSF_ON_00001_1_MONSTER.ClownFish.Hit(&bsMath.Icons.ClownFish)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 4:
		v1, v2 := PSF_ON_00001_1_MONSTER.Dory.Hit(&bsMath.Icons.Dory)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 5:
		v1, v2 := PSF_ON_00001_1_MONSTER.Grammidae.Hit(&bsMath.Icons.Grammidae)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 6:
		v1, v2 := PSF_ON_00001_1_MONSTER.RajahCichlasoma.Hit(&bsMath.Icons.RajahCichlasoma)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 7:
		v1, v2 := PSF_ON_00001_1_MONSTER.PufferFish.Hit(&bsMath.Icons.PufferFish)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 8:
		v1, v2 := PSF_ON_00001_1_MONSTER.LanternFish.Hit(&bsMath.Icons.LanternFish)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 9:
		v1, v2 := PSF_ON_00001_1_MONSTER.LionFish.Hit(&bsMath.Icons.LionFish)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 10:
		v1, v2 := PSF_ON_00001_1_MONSTER.Turtle.Hit(&bsMath.Icons.Turtle)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 11:
		v1, v2 := PSF_ON_00001_1_MONSTER.Lobster.Hit(&bsMath.Icons.Lobster)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 12:
		v1, v2 := PSF_ON_00001_1_MONSTER.Penguin.Hit(&bsMath.Icons.Penguin)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 13:
		v1, v2 := PSF_ON_00001_1_MONSTER.Platypus.Hit(&bsMath.Icons.Platypus)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 14:
		v1, v2 := PSF_ON_00001_1_MONSTER.Manatee.Hit(&bsMath.Icons.Manatee)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 15:
		v1, v2 := PSF_ON_00001_1_MONSTER.Dolphin.Hit(&bsMath.Icons.Dolphin)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 16:
		v1, v2 := PSF_ON_00001_1_MONSTER.Koi.Hit(&bsMath.Icons.Koi)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 17:
		v1, v2 := PSF_ON_00001_1_MONSTER.HermitCrab.Hit(&bsMath.Icons.HermitCrab)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 18:
		v1, v2 := PSF_ON_00001_1_MONSTER.Shark.Hit(&bsMath.Icons.Shark)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 19:
		v1, v2 := PSF_ON_00001_1_MONSTER.AsianArowana.Hit(&bsMath.Icons.AsianArowana)

		if v2 == -1 {
			return newSimplePay(int(fishId), v1, 1)
		}
		return newTriggerBonusPay(int(fishId), v1, v2, -1, 1)

	case 20:
		v1 := PSF_ON_00001_1_MONSTER.HotPot.Hit(&bsMath.Icons.HotPot)
		return newSimplePay(int(fishId), v1, 1)

	case 21:
		v1 := PSF_ON_00001_1_MONSTER.FaceOff1.Hit(&bsMath.Icons.FaceOff1)
		return newSimplePay(int(fishId), v1, 1)

	case 22:
		v1 := PSF_ON_00001_1_MONSTER.FaceOff2.Hit(&bsMath.Icons.FaceOff2)
		return newSimplePay(int(fishId), v1, 1)

	case 23:
		v1 := PSF_ON_00001_1_MONSTER.FaceOff3.Hit(&bsMath.Icons.FaceOff3)
		return newSimplePay(int(fishId), v1, 1)

	case 24:
		v1 := PSF_ON_00001_1_MONSTER.FaceOff4.Hit(&bsMath.Icons.FaceOff4)
		return newSimplePay(int(fishId), v1, 1)

	case 25:
		v1 := PSF_ON_00001_1_MONSTER.FaceOff5.Hit(&bsMath.Icons.FaceOff5)
		return newSimplePay(int(fishId), v1, 1)

	case 26:
		v1, v2 := PSF_ON_00001_1_MONSTER.Drill.Hit(&bsMath.Icons.Drill)

		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 26, -1, 1, v2)
		}

		return newSimplePay(int(fishId), 0, 1)

	case 27:
		v1, v2 := PSF_ON_00001_1_MONSTER.MachineGun.Hit(&bsMath.Icons.MachineGun)

		if v1 > 0 {
			return newOptionBonusPay(int(fishId), v1, 27, -1, 1, v2)
		}

		return newSimplePay(int(fishId), 0, 1)

	case 28:
		v1 := PSF_ON_00001_1_MONSTER.RedEnvelope.Hit(&bsMath.Icons.RedEnvelope)

		if len(v1) > 0 {
			return newOptionBonusPay(int(fishId), 0, 28, -1, 1, v1)
		}

		return newSimplePay(int(fishId), 0, 1)

	case 29:
		v1, v2, v3 := PSF_ON_00001_1_MONSTER.Slot.Hit(
			&bsMath.Icons.Slot,
			&fsMath.MixAward[0],
			&fsMath.StripTable.FsReelMath,
			&fsMath.ChangeLink,
			&fsMath.Icons.StarFish,
			&fsMath.Icons.Dory,
			&fsMath.Icons.LionFish,
			&fsMath.Icons.Lobster,
			&fsMath.Icons.Platypus,
			&fsMath.Icons.Manatee,
			&fsMath.Icons.Dolphin,
			&fsMath.Icons.Koi,
			&fsMath.Icons.HermitCrab,
			&fsMath.Icons.FaceOff1,
			&fsMath.Icons.FaceOff3,
			&fsMath.Icons.FaceOff5,
		)

		if v1 > 0 && len(v2) == 9 && len(v3) == 3 {
			// client don't want to get pay when hit slot
			slot := newExtraOptionBonusPay(int(fishId), 0, 29, -1, 1, make([]interface{}, 2), v3)
			slot.ExtraData[0] = v1
			slot.ExtraData[1] = v2
			return slot
		}
		return newSimplePay(int(fishId), 0, 1)

	case 30:
		v1 := PSF_ON_00001_1_MONSTER.FaceOffChanger.Hit(&bsMath.ChangeLink)
		return newSimplePay(v1, 0, 1)

	case 31:
		v1 := PSF_ON_00001_1_MONSTER.Envelope.Pick(&bsMath.Icons.RedEnvelope)
		return newOptionBonusPay(int(fishId), 0, 31, -1, 1, v1)

	default:
		return newSimplePay(int(fishId), 0, 1)
	}
}
