package api

import (
	"sync"
	"time"
)

type Session struct {
	ExpiredAt time.Time
	Pid       int64
}
type SessionMng struct {
	data map[string]*Session
	mtx  sync.Mutex
}

func (sm *SessionMng) Get(s string) *Session {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	return sm.data[s]
}

func (sm *SessionMng) Set(s string, pid int64) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	sm.data[s] = &Session{
		Pid: pid,
	}
}

var sessionMng = &SessionMng{data: map[string]*Session{}}
