//go:build prod
// +build prod

package PSF_ON_00005_1_MONSTER

func (s *slot) rtpHit(
	rtp *struct{ HitBsSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	giantWhaleChanger1Math *FsGiantWhaleChanger1Math,
	giantWhaleChanger2Math *FsGiantWhaleChanger2Math,
	starFish, clownFish, lionFish, giantJellyFish, lobster, hermitCrab,
	goldGemTurtle, goldCoinToad, goldCrab,
	fruitDish, aPackOfBeer, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := s.hit(fsMath, fsStripMath, giantWhaleChanger1Math, giantWhaleChanger2Math,
			starFish, clownFish, lionFish, giantJellyFish, lobster,
			hermitCrab, goldGemTurtle, goldCoinToad, goldCrab,
			fruitDish, aPackOfBeer, giantWhale,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
