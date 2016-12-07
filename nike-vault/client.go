package vault

import (
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	SECRET_PATH_PREFIX = "v1/secret/"
	AUTH_PATH_PREFIX   = "v1/auth/"
	SYS_PATH_PREFIX    = "v1/sys/"
	TOKEN_LOOKUP_PATH  = "token/lookup-self"

	DEFAULT_MEDIA_TYPE  = "application/json; charset=utf-8"
	VAULT_TOKEN_HEADER  = "X-Vault-Token"
	ACCEPT_HEADER       = "Accept"
	CONTENT_TYPE_HEADER = "Content-Type"

	DEFAULT_TIMEOUT = 15
)

type VaultClient struct {
	vaultUrl   string
	vaultToken string

	httpClient http.Client
}

func NewVaultClient(tokenResolvers *[]StringResolver, urlResolvers *[]StringResolver) (*VaultClient, error) {
	// Check for nil providers and use default providers instead
	if tokenResolvers == nil {
		tokenResolvers = &DefaultTokenResolver
	}

	if urlResolvers == nil {
		urlResolvers = &DefaultUrlResolver
	}

	// Get Vault Token using token providers
	var token, url string
	var valid bool
	for _, resolver := range *tokenResolvers {
		token, valid = resolver()
		if valid {
			break
		}
	}

	if !valid {
		return nil, NoTokenError
	}

	// Get Vault URL using url providers
	for _, resolver := range *urlResolvers {
		url, valid = resolver()
		if valid {
			break
		}
	}
	if !valid {
		return nil, NoUrlError
	}

	c := http.Client{
		Timeout: time.Second * DEFAULT_TIMEOUT,
	}

	return &VaultClient{
		vaultToken: token,
		vaultUrl:   url,
		httpClient: c,
	}, nil
}

func (v *VaultClient) List(path string) (VaultListResponse, error) {
	url, err := v.buildSecretUrl(path, map[string]string{"list": "true"})
	if err != nil {
		return VaultListResponse{}, err
	}

	res, err := v.executeRequest(http.MethodGet, url, nil)
	if err != nil {
		return VaultListResponse{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return VaultListResponse{}, nil
	} else if res.StatusCode != http.StatusOK {
		// TODO parse error and return it
	}

	vlResponse := VaultListResponse{}
	err = json.NewDecoder(res.Body).Decode(&vlResponse)
	if err != nil {
		return VaultListResponse{}, err
	}

	return vlResponse, nil
}

func (v *VaultClient) Read(path string) (VaultResponse, error) {
	url, err := v.buildSecretUrl(path, nil)
	if err != nil {
		return VaultResponse{}, err
	}

	res, err := v.executeRequest(http.MethodGet, url, nil)
	if err != nil {
		return VaultResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO parse error and return it
	}

	vrResponse := VaultResponse{}
	err = json.NewDecoder(res.Body).Decode(&vrResponse)
	if err != nil {
		return VaultResponse{}, err
	}

	return vrResponse, nil
}

func (v *VaultClient) Write(path string, data map[string]string) error {
	url, err := v.buildSecretUrl(path, nil)
	if err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := v.executeRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		// TODO parse error and return it
	}

	return nil
}

func (v *VaultClient) Delete(path string) error {
	url, err := v.buildSecretUrl(path, nil)
	if err != nil {
		return err
	}

	res, err := v.executeRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		// TODO parse error and return it
	}

	return nil
}

func (v *VaultClient) LookupSelf() (VaultClientTokenResponse, error) {
	url, err := v.buildAuthUrl(TOKEN_LOOKUP_PATH, nil)

	res, err := v.executeRequest(http.MethodGet, url, nil)
	if err != nil {
		return VaultClientTokenResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO parse error and return it
	}

	vctResponse := struct {
		Data VaultClientTokenResponse `json:"data"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&vctResponse)
	if err != nil {
		return VaultClientTokenResponse{}, err
	}

	return vctResponse.Data, nil
}

func (v *VaultClient) buildSecretUrl(path string, queryParams map[string]string) (string, error) {
	return v.buildUrl(SECRET_PATH_PREFIX, path, queryParams)
}

func (v *VaultClient) buildAuthUrl(path string, queryParams map[string]string) (string, error) {
	return v.buildUrl(AUTH_PATH_PREFIX, path, queryParams)
}

func (v *VaultClient) buildUrl(pathPrefix, path string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(v.vaultUrl)
	if err != nil {
		return "", nil
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	u.Path = u.Path + SECRET_PATH_PREFIX + "/" + path

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	if queryParams != nil {
		q := u.Query()

		for key, val := range queryParams {
			q.Set(key, val)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil

}

func (v *VaultClient) executeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	h := http.Header{}
	h.Add(VAULT_TOKEN_HEADER, v.vaultToken)
	h.Add(ACCEPT_HEADER, DEFAULT_MEDIA_TYPE)

	if body != nil {
		h.Add(CONTENT_TYPE_HEADER, DEFAULT_MEDIA_TYPE)
	}
	req.Header = h

	return v.httpClient.Do(req)
}
