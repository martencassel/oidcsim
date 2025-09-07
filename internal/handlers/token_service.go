package handlers

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/martencassel/oidcsim/authcode"
	"github.com/martencassel/oidcsim/internal/identity"
)

type RoutesConfig struct {
	Discovery  string `yaml:"discovery"`
	JWKS       string `yaml:"jwks"`
	Authorize  string `yaml:"authorize"`
	Token      string `yaml:"token"`
	Userinfo   string `yaml:"userinfo"`
	Introspect string `yaml:"introspect"`
	Revoke     string `yaml:"revoke"`
	Logout     string `yaml:"logout"`
}

type TokenServiceController struct {
	issuer         string
	routesConfig   *RoutesConfig
	jwks           []byte
	codeStore      *authcode.Store
	privSigningKey *rsa.PrivateKey
	idStore        *identity.CoreIdentityStore
}

type TokenServiceControllerBuilder struct {
	controller *TokenServiceController
}

func NewTokenServiceControllerBuilder() *TokenServiceControllerBuilder {
	return &TokenServiceControllerBuilder{
		controller: &TokenServiceController{},
	}
}

func (b *TokenServiceControllerBuilder) WithIssuer(issuer string) *TokenServiceControllerBuilder {
	b.controller.issuer = issuer
	return b
}

func (b *TokenServiceControllerBuilder) WithRoutesConfig(cfg *RoutesConfig) *TokenServiceControllerBuilder {
	b.controller.routesConfig = cfg
	return b
}

func (b *TokenServiceControllerBuilder) WithJWKS(jwks []byte) *TokenServiceControllerBuilder {
	b.controller.jwks = jwks
	return b
}

func (b *TokenServiceControllerBuilder) WithCodeStore(store *authcode.Store) *TokenServiceControllerBuilder {
	b.controller.codeStore = store
	return b
}

func (b *TokenServiceControllerBuilder) WithSigningKey(key *rsa.PrivateKey) *TokenServiceControllerBuilder {
	b.controller.privSigningKey = key
	return b
}

func (b *TokenServiceControllerBuilder) WithIdentityStore(store *identity.CoreIdentityStore) *TokenServiceControllerBuilder {
	b.controller.idStore = store
	return b
}

func (b *TokenServiceControllerBuilder) Build() *TokenServiceController {
	return b.controller
}

func (ts *TokenServiceController) RegisterRoutes(r gin.IRoutes) {
	r.GET(ts.routesConfig.Discovery, ts.DiscoveryHandler)    // /.well-known/openid-configuration
	r.GET(ts.routesConfig.JWKS, ts.JWKSHandler)              // /.well-known/jwks.json
	r.GET(ts.routesConfig.Authorize, ts.AuthorizeHandler)    // /authorize
	r.POST(ts.routesConfig.Token, ts.TokenHandler)           // /token
	r.GET(ts.routesConfig.Userinfo, ts.UserInfoHandler)      // /userinfo
	r.POST(ts.routesConfig.Introspect, ts.IntrospectHandler) // /introspect
	r.POST(ts.routesConfig.Revoke, ts.RevokeHandler)         // /revoke
	r.POST(ts.routesConfig.Logout, ts.LogoutHandler)         // /logout (RP-Initiated Logout)
}

// JWKSHandler
func (ts *TokenServiceController) JWKSHandler(c *gin.Context) {
	jwks := ts.jwks
	c.Data(http.StatusOK, "application/json", jwks)
}

