// Package client provides a quick+easy http client that'll marshal and unmarshal structs for you
package client

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const timeout = 5

// Client is the connection struct
type Client struct {
	client http.Client
	target string
}

// New accepts options and returns a new http client
func New(opts ...Option) (*Client, error) {
	cfg := options{
		log: logrus.New(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	cli := Client{
		target: cfg.target,
		client: http.Client{
			Timeout: timeout * time.Second,
		},
	}

	return &cli, nil
}
