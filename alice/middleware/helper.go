package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	l = logrus.New()
)

// HandlerFunc converts a handler with an error to a standard handler
func HandlerFunc(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log := GetLoggerFromContext(r.Context())
			log.Printf("Error from handler: %v", err)
		}
	})
}

// GetDefaultLogger gets a default logger to use
func GetDefaultLogger(serviceName, environment string) *logrus.Entry {
	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		FieldMap: FieldMap{
			"service":     serviceName,
			"environment": environment,
		},
		Level: logrus.ErrorLevel,
	}
	return logrus.NewEntry(l)
}

func SetLogLevel(level int) {
	l.Level = level
}
