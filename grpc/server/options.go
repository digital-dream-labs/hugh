package server

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/digital-dream-labs/hugh/log"

	"google.golang.org/grpc"
)

// Option provides a function definition to set options
type Option func(*options)

type options struct {
	certificates            []tls.Certificate
	errs                    []error
	ssInterceptors          []grpc.StreamServerInterceptor
	usInterceptors          []grpc.UnaryServerInterceptor
	log                     log.Logger
	certPool                *x509.CertPool
	clientAuth              tls.ClientAuthType
	port                    int
	tlsCert                 string
	tlsKey                  string
	insecure                bool
	reflect                 bool
	httpPassthrough         bool
	httpPassthroughInsecure bool
}

func (o *options) mustGetCertPool() *x509.CertPool {
	if o.certPool != nil {
		return o.certPool
	}

	return x509.NewCertPool()
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

// WithLogger set the log instance.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.log = l
	}
}

// WithPort sets the listener port.
func WithPort(p int) Option {
	return func(o *options) {
		o.port = p
	}
}

// WithCertPool overrides the system cert pool.
func WithCertPool(p *x509.CertPool) Option {
	return func(o *options) {
		o.certPool = p
	}
}

// WithCertificate specifies the certificates that will be appended to the tls.Config.
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

// WithStreamServerInterceptors sets the streaming middleware.
func WithStreamServerInterceptors(ss ...grpc.StreamServerInterceptor) Option {
	return func(o *options) {
		o.ssInterceptors = ss
	}
}

// WithUnaryServerInterceptors sets the unary middleware.
func WithUnaryServerInterceptors(us ...grpc.UnaryServerInterceptor) Option {
	return func(o *options) {
		o.usInterceptors = us
	}
}

// WithReflectionService starts a reflection service capable to describing services
func WithReflectionService() Option {
	return func(o *options) {
		o.reflect = true
	}
}

// WithHTTPPassthrough starts an additional HTTP passthrough for rest operations
func WithHTTPPassthrough() Option {
	return func(o *options) {
		o.httpPassthrough = true
	}
}

// WithHTTPPassthroughInsecure enables an insecure http passthrough
func WithHTTPPassthroughInsecure() Option {
	return func(o *options) {
		o.httpPassthroughInsecure = true
	}
}

// WithInsecureSkipVerify makes connections not that safe to use.
func WithInsecureSkipVerify() Option {
	return func(o *options) {
		o.insecure = true
	}
}

// WithTLSCert statically sets a TLS certificate.
func WithTLSCert(s string) Option {
	return func(o *options) {
		o.tlsCert = s
	}
}

// WithTLSKey statically sets a TLS key.
func WithTLSKey(s string) Option {
	return func(o *options) {
		o.tlsKey = s
	}
}
