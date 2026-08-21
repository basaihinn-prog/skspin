//go:build dev || staging
// +build dev staging

package PSF_ON_00002_1_MONSTER

func (s *slot) Hit(
	rtpId string,
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	reelMath *FsReelMath,
	faceOffChangerMath31 *FsFaceOffChangerMath_31,
	faceOffChangerMath *FsFaceOffChangerMath,
	starFish, clownfish, lionFish, giantJellyFish, lobster, hermitCrab, goldGemTurtle, goldCoinToad, goldCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (iconPay int, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	switch rtpId {
	case "0":
		if MONSTER.isHit(bsMath.RTP1.TriggerWeight) {
			triggerIconId = bsMath.RTP1.TriggerIconID
			bonusTypeId = bsMath.RTP1.Type
		}

		if MONSTER.isHit(bsMath.RTP1.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "20":
		if MONSTER.isHit(bsMath.RTP2.TriggerWeight) {
			triggerIconId = bsMath.RTP2.TriggerIconID
			bonusTypeId = bsMath.RTP2.Type
		}

		if MONSTER.isHit(bsMath.RTP2.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "40":
		if MONSTER.isHit(bsMath.RTP3.TriggerWeight) {
			triggerIconId = bsMath.RTP3.TriggerIconID
			bonusTypeId = bsMath.RTP3.Type
		}

		if MONSTER.isHit(bsMath.RTP3.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "60":
		if MONSTER.isHit(bsMath.RTP4.TriggerWeight) {
			triggerIconId = bsMath.RTP4.TriggerIconID
			bonusTypeId = bsMath.RTP4.Type
		}

		if MONSTER.isHit(bsMath.RTP4.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "80":
		if MONSTER.isHit(bsMath.RTP5.TriggerWeight) {
			triggerIconId = bsMath.RTP5.TriggerIconID
			bonusTypeId = bsMath.RTP5.Type
		}

		if MONSTER.isHit(bsMath.RTP5.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "901":
		if MONSTER.isHit(bsMath.RTP6.TriggerWeight) {
			triggerIconId = bsMath.RTP6.TriggerIconID
			bonusTypeId = bsMath.RTP6.Type
		}

		if MONSTER.isHit(bsMath.RTP6.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "902":
		if MONSTER.isHit(bsMath.RTP7.TriggerWeight) {
			triggerIconId = bsMath.RTP7.TriggerIconID
			bonusTypeId = bsMath.RTP7.Type
		}

		if MONSTER.isHit(bsMath.RTP7.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "100":
		if MONSTER.isHit(bsMath.RTP8.TriggerWeight) {
			triggerIconId = bsMath.RTP8.TriggerIconID
			bonusTypeId = bsMath.RTP8.Type
		}

		if MONSTER.isHit(bsMath.RTP8.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "150":
		if MONSTER.isHit(bsMath.RTP9.TriggerWeight) {
			triggerIconId = bsMath.RTP9.TriggerIconID
			bonusTypeId = bsMath.RTP9.Type
		}

		if MONSTER.isHit(bsMath.RTP9.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "200":
		if MONSTER.isHit(bsMath.RTP10.TriggerWeight) {
			triggerIconId = bsMath.RTP10.TriggerIconID
			bonusTypeId = bsMath.RTP10.Type
		}

		if MONSTER.isHit(bsMath.RTP10.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	case "300":
		if MONSTER.isHit(bsMath.RTP11.TriggerWeight) {
			triggerIconId = bsMath.RTP11.TriggerIconID
			bonusTypeId = bsMath.RTP11.Type
		}

		if MONSTER.isHit(bsMath.RTP11.HitWeight) {
			v1, v2, v3 := s.hit(
				fsMath,
				reelMath,
				faceOffChangerMath31,
				faceOffChangerMath,
				starFish,
				clownfish,
				lionFish,
				giantJellyFish,
				lobster,
				hermitCrab,
				goldGemTurtle,
				goldCoinToad,
				goldCrab,
				faceOff1,
				faceOff3,
				faceOff5,
			)
			return v1, triggerIconId, bonusTypeId, v2, v3
		}
		return 0, triggerIconId, bonusTypeId, nil, nil

	default:
		return 0, -1, -1, nil, nil
	}
}
