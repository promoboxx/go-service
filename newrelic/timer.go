package newrelic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/healthimation/go-service/alice/middleware"
	"github.com/healthimation/go-service/internal"
	"github.com/justinas/alice"
	nr "github.com/newrelic/go-agent"
)

type newrelicTimer struct {
	app nr.Application
}

// NewTimer returns a timer that logs to newrelic
func NewTimer(environment, serviceName, licenseKey string) (middleware.Timer, error) {
	config := nr.NewConfig(fmt.Sprintf("%s-%s", environment, serviceName), licenseKey)
	config.Enabled = false
	if licenseKey != "" {
		config.Enabled = true
	}
	app, err := nr.NewApplication(config)
	if err != nil {
		return nil, err
	}
	return &newrelicTimer{app: app}, nil
}

func (n *newrelicTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := n.app.StartTransaction(name, w, r)
			// add user-agent to all requests
			txn.AddAttribute(nr.AttributeRequestUserAgent, r.Header.Get("User-Agent"))
			// Put newrelic app in context for custom events
			appCtx := context.WithValue(r.Context(), internal.ContextKeyNewrelicApp, n.app)
			// Put newrelic transaction in context for all other instrumentation
			ctx := context.WithValue(appCtx, internal.ContextKeyNewrelicTransaction, txn)
			r = r.WithContext(ctx)
			h.ServeHTTP(txn, r)
			txn.End()
		})
	}
}
