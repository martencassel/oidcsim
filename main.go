package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	oauth2app "github.com/martencassel/oidcsim/internal/application/oauth2"
	"github.com/martencassel/oidcsim/internal/application/oauth2/flows"
	"github.com/martencassel/oidcsim/internal/config"
	oauth2 "github.com/martencassel/oidcsim/internal/interface/http"
	jwksutil "github.com/martencassel/oidcsim/jwskutil"
)

// type ServerConfig struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	BasePath string `yaml:"base_path"`
// }

// type SigningConfig struct {
// 	PrivateKeyFile string `yaml:"privateKeyFile"`
// 	KeyID          string `yaml:"keyID"`
// }

// type OIDCConfig struct {
// 	Issuer  string        `yaml:"issuer"`
// 	Signing SigningConfig `yaml:"signing"`
// }

// // Default routes for the OIDC provider
// func (c *AppConfig) DefaultRoutes() {
// 	c.Routes.Discovery = "/.well-known/openid-configuration"
// 	c.Routes.JWKS = "/.well-known/jwks.json"
// 	c.Routes.Authorize = "/authorize"
// }

// type RoutesConfig struct {
// 	Discovery  string `yaml:"discovery"`
// 	JWKS       string `yaml:"jwks"`
// 	Authorize  string `yaml:"authorize"`
// 	Token      string `yaml:"token"`
// 	Userinfo   string `yaml:"userinfo"`
// 	Introspect string `yaml:"introspect"`
// 	Revoke     string `yaml:"revoke"`
// 	Logout     string `yaml:"logout"`
// }

// type AppConfig struct {
// 	Server ServerConfig `yaml:"server"`
// 	OIDC   OIDCConfig   `yaml:"oidc"`
// 	Routes RoutesConfig `yaml:"routes"`
// }

// func loadConfig(file string) (config *AppConfig, err error) {
// 	yamlFile, err := os.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config = &AppConfig{}
// 	if err := yaml.Unmarshal(yamlFile, config); err != nil {
// 		return nil, err
// 	}
// 	if config.OIDC.Issuer == "" {
// 		return nil, fmt.Errorf("OIDC issuer is required")
// 	}
// 	if config.Server.Port < 1 || config.Server.Port > 65535 {
// 		return nil, fmt.Errorf("invalid server port: %d", config.Server.Port)
// 	}
// 	return config, nil
// }

func parseRSAPrivateKeyFromPEM(path string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		k, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err2)
		}
		key = k.(*rsa.PrivateKey)
	}
	return key, nil
}

// type AuthorizationRequest struct {
// 	Scope               string `json:"scope"`
// 	State               string `json:"state"`
// 	Nonce               string `json:"nonce"`
// 	CodeChallenge       string `json:"code_challenge"`
// 	CodeChallengeMethod string `json:"code_challenge_method"`
// 	ResponseType        string `json:"response_type"`
// 	ClientID            string `json:"client_id"`
// 	RedirectURI         string `json:"redirect_uri"`
// }

// type AuthorizationResponse struct {
// 	Issuer      string
// 	Code        string
// 	State       string
// 	RedirectURI string
// }

// // RedirectToClient sends a 302 redirect to the client with code and state
// func (r *AuthorizationResponse) RedirectToClient(w http.ResponseWriter) {
// 	redirectURL, err := url.Parse(r.RedirectURI)
// 	if err != nil {
// 		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
// 		return
// 	}
// 	// Add query parameters
// 	q := redirectURL.Query()
// 	q.Set("code", r.Code)
// 	q.Set("state", r.State)
// 	q.Set("iss", r.Issuer)
// 	redirectURL.RawQuery = q.Encode()
// 	http.Redirect(w, &http.Request{}, redirectURL.String(), http.StatusFound)
// }

// // AuthorizeHandler
// func (ts *TokenServiceController) AuthorizeHandler(c *gin.Context) {
// 	// Bind it
// 	authReq := AuthorizationRequest{
// 		Scope:               c.Query("scope"),
// 		State:               c.Query("state"),
// 		Nonce:               c.Query("nonce"),
// 		CodeChallenge:       c.Query("code_challenge"),
// 		CodeChallengeMethod: c.Query("code_challenge_method"),
// 		ResponseType:        c.Query("response_type"),
// 		ClientID:            c.Query("client_id"),
// 		RedirectURI:         c.Query("redirect_uri"),
// 	}
// 	log.Infof("Authorization request: %+v", authReq)

