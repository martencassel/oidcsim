package identitysources

import (
	"context"

	identitysourcesdomain "github.com/martencassel/oidcsim/internal/domain/identitysources"
)

// LDAP bind + attributes based identity source
type ldapProviderImpl struct{}

func (p *ldapProviderImpl) Type() identitysourcesdomain.IdentitySourceType {
	return identitysourcesdomain.SourceLDAP
}

func NewLDAPProvider(settings map[string]interface{}) identitysourcesdomain.IdentityProvider {
	return &ldapProviderImpl{}
}

// Authenticate implements IdentitySource
func (l *ldapProviderImpl) Authenticate(username, password string) (string, *identitysourcesdomain.SubjectID, []identitysourcesdomain.Claim, error) {
	return "", nil, nil, nil
}

func (l *ldapProviderImpl) GetClaims(ctx context.Context, subjectID identitysourcesdomain.SubjectID) (map[string]interface{}, error) {
	return nil, nil
}

func (l *ldapProviderImpl) AuthenticatePassword(ctx context.Context, username, password string) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}

func (l *ldapProviderImpl) AuthenticateExternal(ctx context.Context, assertion interface{}) (identitysourcesdomain.SubjectID, error) {
	return "", nil
}
