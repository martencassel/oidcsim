package identitysources

import (
	"context"

	identitysourcesdomain "github.com/martencassel/oidcsim/internal/domain/identitysources"
)

// External OIDC IdP token validation based identity source

type oidcProviderImpl struct{}

func (p *oidcProviderImpl) Type() identitysourcesdomain.IdentitySourceType {
	return identitysourcesdomain.SourceOIDC
}

func NewOIDCProvider(settings map[string]interface{}) identitysourcesdomain.IdentityProvider {
	return &oidcProviderImpl{}
}

func (o *oidcProviderImpl) AuthenticatePassword(ctx context.Context, username, password string) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (o *oidcProviderImpl) AuthenticateExternal(ctx context.Context, assertion interface{}) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (o *oidcProviderImpl) GetClaims(ctx context.Context, subjectID identitysourcesdomain.SubjectID) (map[string]interface{}, error) {
	return nil, nil
}
