package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/martencassel/oidcsim/internal/errors"
)

func writeOAuthError(w http.ResponseWriter, err error) {
	// Default to internal server error
	status := http.StatusInternalServerError
	code := errors.ErrServerError
	desc := "internal server error"

	if ae, ok := err.(errors.AuthError); ok {
		code = ae
		desc = ae.Error() // or a friendlier description if you have one
		status = http.StatusBadRequest

		if ae == errors.ErrInvalidClient {
			status = http.StatusUnauthorized
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":             code.Error(),
		"error_description": desc,
	})
}

func writeTokenError(w http.ResponseWriter, status int, err errors.AuthError, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":             err.Error(),
		"error_description": description,
	})
}

func writeAuthorizeError(w http.ResponseWriter, redirectURI string, state string, err errors.AuthError, description string) {
	// If we have a redirect URI, redirect with error parameters
	if redirectURI != "" {
		http.Redirect(w, nil, redirectURI+"?error="+err.Error()+"&error_description="+description+"&state="+state, http.StatusFound)
		return
	}
	// Otherwise, just write the error directly
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":             err.Error(),
		"error_description": description,
	})
}
