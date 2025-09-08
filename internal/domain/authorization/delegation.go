package authorization

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/martencassel/oidcsim/internal/domain"
)

type ConsentInput struct {
	Client    ClientID
	Subject   SubjectID
	Scopes    []string
	NotBefore time.Time
	NotAfter  time.Time
	Audiences []string
	IPCIDRs   []string
	Resources []string
}

type ConsentResult struct {
	DelegationID DelegationID
}

type ExchangeInput struct {
	CodeID       string
	Client       string
	ClientSecret string // for confidential clients
	RedirectURI  string
	PKCEVerifier string // if PKCE is used
	Audience     string
	CallerIP     net.IP
}

type ExchangeResult struct {
	Token Token
}

type ConsentType string

const (
	ConsentTypeUser  ConsentType = "user"
	ConsentTypeAdmin ConsentType = "admin"
)

type SubjectSelector struct {
	SubjectIDs []SubjectID // explicit list
	GroupIDs   []string    // or group/org identifiers
	AllUsers   bool
}

type DelegationService struct {
	delegations DelegationRepo
	codes       string
	clients     ClientRepo
	tokens      TokenIssuer2
	now         func() time.Time
	idTokens    domain.IDTokenIssuer
}

type IssueCodeInput struct {
	DelegationID  DelegationID
	Client        ClientID
	RedirectURI   string
	PKCEChallenge string        // optional
	TTL           time.Duration // eg. 60s
	Nonce         string        // optional
}

type IssueCodeResult struct {
	CodeID string
}

func newID() string {
	return "d_" + uuid.New().String()
}

func (s *DelegationService) CreateDelegationFromConsent(ctx context.Context, in ConsentInput) (ConsentResult, error) {
	requested := toScopeSet(in.Scopes)
	if !s.clients.AllowsScopes(ctx, in.Client, requested) {
		return ConsentResult{}, ErrScopeNotAllowed
	}
	now := s.now()
	d := Delegation{
		ID:      DelegationID(newID()), // implement UUID/ulid generator in infra
		Actor:   in.Client,
		Subject: in.Subject,
		Scopes:  requested,
		Window: TimeWindow{
			NotBefore: in.NotBefore,
			NotAfter:  in.NotAfter,
		},
		Constraints: Constraints{
			Audiences: toAudiences(in.Audiences),
			Resources: in.Resources,
		},
		CreatedAt: now,
	}
	if err := s.delegations.Save(ctx, d); err != nil {
		return ConsentResult{}, err
	}
	return ConsentResult{DelegationID: d.ID}, nil
}

func (s *DelegationService) ExchangeCodeForTokens(ctx context.Context, in ExchangeInput) (ExchangeResult, error) {
	// now := s.now()
	// // code, err := s.codes.GetForRedemption(ctx, in.CodeID)
	// // if code.ClientID != in.Client {
	// // 	return ExchangeResult{}, ErrClientBindingFailed
	// // }
	// if subtle.ConstantTimeCompare([]byte(code.RedirectURI), []byte(in.RedirectURI)) != 1 {
	// 	return ExchangeResult{}, errors.New("redirect uri mismatch")
	// }
	// if now.After(code.ExpiresAt) || code.UsedAt != nil {
	// 	return ExchangeResult{}, ErrCodeUsedOrExpired
	// }
	// // Optional: verify client secret for confidential clients.
	// if in.ClientSecret != "" && !s.clients.VerifySecret(ctx, in.Client, in.ClientSecret) {
	// 	return ExchangeResult{}, errors.New("invalid client secret")
	// }
	// // PKCE check if present
	// if code.PKCEChallenge != "" {
	// 	if !verifyPKCE(code.PKCEChallenge, "S256", in.PKCEVerifier) {
	// 		return ExchangeResult{}, errors.New("pkce verification failed")
	// 	}
	// }
	// d, err := s.delegations.FindByID(ctx, code.DelegationID)
	// if err != nil {
	// 	return ExchangeResult{}, err
	// }
	// // Validate the delegation for the requested audience at this moment with caller IP.
	// if err := d.Validate(in.Audience, in.CallerIP, d.Scopes, now); err != nil {
	// 	return ExchangeResult{}, err
	// }
	// // Mark code as used (single-use). Use transactional semantics in infra.
	// if err := s.codes.MarkUsed(ctx, code.ID, now); err != nil {
	// 	return ExchangeResult{}, err
	// }
	// access, exp, err := s.tokens.MintAccess(ctx, d, in.Audience, now)
	// if err != nil {
	// 	return ExchangeResult{}, err
	// }
	// refresh, err := s.tokens.MintRefresh(ctx, d, now)
	// if err != nil {
	// 	return ExchangeResult{}, err
	// }

	// var idToken string
	// if d.Scopes.ContainsAll(NewScopeSet("openid")) && s.idTokens != nil {
	// 	//		idToken, err = s.idTokens.IssueIDToken(ctx, d,
	// 	if err != nil {
	// 		return ExchangeResult{}, err
	// 	}
	// }
	// return ExchangeResult{
	// 	Token: Token{
	// 		AccessToken:  access,
	// 		RefreshToken: refresh,
	// 		IDToken:      idToken, // <-- now included
	// 		ExpiresIn:    exp,
	// 		Scopes:       d.Scopes.ToSlice(),
	// 	},
	// }, nil
	return ExchangeResult{}, nil
}

