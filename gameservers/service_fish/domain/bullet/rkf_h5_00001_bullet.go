package bullet

import "serve/fish_comm/flux"

func rkf_h5_00001_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00004_Shoot(action, bullets, bonusBullets)
}

func rkf_h5_00001_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	psf_on_00004_Bonus(action, bonusBullets)
}

func rkf_h5_00001_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00004_Hit(action, bullets, bonusBullets)
}
