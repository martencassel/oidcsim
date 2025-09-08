package authentication

// type AuthService interface {
// 	Initiate(ctx context.Context, clientID string) (*AuthFlowSpec, error)
// }

// type authService struct {
// 	sessions SessionRepo
// 	engine   FlowEngine
// 	authTTL  time.Duration
// 	stateTTL time.Duration
// }

// func NewAuthService(sessions SessionRepo, engine FlowEngine, authTTL, stateTTL time.Duration) AuthService {
// 	return &authService{sessions: sessions, engine: engine, authTTL: authTTL, stateTTL: stateTTL}
// }

// func (s *authService) NextStep(sessionID string) domauth.StepSpec {
// 	state, exists, err := s.sessions.GetFlowState(context.Background(), sessionID)
// 	if err != nil || !exists {
// 		return domauth.StepSpec{}
// 	}
// 	spec, err := s.engine.Plan(context.Background(), "")
// 	if err != nil {
// 		return domauth.StepSpec{}
// 	}
// 	for _, step := range spec.Steps {
// 		// Simple heuristic: return first step that has not yet been completed
// 		if _, done := state[string(step.Method)]; !done {
// 			return step
// 		}
// 	}
// 	return domauth.StepSpec{}
// }

// func (s *authService) Current(ctx context.Context, sid string) (domauth.Context, bool, error) {
// 	return s.sessions.GetAuthContext(ctx, sid)
// }

// func (s *authService) Initiate(ctx context.Context, sid, clientID string) (domauth.FlowSpec, error) {
// 	spec, err := s.engine.Plan(ctx, clientID)
// 	if err != nil {
// 		return domauth.FlowSpec{}, err
// 	}
// 	// reset any previous flow state
// 	_ = s.sessions.ClearFlowState(ctx, sid)
// 	_ = s.sessions.SetFlowState(ctx, sid, map[string]string{}, s.stateTTL)
// 	return spec, nil
// }

// func (s *authService) StartStep(ctx context.Context, sid string, step domauth.StepSpec) (map[string]string, error) {
// 	state, _, err := s.sessions.GetFlowState(ctx, sid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hints, err := s.engine.StartStep(ctx, step, state)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// persist any state hints if engine used them
// 	if err := s.sessions.SetFlowState(ctx, sid, state, s.stateTTL); err != nil {
// 		return nil, err
// 	}
// 	return hints, nil
// }

// func (s *authService) CompleteStep(ctx context.Context, sid string, step domauth.StepSpec, inputs map[string]string) (bool, domauth.StepSpec, map[string]string, error) {
// 	state, _, err := s.sessions.GetFlowState(ctx, sid)
// 	if err != nil {
// 		return false, domauth.StepSpec{}, nil, err
// 	}

// 	done, updates, err := s.engine.CompleteStep(ctx, step, inputs, state)
// 	if err != nil {
// 		return false, domauth.StepSpec{}, nil, err
// 	}
// 	// merge updates
// 	for k, v := range updates {
// 		state[k] = v
// 	}
// 	if err := s.sessions.SetFlowState(ctx, sid, state, s.stateTTL); err != nil {
// 		return false, domauth.StepSpec{}, nil, err
// 	}
// 	if !done {
// 		// Engine determines next step from state (common pattern) — or caller does
// 		// For simplicity we let caller ask engine.Plan again to know next (not shown)
// 		return false, domauth.StepSpec{}, nil, nil
// 	}
// 	// Build final auth context and persist
// 	authCtx, err := s.engine.BuildAuthContext(ctx, state)
// 	if err != nil {
// 		return false, domauth.StepSpec{}, nil, err
// 	}
// 	if err := s.sessions.SetAuthContext(ctx, sid, authCtx, s.authTTL); err != nil {
// 		return false, domauth.StepSpec{}, nil, err
// 	}
// 	_ = s.sessions.ClearFlowState(ctx, sid)
// 	return true, domauth.StepSpec{}, nil, nil
// }
