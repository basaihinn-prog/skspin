package PSF_ON_00003_1_MONSTER

var Zombie3Drill = &zombie3Drill{}

type zombie3Drill struct {
}

func (z *zombie3Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	return Zombie1Drill.Hit(rtpId, math, bulletMath)
}
