//go:build staging || prod
// +build staging prod

package activity

import (
	"encoding/json"
	errorcode "serve/fish_comm/flux/error-code"
	"serve/fish_comm/flux/logger"
	"serve/fish_comm/flux/redis"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/rank"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"

	redigo "github.com/gomodule/redigo/redis"
)

func init() {
	logger.Service.Zap.Infow("Service Created",
		"Service", Service.Id,
	)
}

func (s *service) Check(secWebSocketKey, hostExtId, hostId, memberId string, isGuest bool) {

	if v, ok := s.activities.Load(secWebSocketKey); ok {
		go v.(*activity).run()
		return
	}

	activities := s.activityOn(secWebSocketKey, hostExtId, hostId, memberId)

	if activities == nil {
		return
	}

	for _, uuid := range activities {

		redisActivity := s.hget(secWebSocketKey, hostExtId, hostId, uuid)

		if redisActivity == nil {
			logger.Service.Zap.Errorw(Activity_REDIS_HGET_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"HostId", hostId,
				"HostExtId", hostExtId,
				"Field", activities,
			)
			errorcode.Service.Fatal(secWebSocketKey, Activity_REDIS_HGET_NOT_FOUND)
			return
		}

		a := &activity{
			secWebSocketKey: secWebSocketKey,
			hostExtId:       hostExtId,
			hostId:          hostId,
			memberId:        memberId,
			isStop:          make(chan bool),
			isGuest:         isGuest,
		}

		if err := json.Unmarshal(redisActivity, a); err != nil {
			logger.Service.Zap.Errorw(Activity_JSON_DECODE_INVALID,
				"GameUser", secWebSocketKey,
				"HostId", hostId,
				"HostExtId", hostExtId,
				"JSON", string(redisActivity),
			)
			errorcode.Service.Fatal(secWebSocketKey, Activity_JSON_DECODE_INVALID)
			return
		}

		logger.Service.Zap.Infow("Activity",
			"GameUser", secWebSocketKey,
			"HostId", hostId,
			"HostExtId", hostExtId,
			"Data", string(redisActivity),
		)

		// activity did not close yet
		if s.nowUTC() < a.Time[3] {
			s.activities.Store(secWebSocketKey, a)
			go a.run()

			break
		}
	}
}

func (s *service) activityOn(secWebSocketKey, hostExtId, hostId, memberId string) []string {
	settings := s.hget(secWebSocketKey, hostExtId, hostId, "activities")

	if settings == nil {
		logger.Service.Zap.Infow("No Activities",
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
			"HostId", hostId,
			"HostExtId", hostExtId,
		)
		return nil
	}

	as := &activities{}

	if err := json.Unmarshal(settings, as); err != nil {
		logger.Service.Zap.Errorw(Activity_JSON_DECODE_INVALID,
			"GameUser", secWebSocketKey,
			"HostId", hostId,
			"HostExtId", hostExtId,
			"JSON", string(settings),
		)
		errorcode.Service.Fatal(secWebSocketKey, Activity_JSON_DECODE_INVALID)
		return nil
	}

	if len(as.On) == 0 {
		return nil
	}

	return as.On
}

func (s *service) Record(secWebSocketKey string, hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability) {
	if hitResult.Pay <= 0 {
		return
	}

	if !s.isLimits(secWebSocketKey, hitFish.GameId, hitFish.TypeId) ||
		!s.isMinBet(secWebSocketKey, hitBullet.Bet*hitBullet.Rate) ||
		!s.isStart(secWebSocketKey) ||
		s.isGuest(secWebSocketKey) {
		return
	}

	// TODO JOHNNY 後續clean code
	if a, ok := s.activities.Load(secWebSocketKey); ok {
		v := a.(*activity)

		rank.Service.New(
			rank.Builder().
				SetActivityUuid(v.Uuid).
				SetSecWebSocketKey(secWebSocketKey).
				SetMemberId(v.memberId).
				SetHostId(v.hostId).
				SetHostExtId(v.hostExtId).
				SetCent(uint64(hitResult.Pay) * hitBullet.Bet * hitBullet.Rate * uint64(hitResult.Multiplier)).
				SetKind(v.Kind).
				SetConn(v.hostExtId).
				Build(),
		)
	}
}

