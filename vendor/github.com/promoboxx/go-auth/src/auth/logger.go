package auth

// Logger can Printf
type Logger interface {
	Printf(format string, v ...interface{})
}

type nullLogger struct{}

func (n *nullLogger) Printf(format string, v ...interface{}) {
	//noop
}
