package PSF_ON_00003_1_MONSTER

var Zombie4Drill = &zombie4Drill{}

type zombie4Drill struct {
}

func (z *zombie4Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	return Zombie1Drill.Hit(rtpId, math, bulletMath)
}
