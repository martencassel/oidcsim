package identitysources

import domIDS "github.com/martencassel/oidcsim/internal/domain/identitysources"

type Registry struct {
	providers map[domIDS.IdentitySourceType]domIDS.IdentityProvider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[domIDS.IdentitySourceType]domIDS.IdentityProvider)}
}

func (r *Registry) Register(t domIDS.IdentitySourceType, p domIDS.IdentityProvider) {
	r.providers[t] = p
}

func (r *Registry) Get(t domIDS.IdentitySourceType) (domIDS.IdentityProvider, bool) {
	p, ok := r.providers[t]
	return p, ok
}
