package PSF_ON_00003_1_MONSTER

var Zombie3 = &zombie3{}

type zombie3 struct {
}

func (z *zombie3) Hit(rtpId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	return Zombie1.Hit(rtpId, math, bulletMath)
}

func (z *zombie3) HitFs(math *FsZombieMath) (iconPays int) {
	return Zombie1.HitFs(math)
}
