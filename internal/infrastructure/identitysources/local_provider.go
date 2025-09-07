package identitysources

import (
	"context"

	identitysourcesdomain "github.com/martencassel/oidcsim/internal/domain/identitysources"
)

// SQL / NoSQL based identity source with bcrypt
type localProviderImpl struct{}

func (p *localProviderImpl) Type() identitysourcesdomain.IdentitySourceType {
	return identitysourcesdomain.SourceLocal
}

func NewLocalProvider(settings map[string]interface{}) identitysourcesdomain.IdentityProvider {
	return &localProviderImpl{}
}

func (l *localProviderImpl) AuthenticatePassword(ctx context.Context, username, password string) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (l *localProviderImpl) AuthenticateExternal(ctx context.Context, assertion interface{}) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (l *localProviderImpl) GetClaims(ctx context.Context, subjectID identitysourcesdomain.SubjectID) (map[string]interface{}, error) {
	return nil, nil
}
