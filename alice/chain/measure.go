package chain

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/promoboxx/go-service/alice/middleware"
)

// Measurer can setup a Measured chain
type Measurer interface {
	Measure(name string, handler http.Handler) http.HandlerFunc
}

type base struct {
	baseChain alice.Chain
	timer     middleware.Timer
	logger    middleware.Logger
}

// NewBase gets a new measurer with the provided base chain
// Expected usage:
// t, err := middleware.NewNewrelicTimer(env, serviceName, nrKey)
// if err != nil {
// 	log.Fatalf("Could not instantiate newrelic timer: %v", err)
// }
// b := chain.NewBase(alice.New(), t, middleware.NewLogrusLogger(logrus.NewEntry(logrus.New())), jwtdecode.NewJWTDecoder())
// router.Get("/user", b.Measure("get users", user.Get()))
//
func NewBase(b alice.Chain, timer middleware.Timer, logger middleware.Logger, jwtDecoder middleware.JWTDecoder) Measurer {
	c := b.Append(middleware.Recovery, middleware.NewUserIDInjector(jwtDecoder).Inject, middleware.RequestID)
	return &base{baseChain: c, timer: timer, logger: logger}
}

// NewBaseWithExtras similar to NewBase but allows users to pass in a set of additional constructors to append the the base chain
func NewBaseWithExtras(b alice.Chain, timer middleware.Timer, logger middleware.Logger, jwtDecoder middleware.JWTDecoder, constructors ...alice.Constructor) Measurer {
	c := b.Append(middleware.Recovery, middleware.NewUserIDInjector(jwtDecoder).Inject, middleware.RequestID)
	c = c.Append(constructors...)
	return &base{baseChain: c, timer: timer, logger: logger}
}

// Measure returns a chain that will have metrics measured
func (b *base) Measure(name string, handler http.Handler) http.HandlerFunc {
	if b.timer != nil {
		return b.baseChain.Append(b.timer.Time(name)).Append(b.logger.Log).Then(handler).ServeHTTP
	}
	return b.baseChain.Then(handler).ServeHTTP
}
