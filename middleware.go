package main

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is is an unexported custom type to prevent collisions with keys defined outside of
// this package
//
// See https://golang.org/pkg/context/#WithValue or https://blog.golang.org/context#TOC_3.2. for
// more info
type contextKey int

type customKey string

const (
	// authedKey is the context key used for storing the auth status of the current request
	authedKey contextKey = iota
	// tokenKey is the context key to use for storing the parsed auth token
	tokenKey

	currentUserID    customKey = "currentUserID"
	currentUserEmail customKey = "currentUserEmail"
)

// CheckJWT is http middleware that reads a JWT from the Authorization header and verifies and
// stores the parsed token to context. In addition to verifying the signature and time claims, it
// ensures that the token has the correct issuer and audience. If no JWT is provided or the token is
// invalid an 401 response is sent
func CheckJWT(verifier IVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get token from authorization header.
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 7 || strings.ToUpper(authHeader[0:6]) != "BEARER" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// strip off "Bearer "
			tokenStr := authHeader[7:]

			// parse and verify token and claims
			token, err := verifier.VerifyToken(tokenStr)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, tokenKey, token)
			ctx = context.WithValue(ctx, authedKey, true)
			ctx = context.WithValue(ctx, currentUserID, token.Claims.UserID)
			ctx = context.WithValue(ctx, currentUserEmail, token.Claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckPermissions makes sure the target permission is found within the token
func CheckPermissions(targetPerm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Grab token from context
			token := r.Context().Value(tokenKey).(*Token)
			logger.Printf("Scopes are: %s", token.Claims.Scope)

			var isAllowed bool
			for _, perm := range token.Claims.Permissions {
				if perm == targetPerm {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// We are allowed move on to next request
			next.ServeHTTP(w, r)
		})
	}
}