func (s *DelegationService) RevokeDelegation(ctx context.Context, id DelegationID) error {
	return s.delegations.Revoke(ctx, id, s.now())
}

// verifyPKCE checks that the provided verifier matches the stored challenge.
// method is usually "S256" or "plain".
func verifyPKCE(storedChallenge, method, verifier string) bool {
	switch strings.ToUpper(method) {
	case "S256":
		h := sha256.Sum256([]byte(verifier))
		// Base64 URL encoding without padding
		encoded := base64.RawURLEncoding.EncodeToString(h[:])
		return subtleConstantTimeCompare(encoded, storedChallenge)
	case "PLAIN":
		return subtleConstantTimeCompare(verifier, storedChallenge)
	default:
		// Unknown method — fail closed
		return false
	}
}

// subtleConstantTimeCompare avoids timing attacks.
func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	// Use crypto/subtle for constant‑time compare
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (s *DelegationService) IssueAuthorizationCode(ctx context.Context, in IssueCodeInput) (IssueCodeResult, error) {
	d, err := s.delegations.FindByID(ctx, in.DelegationID)
	if err != nil {
		return IssueCodeResult{}, err
	}
	if d.Actor != in.Client {
		return IssueCodeResult{}, ErrClientBindingFailed
	}
	now := s.now()
	if err := d.IsActiveAt(now); err != nil {
		return IssueCodeResult{}, err
	}
	if !s.clients.IsRedirectAllowed(ctx, in.Client, in.RedirectURI) {
		return IssueCodeResult{}, errors.New("redirect uri not allowed")
	}
	// code := AuthorizationCode{
	// 	ID:            AuthorizationCodeID(newOpaqueCode()), // short-lived, single-use
	// 	DelegationID:  d.ID,
	// 	ClientID:      in.Client,
	// 	RedirectURI:   in.RedirectURI,
	// 	PKCEChallenge: in.PKCEChallenge,
	// 	ExpiresAt:     now.Add(in.TTL),
	// 	CreatedAt:     now,
	// 	Nonce:         in.Nonce,
	// }
	// if err := s.codes.Issue(ctx, code); err != nil {
	// 	return IssueCodeResult{}, err
	// }
	// return IssueCodeResult{CodeID: code.ID}, nil
	return IssueCodeResult{CodeID: ""}, nil
}

func newOpaqueCode() string { return "c_" + newULID() }

func newULID() string {
	return uuid.NewString()
}

// Helpers
func toScopeSet(in []string) ScopeSet {
	out := make(ScopeSet, len(in))
	for _, v := range in {
		if v != "" {
			out[Scope(v)] = struct{}{}
		}
	}
	return out
}

func toAudiences(in []string) []Audience {
	out := make([]Audience, 0, len(in))
	for _, v := range in {
		if v != "" {
			out = append(out, Audience(v))
		}
	}
	return out
}

// func NewDelegationService(d DelegationRepo, c AuthorizationCodeRepo, cl ClientRepo, t TokenIssuer2) *DelegationService {
// 	return &DelegationService{
// 		delegations: d,

// 		codes:   c,
// 		clients: cl,
// 		tokens:  t,
// 		now:     time.Now,
// 	}
// }

type DelegationID string
type ClientID string
type SubjectID string
type Audience string
type Scope string

