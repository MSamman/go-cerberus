package cerberus

import (
	"errors"
)

var NoAwsAccountId = errors.New("Unable to obtain AWS account ID from instance profile ARN.")
var NoIAMAuthData = errors.New("Success response from IAM role authenticate endpoint missing auth data!")
