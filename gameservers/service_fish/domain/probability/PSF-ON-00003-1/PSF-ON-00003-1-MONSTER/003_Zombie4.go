package PSF_ON_00003_1_MONSTER

var Zombie4 = &zombie4{}

type zombie4 struct {
}

func (z *zombie4) Hit(rtpId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	return Zombie1.Hit(rtpId, math, bulletMath)
}

func (z *zombie4) HitFs(math *FsZombieMath) (iconPays int) {
	return Zombie1.HitFs(math)
}
