package client

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/digital-dream-labs/hugh/log"
	"google.golang.org/grpc"
)

// Option provides a function definition to set options
type Option func(*options)

type options struct {
	log            log.Logger
	disableTLS     bool
	insecure       bool
	certPool       *x509.CertPool
	certificates   []tls.Certificate
	clientAuth     tls.ClientAuthType
	errz           []error
	target         string
	scInterceptors []grpc.StreamClientInterceptor
	ucInterceptors []grpc.UnaryClientInterceptor
	dialOpts       []grpc.DialOption
}

func (o *options) mustGetCertPool() *x509.CertPool {
	if o.certPool != nil {
		return o.certPool
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	return certPool
}

func (o *options) errored() bool {
	return len(o.errz) > 0
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
func WithLogger(l log.Logger) Option {
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

// WithCertPool overrides the system CA pool
func WithCertPool(p *x509.CertPool) Option {
	return func(o *options) {
		o.certPool = p
	}
}

// WithCertificate adds certs for  authentication.
func WithCertificate(c ...tls.Certificate) Option {
	return func(o *options) {
		o.certificates = c
	}
}

// WithClientAuth sets the tls ClientAuthType to control auth behavior.
func WithClientAuth(a tls.ClientAuthType) Option {
	return func(o *options) {
		o.clientAuth = a
	}
}

// WithStreamClientInterceptors sets the streaming client middleware.
func WithStreamClientInterceptors(sc ...grpc.StreamClientInterceptor) Option {
	return func(o *options) {
		o.scInterceptors = append(o.scInterceptors, sc...)
	}
}

// WithUnaryClientInterceptors sets the unary client middleware.
func WithUnaryClientInterceptors(uc ...grpc.UnaryClientInterceptor) Option {
	return func(o *options) {
		o.ucInterceptors = append(o.ucInterceptors, uc...)
	}
}

// WithDialopts adds dial options.
func WithDialopts(uc ...grpc.DialOption) Option {
	return func(o *options) {
		o.dialOpts = append(o.dialOpts, uc...)
	}
}

// WithInsecureSkipVerify makes connections not that safe to use.
func WithInsecureSkipVerify() Option {
	return func(o *options) {
		o.insecure = true
	}
}

// WithDisableTLS turns off tls.
func WithDisableTLS() Option {
	return func(o *options) {
		o.disableTLS = true
	}
}
