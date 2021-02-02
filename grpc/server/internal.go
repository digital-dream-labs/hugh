package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/digital-dream-labs/hugh/log"
)

func (s *Server) changeState(st State) {
	s.mu.Lock()
	s.state = st
	s.mu.Unlock()
	s.log.Debugf("State changed to %s", st)
	s.notifyState(st)
}

func (s *Server) notifyState(st State) {
	s.mu.RLock()
	for _, chs := range s.notifyChan[st] {
		chs <- st
	}
	s.mu.RUnlock()
}

// handleSignals allows for clean exits.
func (s *Server) handleSignals() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// block until signal is Received
	x := <-ch

	s.log.WithFields(log.Fields{
		"signal":              x,
		"numActiveGoRoutines": runtime.NumGoroutine(),
	}).Warn("received os signal")

	s.log.Warn("shutting down")
	s.Stop()
	s.log.Warn("shut down")
}

func (s *Server) appendErr(err error) {
	s.mu.Lock()
	s.errs = append(s.errs, err)
	s.mu.Unlock()
	s.log.Debugf("Registered error: %v", err)
}

func serverTLS(o *options) *tls.Config {
	//nolint -- This is in place to interact with the older services until they're upgraded...
	return &tls.Config{
		//MinVersion:   tls.VersionTLS13,
		ClientCAs:    o.mustGetCertPool(),
		Certificates: o.certificates,
		ClientAuth:   o.clientAuth,
	}
}

func grpcHandlerFunc(grpcServer, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					headers := []string{"Content-Type", "Accept"}
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
					methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
					return
				}
			}
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func (s *Server) getListener(port int, certs []tls.Certificate) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	lis = &internalListener{
		Listener:        lis,
		firstAcceptFunc: func() { s.changeState(Ready) },
	}

	if certs == nil {
		return lis, nil
	}

	tlsListener := tls.NewListener(
		lis,
		&tls.Config{
			MinVersion:   tls.VersionTLS13,
			Certificates: certs,
			NextProtos:   []string{"http/1.1"},
		},
	)

	return tlsListener, nil
}
