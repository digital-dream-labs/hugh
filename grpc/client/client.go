package client

import (
	"crypto/tls"

	"github.com/pkg/errors"

	"fmt"
	"sync"

	log "github.com/digital-dream-labs/hugh/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is an rpc connection
type Client struct {
	target   string
	conn     *grpc.ClientConn
	dialOpts []grpc.DialOption
	mu       sync.RWMutex
}

// New creates a client
func New(opts ...Option) (*Client, error) {
	cfg := options{
		log: log.Base(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	cli := Client{
		target:   cfg.target,
		dialOpts: cfg.dialOpts,
	}

	if cfg.errz != nil {
		err := errors.New("error:")
		for _, v := range cfg.errz {
			err = errors.Wrap(
				err,
				v.Error(),
			)
		}
		return nil, err
	}

	// configure tls options
	switch {
	case cfg.disableTLS:
		cli.dialOpts = append(cli.dialOpts, grpc.WithInsecure())
	case cfg.insecure && cfg.certificates == nil:
		creds := credentials.NewTLS(insecureClientTLS(&cfg))
		cli.dialOpts = append(cli.dialOpts, grpc.WithTransportCredentials(creds))
	case cfg.insecure && cfg.certificates != nil:
		creds := credentials.NewTLS(insecureClientTLSWithClientCerts(&cfg))
		cli.dialOpts = append(cli.dialOpts, grpc.WithTransportCredentials(creds))

	default:
		fmt.Println("default")
		creds := credentials.NewTLS(clientTLS(&cfg))
		cli.dialOpts = append(cli.dialOpts, grpc.WithTransportCredentials(creds))
	}

	if len(cfg.scInterceptors) > 0 {
		cli.dialOpts = append(cli.dialOpts, grpc.WithChainStreamInterceptor(cfg.scInterceptors...))
	}

	if len(cfg.ucInterceptors) > 0 {
		cli.dialOpts = append(cli.dialOpts, grpc.WithChainUnaryInterceptor(cfg.ucInterceptors...))
	}

	return &cli, nil
}

// Connect establishes the GRPC connection to the specified host.
func (c *Client) Connect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return errors.New("connect already called")
	}

	c.conn, err = grpc.Dial(c.target, c.dialOpts...)

	return err
}

// Close closes the connection. Calling close on an empty connection yells at you.
func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("client not connected")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

// Conn returns the connection established to the grpc host.
// This connection is concurrent safe and should be used when
// constructing grpc clients.
func (c *Client) Conn() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil {
		panic("You must connect before you can use this")
	}

	return c.conn
}

func clientTLS(o *options) *tls.Config {
	//nolint -- once older services are upgraded, fix this!
	return &tls.Config{
		RootCAs:      o.mustGetCertPool(),
		Certificates: o.certificates,
	}
}

func insecureClientTLS(o *options) *tls.Config {
	return &tls.Config{
		//nolint -- this is intentional when calling this!
		InsecureSkipVerify: true,
	}
}

func insecureClientTLSWithClientCerts(o *options) *tls.Config {
	return &tls.Config{
		//nolint -- this is intentional when calling this!
		InsecureSkipVerify: true,
		Certificates:       o.certificates,
	}
}
