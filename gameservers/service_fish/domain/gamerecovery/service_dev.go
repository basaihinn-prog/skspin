//go:build dev
// +build dev

package gamerecovery

import (
	game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
)

func (s *service) New(secWebSocketKey, gameId string, subgameId int, hostId, hostExtId, memberId string) {
	// do nothing
}

func (s *service) Open(secWebSocketKey, gameId string, subgameId int, hostId, hostExtId, memberId string) *game_recovery_proto.GameRecovery {
	return nil
}

func (s *service) Save(secWebSocketKey, rtpId string, rtpState int, rtpBudget, denominator, bulletCollection uint64, mercenaryInfo map[uint64]uint64) {
	// do nothing
}

func (s *service) Close(secWebSocketKey string) {
	// do nothing
}
