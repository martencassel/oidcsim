package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	authn "github.com/martencassel/oidcsim/internal/domain/authentication"

	"github.com/martencassel/oidcsim/internal/interface/http/dto"
)

type memorySessionManager struct {
	cookieName    string
	sessions      map[string]*sessionData
	mu            sync.RWMutex
	allowInsecure bool
}

type sessionData struct {
	AuthCtx   authn.Context
	AuthzReq  *dto.AuthorizeRequest
	FlowState map[string]string
}

type SessionOptions func(*memorySessionManager)

func WithAllowInsecure() SessionOptions {
	return func(m *memorySessionManager) {
		m.allowInsecure = true
	}
}

func NewMemorySessionManager(cookieName string, opts ...SessionOptions) *memorySessionManager {
	if cookieName == "" {
		cookieName = "session_id"
	}
	msm := &memorySessionManager{
		sessions:      make(map[string]*sessionData),
		mu:            sync.RWMutex{},
		cookieName:    cookieName,
		allowInsecure: false,
	}
	// apply options
	for _, opt := range opts {
		opt(msm)
	}
	return msm
}

func (m *memorySessionManager) Ensure(w http.ResponseWriter, r *http.Request) (sessionID string, err error) {
	// Check if session cookie exists
	cookie, err := r.Cookie(m.cookieName)
	if err == nil {
		// Session cookie exists, return its value
		return cookie.Value, nil
	}
	// Create a new session ID using crypto/rand
	newSessionID, err := genSessionID()
	if err != nil {
		return "", err
	}

	// Store the new session ID in memory
	m.mu.Lock()
	m.sessions[newSessionID] = &sessionData{}
	m.mu.Unlock()

	// Set the session cookie with secure attributes
	secure := !m.allowInsecure && r.TLS != nil
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // 1 hour default
	})

	return newSessionID, nil
}

func (m *memorySessionManager) Current(r *http.Request) (sessionID string, ok bool) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return "", false
	}
	m.mu.RLock()
	_, exists := m.sessions[cookie.Value]
	m.mu.RUnlock()
	return cookie.Value, exists
}

func (m *memorySessionManager) Rotate(w http.ResponseWriter, r *http.Request) (newID string, err error) {
	oldCookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return "", fmt.Errorf("no existing session to rotate")
	}
	oldSessionID := oldCookie.Value

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.sessions[oldSessionID]; !exists {
		return "", fmt.Errorf("session does not exist")
	}
	// Remove old session
	delete(m.sessions, oldSessionID)

	// Create new session ID using crypto/rand
	newSessionID, err := genSessionID()
	if err != nil {
		return "", err
	}
	m.sessions[newSessionID] = &sessionData{}

	// Set new session cookie with secure attributes
	secure := !m.allowInsecure && r.TLS != nil
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})
	return newSessionID, nil
}

func (m *memorySessionManager) Destroy(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return fmt.Errorf("no session to destroy")
	}
	sessionID := cookie.Value

	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, sessionID)

	// Remove the session cookie (preserve flags)
	secure := !m.allowInsecure && r.TLS != nil
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	return nil
}

// genSessionID creates a cryptographically secure random session id encoded URL-safe base64.
func genSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (m *memorySessionManager) SaveAuthorizeRequest(sid string, req dto.AuthorizeRequest) error {
	m.mu.Unlock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sid]
	if !ok {
		s = &sessionData{}
		m.sessions[sid] = s
	}
	s.AuthzReq = &req
	return nil
}

func (m *memorySessionManager) GetAuthorizeRequest(sid string) (dto.AuthorizeRequest, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sid]
	if !ok || s.AuthzReq == nil {
		return dto.AuthorizeRequest{}, false, nil
	}
	return *s.AuthzReq, true, nil
}

func (m *memorySessionManager) ClearAuthorizeRequest(sid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sid]; ok {
		s.AuthzReq = nil
	}
	return nil
}

func (m *memorySessionManager) GetID(sid string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, exists := m.sessions[sid]; exists {
		return sid
	}
	return ""
}