func (s *service) RecordBonus(secWebSocketKey string, bonus interface{}) {
	var hitFish *fish.Fish
	var hitBullet *bullet.Bullet
	var fishTypeId int32 = -1
	cent := uint64(0)

	switch bonus.(type) {
	case *redenvelope.RedEnvelope:
		re := bonus.(*redenvelope.RedEnvelope)

		hitFish = re.ExtraData[0].(*fish.Fish)
		hitBullet = re.ExtraData[1].(*bullet.Bullet)

		fishTypeId = re.FishTypeId

		cent = re.Pay * hitBullet.Bet * hitBullet.Rate

	case *slot.Slot:
		sl := bonus.(*slot.Slot)

		hitFish = sl.ExtraData[0].(*fish.Fish)
		hitBullet = sl.ExtraData[1].(*bullet.Bullet)

		fishTypeId = sl.FishTypeId

		cent = sl.Pay * hitBullet.Bet * hitBullet.Rate

	default:
		panic("ERROR BONUS TYPE")
	}

	if !s.isLimits(secWebSocketKey, hitFish.GameId, fishTypeId) ||
		!s.isMinBet(secWebSocketKey, hitBullet.Bet*hitBullet.Rate) ||
		!s.isStart(secWebSocketKey) ||
		s.isGuest(secWebSocketKey) {
		return
	}

	// TODO JOHNNY 後續clean code
	if a, ok := s.activities.Load(secWebSocketKey); ok {
		v := a.(*activity)

		rank.Service.New(
			rank.Builder().
				SetActivityUuid(v.Uuid).
				SetSecWebSocketKey(secWebSocketKey).
				SetMemberId(v.memberId).
				SetHostId(v.hostId).
				SetHostExtId(v.hostExtId).
				SetCent(cent).
				SetKind(v.Kind).
				SetConn(v.hostExtId).
				Build(),
		)
	}
}

func (s *service) isGuest(secWebSocketKey string) bool {
	if a, ok := s.activities.Load(secWebSocketKey); ok {
		return a.(*activity).isGuest
	}
	return true
}

func (s *service) isMinBet(secWebSocketKey string, bet uint64) bool {
	if a, ok := s.activities.Load(secWebSocketKey); ok {
		for _, v := range a.(*activity).MinBet {

			if v <= bet {
				return true
			}
		}
	}
	return false
}

func (s *service) isLimits(secWebSocketKey, gameId string, iconId int32) bool {
	if a, ok := s.activities.Load(secWebSocketKey); ok {
		// TODO JOHNNY 後續再來改效能

		for _, games := range a.(*activity).Limits {
			if icons, ok := games[gameId]; ok {
				for _, v := range icons.([]interface{}) {

					if int32(v.(float64)) == iconId {

						return true
					}
				}
			}

		}
	}
	return false
}

func (s *service) isStart(secWebSocketKey string) bool {
	if a, ok := s.activities.Load(secWebSocketKey); ok {

		if a.(*activity).status == status_START { // TODO JOHNNY maybe data race read
			return true
		}
	}
	return false
}

func (s *service) hget(secWebSocketKey, hostExtId, hostId, field string) []uint8 {
	if c := redis.Repository.Conn(hostExtId); c != nil {
		defer c.Close()

		reply, err := redigo.Bytes(c.Do("HGET", hostId, field))

		if err != nil || reply == nil {
			logger.Service.Zap.Warnw(Activity_REDIS_HGET_ERROR,
				"GameUser", secWebSocketKey,
				"HostId", hostId,
				"HostExtId", hostExtId,
				"Field", field,
				"Error", err,
			)
			return nil
		}

		return reply
	}

	logger.Service.Zap.Errorw(Activity_REDIS_CONNECTION_NOT_FOUND,
		"GameUser", secWebSocketKey,
		"HostId", hostId,
		"HostExtId", hostExtId,
		"Field", field,
	)
	errorcode.Service.Fatal(secWebSocketKey, Activity_REDIS_CONNECTION_NOT_FOUND)
	return nil
}
