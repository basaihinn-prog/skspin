package PSF_ON_00003_1_MONSTER

var Zombie2 = &zombie2{}

type zombie2 struct {
}

func (z *zombie2) Hit(rptId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	return Zombie1.Hit(rptId, math, bulletMath)
}

func (z *zombie2) HitFs(math *FsZombieMath) (iconPays int) {
	return Zombie1.HitFs(math)
}
