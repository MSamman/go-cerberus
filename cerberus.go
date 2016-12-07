package cerberus

import (
	vault "github.com/msamman/go-cerberus/nike-vault"
)

const (
	CERBERUS_TOKEN_ENV_PROPERTY  = "CERBERUS_TOKEN"
	CERBERUS_TOKEN_FLAG_PROPERTY = "cerberus_token"

	CERBERUS_ADDR_ENV_PROPERTY  = "CERBERUS_ADDR"
	CERBERUS_ADDR_FLAG_PROPERTY = "cerberus_addr"
)

var DefaultUrlResolver = []vault.StringResolver{
	vault.EnvironmentResolver(CERBERUS_ADDR_ENV_PROPERTY),
	vault.FlagResolver(CERBERUS_ADDR_FLAG_PROPERTY),
}

var DefaultTokenResolver = []vault.StringResolver{
	vault.EnvironmentResolver(CERBERUS_TOKEN_ENV_PROPERTY),
	vault.FlagResolver(CERBERUS_TOKEN_FLAG_PROPERTY),
	InstanceRoleTokenResolver(&DefaultUrlResolver),
}

func NewClient(tokenResolvers *[]vault.StringResolver, urlResolvers *[]vault.StringResolver) (*vault.VaultClient, error) {
	// Check for nil providers and use default providers instead
	if tokenResolvers == nil {
		tokenResolvers = &DefaultTokenResolver
	}

	if urlResolvers == nil {
		urlResolvers = &DefaultUrlResolver
	}

	return vault.NewVaultClient(tokenResolvers, urlResolvers)
}
