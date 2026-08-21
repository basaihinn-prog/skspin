package PSF_ON_00004_1_MONSTER

import (
	"os"
	"strings"
)

var Clam = &clam{
	bonusMode: strings.Contains(os.Getenv("ENABLED_BONUS"), "31"),
}

type clam struct {
	bonusMode bool
}

type BsFishMath struct {
	IconID     int                     `json:"IconID"`
	BonusID    string                  `json:"BonusID"`
	BonusTimes []int                   `json:"BonusTimes"`
	IconPays   []int                   `json:"IconPays"`
	RTP1       struct{ HitBsFishMath } `json:"RTP1"`
	RTP2       struct{ HitBsFishMath } `json:"RTP2"`
	RTP3       struct{ HitBsFishMath } `json:"RTP3"`
	RTP4       struct{ HitBsFishMath } `json:"RTP4"`
	RTP5       struct{ HitBsFishMath } `json:"RTP5"`
}

type HitBsFishMath struct {
	ID        string `json:"ID"`
	HitWeight []int  `json:"HitWeight"`
}

type FsFishMath struct {
	IconID   int
	IconPays []int
}

func (c *clam) Hit(rtpId string, math *BsFishMath) (iconPay int) {
	iconPay = 0

	switch rtpId {
	case math.RTP1.ID:
		iconPay = c.rtpHit(math.IconPays[0], &math.RTP1)

	case math.RTP2.ID:
		iconPay = c.rtpHit(math.IconPays[0], &math.RTP2)

	case math.RTP3.ID:
		iconPay = c.rtpHit(math.IconPays[0], &math.RTP3)

	case math.RTP4.ID:
		iconPay = c.rtpHit(math.IconPays[0], &math.RTP4)

	case math.RTP5.ID:
		iconPay = c.rtpHit(math.IconPays[0], &math.RTP5)

	default:
		return 0
	}

	return iconPay
}

func (c *clam) HitFs(math *FsFishMath) (iconPays int) {
	return math.IconPays[0]
}

func (c *clam) Rtp(order int, math *BsFishMath) (rtpId string) {
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
func (c *clam) rtpHit(iconPays int, bsFish *struct{ HitBsFishMath }) (iconPay int) {
	if MONSTER.isHit(bsFish.HitWeight) {
		iconPay = iconPays
	}

	return iconPay
}
