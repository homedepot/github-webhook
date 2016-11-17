package zest

import (
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/contraband/holler"
)

// Sink is the type that represents the sink that will emit errors to Yeller.
type Sink struct {
	yeller *holler.Yeller
}

// NewYellerSink creates a new Sink for use with Lager.
func NewYellerSink(token, env string) *Sink {
	return &Sink{
		yeller: holler.NewYeller(token, env),
	}
}

// Log will send any error log lines up to Yeller.
func (s *Sink) Log(line lager.LogFormat) {
	if line.LogLevel < lager.ERROR {
		return
	}

	if errStr, ok := line.Data["error"].(string); ok {
		delete(line.Data, "message")
		delete(line.Data, "error")

		s.yeller.Notify(
			line.Message,
			errors.New(errStr),
			line.Data,
		)
	}
}
