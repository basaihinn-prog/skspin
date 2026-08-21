package bullet

import common_fish_proto "serve/service_fish/models/proto"

type Bullet struct {
	BonusUuid       string
	Uuid            string
	SecWebSocketKey string
	RoomUuid        string
	RtpId           string
	Bet             uint64
	Rate            uint64
	PlayerCent      uint64
	LineIndex       int32
	RateIndex       int32
	BetIndex        int32
	BetLevelIndex   int32
	TypeId          int32
	Bullets         int32
	ShootBullets    int32 // count client side bullets
	UsedBullets     int32 // count server side bullets
	ShootMode       int32
	Target          *common_fish_proto.Point
	ExtraData       []interface{}
	hitLock         bool
}

func New(bonusUuid, uuid, secWebSocketKey, roomUuid string,
	betIndex, lineIndex, rateIndex, betLevelIndex, typeId, bullets, shootBullets, usedBullets, shootMode int32,
	target *common_fish_proto.Point) *Bullet {
	return &Bullet{
		BonusUuid:       bonusUuid,
		Uuid:            uuid,
		SecWebSocketKey: secWebSocketKey,
		RoomUuid:        roomUuid,
		BetIndex:        betIndex,
		LineIndex:       lineIndex,
		RateIndex:       rateIndex,
		BetLevelIndex:   betLevelIndex,
		TypeId:          typeId,
		Bullets:         bullets,
		ShootBullets:    shootBullets,
		UsedBullets:     usedBullets,
		ShootMode:       shootMode,
		Target:          target,
		ExtraData:       make([]interface{}, 4),
		hitLock:         false,
	}
}
