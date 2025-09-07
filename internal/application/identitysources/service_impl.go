package identitysources

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/configuration"
	domIDS "github.com/martencassel/oidcsim/internal/domain/identitysources"
)

type Service struct {
	registry *Registry
}

func NewService(cfg configuration.IdentitySourcesConfig) (*Service, error) {
	reg := NewRegistry()
	for _, srcCfg := range cfg.Sources {
		if !srcCfg.Enabled {
			continue
		}
		provider, err := BuildProvider(srcCfg)
		if err != nil {
			return nil, err
		}
		reg.Register(domIDS.IdentitySourceType(srcCfg.Type), provider)
	}
	return &Service{registry: reg}, nil
}

func (s *Service) Authenticate(ctx context.Context, srcType domIDS.IdentitySourceType, creds map[string]string) (domIDS.SubjectID, error) {
	provider, ok := s.registry.Get(srcType)
	if !ok {
		return "", domIDS.ErrUnknownSource
	}
	return provider.AuthenticatePassword(ctx, creds["username"], creds["password"])
}
