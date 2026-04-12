---
title: OIDC Integration Overview
---

# OIDC Integration

AuthForest is a standards-compliant **OpenID Connect** provider. Any OIDC-compatible client library can be used to integrate with it.

## Supported Flows

| Flow                          | Use Case                                      |
| ----------------------------- | --------------------------------------------- |
| **Authorization Code + PKCE** | Single-page apps, mobile apps, public clients |
| **Authorization Code**        | Server-side apps, confidential clients        |

## Configuration

All OIDC clients need these settings:

| Setting                    | Value                                         |
| -------------------------- | --------------------------------------------- |
| **Authority / Issuer URL** | `http://localhost:8080` (your AuthForest URL) |
| **Client ID**              | From the application you registered           |
| **Redirect URI**           | Must match a registered redirect URI          |
| **Scopes**                 | `openid profile email`                        |
| **Response Type**          | `code`                                        |

## Discovery Endpoint

AuthForest serves a standard discovery document:

```
GET /.well-known/openid-configuration
```

Most OIDC client libraries will auto-discover all endpoints from this URL.

## Token Response

After a successful authorization code exchange, you'll receive:

```json
{
	"access_token": "eyJhbGciOi...",
	"id_token": "eyJhbGciOi...",
	"token_type": "Bearer",
	"expires_in": 3600,
	"refresh_token": "eyJhbGciOi..."
}
```

The `id_token` is a JWT containing user claims (sub, email, name, roles).

## Framework Guides

Choose your framework for a step-by-step integration example:

- [JavaScript / TypeScript](/docs/oidc-javascript) — using `oidc-client-ts`
- [React](/docs/oidc-react) — using `react-oidc-context`
- [NextAuth / Auth.js](/docs/oidc-nextauth) — using the OIDC provider
- [Python](/docs/oidc-python) — using `authlib`
- [Go](/docs/oidc-go) — using `coreos/go-oidc`
