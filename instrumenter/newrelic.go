package instrumenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/promoboxx/go-service/internal"
	newrelic "github.com/newrelic/go-agent"
)

type newrelicClient struct {
	app newrelic.Application
}

func (c *newrelicClient) StartTransaction(ctx context.Context, name string, w http.ResponseWriter, r *http.Request) (context.Context, Transaction) {
	txn := c.app.StartTransaction(name, w, r)
	txnCtx := context.WithValue(ctx, internal.ContextKeyNewrelicTransaction, txn)
	return txnCtx, &newrelicTxn{txn: txn}
}

type newrelicTxn struct {
	txn newrelic.Transaction
}

func (n *newrelicTxn) End() error {
	return n.txn.End()
}

func (n *newrelicTxn) NoticeError(err error) error {
	return n.txn.NoticeError(err)
}

func newrelicStartTimer(ctx context.Context, name string) newrelic.Segment {
	transaction := ctx.Value(internal.ContextKeyNewrelicTransaction)
	s := newrelic.Segment{}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		s = newrelic.StartSegment(txn, name)
	}
	return s // newrelic.Segment implements End(), and safely handles nil segments
}

// externalCallTimer wrapes newrelic.ExternalSegment to make the implementation
// details a bit more abstract.
type externalCallTimer struct {
	segment newrelic.ExternalSegment
}

func newrelicStartExternalTimer(ctx context.Context, serviceName string, path string) *externalCallTimer {
	c := &externalCallTimer{}
	transaction := ctx.Value(internal.ContextKeyNewrelicTransaction)
	if txn, ok := transaction.(newrelic.Transaction); ok {
		c.segment = newrelic.ExternalSegment{
			StartTime: newrelic.StartSegmentNow(txn),
			URL:       fmt.Sprintf("http://%s/%s", serviceName, path),
		}
	}
	return c // newrelic.ExternalSegment implements End(), and safely handles nil segments
}

func (e *externalCallTimer) End() {
	e.segment.End()
}

func newrelicAddAddtribute(ctx context.Context, key string, value interface{}) error {
	transaction := ctx.Value(internal.ContextKeyNewrelicTransaction)
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return txn.AddAttribute(key, value)
	}
	return nil // i don't think we want an error if we don't have a transaction
}
