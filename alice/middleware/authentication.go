package middleware

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/promoboxx/go-service/service"
)

func Auth(token auth.Token) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := token.ValidateJWT(r.Header.Get("Authorization"))
			if err != nil {
				service.WriteProblem(w, "Could not validate JWT", "NOT_AUTHORIZED", http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.ContextKeyClaims, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
