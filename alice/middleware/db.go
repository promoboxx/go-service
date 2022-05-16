package middleware

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/promoboxx/go-service/database"
	"github.com/promoboxx/go-service/service"
)

func InjectDBConn(connector database.SQLDBConnector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbConn, err := connector.GetConnection()
			if err != nil {
				service.WriteProblem(w, "database error", "DATABASE_ERROR", http.StatusInternalServerError, err)
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.ContextKeyDBConn, dbConn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
