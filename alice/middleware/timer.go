package middleware

import (
	"context"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

	"github.com/justinas/alice"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
			var ctx context.Context
			var span ddtrace.Span
			defer func() {
				if span != nil {
					span.Finish()
				}
			}()

			spanCtx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
			if err != nil && err != tracer.ErrSpanContextNotFound {
				log.Printf("Error extracting tracing headers: %v", err)
			}

			if spanCtx != nil {
				span = tracer.StartSpan(name, tracer.ChildOf(spanCtx))
			} else {
				span, ctx = tracer.StartSpanFromContext(r.Context(), name)
				err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(r.Header))
				if err != nil {
					log.Printf("Error injecting tracing headers: %v", err)
				}
				r.WithContext(ctx)
			}

			h.ServeHTTP(w, r)
		})
	}
}
