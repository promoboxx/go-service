package trace

import (
	"log"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"

	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	"github.com/promoboxx/go-service/alice/middleware"
	"github.com/promoboxx/go-service/alice/middleware/lrw"
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
			// make the response writer a logging response writer so we can access the status code
			w = &lrw.LoggingResponseWriter{w, http.StatusOK, nil}
			var span opentracing.Span
			ctx := r.Context()

			// attempt to extract any span information in case this request is coming from somewhere
			// that may have already set one
			sctx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil && err != tracer.ErrSpanContextNotFound {
				log.Printf("error extracting span data: %v", err)
			}

			// if there was no span context start one
			if sctx == nil {
				span, ctx = opentracing.StartSpanFromContext(r.Context(), "http.request")
			} else {
				// start a new span as a child of the existing span
				span = opentracing.StartSpan("http.request", opentracing.ChildOf(sctx))
			}

			// required to get trace search and analytics
			span.SetTag(ext.AnalyticsEvent, true)
			// sets the name of the request that gets passed into measure
			span.SetTag(ext.ResourceName, name)
			span.SetTag(ext.HTTPURL, r.URL.String())
			span.SetTag(ext.HTTPMethod, r.Method)
			defer span.Finish()

			// inject the span information into the http headers so when we send off to another internal
			// service it will have the information to chain everything together
			err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Printf("error injecting span data: %v", err)
			}

			h.ServeHTTP(w, r.WithContext(ctx))

			// set some tags based on status code
			if lrw, ok := w.(*lrw.LoggingResponseWriter); ok {
				otext.HTTPStatusCode.Set(span, uint16(lrw.StatusCode))

				if lrw.StatusCode > 299 || lrw.StatusCode < 200 {
					otext.Error.Set(span, true)

					if lrw.InnerError != nil {
						span.SetTag(ext.Error, lrw.InnerError)
					}
				}
			}
		})
	}
}
