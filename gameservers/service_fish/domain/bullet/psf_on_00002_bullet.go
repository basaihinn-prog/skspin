package bullet

import (
	"serve/fish_comm/flux"
)

func psf_on_00002_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00001_Shoot(action, bullets, bonusBullets)
}

func psf_on_00002_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	psf_on_00001_Bonus(action, bonusBullets)
}

func psf_on_00002_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00001_Hit(action, bullets, bonusBullets)
}

func psf_on_00002_CheckDrill(bulletType int32) bool {
	return psf_on_00001_CheckDrill(bulletType)
}
