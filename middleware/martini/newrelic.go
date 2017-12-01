package martini

import (
	"context"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/healthimation/go-service/internal"
	newrelic "github.com/newrelic/go-agent"
)

type newrelicMiddlewareTimer struct {
	app newrelic.Application
}

func (n *newrelicMiddlewareTimer) Measure(name string) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		txn := n.app.StartTransaction(name, res, req)
		// add user-agent to all requests
		txn.AddAttribute(newrelic.AttributeRequestUserAgent, req.Header.Get("User-Agent"))
		// Put newrelic app in context for custom events
		appCtx := context.WithValue(req.Context(), internal.ContextKeyNewrelicApp, n.app)
		// Put newrelic transaction in context for all other instrumentation
		ctx := context.WithValue(appCtx, internal.ContextKeyNewrelicTransaction, txn)
		c.Map(req.WithContext(ctx))
		c.MapTo(txn, (*http.ResponseWriter)(nil))
		c.Next()
		txn.End()
	}
}
