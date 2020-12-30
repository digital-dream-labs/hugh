package echopb

import (
	context "golang.org/x/net/context"
)

// Echo is the config struct for the echo service
type Echo struct{}

// Echo is a function for the echo service
func (t Echo) Echo(cxt context.Context, in *EchoMessage) (*EchoMessage, error) {
	return in, nil
}