var (
	ErrInvalidDelegation   = errors.New("invalid delegation")
	ErrExpiredDelegation   = errors.New("delegation expired")
	ErrRevokedDelegation   = errors.New("delegation revoked")
	ErrScopeNotAllowed     = errors.New("requested scopes exceed allowed set")
	ErrAudienceNotAllowed  = errors.New("audience not allowed")
	ErrIPNotAllowed        = errors.New("ip not allowed by constraints")
	ErrCodeInvalid         = errors.New("authorization code invalid")
	ErrCodeUsedOrExpired   = errors.New("authorization code used or expired")
	ErrClientBindingFailed = errors.New("authorization code not bound to this client")
)

type TimeWindow struct {
	NotBefore time.Time
	NotAfter  time.Time
}

func (tw TimeWindow) Contains(t time.Time) bool {
	if !tw.NotBefore.IsZero() && t.Before(tw.NotBefore) {
		return false
	}
	if !tw.NotAfter.IsZero() && t.After(tw.NotAfter) {
		return false
	}
	return true
}

type ScopeSet map[Scope]struct{}

func NewScopeSet(scopes ...Scope) ScopeSet {
	s := make(ScopeSet, len(scopes))
	for _, sc := range scopes {
		if sc != "" {
			s[sc] = struct{}{}
		}
	}
	return s
}

func (s ScopeSet) ContainsAll(other ScopeSet) bool {
	for sc := range other {
		if _, ok := s[sc]; !ok {
			return false
		}
	}
	return true
}

func (s ScopeSet) ToSlice() []string {
	out := make([]string, 0, len(s))
	for sc := range s {
		out = append(out, string(sc))
	}
	return out
}

type Delegation struct {
	ID          DelegationID
	Actor       ClientID
	Subject     SubjectID       // for user consent
	Subjects    SubjectSelector // for admin consent to multiple users
	Scopes      ScopeSet
	Window      TimeWindow
	Constraints Constraints
	ConsentType ConsentType
	GrantedBy   string // admin ID or system
	RevokedAt   *time.Time
	CreatedAt   time.Time
}

// Constraints narrow where and how the delegation is valid.
type Constraints struct {
	Audiences []Audience   // allowed token audiences
	IPRanges  []*net.IPNet // allowed caller IPs
	Resources []string     // resource-specific rules (e.g., doc:123, bucket:foo)
}

func (c Constraints) AllowsAudience(a Audience) bool {
	if len(c.Audiences) == 0 {
		return true
	}
	for _, allowed := range c.Audiences {
		if allowed == a {
			return true
		}
	}
	return false
}

func (c Constraints) AllowsIP(ip net.IP) bool {
	if len(c.IPRanges) == 0 || ip == nil {
		return true
	}
	for _, cidr := range c.IPRanges {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func (d Delegation) IsActiveAt(t time.Time) error {
	if d.RevokedAt != nil && !d.RevokedAt.After(t) {
		return ErrRevokedDelegation
	}
	if !d.Window.Contains(t) {
		return ErrExpiredDelegation
	}
	return nil
}

func (d Delegation) Validate(aud Audience, ip net.IP, required ScopeSet, now time.Time) error {
	if err := d.IsActiveAt(now); err != nil {
		return err
	}
	if !d.Scopes.ContainsAll(required) {
		return ErrScopeNotAllowed
	}
	if !d.Constraints.AllowsAudience(aud) {
		return ErrAudienceNotAllowed
	}
	if !d.Constraints.AllowsIP(ip) {
		return ErrIPNotAllowed
	}
	return nil
}

type DelegationRepo interface {
	Save(ctx context.Context, d Delegation) error
	FindByID(ctx context.Context, id DelegationID) (Delegation, error)
	Revoke(ctx context.Context, id DelegationID, at time.Time) error
}

type Token struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int // seconds
	Scopes       []string
}

type TokenIssuer2 interface {
	// MintAccess should embed subject, client, scopes, aud, and delegation id as claims.
	MintAccess(ctx context.Context, d Delegation, aud Audience, now time.Time) (string, int, error)
	// MintRefresh may bind to delegation and client; consider rotation.
	MintRefresh(ctx context.Context, d Delegation, now time.Time) (string, error)
}

type ClientRepo interface {
	// Used to validate client and its allowed scopes/redirect URIs if you enforce them here.
	IsRedirectAllowed(ctx context.Context, client ClientID, redirectURI string) bool
	AllowsScopes(ctx context.Context, client ClientID, scopes ScopeSet) bool
	VerifySecret(ctx context.Context, client ClientID, secret string) bool
}
