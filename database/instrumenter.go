package database

import "context"

type DBInstrumentTimer interface {
	End() error
}

// StartDBTimer starts a special DB timer and returns something that must be End()-ed
func StartDBTimer(ctx context.Context, collection, operation, query string) DBInstrumentTimer {
	return newrelicStartDBTimer(ctx, collection, operation, query)
}
