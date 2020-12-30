package log

import (
	"context"
	"time"

	"github.com/digital-dream-labs/hugh/log"
	"google.golang.org/grpc"
)

// SafeLoggingModifier provides functions needed for this
// middleware to safely modify copies of data on the chain
type SafeLoggingModifier interface {
	Handles(interface{}) bool
	Copy(interface{}) interface{}
	Modify(interface{}) interface{}
}

// InterceptOption is used to configure interceptors
type InterceptOption func(*interceptConfig)

// InterceptWithReqModifier returns an interceptOption to pass to Intercept* functions
func InterceptWithReqModifier(method string, modifier SafeLoggingModifier) InterceptOption {
	return func(cfg *interceptConfig) {
		cfg.reqModifiers[method] = modifier
	}
}

// InterceptWithRespModifier returns an interceptOption to pass to Intercept* functions
func InterceptWithRespModifier(method string, modifier SafeLoggingModifier) InterceptOption {
	return func(cfg *interceptConfig) {
		cfg.respModifiers[method] = modifier
	}
}

type interceptConfig struct {
	reqModifiers  map[string]SafeLoggingModifier
	respModifiers map[string]SafeLoggingModifier
}

func newInterceptConfig() *interceptConfig {
	return &interceptConfig{
		reqModifiers:  make(map[string]SafeLoggingModifier),
		respModifiers: make(map[string]SafeLoggingModifier),
	}
}

type loggingServerStream struct {
	grpc.ServerStream
	ctx             context.Context
	method          string
	rcvd, snt       int
	reqMod, respMod SafeLoggingModifier
}

func (l *loggingServerStream) Context() context.Context {
	return l.ctx
}

func (l *loggingServerStream) SendMsg(m interface{}) error {
	start := time.Now()
	err := l.ServerStream.SendMsg(m)

	resp := prepareForLog(l.respMod, m)

	log.WithFields(log.Fields{
		"Method":   l.method,
		"Response": resp,
		"Error":    err,
		"Duration": time.Since(start).String(),
	}).Debug("send message")

	l.snt++

	return err
}

func (l *loggingServerStream) RecvMsg(m interface{}) error {
	start := time.Now()
	err := l.ServerStream.RecvMsg(m)

	req := prepareForLog(l.reqMod, m)

	log.WithFields(log.Fields{
		"Method":   l.method,
		"Request":  req,
		"Error":    err,
		"Duration": time.Since(start).String(),
	}).Debug("receive message")

	l.rcvd++

	return err
}

// StreamServerInterceptor returns an interceptor for logging grpc streaming messages
func StreamServerInterceptor(opts ...InterceptOption) grpc.StreamServerInterceptor {
	cfg := newInterceptConfig()
	for _, o := range opts {
		o(cfg)
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		newCtx := log.AddToContext(ss.Context())
		lss := &loggingServerStream{ServerStream: ss, method: info.FullMethod, ctx: newCtx}

		if f, ok := cfg.reqModifiers[info.FullMethod]; ok {
			lss.reqMod = f
		}
		if f, ok := cfg.respModifiers[info.FullMethod]; ok {
			lss.respMod = f
		}

		err := handler(srv, lss)

		log.FromContext(newCtx).WithFields(log.Fields{
			"Method":   info.FullMethod,
			"Sent":     lss.snt,
			"Received": lss.rcvd,
			"Error":    err,
			"Duration": time.Since(start).String(),
		}).Info("stream info")

		return err
	}
}

// UnaryServerInterceptor returns an interceptor for logging grpc unary requests
func UnaryServerInterceptor(opts ...InterceptOption) grpc.UnaryServerInterceptor {
	cfg := newInterceptConfig()
	for _, o := range opts {
		o(cfg)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		newCtx := log.AddToContext(ctx)
		resp, err := handler(newCtx, req)

		reqCopy := prepareForLog(cfg.reqModifiers[info.FullMethod], req)
		respCopy := prepareForLog(cfg.respModifiers[info.FullMethod], resp)

		log.FromContext(newCtx).WithFields(log.Fields{
			"Method":   info.FullMethod,
			"Request":  reqCopy,
			"Response": respCopy,
			"Error":    err,
			"Duration": time.Since(start).String(),
		}).Info("request info")

		return resp, err
	}
}

func prepareForLog(modifier SafeLoggingModifier, src interface{}) interface{} {
	cp := src
	if modifier != nil && modifier.Handles(src) {
		cp = modifier.Copy(src)
		if src == cp {
			return src
		}
		cp = modifier.Modify(cp)
	}
	return cp
}
