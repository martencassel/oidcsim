package session

import (
	"net/http"

	"github.com/martencassel/oidcsim/internal/interface/http/dto"
)

type SessionManager interface {
	GetID(sid string) string
	Current(r *http.Request) (sessionID string, ok bool)
	Ensure(w http.ResponseWriter, r *http.Request) (sessionID string, err error)
	Rotate(w http.ResponseWriter, r *http.Request) (newID string, err error)
	Destroy(w http.ResponseWriter, r *http.Request) error

	SaveAuthorizeRequest(sid string, req dto.AuthorizeRequest) error
	GetAuthorizeRequest(sid string) (dto.AuthorizeRequest, bool, error)
	ClearAuthorizeRequest(sid string) error
}
