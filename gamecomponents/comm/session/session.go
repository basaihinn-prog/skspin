package session

import (
	"errors"
	"game/duck/lang"
	"game/duck/ut2"
	"reflect"
)

type IManager[T ISession] interface {
	Get(pid int64) T
}

// 这是一个通用的PlayerManager
type Manager[T ISession] struct {
	new func() T
	mp  ut2.IMap[int64, T]
}

func NewManager[T ISession]() *Manager[T] {
	var val T
	v := reflect.TypeOf(val)
	v = v.Elem()

	m := &Manager[T]{}

	m.new = func() T {
		return reflect.New(v).Interface().(T)
	}

	m.mp = ut2.NewSyncMap[int64, T]()

	return m
}

func (m *Manager[T]) Remove(pid int64) {
	m.mp.Delete(pid)
}

func (m *Manager[T]) Get(pid int64) (t T) {
	r, ok := m.mp.Load(pid)
	if !ok {
		return
	}
	return r
}

func (m *Manager[T]) GetOrCreate(pid int64) T {
	r, ok := m.mp.Load(pid)
	if !ok {
		r = m.new()
		s := r.mustEmbedSession()
		s.Pid = pid

		// TODO: Language 应该从登录链接里面拿到
		s.Language = "th"

		m.mp.Store(pid, r)
	}
	return r
}

//

type ISession interface {
	mustEmbedSession() *Session
}

type Session struct {
	Pid      int64
	Language string
}

func (ps *Session) Err(err string) error {
	return errors.New(lang.Get(ps.Language, err))
}

func (ps *Session) Errf(id string, data any, plural int) error {
	return errors.New(lang.Translate(ps.Language, id, data, plural))
}

func (ps *Session) mustEmbedSession() *Session {
	return ps
}
