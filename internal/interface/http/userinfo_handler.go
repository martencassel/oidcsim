package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martencassel/oidcsim/internal/interface/http/dto"
)

func (h *Handler) UserInfo(c *gin.Context) {
	token := extractBearerToken(c.Request)
	resp, err := h.UserInfoAppService.HandleUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:            "invalid_token",
			ErrorDescription: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
