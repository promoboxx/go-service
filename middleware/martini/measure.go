package martini

import (
	"fmt"

	"github.com/go-martini/martini"
	newrelic "github.com/newrelic/go-agent"
)

type MiddlewareTimer interface {
	Measure(name string) martini.Handler
}

// NewMiddlewareTimer returns an instance of something that can Measure() an HTTP handler via middleware
func NewMiddlewareTimer(environment, serviceName, licenseKey string) (MiddlewareTimer, error) {
	config := newrelic.NewConfig(fmt.Sprintf("%s-%s", environment, serviceName), licenseKey)
	config.Enabled = false
	if licenseKey != "" {
		config.Enabled = true
	}
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}
	return &newrelicMiddlewareTimer{app: app}, nil
}

func NewMiddlewareTimerFromNR(app newrelic.Application) MiddlewareTimer {
	return &newrelicMiddlewareTimer{app: app}
}
