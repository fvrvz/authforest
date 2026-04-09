package dto

type OIDCDiscovery struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	ResponseModesSupported            []string `json:"response_modes_supported"`
}

type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	CodeVerifier string `form:"code_verifier"`
	RefreshToken string `form:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type UserInfoResponse struct {
	Sub               string `json:"sub"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope"`
	State               string `form:"state"`
	Nonce               string `form:"nonce"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

type RegisterClientRequest struct {
	ClientName               string   `json:"client_name" binding:"required"`
	RedirectURIs             []string `json:"redirect_uris" binding:"required,min=1"`
	ClientType               string   `json:"client_type" binding:"required,oneof=public confidential"`
	Scopes                   string   `json:"scopes"`
	GrantTypes               string   `json:"grant_types"`
	AccessTokenExpiryMinutes *int     `json:"access_token_expiry_minutes"`
	RefreshTokenExpiryHours  *int     `json:"refresh_token_expiry_hours"`
	IDTokenExpiryMinutes     *int     `json:"id_token_expiry_minutes"`
}

type RegisterClientResponse struct {
	ClientID                 string   `json:"client_id"`
	ClientSecret             string   `json:"client_secret,omitempty"`
	ClientName               string   `json:"client_name"`
	ClientType               string   `json:"client_type"`
	RedirectURIs             []string `json:"redirect_uris"`
	Scopes                   string   `json:"scopes"`
	GrantTypes               string   `json:"grant_types"`
	AccessTokenExpiryMinutes int      `json:"access_token_expiry_minutes"`
	RefreshTokenExpiryHours  int      `json:"refresh_token_expiry_hours"`
	IDTokenExpiryMinutes     int      `json:"id_token_expiry_minutes"`
	CreatedAt                string   `json:"created_at,omitempty"`
}

type UpdateClientRequest struct {
	ClientName               string   `json:"client_name"`
	RedirectURIs             []string `json:"redirect_uris"`
	Scopes                   string   `json:"scopes"`
	GrantTypes               string   `json:"grant_types"`
	AccessTokenExpiryMinutes *int     `json:"access_token_expiry_minutes"`
	RefreshTokenExpiryHours  *int     `json:"refresh_token_expiry_hours"`
	IDTokenExpiryMinutes     *int     `json:"id_token_expiry_minutes"`
}

type DashboardStats struct {
	TotalUsers   int64 `json:"total_users"`
	TotalClients int64 `json:"total_clients"`
}
