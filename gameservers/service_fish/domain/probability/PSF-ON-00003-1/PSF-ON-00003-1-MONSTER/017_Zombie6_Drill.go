package PSF_ON_00003_1_MONSTER

var Zombie6Drill = &zombie6Drill{}

type zombie6Drill struct {
}

func (z *zombie6Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	return Zombie1Drill.Hit(rtpId, math, bulletMath)
}
