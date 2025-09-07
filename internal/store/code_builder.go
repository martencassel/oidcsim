package store

// import (
// 	"crypto/rand"
// 	"encoding/base64"
// 	"time"

// 	"github.com/martencassel/oidcsim/internal/dto"
// )

// type AuthorizationCodeBuilder struct {
// 	code AuthorizationCode
// }

// func NewAuthorizationCodeBuilder() *AuthorizationCodeBuilder {
// 	return &AuthorizationCodeBuilder{code: AuthorizationCode{}}
// }

// func (b *AuthorizationCodeBuilder) FromRequest(req dto.AuthorizeRequest) *AuthorizationCodeBuilder {
// 	b.code.ClientID = req.ClientID
// 	b.code.RedirectURI = req.RedirectURI
// 	b.code.Scope = req.Scope
// 	b.code.State = req.State
// 	b.code.CodeChallenge = CodeChallenge(req.CodeChallenge)
// 	b.code.CodeChallengeMethod = CodeChallengeMethod(req.CodeChallengeMethod)
// 	b.code.Nonce = Nonce(req.Nonce)
// 	return b
// }

// func (b *AuthorizationCodeBuilder) WithUser(userID, sessionID string, authTime time.Time) *AuthorizationCodeBuilder {
// 	b.code.UserID = userID
// 	b.code.SessionID = sessionID
// 	b.code.AuthTime = authTime
// 	return b
// }

// func (b *AuthorizationCodeBuilder) WithCode(code string) *AuthorizationCodeBuilder {
// 	b.code.Code = code
// 	return b
// }

// func (b *AuthorizationCodeBuilder) WithExpiry(ttl time.Duration) *AuthorizationCodeBuilder {
// 	b.code.ExpiresAt = time.Now().Add(ttl)
// 	return b
// }

// func (b *AuthorizationCodeBuilder) Build() AuthorizationCode {
// 	return b.code
// }

// func generateRandomString() string {
// 	b := make([]byte, 32)
// 	if _, err := rand.Read(b); err != nil {
// 		panic(err) // or handle gracefully
// 	}
// 	return base64.RawURLEncoding.EncodeToString(b)
// }