// 	issuedCode, err := ts.codeStore.Generate(authReq.ClientID, authReq.RedirectURI)
// 	if err != nil {
// 		log.Errorf("Failed to generate authorization code: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate authorization code"})
// 		return
// 	}
// 	response := AuthorizationResponse{
// 		Issuer:      ts.issuer,
// 		Code:        issuedCode,
// 		State:       authReq.State,
// 		RedirectURI: authReq.RedirectURI,
// 	}
// 	response.RedirectToClient(c.Writer)
// }

// // TokenRequest represents the body of a POST /token request
// type TokenRequest struct {
// 	GrantType    string `json:"grant_type"`
// 	Code         string `json:"code"`
// 	RedirectURI  string `json:"redirect_uri"`
// 	ClientID     string `json:"client_id,omitempty"`     // optional if using Basic Auth
// 	ClientSecret string `json:"client_secret,omitempty"` // optional if using Basic Auth
// }

// // ParseTokenRequest parses a POST /token request body into a TokenRequest struct
// func ParseTokenRequest(r *http.Request) (*TokenRequest, error) {
// 	if err := r.ParseForm(); err != nil {
// 		return nil, err
// 	}

// 	clientID := r.FormValue("client_id")
// 	clientSecret := r.FormValue("client_secret")

// 	// Check Authorization header for Basic Auth
// 	auth := r.Header.Get("Authorization")
// 	if auth != "" && len(auth) > 6 && auth[:6] == "Basic " {
// 		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
// 		if err == nil {
// 			parts := strings.SplitN(string(decoded), ":", 2)
// 			if len(parts) == 2 {
// 				clientID = parts[0]
// 				clientSecret = parts[1]
// 			}
// 		}
// 	}

// 	return &TokenRequest{
// 		GrantType:    r.FormValue("grant_type"),
// 		Code:         r.FormValue("code"),
// 		RedirectURI:  r.FormValue("redirect_uri"),
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 	}, nil
// }

// // TokenResponse represents the JSON response from the token endpoint
// type TokenResponse struct {
// 	AccessToken  string `json:"access_token"`
// 	TokenType    string `json:"token_type"`              // usually "Bearer"
// 	ExpiresIn    int    `json:"expires_in"`              // in seconds
// 	IDToken      string `json:"id_token,omitempty"`      // OIDC-specific
// 	RefreshToken string `json:"refresh_token,omitempty"` // optional
// 	Scope        string `json:"scope,omitempty"`         // optional
// }

// // WriteTokenResponse writes the token response as JSON to the http.ResponseWriter
// func WriteTokenResponse(w http.ResponseWriter, resp TokenResponse) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	return json.NewEncoder(w).Encode(resp)
// }

// // TokenHandler
// func (ts *TokenServiceController) TokenHandler(c *gin.Context) {
// 	// Parse the token request
// 	tokenReq, err := ParseTokenRequest(c.Request)
// 	if err != nil {
// 		log.Errorf("Failed to parse token request: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token request"})
// 		return
// 	}

// 	userGroups, err := ts.idStore.GetUserGroups(c.Request.Context(), "alice")
// 	if err != nil {
// 		log.Infof("Failed to get user groups: %v", err)
// 	}
// 	log.Infof("User groups: %+v", userGroups)

// 	claims := jwt.MapClaims{
// 		"iss": ts.issuer,
// 		"sub": "alice",
// 		"aud": tokenReq.ClientID,
// 		"exp": time.Now().Add(time.Hour).Unix(),
// 		"iat": time.Now().Unix(),
// 	}
// 	log.Infof("Pretty printed claims: %+v", claims)
// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	if ts.privSigningKey == nil {
// 		log.Errorf("Private signing key is not configured")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "private signing key not configured"})
// 		return
// 	}
// 	tokenString, err := token.SignedString(ts.privSigningKey)
// 	if err != nil {
// 		log.Errorf("Failed to sign ID token: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign ID token"})
// 		return
// 	}
// 	log.Infof("Parsed token request: %+v", tokenReq)
// 	resp := TokenResponse{
// 		AccessToken:  "abc123",
// 		TokenType:    "Bearer",
// 		ExpiresIn:    3600,
// 		IDToken:      tokenString,
// 		RefreshToken: "xyz789",
// 		Scope:        "openid profile email",
// 	}
// 	if err := WriteTokenResponse(c.Writer, resp); err != nil {
// 		http.Error(c.Writer, "Failed to write response", http.StatusInternalServerError)
// 	}
// }

// // UserInfoHandler
// func (ts *TokenServiceController) UserInfoHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "user info"})
// }

// // IntrospectHandler
// func (ts *TokenServiceController) IntrospectHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "token introspected"})
// }

