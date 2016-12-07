package vault

type VaultAuthResponse struct {
	ClientToken   string            `json:"clientToken"`
	Policies      []string          `json:"policies"`
	MetaData      map[string]string `json:"metaData"`
	LeaseDuration int               `json:"leaseDuration"`
	Renewable     bool              `json:"renewable"`
}

type VaultClientTokenResponse struct {
	ID          string            `json:"id"`
	Policies    []string          `json:"policies"`
	Path        string            `json:"path"`
	Meta        map[string]string `json:"meta"`
	DisplayName string            `json:"displayName"`
	NumUses     int               `json:"numUses"`
}

type VaultEnableAuditBackendRequest struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
}

type VaultHealthResponse struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
	Standby     bool `json:"standby"`
}

type VaultInitResponse struct {
	Keys      []string `json:"keys"`
	RootToken string   `json:"rootToken"`
}

type VaultListResponse struct {
	Keys []string `json:"keys"`
}

type VaultPolicy struct {
	Rules string `json:"rules"`
}

type VaultResponse struct {
	Data map[string]string `json:"data"`
}

type VaultRevokeTokenRequest struct {
	token string `json:"token"`
}

type VaultSealStatusResponse struct {
	Sealed   bool `json:"sealed"`
	T        int  `json:"t"`
	N        int  `json:"n"`
	Progress int  `json:"progress"`
}

type VaultTokenAuthRequest struct {
	ID              string            `json:"id"`
	Policies        []string          `json:"policies"`
	Meta            map[string]string `json:"meta"`
	NoParent        bool              `json:"noParent"`
	NoDefaultPolicy bool              `json:"noDefaultPolicy"`
	TTL             string            `json:"ttl"`
	DisplayName     string            `json:"displayName"`
	NumUses         int               `json:"numUses"`
}

type VaultUnsealRequest struct {
	Key   string `json:"key"`
	Reset bool   `json:"reset"`
}

type VaultErrorResponse struct {
	Errors []string `json:"errors"`
}
