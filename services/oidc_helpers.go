package services

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// isValidRedirectURI checks if a given URI exactly matches one of the registered redirect URIs.
// Per OIDC Core §3.1.2.1, matching uses Simple String Comparison (RFC 3986 §6.2.1).
func isValidRedirectURI(allowed []string, uri string) bool {
	for _, u := range allowed {
		if u == uri {
			return true
		}
	}
	return false
}

// verifyPKCE verifies the code_verifier against the stored code_challenge using S256.
// Per RFC 7636 §4.6: BASE64URL-ENCODE(SHA256(ASCII(code_verifier))) == code_challenge
func verifyPKCE(codeChallenge, codeVerifier string) bool {
	h := sha256.Sum256([]byte(codeVerifier))
	computed := base64.RawURLEncoding.EncodeToString(h[:])
	return computed == codeChallenge
}

// redirectWithError redirects the user agent with properly URL-encoded OAuth error parameters.
// Per OIDC Core §3.1.2.6, error parameters are added to the redirect_uri query component.
func redirectWithError(ctx *gin.Context, redirectURI, errorCode, description, state string) {
	params := url.Values{}
	params.Set("error", errorCode)
	params.Set("error_description", description)
	if state != "" {
		params.Set("state", state)
	}
	ctx.Redirect(http.StatusFound, redirectURI+"?"+params.Encode())
}
