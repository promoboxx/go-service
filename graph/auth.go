package graph

import (
	"context"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware"
)

// This is a bit of a hack workaround to make the `middleware.GetClaimsFromContext` function return a pointer,
// to gracefully handle situations where no auth claims were injected by any middleware. This change is necessary
// because the GraphQL middleware does not prevent requests from reaching the resolvers if no token is present,
// so they could panic with an unhelpful error message if they tried to call the standard function
func EnsureAuthenticated(ctx context.Context) (c *auth.Claim) {
	defer func() {
		if err := recover(); err != nil {
			c = nil // the return value is named so it can be modified here, pretty cool
		}
	}()

	claim := middleware.GetClaimsFromContext(ctx)
	return &claim
}
