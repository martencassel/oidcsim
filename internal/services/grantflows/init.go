package grantflows

import (
	"github.com/martencassel/oidcsim/internal/registry"
	"github.com/martencassel/oidcsim/internal/services/granthandlers"
	"github.com/martencassel/oidcsim/internal/services/grantvalidators"
)

func NewRegistry() *registry.Registry[GrantFlow] {
	r := registry.New[GrantFlow]()
	r.Register("authorization_code", GrantFlow{
		Validator: &grantvalidators.AuthCodeValidator{},
		Handler:   &granthandlers.AuthCodeHandler{},
	})
	return r
}
