package middleware

import (
	"context"
	"database/sql"
	"errors"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/sirupsen/logrus"
)

// GetRequestIDFromCtx returns the requestID and an error if it's not present
func GetRequestIDFromCtx(ctx context.Context) (string, error) {
	requestID := ctx.Value(contextkey.ContextKeyRequestID)
	if requestID == nil {
		return "", errors.New("no request ID in context")
	}

	strRequestID, ok := requestID.(string)
	if !ok {
		return "", errors.New("invalid request ID type in context")
	}

	return strRequestID, nil
}

// MustGetRequestIDFromContext returns the requestID and panics if it's not present
func MustGetRequestIDFromContext(ctx context.Context) string {
	requestID, err := GetRequestIDFromCtx(ctx)
	if err != nil {
		panic(err)
	}
	return requestID
}

// GetInsecureUserIDFromCtx returns the user ID and an error if it's not present
func GetInsecureUserIDFromCtx(ctx context.Context) (string, error) {
	userID := ctx.Value(contextkey.ContextKeyInsecureUserID)
	if userID == nil {
		return "", errors.New("no user ID in context")
	}

	strUserID, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID type in context")
	}

	return strUserID, nil
}

// MustGetInsecureUserIDFromContext returns the user ID and panics if it's not present
func MustGetInsecureUserIDFromContext(ctx context.Context) string {
	userID, err := GetInsecureUserIDFromCtx(ctx)
	if err != nil {
		panic(err)
	}
	return userID
}

// GetLoggerFromCtx returns a logrus entry from the context and an error if it's not present
func GetLoggerFromCtx(ctx context.Context) (*logrus.Entry, error) {
	logger := ctx.Value(contextkey.ContextKeyLogger)
	if logger == nil {
		return nil, errors.New("no logger in context")
	}

	entry, ok := logger.(*logrus.Entry)
	if !ok {
		return nil, errors.New("invalid logger type in context")
	}

	return entry, nil
}

// MustGetLoggerFromContext returns a logrus entry from the context and panics if it's not present
// If no logger is found, it returns a default logger
func MustGetLoggerFromContext(ctx context.Context) *logrus.Entry {
	logger, err := GetLoggerFromCtx(ctx)
	if err != nil {
		return logrus.NewEntry(logrus.New()) // default to a new entry
	}
	return logger
}

// GetClaimsFromCtx returns the auth claims from the context and an error if they are not present
func GetClaimsFromCtx(ctx context.Context) (auth.Claim, error) {
	claims := ctx.Value(contextkey.ContextKeyClaims)
	if claims == nil {
		var nilClaim auth.Claim
		return nilClaim, errors.New("no auth claims in context")
	}

	authClaims, ok := claims.(auth.Claim)
	if !ok {
		var nilClaim auth.Claim
		return nilClaim, errors.New("invalid auth claims type in context")
	}

	return authClaims, nil
}

// MustGetClaimsFromContext returns the auth claims from the context and panics if they are not present
func MustGetClaimsFromContext(ctx context.Context) auth.Claim {
	claims, err := GetClaimsFromCtx(ctx)
	if err != nil {
		panic(err)
	}
	return claims
}

// GetDBConnFromCtx returns the DB connection from the context and an error if it's not present
func GetDBConnFromCtx(ctx context.Context) (*sql.DB, error) {
	dbConn := ctx.Value(contextkey.ContextKeyDBConn)
	if dbConn == nil {
		return nil, errors.New("no DB connection in context")
	}

	db, ok := dbConn.(*sql.DB)
	if !ok {
		return nil, errors.New("invalid DB connection type in context")
	}

	return db, nil
}

// MustGetDBConnFromContext returns the DB connection from the context and panics if it's not present
func MustGetDBConnFromContext(ctx context.Context) *sql.DB {
	db, err := GetDBConnFromCtx(ctx)
	if err != nil {
		panic(err)
	}
	return db
}

// GetDBFromCtx returns a database object from the context and an error if it's not present
func GetDBFromCtx[T any](ctx context.Context, f func(*sql.DB) T) (T, error) {
	if ctx.Value(contextkey.ContextKeyDB) == nil {
		dbConn, err := GetDBConnFromCtx(ctx)
		if err != nil {
			var zero T
			return zero, err
		}
		return f(dbConn), nil
	}

	db := ctx.Value(contextkey.ContextKeyDB)
	result, ok := db.(T)
	if !ok {
		var zero T
		return zero, errors.New("invalid DB type in context")
	}

	return result, nil
}

// MustGetDBFromContext returns a database object from the context and panics if it's not present
func MustGetDBFromContext[T any](ctx context.Context, f func(*sql.DB) T) T {
	result, err := GetDBFromCtx(ctx, f)
	if err != nil {
		panic(err)
	}
	return result
}

// GetCanaryVersionFromCtx returns the canary version and an error if it's not present
func GetCanaryVersionFromCtx(ctx context.Context) (string, error) {
	canary := ctx.Value(contextkey.ContextKeyCanary)
	if canary == nil {
		return "", errors.New("no canary version in context")
	}

	strCanary, ok := canary.(string)
	if !ok {
		return "", errors.New("invalid canary version type in context")
	}

	return strCanary, nil
}

// MustGetCanaryVersionFromContext returns the canary version and panics if it's not present
func MustGetCanaryVersionFromContext(ctx context.Context) string {
	canary, err := GetCanaryVersionFromCtx(ctx)
	if err != nil {
		return "" // Return empty string if not found
	}
	return canary
}

// GetJWTFromCtx returns the JWT and an error if it's not present
func GetJWTFromCtx(ctx context.Context) (string, error) {
	jwt := ctx.Value(contextkey.ContextKeyJWT)
	if jwt == nil {
		return "", errors.New("no JWT in context")
	}

	strJWT, ok := jwt.(string)
	if !ok {
		return "", errors.New("invalid JWT type in context")
	}

	return strJWT, nil
}

// MustGetJWTFromContext returns the JWT and panics if it's not present
func MustGetJWTFromContext(ctx context.Context) string {
	jwt, err := GetJWTFromCtx(ctx)
	if err != nil {
		panic(err)
	}
	return jwt
}
