package PSF_ON_00003_1_MONSTER

import "serve/fish_comm/rng"

var Zombie1Drill = &zombie1Drill{}

type zombie1Drill struct {
}

type BsZombieDrillMath struct {
	IconID           int                            `json:"IconID"`
	BonusID          string                         `json:"BonusID"`
	UseRtp           string                         `json:"UseRTP"`
	BonusTimes       []int                          `json:"BonusTimes"`
	BonusTimesWeight []int                          `json:"BonusTimesWeight"`
	BonusIconID      []int                          `json:"BonusIconID"`
	IconPays         []int                          `json:"IconPays"`
	RTP1             struct{ HitBsZombieDrillMath } `json:"RTP1"`
	RTP2             struct{ HitBsZombieDrillMath } `json:"RTP2"`
	RTP3             struct{ HitBsZombieDrillMath } `json:"RTP3"`
	RTP4             struct{ HitBsZombieDrillMath } `json:"RTP4"`
	RTP5             struct{ HitBsZombieDrillMath } `json:"RTP5"`
	RTP6             struct{ HitBsZombieDrillMath } `json:"RTP6"`
	RTP7             struct{ HitBsZombieDrillMath } `json:"RTP7"`
	RTP8             struct{ HitBsZombieDrillMath } `json:"RTP8"`
	RTP9             struct{ HitBsZombieDrillMath } `json:"RTP9"`
	RTP10            struct{ HitBsZombieDrillMath } `json:"RTP10"`
	RTP11            struct{ HitBsZombieDrillMath } `json:"RTP11"`
}

type HitBsZombieDrillMath struct {
	ID            string `json:"ID"`
	HitWeight     []int  `json:"HitWeight"`
	TriggerIconID int    `json:"TriggerIconID"`
	TriggerWeight []int  `json:"TriggerWeight"`
}

func (z *zombie1Drill) Hit(rtpId string, math *BsZombieDrillMath, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	iconPay = 0
	bonusTimes = 0

	bullets = -1

	switch rtpId {
	case math.RTP1.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP1, bulletMath)
	case math.RTP2.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP2, bulletMath)
	case math.RTP3.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP3, bulletMath)
	case math.RTP4.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP4, bulletMath)
	case math.RTP5.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP5, bulletMath)
	case math.RTP6.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP6, bulletMath)
	case math.RTP7.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP7, bulletMath)
	case math.RTP8.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP8, bulletMath)
	case math.RTP9.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP9, bulletMath)
	case math.RTP10.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP10, bulletMath)
	case math.RTP11.ID:
		iconPay, bonusTimes, bullets = z.rtpHit(math, &math.RTP11, bulletMath)
	default:
		return 0, 0, -1
	}

	return iconPay, bonusTimes, bullets
}

func (z *zombie1Drill) pick(math *BsZombieDrillMath) int {
	bonusTimes := make([]rng.Option, 0, len(math.BonusTimes))

	for i := 0; i < len(math.BonusTimesWeight); i++ {
		bonusTimes = append(bonusTimes, rng.Option{
			math.BonusTimesWeight[i],
			math.BonusTimes[i],
		})
	}

	return MONSTER.rng(bonusTimes).(int)
}

func (z *zombie1Drill) UseRTP(math *BsZombieDrillMath) (rtpId string) {
	return math.UseRtp
}

func (z *zombie1Drill) rtpHit(math *BsZombieDrillMath, rtp *struct{ HitBsZombieDrillMath }, bulletMath *BsBullet) (iconPay, bonusTimes, bullets int) {
	if MONSTER.isHit(rtp.HitWeight) {
		iconPay = math.IconPays[0]
		bonusTimes = z.pick(math)
	}

	if MONSTER.isHit(rtp.TriggerWeight) {
		bullets = Bullet.RngBonusBullets(bulletMath)
	}

	return iconPay, bonusTimes, bullets
}
