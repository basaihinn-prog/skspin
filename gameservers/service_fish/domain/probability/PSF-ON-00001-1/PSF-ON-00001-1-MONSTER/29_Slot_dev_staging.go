//go:build dev || staging
// +build dev staging

package PSF_ON_00001_1_MONSTER

func (s *slot) Hit(
	bsMath *BsSlotMath,
	fsMath *FsSlotMath,
	reelMath *FsReelMath,
	faceOffChangerMath *FsFaceOffChangerMath,
	starFish, dory, lionFish, lobster, platypus, manatee, dolphin, koi, hermitCrab *FsFishMath,
	faceOff1, faceOff3, faceOff5 *FsFaceOffMath,
) (pay int, iconIds []int32, iconPays []int64) {
	if MONSTER.isHit(bsMath.HitWeight) {
		blockReelIndex := s.rngBlockReelIndex()

		iconIds := make([]int32, 0)
		iconPays := make([]int64, 3)

		for i := 0; i < fsMath.MixAwardId; i++ {
			if i == blockReelIndex {
				v1, v2 := s.rngBlockReel(fsMath, reelMath, faceOffChangerMath)

				for _, v := range v1 {
					iconIds = append(iconIds, v)
				}

				iconPays[i] = int64(v2)
			} else {
				v1, v2 := s.rngNormalReel(fsMath, reelMath, faceOffChangerMath)

				for _, v := range v1 {
					iconIds = append(iconIds, v)
				}

				iconPays[i] = int64(v2)
			}
		}

		// big win
		if iconIds[1] == iconIds[4] && iconIds[1] == iconIds[7] {
			bonusPay := s.bonusPay(
				iconIds[1],
				starFish, dory, lionFish, lobster, platypus, manatee, dolphin, koi, hermitCrab,
				faceOff1, faceOff3, faceOff5,
			)

			if s.contain(fsMath, bonusPay) {
				return bonusPay, iconIds, iconPays
			}
			return 0, nil, nil
		}

		if s.contain(fsMath, int(iconPays[0]+iconPays[1]+iconPays[2])) {
			return int(iconPays[0] + iconPays[1] + iconPays[2]), iconIds, iconPays
		}
	}
	return 0, nil, nil
}
