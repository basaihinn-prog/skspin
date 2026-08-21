package PSF_ON_00003_1_MONSTER

var Zombie5 = &zombie5{}

type zombie5 struct {
}

func (z *zombie5) Hit(rtpId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	return Zombie1.Hit(rtpId, math, bulletMath)
}
