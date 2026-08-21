package PSF_ON_00003_1_MONSTER

import "serve/fish_comm/rng"

var BonusZombie = &bonusZombie{}

type bonusZombie struct {
}

type BsBonusZombieMath struct {
	IconID     int                            `json:"IconID"`
	BonusID    string                         `json:"BonusID"`
	BonusTimes []int                          `json:"BonusTimes"`
	RTP1       struct{ HitBsBonusZombieMath } `json:"RTP1"`
	RTP2       struct{ HitBsBonusZombieMath } `json:"RTP2"`
	RTP3       struct{ HitBsBonusZombieMath } `json:"RTP3"`
	RTP4       struct{ HitBsBonusZombieMath } `json:"RTP4"`
	RTP5       struct{ HitBsBonusZombieMath } `json:"RTP5"`
	RTP6       struct{ HitBsBonusZombieMath } `json:"RTP6"`
	RTP7       struct{ HitBsBonusZombieMath } `json:"RTP7"`
	RTP8       struct{ HitBsBonusZombieMath } `json:"RTP8"`
	RTP9       struct{ HitBsBonusZombieMath } `json:"RTP9"`
	RTP10      struct{ HitBsBonusZombieMath } `json:"RTP10"`
}

type HitBsBonusZombieMath struct {
	ID                string `json:"ID"`
	HitWeight         []int  `json:"HitWeight"`
	IconBullets       []int  `json:"IconBullets"`
	IconBulletsWeight []int  `json:"IconBulletsWeight"`
	TriggerIconID     int    `json:"TriggerIconID"`
}

func (b *bonusZombie) Hit(rtpId string, math *BsBonusZombieMath) (iconBullets, triggerIconId int) {
	iconBullets = 0
	triggerIconId = -1

	switch rtpId {
	case math.RTP1.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP1)
	case math.RTP2.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP2)
	case math.RTP3.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP3)
	case math.RTP4.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP4)
	case math.RTP5.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP5)
	case math.RTP6.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP6)
	case math.RTP7.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP7)
	case math.RTP8.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP8)
	case math.RTP9.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP9)
	case math.RTP10.ID:
		iconBullets, triggerIconId = b.rtpHit(&math.RTP10)
	default:
		return 0, -1
	}

	return iconBullets, triggerIconId
}

func (b *bonusZombie) pick(math *struct{ HitBsBonusZombieMath }) int {
	iconBullets := make([]rng.Option, 0, len(math.IconBullets))

	for i := 0; i < len(math.IconBulletsWeight); i++ {
		iconBullets = append(iconBullets, rng.Option{
			math.IconBulletsWeight[i],
			math.IconBullets[i],
		})
	}

	return MONSTER.rng(iconBullets).(int)
}

func (b *bonusZombie) Rtp(order int, math *BsBonusZombieMath) (rtpId string) {
	switch order {
	case 1:
		return math.RTP1.ID
	case 2:
		return math.RTP2.ID
	case 3:
		return math.RTP3.ID
	case 4:
		return math.RTP4.ID
	case 5:
		return math.RTP5.ID
	default:
		return ""
	}
}
func (b *bonusZombie) rtpHit(math *struct{ HitBsBonusZombieMath }) (iconBullets, triggerIconId int) {
	if MONSTER.isHit(math.HitWeight) {
		iconBullets = b.pick(math)
		triggerIconId = math.TriggerIconID
	}

	return iconBullets, triggerIconId
}
