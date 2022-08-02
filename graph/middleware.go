package graph

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/promoboxx/go-service/service"
)

// Attempts to get the auth token from the HTTP Authorization header, decrypt it, and place it on the context, but allows
// requests through even if they don't have this header field
func AuthOptional(token auth.Token) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authRaw := r.Header.Get("Authorization")

			if len(authRaw) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := token.ValidateJWT(authRaw)
			if err != nil {
				service.WriteProblem(w, "Could not validate JWT", "NOT_AUTHORIZED", http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.ContextKeyClaims, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