// // RevokeHandler
// func (ts *TokenServiceController) RevokeHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "token revoked"})
// }

// func (ts *TokenServiceController) LogoutHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
// }

// 1) Domain objects (what each one is and why it exists)
// Adapter / DTO
// adapter.BuildDTO(*http.Request) -> dto.ClientAuthDTO`
// Extract what the client **presented** (headers, form fields, TLS peer cert) and normalize to an immutable payload for the domain layer.
// No policy decisions here.

func mustLoadConfig(path string) *config.AppConfig {
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func main() {
	registry := oauth2app.NewFlowRegistry()
	codeRepo := infraoauth2.
		registry.Register("code", flows.NewCodeFlow(codeRepo, clientRepo, clock))
	authorizeSvc := oauth2app.NewAuthorizeService(delegationSvc, registry)
	mux := gin.Default()
	authHandler := &oauth2.Handler{}
	mux.GET("/authorize", authHandler.Authorize)
	if err := mux.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// 	cfg := mustLoadConfig("config.yaml")
// 	log.Infof("Loaded config: %+v", cfg)
// 	priv, pub := mustLoadKeys(cfg.OIDC.Signing.PrivateKeyFile)
// 	jwks := mustGenerateJWKS(pub)
// 	addr, port, tlsCert, tlsKey := parseFlags(cfg)
// 	app := bootstrap.BuildApp(cfg, jwks, priv)
// 	server := NewServer(addr, port, tlsCert, tlsKey, app.Router)
// 	if err := server.ListenAndServe(); err != nil {
// 		log.Fatalf("server failed: %v", err)
// 	}
// }

// 	cfg := mustLoadConfig("config.yaml")
// 	priv, pub := mustLoadKeys(cfg.OIDC.Signing.PrivateKeyFile)
// 	jwks := mustGenerateJWKS(pub)
// 	log.SetFormatter(&log.TextFormatter{
// 		FullTimestamp: true,
// 	})
// 	log.SetOutput(os.Stdout)
// 	log.SetLevel(log.DebugLevel)
// 	fmt.Println(jwks)

// 	addr, port, tlsCert, tlsKey := parseFlags(cfg)
// 	// addr := flag.String("listen", "0.0.0.0", "address to listen on")
// 	// port := flag.Int("port", 8080, "port to listen on")
// 	// tlsCert := flag.String("tls-cert", "", "path to TLS certificate file (PEM)")
// 	// tlsKey := flag.String("tls-key", "", "path to TLS key file (PEM)")
// 	flag.Parse()

// 	r := setupRouter(cfg, jwks, priv)

// 	// r := gin.Default()
// 	// r.Use(RequestResponseLogger())
// 	// if cfg.OIDC.Issuer == "" {
// 	// 	cfg.OIDC.Issuer = "https://idp.local"
// 	// }

// 	// authReg := registry.New[clientauth.Authenticator]()
// 	// authReg.Register("client_secret_basic", &clientauth.ClientSecretBasic{})
// 	// authReg.Register("client_secret_post", &clientauth.ClientSecretPost{})
// 	// authReg.Register("private_key_jwt", &clientauth.ClientPrivateJWT{})
// 	// authReg.Register("none", &clientauth.ClientNone{})
// 	// //	clientStore := store.NewInMemoryClientStore()

// 	// grantHandlerReg := registry.New[granthandlers.GrantHandler]()
// 	// grantFlowReg := registry.New[grantflows.GrantFlow]()
// 	// clientStore := store.NewInMemoryClientStore()

// 	// deps := services.TokenServiceDeps{
// 	// 	AuthRegistry:      clientauth.BuildAuthRegistry(),
// 	// 	GrantValidatorReg: grantvalidators.BuildGrantValidatorRegistry(),
// 	// 	GrantHandlerReg:   grantHandlerReg,
// 	// 	GrantFlowReg:      grantFlowReg,
// 	// 	ClientStore:       clientStore,
// 	// }
// 	// tokenService := services.NewTokenServiceImpl(deps)
// 	// log.Infof("Token service: %+v", tokenService)

// 	// codeStore := authcode.NewStore(360 * time.Second)
// 	// privSigningKey := priv

// 	// // Identity Store API group
// 	// store := identity.NewCoreIdentityStore("http://localhost:8080/identity")
// 	// handler := identity.NewIdentityStoreHandler(store)
// 	// handler.SeedDefault()
// 	// handler.RegisterRoutes(r)

// 	// tokenServiceController := NewTokenServiceController(cfg.Routes, cfg.OIDC.Issuer, jwks, privSigningKey, codeStore, store)
// 	// tokenServiceController.RegisterRoutes(r)

// 	// // Example route
// 	// r.GET("/ping", func(c *gin.Context) {
// 	// 	c.JSON(http.StatusOK, gin.H{"message": "pong"})
// 	// })

// 	// // Catch-all route for unmatched paths
// 	// r.NoRoute(func(c *gin.Context) {
// 	// 	c.JSON(http.StatusNotFound, gin.H{
// 	// 		"error": "Route not found",
// 	// 		"path":  c.Request.URL.Path,
// 	// 	})
// 	// })

// 	// Determine address and port: command line flags take precedence, else use config.yaml
// 	if flag.Lookup("listen").Value.String() != "0.0.0.0" || flag.Lookup("port").Value.String() != "8080" {
// 		addr = fmt.Sprintf("%s:%d", addr, port)
// 	} else {
// 		addr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
// 	}

// 	if tlsCert != "" && tlsKey != "" {
// 		// Run with TLS
// 		if err := r.RunTLS(addr, tlsCert, tlsKey); err != nil {
// 			panic(fmt.Sprintf("failed to run server with TLS: %v", err))
// 		}
// 	} else {
// 		// Run without TLS
// 		if err := r.Run(addr); err != nil {
// 			panic(fmt.Sprintf("failed to run server: %v", err))
// 		}
// 	}
// }

type Server struct {
	Address string
	Port    int
	TLSCert string
	TLSKey  string
	Router  *gin.Engine
}

func NewServer(address string, port int, tlsCert, tlsKey string, router *gin.Engine) *Server {
	return &Server{
		Address: address,
		Port:    port,
		TLSCert: tlsCert,
		TLSKey:  tlsKey,
		Router:  router,
	}
}

func (s *Server) ListenAndServe() error {
	listenAddr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	if s.TLSCert != "" && s.TLSKey != "" {
		// Run with TLS
		return s.Router.RunTLS(listenAddr, s.TLSCert, s.TLSKey)
	} else {
		// Run without TLS
		return s.Router.Run(listenAddr)
	}
}

// func mustLoadConfig(file string) *AppConfig {
// 	cfg, err := loadConfig(file)
// 	if err != nil {
// 		panic(fmt.Sprintf("failed to load config: %v", err))
// 	}
// 	return cfg
// }

// func setupLogger() {
// 	log.SetFormatter(&log.TextFormatter{
// 		FullTimestamp: true,
// 	})parseRSAPrivateKeyFromPEM
// 	log.SetOutput(os.Stdout)
// 	log.SetLevel(log.DebugLevel)
// }

func mustLoadKeys(path string) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := parseRSAPrivateKeyFromPEM(path)
	if err != nil {
		panic(fmt.Sprintf("failed to parse private key: %v", err))
	}
	pub := &priv.PublicKey
	return priv, pub
}

func mustGenerateJWKS(pub *rsa.PublicKey) []byte {
	jwks, err := jwksutil.GenerateJWKS(pub, "idp-key")
	if err != nil {
		panic(fmt.Sprintf("failed to generate JWKS: %v", err))
	}
	return jwks
}

func parseFlags(cfg *config.AppConfig) (string, int, string, string) {
	addr := flag.String("listen", "0.0.0.0", "Address to listen on")
	port := flag.Int("port", 8080, "Port to listen on")
	tlsCert := flag.String("tls-cert", "", "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "Path to TLS key file")
	flag.Parse()
	return *addr, *port, *tlsCert, *tlsKey
}

// func setupRouter(cfg *AppConfig, jwks []byte, priv *rsa.PrivateKey) *gin.Engine {
// 	r := gin.Default()
// 	r.Use(RequestResponseLogger())
// 	if cfg.OIDC.Issuer == "" {
// 		cfg.OIDC.Issuer = "https://idp.local"
// 	}

// 	// Identity Store API group
// 	store := identity.NewCoreIdentityStore("http://localhost:8080/identity")
// 	handler := identity.NewIdentityStoreHandler(store)
// 	handler.SeedDefault()
// 	handler.RegisterRoutes(r)

// 	codeStore := authcode.NewStore(360 * time.Second)
// 	privSigningKey := priv

// 	tokenServiceController := NewTokenServiceController(cfg.Routes, cfg.OIDC.Issuer, jwks, privSigningKey, codeStore, store)
// 	tokenServiceController.RegisterRoutes(r)

// 	// Example route
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "pong"})
// 	})

// 	// Catch-all route for unmatched paths
// 	r.NoRoute(func(c *gin.Context) {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "Route not found",
// 			"path":  c.Request.URL.Path,
// 		})
// 	})
// 	return r
// }
