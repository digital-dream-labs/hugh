package log

import (
	"context"
	"time"

	"github.com/digital-dream-labs/hugh/log"
	"google.golang.org/grpc"
)

type loggingServerStream struct {
	grpc.ServerStream
	ctx             context.Context
	method          string
	rcvd, snt       int
	reqMod, respMod SafeLoggingModifier
}

const (
	methodTag   = "method"
	sentTag     = "sent"
	receivedTag = "received"
	responseTag = "response"
	errorTag    = "error"
	durationTag = "duration"
)

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
			methodTag:   info.FullMethod,
			sentTag:     lss.snt,
			receivedTag: lss.rcvd,
			errorTag:    err,
			durationTag: time.Since(start).Seconds(),
		}).Info("stream info")

		return err
	}
}

func (l *loggingServerStream) Context() context.Context {
	return l.ctx
}

func (l *loggingServerStream) SendMsg(m interface{}) error {
	start := time.Now()
	err := l.ServerStream.SendMsg(m)

	resp := prepareForLog(l.respMod, m)

	log.WithFields(log.Fields{
		methodTag:   l.method,
		sentTag:     resp,
		receivedTag: err,
		durationTag: time.Since(start).Seconds(),
	}).Debug("send message")

	l.snt++

	return err
}

func (l *loggingServerStream) RecvMsg(m interface{}) error {
	start := time.Now()
	err := l.ServerStream.RecvMsg(m)

	req := prepareForLog(l.reqMod, m)

	log.WithFields(log.Fields{
		methodTag:   l.method,
		sentTag:     req,
		receivedTag: err,
		durationTag: time.Since(start).Milliseconds(),
	}).Debug("receive message")

	l.rcvd++

	return err
}
