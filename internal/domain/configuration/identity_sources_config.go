package configuration

import "time"

// IdentitySourceConfig is the configuration for an identity source.
type IdentitySourceConfig struct {
	Name         string // Logical name, e.g. "local", "ldap", "oidc"
	Type         string // Type of source, e.g. "local", "ldap", "oidc"
	Enabled      bool
	Priority     int                    // For selection if multiple sources are enabled
	Settings     map[string]interface{} // Provider-specific settings (host, bindDN, etc)
	ClaimMapping []ClaimMappingConfig
	AuthPolicy   AuthPolicyConfig
}

// ClaimMappingConfig defines how source attributes map to internal claims.
type ClaimMappingConfig struct {
	SourceAttr  string // e.g., "mail"
	TargetClaim string // e.g., "email"
}

// AuthPolicyConfig defines authentication rules for this source.
type AuthPolicyConfig struct {
	MinACR        bool // Minimum assurance level
	RequireMFA    bool
	SessionMaxAge time.Duration // Force re-auth after this time
}

// IdentitySourcesConfig is the top-level config for all sources.
type IdentitySourcesConfig struct {
	Sources []IdentitySourceConfig
}
