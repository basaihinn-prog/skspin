//go:build dev || staging
// +build dev staging

package PSF_ON_00002_1_MONSTER

import "serve/fish_comm/rng"

func (f *faceOff1) Hit(rtpId string, math *BsFaceOffMath) (iconPay int, triggerIconId int, bonusTypeId int) {
	iconPay = 0
	triggerIconId = -1
	bonusTypeId = -1
	var FaceOffOptions []rng.Option
	for i := 0; i < len(math.IconPays); i++ {
		FaceOffOptions = append(FaceOffOptions, rng.Option{math.IconPaysWeight[i], math.IconPays[i]})
	}

	switch rtpId {
	case "0":
		if MONSTER.isHit(math.RTP1.TriggerWeight) {
			triggerIconId = math.RTP1.TriggerIconID
			bonusTypeId = math.RTP1.Type
		}

		if MONSTER.isHit(math.RTP1.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "20":
		if MONSTER.isHit(math.RTP2.TriggerWeight) {
			triggerIconId = math.RTP2.TriggerIconID
			bonusTypeId = math.RTP2.Type
		}

		if MONSTER.isHit(math.RTP2.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "40":
		if MONSTER.isHit(math.RTP3.TriggerWeight) {
			triggerIconId = math.RTP3.TriggerIconID
			bonusTypeId = math.RTP3.Type
		}

		if MONSTER.isHit(math.RTP3.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "60":
		if MONSTER.isHit(math.RTP4.TriggerWeight) {
			triggerIconId = math.RTP4.TriggerIconID
			bonusTypeId = math.RTP4.Type
		}

		if MONSTER.isHit(math.RTP4.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)

		}

	case "80":
		if MONSTER.isHit(math.RTP5.TriggerWeight) {
			triggerIconId = math.RTP5.TriggerIconID
			bonusTypeId = math.RTP5.Type
		}

		if MONSTER.isHit(math.RTP5.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)

		}

	case "901":
		if MONSTER.isHit(math.RTP6.TriggerWeight) {
			triggerIconId = math.RTP6.TriggerIconID
			bonusTypeId = math.RTP6.Type
		}

		if MONSTER.isHit(math.RTP6.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "902":
		if MONSTER.isHit(math.RTP7.TriggerWeight) {
			triggerIconId = math.RTP7.TriggerIconID
			bonusTypeId = math.RTP7.Type
		}

		if MONSTER.isHit(math.RTP7.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "100":
		if MONSTER.isHit(math.RTP8.TriggerWeight) {
			triggerIconId = math.RTP8.TriggerIconID
			bonusTypeId = math.RTP8.Type
		}

		if MONSTER.isHit(math.RTP8.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "150":
		if MONSTER.isHit(math.RTP9.TriggerWeight) {
			triggerIconId = math.RTP9.TriggerIconID
			bonusTypeId = math.RTP9.Type
		}

		if MONSTER.isHit(math.RTP9.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "200":
		if MONSTER.isHit(math.RTP10.TriggerWeight) {
			triggerIconId = math.RTP10.TriggerIconID
			bonusTypeId = math.RTP10.Type
		}

		if MONSTER.isHit(math.RTP10.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	case "300":
		if MONSTER.isHit(math.RTP11.TriggerWeight) {
			triggerIconId = math.RTP11.TriggerIconID
			bonusTypeId = math.RTP11.Type
		}

		if MONSTER.isHit(math.RTP11.HitWeight) {
			iconPay = MONSTER.rng(FaceOffOptions).(int)
		}

	default:
		return iconPay, triggerIconId, bonusTypeId
	}

	return iconPay, triggerIconId, bonusTypeId
}

func (f *faceOff1) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
