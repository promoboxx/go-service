package middleware

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/promoboxx/go-service/alice/middleware/lrw"
)

const (
	logFieldRequestID  = "request_id"
	logFieldUserID     = "user_id"
	logFieldMethod     = "method"
	logFieldStatusCode = "status_code"
	logFieldPath       = "path"
	logFieldError      = "error"
)

// Logger injects a logger into the context
type Logger interface {
	Log(h http.Handler) http.Handler
}

type logger struct {
	entry       *logrus.Entry
	logRequests bool
}

// NewLogrusLogger allows you to setup a base entry to use for logging
func NewLogrusLogger(baseEntry *logrus.Entry, logRequests bool) Logger {
	return &logger{entry: baseEntry, logRequests: logRequests}
}

func (l *logger) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingResponseWriter := &lrw.LoggingResponseWriter{w, http.StatusOK, nil}

		// add logger to the context
		ctx := context.WithValue(r.Context(), contextKeyLogger, entry)
		r = r.WithContext(ctx)

		h.ServeHTTP(loggingResponseWriter, r)

		entry := l.entry.WithFields(logrus.Fields{
			logFieldRequestID:  getRequestIDFromContext(r.Context()),
			logFieldUserID:     getInsecureUserIDFromContext(r.Context()),
			logFieldMethod:     r.Method,
			logFieldStatusCode: loggingResponseWriter.StatusCode,
			logFieldPath:       r.URL.String(),
			logFieldError:      loggingResponseWriter.InnerError,
		})

		if l.logRequests {
			entry.Printf("Finished request")
		}
	})
}
