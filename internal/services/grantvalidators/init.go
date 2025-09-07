package grantvalidators

import "github.com/martencassel/oidcsim/internal/registry"

func NewRegistry() *registry.Registry[GrantValidator] {
	r := registry.New[GrantValidator]()
	grantValidators := registry.New[GrantValidator]()
	grantValidators.Register((&AuthCodeValidator{}).GrantType(), &AuthCodeValidator{})
	grantValidators.Register((&RefreshTokenValidator{}).GrantType(), &RefreshTokenValidator{})
	return r
}

func BuildGrantValidatorRegistry() *registry.Registry[GrantValidator] {
	return NewRegistry()
}
