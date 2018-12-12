package database

import (
	"context"

	newrelic "github.com/newrelic/go-agent"
	"github.com/promoboxx/go-service/internal"
)

func newrelicStartDBTimer(ctx context.Context, collection, operation, query string) newrelic.DatastoreSegment {
	transaction := ctx.Value(internal.ContextKeyNewrelicTransaction)
	s := newrelic.DatastoreSegment{}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		s = newrelic.DatastoreSegment{
			StartTime:          newrelic.StartSegmentNow(txn),
			Product:            newrelic.DatastorePostgres,
			Collection:         collection,
			Operation:          operation,
			ParameterizedQuery: query,
		}
	}
	return s // newrelic.DatastoreSegment implements End(), and safely handles nil segments
}
