---
title: Go (go-oidc)
---

# Go Integration

Use [`coreos/go-oidc`](https://github.com/coreos/go-oidc) — the standard OIDC library for Go.

## Install

```bash
go get github.com/coreos/go-oidc/v3
go get golang.org/x/oauth2
```

## Configuration

```go
package main

import (
    "context"
    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
)

var (
    provider     *oidc.Provider
    oauth2Config oauth2.Config
)

func init() {
    ctx := context.Background()

    var err error
    provider, err = oidc.NewProvider(ctx, "http://localhost:8080")
    if err != nil {
        panic(err)
    }

    oauth2Config = oauth2.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        RedirectURL:  "http://localhost:3000/callback",
        Endpoint:     provider.Endpoint(),
        Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
    }
}
```

## Login Handler

```go
func handleLogin(w http.ResponseWriter, r *http.Request) {
    state := generateRandomState() // Implement CSRF state
    http.SetCookie(w, &http.Cookie{
        Name:  "oauth_state",
        Value: state,
    })
    http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}
```

## Callback Handler

```go
func handleCallback(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Verify state
    cookie, _ := r.Cookie("oauth_state")
    if r.URL.Query().Get("state") != cookie.Value {
        http.Error(w, "Invalid state", http.StatusBadRequest)
        return
    }

    // Exchange code for tokens
    token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
    if err != nil {
        http.Error(w, "Token exchange failed", http.StatusInternalServerError)
        return
    }

    // Extract and verify ID token
    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "No id_token in response", http.StatusInternalServerError)
        return
    }

    verifier := provider.Verifier(&oidc.Config{ClientID: "your-client-id"})
    idToken, err := verifier.Verify(ctx, rawIDToken)
    if err != nil {
        http.Error(w, "ID token verification failed", http.StatusInternalServerError)
        return
    }

    // Extract claims
    var claims struct {
        Email string   `json:"email"`
        Name  string   `json:"name"`
        Roles []string `json:"roles"`
    }
    idToken.Claims(&claims)

    fmt.Fprintf(w, "Hello %s! Roles: %v", claims.Name, claims.Roles)
}
```

## Token Validation Middleware

For resource servers that need to validate access tokens:

```go
func authMiddleware(next http.Handler) http.Handler {
    verifier := provider.Verifier(&oidc.Config{ClientID: "your-client-id"})

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := verifier.Verify(r.Context(), tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add claims to request context
        var claims struct {
            Sub   string   `json:"sub"`
            Email string   `json:"email"`
            Roles []string `json:"roles"`
        }
        token.Claims(&claims)

        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```
