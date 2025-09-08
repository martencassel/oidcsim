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
	delegationapp "github.com/martencassel/oidcsim/internal/application/delegation"
	"github.com/martencassel/oidcsim/internal/config"
	"github.com/martencassel/oidcsim/internal/domain/delegation"
	infradelegation "github.com/martencassel/oidcsim/internal/infrastructure/delegation"
	"github.com/martencassel/oidcsim/internal/infrastructure/session"
	infrasession "github.com/martencassel/oidcsim/internal/infrastructure/session"
	oauth2 "github.com/martencassel/oidcsim/internal/interface/http"
	jwksutil "github.com/martencassel/oidcsim/jwskutil"
)

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

func mustLoadConfig(path string) *config.AppConfig {
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func TestDelegation(repo delegationapp.Repository) {
	repo.Save(nil, delegation.Delegation{
		UserID:   "alice",
		ClientID: "client1",
		Scopes:   []string{"openid", "profile", "email"},
	})
	repo.FindByID(nil, "alice|client1")
	repo.FindByUserAndClient(nil, "alice", "client1")
	repo.Delete(nil, "alice", "client1")
	repo.FindByUserAndClient(nil, "alice", "client1")
}

func main() {
	delegRepo := infradelegation.NewMemoryRepo()
	//cliRepo := infradelegation.NewMemoryRepo()
	mux := gin.Default()
	authHandler := &oauth2.Handler{}
	authHandler.DelegationSvc = delegationapp.NewDelegationService(delegRepo)
	sessionManager := infrasession.NewMemorySessionManager("oidcsim_session", session.WithAllowInsecure())
	authHandler.Sessions = sessionManager

	mux.GET("/authorize", authHandler.Authorize)
	// if err := mux.Run(":8080"); err != nil {
	// 	log.Fatalf("failed to run server: %v", err)
	// }
}

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
