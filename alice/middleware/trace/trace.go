package trace

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/promoboxx/go-service/alice/middleware"

	"github.com/justinas/alice"
)

type openTracingTimer struct {
}

// NewOpenTracingTimer creates a new Tracer that uses opentracing spans
func NewOpenTracingTimer() middleware.Timer {
	return &openTracingTimer{}
}

// Time returns a middleware that will handle creating spans and tracing
// for each request
func (ott *openTracingTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := middleware.GetLoggerFromContext(r.Context())
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

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
