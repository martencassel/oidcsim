package handlers

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/martencassel/oidcsim/internal/dto"
// 	"github.com/martencassel/oidcsim/internal/errors"
// 	"github.com/martencassel/oidcsim/internal/services"
// )

// type TokenHandler struct {
// 	TokenService services.TokenService
// }

// func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	req, err := dto.BuildTokenRequest(r)
// 	if err != nil {
// 		writeOAuthError(w, errors.ErrInvalidRequest.WithDescription(err.Error()))
// 		return
// 	}
// 	resp, err := h.TokenService.IssueToken(r.Context(), req)
// 	if err != nil {
// 		writeOAuthError(w, err)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_ = json.NewEncoder(w).Encode(resp) // resp is a dto.TokenResponse
// }
