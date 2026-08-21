package PSF_ON_00003_1_MONSTER

import "serve/fish_comm/rng"

var Bullet = &bullet{}

type bullet struct {
}

type BsBullet struct {
	IconID            int
	UseRtp            string `json:"UseRTP"`
	IconBullets       []int
	IconBulletsWeight []int
}

func (b *bullet) RngBonusBullets(bulletMath *BsBullet) (bullet int) {
	bulletOption := make([]rng.Option, 0, len(bulletMath.IconBulletsWeight))

	for i := 0; i < len(bulletMath.IconBulletsWeight); i++ {
		bulletOption = append(bulletOption, rng.Option{
			Weight: bulletMath.IconBulletsWeight[i],
			Item:   bulletMath.IconBullets[i],
		})
	}

	return MONSTER.rng(bulletOption).(int)
}

func (b *bullet) UseRTP(math *BsBullet) (rtpId string) {
	return math.UseRtp
}
