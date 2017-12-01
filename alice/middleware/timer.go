package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/alice"

	"github.com/healthimation/go-service/internal"
	newrelic "github.com/newrelic/go-agent"
)

type newrelicTimer struct {
	app newrelic.Application
}

// Timer can time a handler and log it
type Timer interface {
	Time(name string) alice.Constructor
}

// NewNewrelicTimer returns a timer that logs to newrelic
func NewNewrelicTimer(environment, serviceName, licenseKey string) (Timer, error) {
	config := newrelic.NewConfig(fmt.Sprintf("%s-%s", environment, serviceName), licenseKey)
	config.Enabled = false
	if licenseKey != "" {
		config.Enabled = true
	}
	app, err := newrelic.NewApplication(config)
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
			txn.AddAttribute(newrelic.AttributeRequestUserAgent, r.Header.Get("User-Agent"))
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
