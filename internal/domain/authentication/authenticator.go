package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/martencassel/oidcsim/internal/domain/user"
)

const (
	ErrSessionNotFound = "session not found"
)

type AuthSession struct {
	ID            string
	SubjectID     string
	CreatedAt     int64
	Authenticated bool
	AuthTime      int64
	Flow          []AuthStep
	Completed     []AuthMethod
	Methods       []AuthMethod
}

type SessionStore interface {
	IsAuthenticated(sessionID string) bool
	GetSubjectID(sessionID string) string
	GetAuthTime(sessionID string) int64
	GetAuthFlow(sessionID string) []AuthStep
	GetCompletedSteps(sessionID string) []AuthMethod
	Save(session AuthSession) error
	Get(id string) (AuthSession, bool)
	Delete(id string) error
	SetSubject(sessionID, subjectID string) error
	MarkAuthenticated(sessionID string) error
}

type DefaultAuthService struct {
	sessionStore SessionStore
	userRepo     user.UserRepository
}

func NewDefaultAuthService(store SessionStore, repo user.UserRepository) *DefaultAuthService {
	return &DefaultAuthService{
		sessionStore: store,
		userRepo:     repo,
	}
}

func (s *DefaultAuthService) Initiate(ctx context.Context, clientID string) (*AuthFlowSpec, error) {
	return &AuthFlowSpec{
		Steps: []AuthStep{
			{Method: MethodPassword},
			{Method: MethodOTP},
		},
	}, nil
}

func (s *DefaultAuthService) StartStep(ctx context.Context, sessionID string, method AuthMethod) (*StepUI, error) {
	switch method {
	case MethodPassword:
		return &StepUI{Prompt: "Enter your username and password", Fields: []string{"username", "password"}}, nil
	case MethodOTP:
		return &StepUI{Prompt: "Enter the OTP sent to your device", Fields: []string{"otp"}}, nil
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

func (s *DefaultAuthService) CompleteStep(ctx context.Context, sessionID string, method AuthMethod, inputs map[string]string) (bool, error) {
	switch method {
	case MethodPassword:
		user, err := s.userRepo.Authenticate(ctx, inputs["username"], inputs["password"])
		if err != nil {
			return false, err
		}
		s.sessionStore.SetSubject(sessionID, user.ID)
		return false, nil // continue to OTP
	case MethodOTP:
		if inputs["otp"] != "123456" {
			return false, fmt.Errorf("invalid OTP")
		}
		s.sessionStore.MarkAuthenticated(sessionID)
		return true, nil
	default:
		return false, fmt.Errorf("unsupported method")
	}
}

func contains(slice []AuthMethod, item AuthMethod) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s *DefaultAuthService) NextStep(sessionID string) AuthStep {
	flow := s.sessionStore.GetAuthFlow(sessionID) // e.g. slice of AuthStep
	completed := s.sessionStore.GetCompletedSteps(sessionID)
	for _, step := range flow {
		if !contains(completed, step.Method) {
			return step
		}
	}
	return AuthStep{} // or panic if flow is exhausted
}

func (s *DefaultAuthService) Current(ctx context.Context, sessionID string) (AuthContext, bool, error) {
	if !s.sessionStore.IsAuthenticated(sessionID) {
		return AuthContext{}, false, nil
	}
	subjectID := s.sessionStore.GetSubjectID(sessionID)
	authTime := s.sessionStore.GetAuthTime(sessionID)
	return AuthContext{
		SubjectID: subjectID,
		AuthTime:  time.Unix(authTime, 0),
	}, true, nil
}
