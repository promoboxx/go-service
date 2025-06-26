package graph

import (
	"context"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware"
)

// EnsureAuthenticated checks if the context contains valid authentication claims
// and returns them along with an error if they are not present.
func EnsureAuthenticated(ctx context.Context) (auth.Claim, error) {
	return middleware.GetClaimsFromCtx(ctx)
}
