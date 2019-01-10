package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
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
func GetDefaultLogger(serviceName, environment string, level int) *logrus.Entry {
	l := logrus.New()

	if level < 0 {
		l.Level = 0
	} else if level > len(logrus.AllLevels)-1 {
		l.Level = logrus.AllLevels[len(logrus.AllLevels)-1]
	}

	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		FieldMap: logrus.FieldMap{
			"service":     serviceName,
			"environment": environment,
		},
	}
	return logrus.NewEntry(l)
}
