package dto

import "github.com/gin-gonic/gin"

// For /introspect endpoint (RFC 7662).

type IntrospectionRequest struct {
	Token string `form:"token" binding:"required"`
}

type IntrospectionResponse struct {
	Active   bool   `json:"active"`
	Scope    string `json:"scope,omitempty"`
	ClientID string `json:"client_id,omitempty"`
	Username string `json:"username,omitempty"`
	Exp      int64  `json:"exp,omitempty"`
	Sub      string `json:"sub,omitempty"`
	Aud      string `json:"aud,omitempty"`
	Iss      string `json:"iss,omitempty"`
	Iat      int64  `json:"iat,omitempty"`
}

func (ir *IntrospectionRequest) Bind(c *gin.Context) error {
	return c.ShouldBind(ir)
}

type IntrospectionErrorResponse struct {
	Active           bool   `json:"active"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}
