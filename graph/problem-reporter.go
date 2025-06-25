package graph

import (
	"context"

	"github.com/promoboxx/go-glitch/glitch"
	"github.com/promoboxx/go-service/alice/middleware"
)

func ReportProblem(ctx context.Context, prob *glitch.GQLProblem, logMsg string) error {
	middleware.GetLoggerFromContext(ctx).WithError(prob).Error(logMsg)
	return prob
}

func ReportError(ctx context.Context, prob *glitch.GQLProblem, err error) error {
	middleware.GetLoggerFromContext(ctx).WithError(err).Error(prob.Error())
	return prob
}
