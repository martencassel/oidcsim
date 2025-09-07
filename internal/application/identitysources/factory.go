package identitysources

import (
	"fmt"

	"github.com/martencassel/oidcsim/internal/domain/configuration"
	domIDS "github.com/martencassel/oidcsim/internal/domain/identitysources"
	identitysourcesdomain "github.com/martencassel/oidcsim/internal/domain/identitysources"
	infraIDS "github.com/martencassel/oidcsim/internal/infrastructure/identitysources"
)

func BuildProvider(cfg configuration.IdentitySourceConfig) (identitysourcesdomain.IdentityProvider, error) {
	switch cfg.Type {
	case string(domIDS.SourceLocal):
		return infraIDS.NewLocalProvider(cfg.Settings), nil
	case string(domIDS.SourceLDAP):
		return infraIDS.NewLDAPProvider(cfg.Settings), nil
	case string(domIDS.SourceOIDC):
		return infraIDS.NewOIDCProvider(cfg.Settings), nil
	default:
		return nil, fmt.Errorf("unsupported identity source type: %s", cfg.Type)
	}
}
