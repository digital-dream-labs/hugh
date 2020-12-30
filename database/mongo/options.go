package mongo

// Option provides a function definition to set opts
type Option func(*opts)

type opts struct {
	username     string
	password     string
	name         string
	host         string
	authdbname   string
	port         int
	errs         []error
	writeConcern string
	direct       bool
	cluster      bool
	retryWrites  bool
}

/*
	TODO:  This needs a kv parser so we can pass arbitrary options to mongo.
	See: https://docs.mongodb.com/manual/reference/connection-string/#connections-connection-options
*/

// WithViper specifies that the server should construct its opts using viper.
func WithViper(args ...string) Option {
	return func(o *opts) {
		if err := o.viperize(args...); err != nil {
			o.errs = append(o.errs, err)
		}
	}
}

// WithUsername sets the db username.
func WithUsername(p string) Option {
	return func(o *opts) {
		o.username = p
	}
}

// WithPassword sets the db password.
func WithPassword(p string) Option {
	return func(o *opts) {
		o.password = p
	}
}

// WithName sets the db name.
func WithName(p string) Option {
	return func(o *opts) {
		o.name = p
	}
}

// WithName sets the db name.
func WithAuthName(p string) Option {
	return func(o *opts) {
		o.authdbname = p
	}
}

// WithHost sets the db host.
func WithHost(p string) Option {
	return func(o *opts) {
		o.host = p
	}
}

// WithPort sets the listener port.
func WithPort(p int) Option {
	return func(o *opts) {
		o.port = p
	}
}

// WithDirect tells mongo this is a direct connection.
func WithDirect(p bool) Option {
	return func(o *opts) {
		o.direct = p
	}
}

// WithCluster sets this up as a clustered (srv) connection
func WithCluster(p bool) Option {
	return func(o *opts) {
		o.cluster = p
	}
}

// WithRetryWrites turns retryWrites on or off
func WithRetryWrites(p bool) Option {
	return func(o *opts) {
		o.retryWrites = p
	}
}

// WithWriteConcern sets the write concern (0, 1, or majority)
func WithWriteConcern(p string) Option {
	return func(o *opts) {
		o.writeConcern = p
	}
}
