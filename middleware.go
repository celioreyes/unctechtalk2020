package main

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		modCtx := context.WithValue(r.Context(), "currentUser", claims["userId"])
		next.ServeHTTP(w, r.WithContext(modCtx))
	})
}
