package middleware

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/justinas/alice"
)

// Timer can time a handler and log it
type Timer interface {
	Time(name string) alice.Constructor
}

type nullTimer struct{}

func (n *nullTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

func NewNullTimer() Timer {
	return &nullTimer{}
}

type openTracingTimer struct {
}

// NewOpenTracingTimer creates a new timer that uses opentracing spans
func NewOpenTracingTimer() Timer {
	return &openTracingTimer{}
}

// Time returns a middleware that will handle creating spans and tracing
// for each request
func (ott *openTracingTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := GetLoggerFromContext(r.Context())
			var span opentracing.Span
			var ctx context.Context
			sctx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				log.Errorf("Error extracting headers from context: %v", err)
			}

			// if there was no span context start one
			if sctx == nil {
				span, ctx = opentracing.StartSpanFromContext(r.Context(), name)
			} else {
				span = opentracing.StartSpan(name, opentracing.ChildOf(sctx))
				err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
				if err != nil {
					log.Errorf("Error injecting headers: %v", err)
				}
			}

			defer span.Finish()
			newReq := r.WithContext(ctx)

			h.ServeHTTP(w, newReq)
		})
	}
}
