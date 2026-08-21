package bullet

import "serve/fish_comm/flux"

func psf_on_00007_Shoot(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00004_Shoot(action, bullets, bonusBullets)
}

func psf_on_00007_Bonus(action *flux.Action, bonusBullets map[string]*Bullet) {
	psf_on_00004_Bonus(action, bonusBullets)
}

func psf_on_00007_Hit(action *flux.Action, bullets map[string]*Bullet, bonusBullets map[string]*Bullet) {
	psf_on_00004_Hit(action, bullets, bonusBullets)
}
