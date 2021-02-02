package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/digital-dream-labs/hugh/log"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	idletimeout       = 3
	readheadertimeout = 1
	maxheaderbytes    = 10 << 20
)

// Server is a server struct
type Server struct {
	transport     *grpc.Server
	listener      net.Listener
	httpMux       *grpc_runtime.ServeMux
	httpTransport *http.Server
	httpListener  net.Listener
	httpConfig    *tls.Config
	state         State
	log           log.Logger
	notifyChan    map[State][]chan<- State
	shutdown      func()
	mu            sync.RWMutex
	errs          []error
}

// New constructs a new Server
func New(opts ...Option) (*Server, error) {
	cfg := options{
		log: log.Base(),
	}

	var srvOpts []grpc.ServerOption

	for _, o := range opts {
		o(&cfg)
	}

	if cfg.errored() {
		return nil, fmt.Errorf("error during server setup: %v", cfg.errs)
	}

	if !cfg.insecure {
		if (cfg.tlsCert == "" || cfg.tlsKey == "") && cfg.certificates == nil {
			return nil, errors.New("either set insecure or define TLS certificates appropriately")
		}

		if cfg.certificates == nil {
			crt, err := tls.X509KeyPair(
				[]byte(cfg.tlsCert),
				[]byte(cfg.tlsKey),
			)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			cfg.certificates = append(cfg.certificates, crt)
		}

		creds := credentials.NewTLS(serverTLS(&cfg))
		srvOpts = append(srvOpts, grpc.Creds(creds))
	}

	if len(cfg.ssInterceptors) > 0 {
		c := middleware.ChainStreamServer(cfg.ssInterceptors...)
		srvOpts = append(srvOpts, grpc.StreamInterceptor(c))
	}

	if len(cfg.usInterceptors) > 0 {
		c := middleware.ChainUnaryServer(cfg.usInterceptors...)
		srvOpts = append(srvOpts, grpc.UnaryInterceptor(c))
	}

	srv := Server{
		state:      Init,
		transport:  grpc.NewServer(srvOpts...),
		log:        cfg.log,
		notifyChan: make(map[State][]chan<- State),
	}

	srv.shutdown = srv.transport.GracefulStop

	if cfg.httpPassthrough || cfg.httpPassthroughInsecure {
		srv.httpMux = grpc_runtime.NewServeMux(
			grpc_runtime.WithMarshalerOption(
				grpc_runtime.MIMEWildcard,
				&grpc_runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						Indent:          "",
						UseProtoNames:   true,
						EmitUnpopulated: true,
						Multiline:       true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				},
			),
		)

		mux := http.NewServeMux()
		mux.Handle("/", srv.httpMux)

		if cfg.certificates != nil {
			srv.httpConfig = &tls.Config{
				Certificates: cfg.certificates,
				NextProtos:   []string{"h2"},
				ClientCAs:    cfg.mustGetCertPool(),
				ClientAuth:   cfg.clientAuth,
				//nolint -- the only way to make this proper is to have a SAN in the cert, which may expose
				// some of the internals.  I could go either way on this one..
				InsecureSkipVerify: true,
			}
		}

		srv.httpTransport = &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.port),
			Handler:           grpcHandlerFunc(srv.Transport(), mux),
			TLSConfig:         srv.httpConfig,
			IdleTimeout:       time.Second * idletimeout,
			ReadHeaderTimeout: time.Second * readheadertimeout,
			MaxHeaderBytes:    maxheaderbytes,
		}

		var err error

		switch {
		case cfg.httpPassthrough:
			srv.httpListener, err = srv.getListener(cfg.port, cfg.certificates)
			if err != nil {
				return nil, err
			}
		case cfg.httpPassthroughInsecure:
			srv.httpListener, err = srv.getListener(cfg.port, nil)
			if err != nil {
				return nil, err
			}
		}

		srv.listener, err = srv.getListener(internalServerPort, nil)
		if err != nil {
			return nil, err
		}
	} else {
		rpcLis, _ := srv.getListener(cfg.port, nil)
		srv.listener = rpcLis
	}

	if cfg.reflect {
		reflection.Register(srv.transport)
	}

	return &srv, nil
}

// Address returns the address
func (s *Server) Address() net.Addr {
	return s.listener.Addr()
}

// HTTPAddress returns the address
func (s *Server) HTTPAddress() net.Addr {
	if s.httpListener != nil {
		return s.httpListener.Addr()
	}
	return nil
}

// Start calls the underlying grpc.Server.Serve
func (s *Server) Start() {
	if s.State() != Init {
		s.log.Errorf("server is not in a valid state, want: %q have: %q", "NEW", s.State().String())
		return
	}

	log.WithFields(log.Fields{
		"grpc-address": s.Address(),
		"http-address": s.HTTPAddress(),
	}).Infof("server starting")

	go s.handleSignals()

	s.changeState(Starting)

	go func() {
		if err := s.transport.Serve(s.listener); err != nil {
			s.appendErr(err)
			s.changeState(Error)
		}
	}()
	if s.httpTransport != nil {
		go func() {
			if err := s.httpTransport.Serve(s.httpListener); err != nil {
				s.appendErr(err)
				s.changeState(Error)
			}
		}()
	}
}

// Stop shuts down the service
func (s *Server) Stop() {
	s.changeState(Stopping)
	s.shutdown()
	s.changeState(Stopped)
}

// Notify will send requested rpc.State changes over the channel returned by this function.
func (s *Server) Notify(states ...State) <-chan State {
	ch := make(chan State, len(states))

	s.mu.Lock()
	for _, v := range states {
		s.notifyChan[v] = append(s.notifyChan[v], ch)
		if v == s.state {
			ch <- s.state
		}
	}
	s.mu.Unlock()

	return ch
}

// State returns the current state.
func (s *Server) State() State {
	s.mu.RLock()
	st := s.state
	s.mu.RUnlock()
	return st
}

// Errors returns any errors registered during construction, starting, or stopping.
func (s *Server) Errors() []error {
	s.mu.RLock()
	e := s.errs
	s.mu.RUnlock()
	return e
}

// Transport returns the underlying grpc server.
func (s *Server) Transport() *grpc.Server {
	return s.transport
}

// RegisterHTTPService accepts a slice of http gateway services to register
func (s *Server) RegisterHTTPService(in []func(context.Context, *grpc_runtime.ServeMux, string, []grpc.DialOption) error) error {
	dialOpts := []grpc.DialOption{}

	if s.httpConfig != nil {
		log.Info("starting http proxy with tls certificates")
		dialOpts = append(dialOpts,
			grpc.WithTransportCredentials(
				credentials.NewTLS(s.httpConfig),
			),
		)
	} else {
		log.Info("starting http proxy without tls certificates")
		dialOpts = append(dialOpts,
			grpc.WithInsecure(),
		)
	}

	for _, v := range in {
		log.Debug("registering ", runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name())
		if err := v(context.Background(), s.httpMux, fmt.Sprintf(":%d", internalServerPort), dialOpts); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
