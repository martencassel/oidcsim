package handlers

import (
	"net/http"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/services"
)

type AuthorizeHandler struct {
	AuthorizeService services.AuthorizeService
}

// ServeHTTP handles the /authorize endpoint.
func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := dto.BuildAuthorizeRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirectURL, err := h.AuthorizeService.HandleAuthorize(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
