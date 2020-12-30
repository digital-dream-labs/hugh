package server

import (
	"net"
	"sync"
)

const (
	internalServerPort = 65533
)

type internalListener struct {
	net.Listener
	firstAccept     sync.Once
	firstAcceptFunc func()
}

// Accept hijacks the Accept method from the internal listener.
func (l *internalListener) Accept() (net.Conn, error) {
	l.firstAccept.Do(l.firstAcceptFunc)

	return l.Listener.Accept()
}

// Addr hijacks the Addr method from the internal listener.
func (l *internalListener) Addr() net.Addr {
	return l.Listener.Addr()
}

// Close hijacks the Close method from the internal listener.
func (l *internalListener) Close() error {
	return l.Listener.Close()
}
