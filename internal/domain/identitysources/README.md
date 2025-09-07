Application Layer: Orhcestration

In internal/application/identitysources/:


Registry of available providers (local, LDAP, OIDC, etc.)

Selection logic based on configuration (tenant, client, requested method)

Fallback / Aggregation if multiple sources cant satisfy a request.

