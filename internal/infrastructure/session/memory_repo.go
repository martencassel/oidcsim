package session

import (
	"context"
	"sync"
	"time"

	appauth "github.com/martencassel/oidcsim/internal/application/authentication"
	domauth "github.com/martencassel/oidcsim/internal/domain/authentication"
)

type memoryRepo struct {
	mu   sync.RWMutex
	auth map[string]memAuth
	flow map[string]memState
}

type memAuth struct {
	v   domauth.Context
	exp time.Time
}
type memState struct {
	v   map[string]string
	exp time.Time
}

func NewMemoryRepo() appauth.SessionRepo {
	return &memoryRepo{
		auth: make(map[string]memAuth),
		flow: make(map[string]memState),
	}
}

func (m *memoryRepo) GetAuthContext(_ context.Context, sid string) (domauth.Context, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.auth[sid]
	if !ok || time.Now().After(a.exp) {
		return domauth.Context{}, false, nil
	}
	return a.v, true, nil
}

func (m *memoryRepo) SetAuthContext(_ context.Context, sid string, v domauth.Context, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auth[sid] = memAuth{v: v, exp: time.Now().Add(ttl)}
	return nil
}

func (m *memoryRepo) ClearAuthContext(_ context.Context, sid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.auth, sid)
	return nil
}

func (m *memoryRepo) GetFlowState(_ context.Context, sid string) (map[string]string, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.flow[sid]
	if !ok || time.Now().After(s.exp) {
		return map[string]string{}, false, nil
	}
	// shallow copy to avoid aliasing
	out := make(map[string]string, len(s.v))
	for k, v := range s.v {
		out[k] = v
	}
	return out, true, nil
}

func (m *memoryRepo) SetFlowState(_ context.Context, sid string, state map[string]string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cpy := make(map[string]string, len(state))
	for k, v := range state {
		cpy[k] = v
	}
	m.flow[sid] = memState{v: cpy, exp: time.Now().Add(ttl)}
	return nil
}

func (m *memoryRepo) ClearFlowState(_ context.Context, sid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.flow, sid)
	return nil
}
