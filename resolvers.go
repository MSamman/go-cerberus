package cerberus

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	vault "github.com/msamman/go-cerberus/nike-vault"
	"net/http"
)

const (
	IAM_ARN_PATTERN = "(arn\\:aws\\:iam\\:\\:)(?P<accountId>[0-9].*)(\\:.*/)(?P<iamRole>[a-zA-Z0-9_].*)"

	PADDING_TIME = 60

	CERBERUS_AUTH_PATH = "v1/auth/iam-role"
)

type InstanceMetaData struct {
	IntanceAccountID string `json:"account_id"`
	InstanceIAMRole  string `json:"role_name"`
	InstanceRegion   string `json:"region"`
}

func InstanceRoleTokenResolver(urlResolvers *[]vault.StringResolver) func() (string, bool) {
	errFunc := func() (string, bool) {
		return "", false
	}

	awsSession, err := session.NewSession()
	if err != nil {
		return errFunc
	}

	instanceMeta, err := lookupEC2Info(awsSession)
	if err != nil {
		return errFunc
	}

	if urlResolvers == nil {
		urlResolvers = &DefaultUrlResolver
	}

	data, err := getEncrypthedAuthData(instanceMeta, urlResolvers)
	if err != nil {
		return errFunc
	}

	vaultAuth, err := decryptToken(awsSession, instanceMeta, data)
	if err != nil {
		return errFunc
	}

	return func() (string, bool) {
		return vaultAuth.ClientToken, true
	}
}

func lookupEC2Info(s *session.Session) (InstanceMetaData, error) {
	ec2Meta := ec2metadata.New(s)

	ec2IAMInfo, err := ec2Meta.IAMInfo()
	if err != nil {
		return InstanceMetaData{}, err
	}

	accountIdRegex, err := regexp.Compile(IAM_ARN_PATTERN)
	if err != nil {
		return InstanceMetaData{}, err
	}

	matches := accountIdRegex.FindStringSubmatch(ec2IAMInfo.InstanceProfileArn)

	if len(matches) == 0 {
		return InstanceMetaData{}, NoAwsAccountId
	}

	region, err := ec2Meta.Region()
	if err != nil {
		return InstanceMetaData{}, err
	}

	return InstanceMetaData{matches[2], matches[4], region}, nil
}

func getEncrypthedAuthData(meta InstanceMetaData, urlResolvers *[]vault.StringResolver) (string, error) {
	var url string
	var valid bool
	for _, resolver := range *urlResolvers {
		url, valid = resolver()
		if valid {
			break
		}
	}
	if !valid {
		return "", vault.NoUrlError
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(meta)

	req, err := http.NewRequest(http.MethodPost, url+CERBERUS_AUTH_PATH, &b)
	if err != nil {
		return "", err
	}

	h := http.Header{}
	h.Add(vault.ACCEPT_HEADER, vault.DEFAULT_MEDIA_TYPE)
	h.Add(vault.CONTENT_TYPE_HEADER, vault.DEFAULT_MEDIA_TYPE)

	req.Header = h

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// TODO parse error and return it
	}

	data := struct {
		AuthData string `json:"auth_data"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	if data.AuthData == "" {
		return "", NoIAMAuthData
	}

	return data.AuthData, nil
}

func decryptToken(s *session.Session, meta InstanceMetaData, encryptedData string) (vault.VaultAuthResponse, error) {
	kmsService := kms.New(s, aws.NewConfig().WithRegion(meta.InstanceRegion))

	decryptInput := &kms.DecryptInput{}
	decryptInput.SetCiphertextBlob([]byte(encryptedData))

	decryptRes, err := kmsService.Decrypt(decryptInput)
	if err != nil {
		return vault.VaultAuthResponse{}, err
	}

	buf := bytes.NewBuffer([]byte(decryptRes.GoString()))
	authRes := vault.VaultAuthResponse{}
	err = json.NewDecoder(buf).Decode(&authRes)
	if err != nil {
		return vault.VaultAuthResponse{}, err
	}

	return authRes, nil
}