type AuthorizationRequest struct {
	Scope               string `json:"scope"`
	State               string `json:"state"`
	Nonce               string `json:"nonce"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	ResponseType        string `json:"response_type"`
	ClientID            string `json:"client_id"`
	RedirectURI         string `json:"redirect_uri"`
}

type AuthorizationResponse struct {
	Issuer      string
	Code        string
	State       string
	RedirectURI string
}

// RedirectToClient sends a 302 redirect to the client with code and state
func (r *AuthorizationResponse) RedirectToClient(w http.ResponseWriter) {
	redirectURL, err := url.Parse(r.RedirectURI)
	if err != nil {
		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}
	// Add query parameters
	q := redirectURL.Query()
	q.Set("code", r.Code)
	q.Set("state", r.State)
	q.Set("iss", r.Issuer)
	redirectURL.RawQuery = q.Encode()
	http.Redirect(w, &http.Request{}, redirectURL.String(), http.StatusFound)
}

// AuthorizeHandler
func (ts *TokenServiceController) AuthorizeHandler(c *gin.Context) {
	// Bind it
	authReq := AuthorizationRequest{
		Scope:               c.Query("scope"),
		State:               c.Query("state"),
		Nonce:               c.Query("nonce"),
		CodeChallenge:       c.Query("code_challenge"),
		CodeChallengeMethod: c.Query("code_challenge_method"),
		ResponseType:        c.Query("response_type"),
		ClientID:            c.Query("client_id"),
		RedirectURI:         c.Query("redirect_uri"),
	}
	log.Infof("Authorization request: %+v", authReq)

	issuedCode, err := ts.codeStore.Generate(authReq.ClientID, authReq.RedirectURI)
	if err != nil {
		log.Errorf("Failed to generate authorization code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate authorization code"})
		return
	}
	response := AuthorizationResponse{
		Issuer:      ts.issuer,
		Code:        issuedCode,
		State:       authReq.State,
		RedirectURI: authReq.RedirectURI,
	}
	response.RedirectToClient(c.Writer)
}

// TokenRequest represents the body of a POST /token request
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id,omitempty"`     // optional if using Basic Auth
	ClientSecret string `json:"client_secret,omitempty"` // optional if using Basic Auth
}

// ParseTokenRequest parses a POST /token request body into a TokenRequest struct
func ParseTokenRequest(r *http.Request) (*TokenRequest, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	// Check Authorization header for Basic Auth
	auth := r.Header.Get("Authorization")
	if auth != "" && len(auth) > 6 && auth[:6] == "Basic " {
		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		if err == nil {
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) == 2 {
				clientID = parts[0]
				clientSecret = parts[1]
			}
		}
	}

	return &TokenRequest{
		GrantType:    r.FormValue("grant_type"),
		Code:         r.FormValue("code"),
		RedirectURI:  r.FormValue("redirect_uri"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

// TokenResponse represents the JSON response from the token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`              // usually "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // in seconds
	IDToken      string `json:"id_token,omitempty"`      // OIDC-specific
	RefreshToken string `json:"refresh_token,omitempty"` // optional
	Scope        string `json:"scope,omitempty"`         // optional
}

// WriteTokenResponse writes the token response as JSON to the http.ResponseWriter
func WriteTokenResponse(w http.ResponseWriter, resp TokenResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// TokenHandler
func (ts *TokenServiceController) TokenHandler(c *gin.Context) {
	// Parse the token request
	tokenReq, err := ParseTokenRequest(c.Request)
	if err != nil {
		log.Errorf("Failed to parse token request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token request"})
		return
	}

	userGroups, err := ts.idStore.GetUserGroups(c.Request.Context(), "alice")
	if err != nil {
		log.Infof("Failed to get user groups: %v", err)
	}
	log.Infof("User groups: %+v", userGroups)

	claims := jwt.MapClaims{
		"iss": ts.issuer,
		"sub": "alice",
		"aud": tokenReq.ClientID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	log.Infof("Pretty printed claims: %+v", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if ts.privSigningKey == nil {
		log.Errorf("Private signing key is not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "private signing key not configured"})
		return
	}
	tokenString, err := token.SignedString(ts.privSigningKey)
	if err != nil {
		log.Errorf("Failed to sign ID token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign ID token"})
		return
	}
	log.Infof("Parsed token request: %+v", tokenReq)
	resp := TokenResponse{
		AccessToken:  "abc123",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		IDToken:      tokenString,
		RefreshToken: "xyz789",
		Scope:        "openid profile email",
	}
	if err := WriteTokenResponse(c.Writer, resp); err != nil {
		http.Error(c.Writer, "Failed to write response", http.StatusInternalServerError)
	}
}

// UserInfoHandler
func (ts *TokenServiceController) UserInfoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user info"})
}

// IntrospectHandler
func (ts *TokenServiceController) IntrospectHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "token introspected"})
}

// RevokeHandler
func (ts *TokenServiceController) RevokeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "token revoked"})
}

func (ts *TokenServiceController) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
