package accountingmanager

import "sync"

var Service = &service{
	id:          "AccountingManagerService",
	accountings: make(map[string]*AccountingManager),
	rwMutex:     &sync.RWMutex{},
}

type service struct {
	id          string
	accountings map[string]*AccountingManager
	rwMutex     *sync.RWMutex
}

func (s *service) New(a *AccountingManager) *AccountingManager {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if a == nil {
		return nil
	}

	s.accountings[a.secWebSocketKey] = a
	return a
}

func (s *service) Get(secWebSocketKey string) *AccountingManager {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if v, ok := s.accountings[secWebSocketKey]; ok {
		return v
	}
	return nil
}

func (s *service) GetAccountingSn(secWebSocketKey string) uint64 {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if v, ok := s.accountings[secWebSocketKey]; ok {
		return v.accountingSn
	}
	return 0
}

func (s *service) hashId(secWebSocketKey string) string {
	return secWebSocketKey + s.id
}

func (s *service) Delete(secWebSocketKey string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if v, ok := s.accountings[secWebSocketKey]; ok {
		v.destroy()
		delete(s.accountings, secWebSocketKey)
	}

}
