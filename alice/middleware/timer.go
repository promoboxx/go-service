package middleware

import (
	"github.com/justinas/alice"
)

// Timer can time a handler and log it
type Timer interface {
	Time(name string) alice.Constructor
}
