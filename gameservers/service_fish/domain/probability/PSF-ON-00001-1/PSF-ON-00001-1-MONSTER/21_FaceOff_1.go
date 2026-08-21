package PSF_ON_00001_1_MONSTER

import (
	"os"
	"strings"
)

var FaceOff1 = &faceOff1{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "21"),
}

type faceOff1 struct {
	bonusMode bool
}

type BsFaceOffMath struct {
	IconID           int
	IsWildSubstitute bool
	BonusID          string
	BonusTimes       []int
	HitWeight        []int
	IconPays         []int
}

type FsFaceOffMath struct {
	IconID   int
	IconPays []int
}

func (f *faceOff1) Hit(math *BsFaceOffMath) (iconPay int) {
	if MONSTER.isHit(math.HitWeight) {
		return math.IconPays[0]
	}
	return 0
}

func (f *faceOff1) HitFs(math *FsFaceOffMath) (iconPay int) {
	return math.IconPays[0]
}
