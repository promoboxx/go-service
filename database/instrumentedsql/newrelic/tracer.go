package newrelic

import (
	"context"
	"fmt"

	"github.com/ExpansiveWorlds/instrumentedsql"
	"github.com/healthimation/go-service/internal"
	newrelic "github.com/newrelic/go-agent"
)

type contextKey string

type tracer struct {
}

type span struct {
	txn     newrelic.Transaction
	segment *newrelic.DatastoreSegment
}

// NewTracer returns a tracer that will create segments on a Newrelic Transaction
func NewTracer() instrumentedsql.Tracer {
	return tracer{}
}

// GetSpan will always return a span with an EXISTING transaction.  TODO: start transaction?
func (tracer) GetSpan(ctx context.Context) instrumentedsql.Span {
	transaction := ctx.Value(internal.ContextKeyNewrelicTransaction)
	if t, ok := transaction.(newrelic.Transaction); ok {
		return span{txn: t}
	}
	return span{}
}

func (s span) NewChild(name string) instrumentedsql.Span {
	if s.txn != nil {
		s.segment = &newrelic.DatastoreSegment{
			StartTime: newrelic.StartSegmentNow(s.txn),
			Product:   newrelic.DatastorePostgres,
			Operation: name,
		}
	}

	return s
}

func (s span) SetLabel(k, v string) {
	if k == "query" {
		s.segment.Operation = fmt.Sprintf("(%s) %s", s.segment.Operation, v)
	}
	return
}

func (s span) Finish() {
	s.segment.End()
}
