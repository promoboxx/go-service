package instrumenter

import (
	"context"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// Client is a weak wrapper around the NR application 'start transaction' functionality, so we don't
// need to pull the NR lib directly into all the things.
type Client interface {
	// Start a top level transaction, omit w and r for a background transaction
	StartTransaction(ctx context.Context, name string, w http.ResponseWriter, r *http.Request) (context.Context, Transaction)
}

// Txn wraps the bits of the Newrelic Transaction interface that we use
type Transaction interface {
	End() error
	NoticeError(err error) error
}

// InstrumentTimer must be used within a Transaction, it maps to a Newrelic 'segment'
type InstrumentTimer interface {
	End() error
}

// ExternalCallTimer is a timer for service a to instrument outgoing calls within a Transaction
type ExternalCallTimer interface {
	End() // Stop the timer
}

func NewClientFromNR(app newrelic.Application) Client {
	return &newrelicClient{app: app}
}

// StartTimer starts a timer and returns something that must be End()-ed
func StartTimer(ctx context.Context, name string) InstrumentTimer {
	return newrelicStartTimer(ctx, name)
}

// StartExternalCallTimer starts a timer for an outgoing HTTP call
// and returns something that must be End()-ed
func StartExternalCallTimer(ctx context.Context, serviceName string, path string) ExternalCallTimer {
	return newrelicStartExternalTimer(ctx, serviceName, path)
}

// Add an attribute to the transaction stored in the context, if there is one.
func AddTxnAttribute(ctx context.Context, key string, value interface{}) error {
	return newrelicAddAddtribute(ctx, key, value)
}

// TODO: add metric/event
// https://github.com/newrelic/go-agent/blob/master/GUIDE.md#custom-events
