package services

import (
	"encoding/base64"
	"math/big"
	"net/http"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/gin-gonic/gin"
)

// OIDCDiscovery returns the OpenID Connect Discovery document
// per OpenID Connect Discovery 1.0 specification.
func OIDCDiscovery(ctx *gin.Context) {
	issuer := config.GetConfig().OIDC.Issuer

	discovery := dto.OIDCDiscovery{
		Issuer:                            issuer,
		AuthorizationEndpoint:             issuer + "/oauth2/authorize",
		TokenEndpoint:                     issuer + "/oauth2/token",
		UserinfoEndpoint:                  issuer + "/oauth2/userinfo",
		JwksURI:                           issuer + "/.well-known/jwks.json",
		ScopesSupported:                   []string{"openid", "profile", "email", "offline_access"},
		ResponseTypesSupported:            []string{"code"},
		GrantTypesSupported:               []string{"authorization_code", "refresh_token"},
		SubjectTypesSupported:             []string{"public"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post", "none"},
		ClaimsSupported:                   []string{"sub", "iss", "aud", "exp", "iat", "auth_time", "nonce", "name", "given_name", "family_name", "preferred_username", "email", "email_verified", "updated_at"},
		CodeChallengeMethodsSupported:     []string{"S256"},
		ResponseModesSupported:            []string{"query"},
	}

	ctx.JSON(http.StatusOK, discovery)
}

// JWKS returns the JSON Web Key Set containing the public key used to sign tokens.
func JWKS(ctx *gin.Context) {
	pubKey := helpers.GetRSAPublicKey()

	// Derive exponent from the actual public key (not hardcoded)
	e := big.NewInt(int64(pubKey.E))

	jwk := dto.JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: helpers.GetRSAKeyID(),
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(e.Bytes()),
	}

	ctx.JSON(http.StatusOK, dto.JWKSResponse{Keys: []dto.JWK{jwk}})
}
