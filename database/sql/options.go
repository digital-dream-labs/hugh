package sql

// Option provides a function definition to set options
type Option func(*options)

type options struct {
	databaseType string
	username     string
	password     string
	name         string
	host         string
	port         int
	tlsMode      string
	errs         []error
}

func (o *options) errored() bool {
	return len(o.errs) > 0
}

// WithViper specifies that the server should construct its options using viper.
func WithViper(args ...string) Option {
	return func(o *options) {
		if err := o.viperize(args...); err != nil {
			o.errs = append(o.errs, err)
		}
	}
}

// WithDatabaseType sets the db type.
func WithDatabaseType(p string) Option {
	return func(o *options) {
		o.databaseType = p
	}
}

// WithUsername sets the db username.
func WithUsername(p string) Option {
	return func(o *options) {
		o.username = p
	}
}

// WithPassword sets the db password.
func WithPassword(p string) Option {
	return func(o *options) {
		o.password = p
	}
}

// WithName sets the db name.
func WithName(p string) Option {
	return func(o *options) {
		o.name = p
	}
}

// WithHost sets the db host.
func WithHost(p string) Option {
	return func(o *options) {
		o.host = p
	}
}

// WithPort sets the listener port.
func WithPort(p int) Option {
	return func(o *options) {
		o.port = p
	}
}

// WithTLSMode sets the tls mode
func WithTLSMode(p string) Option {
	return func(o *options) {
		o.tlsMode = p
	}
}
