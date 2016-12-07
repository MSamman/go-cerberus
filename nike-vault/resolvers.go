package vault

import (
	//	"flag"
	"os"
)

const (
	VAULT_TOKEN_ENV_PROPERTY  = "VAULT_TOKEN"
	VAULT_TOKEN_FLAG_PROPERTY = "vault_token"

	VAULT_ADDR_ENV_PROPERTY  = "VAULT_ADDR"
	VAULT_ADDR_FLAG_PROPERTY = "vault_addr"
)

type StringResolver func() (string, bool)

var DefaultTokenResolver = []StringResolver{
	EnvironmentResolver(VAULT_TOKEN_ENV_PROPERTY),
	FlagResolver(VAULT_TOKEN_FLAG_PROPERTY),
}

var DefaultUrlResolver = []StringResolver{
	EnvironmentResolver(VAULT_ADDR_ENV_PROPERTY),
	FlagResolver(VAULT_ADDR_FLAG_PROPERTY),
}

func EnvironmentResolver(envvarName string) StringResolver {
	return func() (string, bool) {
		return os.LookupEnv(envvarName)
	}
}

func FlagResolver(flagName string) StringResolver {
	return func() (string, bool) {
		return flagName, true
	}
}

func TokenResolver(token string) StringResolver {
	return func() (string, bool) {
		return token, true
	}
}
