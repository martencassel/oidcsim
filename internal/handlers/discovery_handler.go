package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DiscoveryResponse struct {
	Issuer                 string   `json:"issuer"`
	AuthURL                string   `json:"authorization_endpoint"`
	TokenURL               string   `json:"token_endpoint"`
	JWKSURL                string   `json:"jwks_uri"`
	ResponseTypesSupported []string `json:"response_types_supported"`
}

// DiscoveryHandler
func (ts *TokenServiceController) DiscoveryHandler(c *gin.Context) {
	issuer := ts.issuer
	c.JSON(http.StatusOK, DiscoveryResponse{
		Issuer:                 issuer,
		AuthURL:                issuer + ts.routesConfig.Authorize,
		TokenURL:               issuer + ts.routesConfig.Token,
		JWKSURL:                issuer + ts.routesConfig.JWKS,
		ResponseTypesSupported: []string{"code", "token", "code id_token"},
	})
}
