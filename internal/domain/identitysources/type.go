package identitysources

// IdentitySourceType is a label for the kind of source (local, ldap, oidc, saml, etc.).
type IdentitySourceType string

const (
	SourceLocal IdentitySourceType = "local"
	SourceLDAP  IdentitySourceType = "ldap"
	SourceOIDC  IdentitySourceType = "oidc"
	SourceSAML  IdentitySourceType = "saml"
)
