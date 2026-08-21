package PSF_ON_00003_1_MONSTER

var Zombie2Drill = &zombie2Drill{}

type zombie2Drill struct {
}

func (z *zombie2Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	return Zombie1Drill.Hit(rtpId, math, bulletMath)
}
