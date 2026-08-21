package PSF_ON_00003_1_MONSTER

var Zombie5Drill = &zombie5Drill{}

type zombie5Drill struct {
}

func (z *zombie5Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	return Zombie1Drill.Hit(rtpId, math, bulletMath)
}
