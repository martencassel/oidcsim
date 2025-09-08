package bootstrap

// import (
// 	"bytes"
// 	"crypto/rsa"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/fatih/color"
// 	"github.com/gin-gonic/gin"
// 	"github.com/martencassel/oidcsim/authcode"
// 	"github.com/martencassel/oidcsim/internal/clientauth"
// 	"github.com/martencassel/oidcsim/internal/config"
// 	"github.com/martencassel/oidcsim/internal/handlers"
// 	"github.com/martencassel/oidcsim/internal/identity"
// 	"github.com/martencassel/oidcsim/internal/registry"
// 	"github.com/martencassel/oidcsim/internal/services"
// 	"github.com/martencassel/oidcsim/internal/services/grantflows"
// 	"github.com/martencassel/oidcsim/internal/services/granthandlers"
// 	"github.com/martencassel/oidcsim/internal/services/grantvalidators"
// 	"github.com/martencassel/oidcsim/internal/store"
// 	log "github.com/sirupsen/logrus"
// )

// type bodyLogWriter struct {
// 	gin.ResponseWriter
// 	body *bytes.Buffer
// }

// func (w *bodyLogWriter) Write(b []byte) (int, error) {
// 	w.body.Write(b)
// 	return w.ResponseWriter.Write(b)
// }

// type App struct {
// 	Router *gin.Engine
// }

// func BuildApp(cfg *config.AppConfig, jwks []byte, priv *rsa.PrivateKey) *App {
// 	// Initialize Gin router
// 	router := gin.Default()
// 	// Setup routes, middleware, handlers, etc. using cfg, jwks, and priv
// 	router.Use(RequestResponseLogger())
// 	if cfg.OIDC.Issuer == "" {
// 		cfg.OIDC.Issuer = "https://idp.local"
// 	}
// 	authReg := registry.New[clientauth.Authenticator]()
// 	authReg.Register("client_secret_basic", &clientauth.ClientSecretBasic{})
// 	authReg.Register("client_secret_post", &clientauth.ClientSecretPost{})
// 	authReg.Register("private_key_jwt", &clientauth.ClientPrivateJWT{})
// 	authReg.Register("none", &clientauth.ClientNone{})
// 	clientStore := store.NewInMemoryClientStore()
// 	grantHandlerReg := registry.New[granthandlers.GrantHandler]()
// 	grantFlowReg := registry.New[grantflows.GrantFlow]()
// 	deps := services.TokenServiceDeps{
// 		AuthRegistry:      clientauth.BuildAuthRegistry(),
// 		GrantValidatorReg: grantvalidators.BuildGrantValidatorRegistry(),
// 		GrantHandlerReg:   grantHandlerReg,
// 		GrantFlowReg:      grantFlowReg,
// 		ClientStore:       clientStore,
// 	}
// 	tokenService := services.NewTokenServiceImpl(deps)
// 	log.Infof("Token service: %+v", tokenService)
// 	codeStore := authcode.NewStore(360 * time.Second)
// 	privSigningKey := priv
// 	// Identity Store API group
// 	store := identity.NewCoreIdentityStore("http://localhost:8080/identity")
// 	handler := identity.NewIdentityStoreHandler(store)
// 	handler.SeedDefault()
// 	handler.RegisterRoutes(router)

// 	routesConfig := &handlers.RoutesConfig{
// 		Discovery:  "/.well-known/openid-configuration",
// 		JWKS:       "/.well-known/jwks.json",
// 		Authorize:  "/authorize",
// 		Token:      "/token",
// 		Userinfo:   "/userinfo",
// 		Introspect: "/introspect",
// 		Revoke:     "/revoke",
// 		Logout:     "/logout",
// 	}
// 	// Token Service API group
// 	controller := handlers.NewTokenServiceControllerBuilder().
// 		WithIssuer(cfg.OIDC.Issuer).
// 		WithRoutesConfig(routesConfig).
// 		WithJWKS(jwks).
// 		WithCodeStore(codeStore).
// 		WithSigningKey(privSigningKey).
// 		WithIdentityStore(store).
// 		Build()
// 	controller.RegisterRoutes(router)
// 	return &App{Router: router}
// }

// // Middleware to log requests and responses with colors and pretty-printed structured info
// func RequestResponseLogger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Log request line and query params
// 		color.Set(color.FgCyan)
// 		fmt.Printf("\n→ %s %s %s\n", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
// 		color.Unset()
// 		if len(c.Request.URL.Query()) > 0 {
// 			color.Set(color.FgHiCyan)
// 			fmt.Println("Query Params:")
// 			for k, v := range c.Request.URL.Query() {
// 				fmt.Printf("  %s: %v\n", k, v)
// 			}
// 			color.Unset()
// 		}
// 		for k, v := range c.Request.Header {
// 			color.Set(color.FgBlue)
// 			fmt.Printf("%s: %v\n", k, v)
// 			color.Unset()
// 		}
// 		// Capture response
// 		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
// 		c.Writer = writer
// 		c.Next()
// 		// Log response
// 		color.Set(color.FgGreen)
// 		fmt.Printf("← %d %s\n", c.Writer.Status(), http.StatusText(c.Writer.Status()))
// 		color.Unset()
// 		for k, v := range c.Writer.Header() {
// 			color.Set(color.FgMagenta)
// 			fmt.Printf("%s: %v\n", k, v)
// 			color.Unset()
// 		}
// 		if ct := c.Writer.Header().Get("Content-Type"); ct != "" {
// 			color.Set(color.FgYellow)
// 			fmt.Printf("Content-Type: [%s]\n", ct)
// 			color.Unset()
// 		}
// 		if c.Writer.Header().Get("Content-Type") != "" &&
// 			(c.Writer.Header().Get("Content-Type") == "application/json" ||
// 				c.Writer.Header().Get("Content-Type") == "application/json; charset=utf-8") {
// 			var prettyJSON bytes.Buffer
// 			if err := json.Indent(&prettyJSON, writer.body.Bytes(), "", "  "); err == nil {
// 				color.Set(color.FgYellow)
// 				fmt.Println(prettyJSON.String())
// 				color.Unset()
// 			} else {
// 				fmt.Println(writer.body.String())
// 			}
// 		} else {
// 			fmt.Println(writer.body.String())
// 		}
// 	}
// }
