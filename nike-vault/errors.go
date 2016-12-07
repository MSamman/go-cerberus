package vault

import (
	"errors"
)

var NoTokenError = errors.New("No vault token found using the provided resolvers.")
var NoUrlError = errors.New("No vault URL found using the provided resolvers.")
