package identitysources

import (
	"context"

	identitysourcesdomain "github.com/martencassel/oidcsim/internal/domain/identitysources"
)

// SAML assertion based identity source
type samlProviderImpl struct{}

func (p *samlProviderImpl) Type() identitysourcesdomain.IdentitySourceType {
	return identitysourcesdomain.SourceSAML
}

func NewSAMLProvider(settings map[string]interface{}) identitysourcesdomain.IdentityProvider {
	return &samlProviderImpl{}
}

func (s *samlProviderImpl) AuthenticatePassword(ctx context.Context, username, password string) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (s *samlProviderImpl) AuthenticateExternal(ctx context.Context, assertion interface{}) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (s *samlProviderImpl) GetClaims(ctx context.Context, subjectID identitysourcesdomain.SubjectID) (map[string]interface{}, error) {
	return nil, nil
}
