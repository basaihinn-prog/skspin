//go:build dev
// +build dev

package activity

import (
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
)

func init() {
	// do nothing
}

func (s *service) Check(secWebSocketKey, hostExtId, hostId, memberId string, isGuest bool) {
	// do nothing
}

func (s *service) Record(secWebSocketKey string, hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	// do nothing
}

func (s *service) RecordBonus(secWebSocketKey string, bonus interface{}) {
	// do nothing
}
