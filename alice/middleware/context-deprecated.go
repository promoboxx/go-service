package middleware

import (
	"context"
	"database/sql"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/sirupsen/logrus"
)

// Deprecated: Use GetRequestIDFromCtx instead, which handles missing values more gracefully.
func GetRequestIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyRequestID)
}

// Deprecated: Use GetInsecureUserIDFromCtx instead, which handles missing values more gracefully.
func GetInsecureUserIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyInsecureUserID)
}

// Deprecated: Use GetLoggerFromCtx instead, which returns an error when logger is not present.
func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	inter := ctx.Value(contextkey.ContextKeyLogger)
	if entry, ok := inter.(*logrus.Entry); ok {
		return entry
	}
	return logrus.NewEntry(logrus.New()) // default to a new entry
}

// Deprecated: Use GetClaimsFromCtx instead, which returns an error when claims are not present.
// This function will panic if claims are not present in the context.
func GetClaimsFromContext(ctx context.Context) auth.Claim {
	return ctx.Value(contextkey.ContextKeyClaims).(auth.Claim)
}

// Deprecated: Use GetDBConnFromCtx instead, which returns an error when DB connection is not present.
// This function will panic if DB connection is not present in the context.
func GetDBConnFromContext(ctx context.Context) *sql.DB {
	return ctx.Value(contextkey.ContextKeyDBConn).(*sql.DB)
}

// Deprecated: Use GetDBFromCtx instead, which returns an error when DB is not present.
func GetDBFromContext[T any](ctx context.Context, f func(*sql.DB) T) T {
	if ctx.Value(contextkey.ContextKeyDB) == nil {
		return f(GetDBConnFromContext(ctx))
	}

	return ctx.Value(contextkey.ContextKeyDB).(T)
}

// Deprecated: Use GetCanaryVersionFromCtx instead, which handles missing values more gracefully.
func GetCanaryVersionFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyCanary)
}

// Deprecated: Use GetJWTFromCtx instead, which handles missing values more gracefully.
func GetJWTFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyJWT)
}

// getStringFromContext returns a string from the context or empty string if missing
func getStringFromContext(ctx context.Context, key contextkey.ContextKey) string {
	stri := ctx.Value(key)
	if str, ok := stri.(string); ok {
		return str
	}
	return ""
}
