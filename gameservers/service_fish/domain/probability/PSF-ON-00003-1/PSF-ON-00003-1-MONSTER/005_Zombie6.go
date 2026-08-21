package PSF_ON_00003_1_MONSTER

var Zombie6 = &zombie6{}

type zombie6 struct {
}

func (z *zombie6) Hit(rtpId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	return Zombie1.Hit(rtpId, math, bulletMath)
}
