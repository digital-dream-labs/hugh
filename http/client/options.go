package client

import (
	"github.com/sirupsen/logrus"
)

// Option provides a function definition to set options
type Option func(*options)

type options struct {
	log    *logrus.Logger
	errz   []error
	target string
}

// WithViper specifies that the client or server should construct its options
// using viper
func WithViper(args ...string) Option {
	return func(o *options) {
		if err := o.viperize(args...); err != nil {
			o.errz = append(o.errz, err)
		}
	}
}

// WithLogger set the log instance for client and server instances.
func WithLogger(l *logrus.Logger) Option {
	return func(o *options) {
		o.log = l
	}
}

// WithTarget sets the target host.
func WithTarget(t string) Option {
	return func(o *options) {
		o.target = t
	}
}
