package PSF_ON_00003_1_MONSTER

import (
	"os"
	"serve/fish_comm/rng"
	"strings"
)

var LittleZombie = &littleZombie{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "303"),
}

type littleZombie struct {
	bonusMode bool
}

type BsLittleZombieMath struct {
	IconID     int                             `json:"IconID"`
	BonusID    string                          `json:"BonusID"`
	UseRtp     string                          `json:"UseRTP"`
	BonusTimes []int                           `json:"BonusTimes"`
	RTP1       struct{ HitBsLittleZombieMath } `json:"RTP1"`
	RTP2       struct{ HitBsLittleZombieMath } `json:"RTP2"`
	RTP3       struct{ HitBsLittleZombieMath } `json:"RTP3"`
	RTP4       struct{ HitBsLittleZombieMath } `json:"RTP4"`
	RTP5       struct{ HitBsLittleZombieMath } `json:"RTP5"`
	RTP6       struct{ HitBsLittleZombieMath } `json:"RTP6"`
	RTP7       struct{ HitBsLittleZombieMath } `json:"RTP7"`
	RTP8       struct{ HitBsLittleZombieMath } `json:"RTP8"`
	RTP9       struct{ HitBsLittleZombieMath } `json:"RTP9"`
	RTP10      struct{ HitBsLittleZombieMath } `json:"RTP10"`
	RTP11      struct{ HitBsLittleZombieMath } `json:"RTP11"`
}

type HitBsLittleZombieMath struct {
	ID                string `json:"ID"`
	HitWeight         []int  `json:"HitWeight"`
	PayOrBullet       []int  `json:"PayOrBullet"`
	IconPays          []int  `json:"IconPays"`
	IconPaysWeight    []int  `json:"IconPaysWeight"`
	IconBullets       []int  `json:"IconBullets"`
	IconBulletsWeight []int  `json:"IconBulletsWeight"`
}

func (l *littleZombie) HitFs(math *FsZombieMath) (iconPays int) {
	return math.IconPays[0]
}

func (l *littleZombie) Hit(rtpId string, math *BsLittleZombieMath) (iconPay, iconBullet int) {
	iconPay = 0
	iconBullet = 0

	switch rtpId {
	case math.RTP1.ID:
		return l.rtpHit(&math.RTP1)
	case math.RTP2.ID:
		return l.rtpHit(&math.RTP2)
	case math.RTP3.ID:
		return l.rtpHit(&math.RTP3)
	case math.RTP4.ID:
		return l.rtpHit(&math.RTP4)
	case math.RTP5.ID:
		return l.rtpHit(&math.RTP5)
	case math.RTP6.ID:
		return l.rtpHit(&math.RTP6)
	case math.RTP7.ID:
		return l.rtpHit(&math.RTP7)
	case math.RTP8.ID:
		return l.rtpHit(&math.RTP8)
	case math.RTP9.ID:
		return l.rtpHit(&math.RTP9)
	case math.RTP10.ID:
		return l.rtpHit(&math.RTP10)
	case math.RTP11.ID:
		return l.rtpHit(&math.RTP11)
	default:
		return 0, 0
	}
}

func (l *littleZombie) pick(rtp *struct{ HitBsLittleZombieMath }) int {
	options := make([]rng.Option, 0, len(rtp.PayOrBullet))

	for i := 0; i < len(rtp.PayOrBullet); i++ {
		options = append(options, rng.Option{rtp.PayOrBullet[i], i})
	}

	return MONSTER.rng(options).(int)
}

func (l *littleZombie) pickPay(rtp *struct{ HitBsLittleZombieMath }) int {
	options := make([]rng.Option, 0, len(rtp.IconPaysWeight))

	for i := 0; i < len(rtp.IconPaysWeight); i++ {
		options = append(options, rng.Option{
			rtp.IconPaysWeight[i],
			rtp.IconPays[i],
		})
	}

	return MONSTER.rng(options).(int)
}

func (l *littleZombie) pickBullet(rtp *struct{ HitBsLittleZombieMath }) int {
	options := make([]rng.Option, 0, len(rtp.IconBulletsWeight))

	for i := 0; i < len(rtp.IconBulletsWeight); i++ {
		options = append(options, rng.Option{
			rtp.IconBulletsWeight[i],
			rtp.IconBullets[i],
		})
	}

	return MONSTER.rng(options).(int)
}
func (l *littleZombie) rtpHit(rtp *struct{ HitBsLittleZombieMath }) (iconPay, iconBullet int) {
	if MONSTER.isHit(rtp.HitWeight) {
		payOrBullet := l.pick(rtp)

		switch payOrBullet {
		case 0: // IconPay
			iconPay = l.pickPay(rtp)
		case 1: // Bullets
			iconBullet = l.pickBullet(rtp)
		}
	}

	return iconPay, iconBullet
}
