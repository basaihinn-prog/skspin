package PSF_ON_00003_1_MONSTER

var Zombie1 = &zombie1{}

type zombie1 struct {
}

type BsZombieMath struct {
	IconID     int                       `json:"IconID"`
	BonusID    string                    `json:"BonusID"`
	BonusTimes []int                     `json:"BonusTimes"`
	IconPays   []int                     `json:"IconPays"`
	RTP1       struct{ HitBsZombieMath } `json:"RTP1"`
	RTP2       struct{ HitBsZombieMath } `json:"RTP2"`
	RTP3       struct{ HitBsZombieMath } `json:"RTP3"`
	RTP4       struct{ HitBsZombieMath } `json:"RTP4"`
	RTP5       struct{ HitBsZombieMath } `json:"RTP5"`
	RTP6       struct{ HitBsZombieMath } `json:"RTP6"`
	RTP7       struct{ HitBsZombieMath } `json:"RTP7"`
	RTP8       struct{ HitBsZombieMath } `json:"RTP8"`
	RTP9       struct{ HitBsZombieMath } `json:"RTP9"`
	RTP10      struct{ HitBsZombieMath } `json:"RTP10"`
	RTP11      struct{ HitBsZombieMath } `json:"RTP11"`
}

type HitBsZombieMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
}

type FsZombieMath struct {
	IconID   int
	IconPays []int
}

func (z *zombie1) Hit(rtpId string, math *BsZombieMath, bulletMath *BsBullet) (iconPay, bullets int) {
	iconPay = 0
	bullets = -1

	switch rtpId {
	case math.RTP1.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP1, bulletMath)
	case math.RTP2.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP2, bulletMath)
	case math.RTP3.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP3, bulletMath)
	case math.RTP4.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP4, bulletMath)
	case math.RTP5.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP5, bulletMath)
	case math.RTP6.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP6, bulletMath)
	case math.RTP7.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP7, bulletMath)
	case math.RTP8.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP8, bulletMath)
	case math.RTP9.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP9, bulletMath)
	case math.RTP10.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP10, bulletMath)
	case math.RTP11.ID:
		iconPay, bullets = z.rtpHit(math.IconPays[0], &math.RTP11, bulletMath)
	default:
		return 0, -1
	}

	return iconPay, bullets
}

func (z *zombie1) HitFs(math *FsZombieMath) int {
	return math.IconPays[0]
}

func (z *zombie1) rtpHit(iconPays int, bsZombie *struct{ HitBsZombieMath }, bulletMath *BsBullet) (iconPay, bullets int) {
	if MONSTER.isHit(bsZombie.HitWeight) {
		iconPay = iconPays
	}

	if MONSTER.isHit(bsZombie.TriggerWeight) {
		bullets = Bullet.RngBonusBullets(bulletMath)
	}

	return iconPay, bullets
}
