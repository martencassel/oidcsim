package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/martencassel/oidcsim/internal/domain/authentication"
)

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]authentication.AuthSession
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]authentication.AuthSession),
	}
}

func (s *InMemorySessionStore) Save(session authentication.AuthSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	return nil
}

func (s *InMemorySessionStore) Get(id string) (authentication.AuthSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	return session, ok
}

func (s *InMemorySessionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func (s *InMemorySessionStore) IsAuthenticated(sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[sessionID]
	return ok && sess.Authenticated
}

func (s *InMemorySessionStore) GetSubjectID(sessionID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[sessionID].SubjectID
}

func (s *InMemorySessionStore) GetAuthTime(sessionID string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[sessionID].AuthTime
}

func (s *InMemorySessionStore) GetAuthFlow(sessionID string) []authentication.AuthStep {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[sessionID].Flow
}

func (s *InMemorySessionStore) GetCompletedSteps(sessionID string) []authentication.AuthMethod {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[sessionID].Completed
}

func (s *InMemorySessionStore) SetSubject(sessionID, subjectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	sess.SubjectID = subjectID
	s.sessions[sessionID] = sess
	return nil
}

func (s *InMemorySessionStore) MarkAuthenticated(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	sess.Authenticated = true
	sess.AuthTime = time.Now().Unix()
	s.sessions[sessionID] = sess
	return nil
}
