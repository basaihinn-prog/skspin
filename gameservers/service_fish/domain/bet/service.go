package bet

import (
	common_proto "serve/service_fish/models/proto"
	"sync"
)

var Service = &service{
	Id:      "BetService",
	bets:    make(map[string]*common_proto.BetRateLine),
	rwMutex: &sync.RWMutex{},
}

type service struct {
	Id      string
	bets    map[string]*common_proto.BetRateLine
	rwMutex *sync.RWMutex
}

func (s *service) Get(secWebSocketKey string) *common_proto.BetRateLine {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if v, ok := s.bets[secWebSocketKey]; ok {
		return v
	}
	return nil
}

func (s *service) Set(secWebSocketKey string, betRateLine *common_proto.BetRateLine) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.bets[secWebSocketKey] = betRateLine
}

func (s *service) Delete(secWebSocketKey string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	delete(s.bets, secWebSocketKey)
}
