//go:build dev || staging
// +build dev staging

package PSF_ON_00004_1_MONSTER

func (s *slot) rtpHit(
	weights []int,
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	faceOffChangerMath1 *FsFaceOffChanger1Math,
	faceOffChangerMath2 *FsFaceOffChanger2Math,
	clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
	snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef *FsFishMath,
) (iconPay int, iconIds []int32, iconPays []int64) {
	if MONSTER.isHit(weights) {
		v1, v2, v3 := s.hit(
			fsMath,
			fsStripMath,
			faceOffChangerMath1,
			faceOffChangerMath2,
			clam, flyingFish, mahiMahi, lobsterWaiter, platypusSeniorChef, seaLionChef, hairtCrab,
			snaggletoothShark, swordfish, lobsterDish, fruitDish, whiteTigerChef,
		)

		return v1, v2, v3
	}
	return 0, nil, nil
}
