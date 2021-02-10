package log

import (
	"context"
	"time"

	"github.com/digital-dream-labs/hugh/log"
	"google.golang.org/grpc"
)

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
			methodTag:   info.FullMethod,
			sentTag:     reqCopy,
			responseTag: respCopy,
			errorTag:    err,
			durationTag: time.Since(start).Seconds(),
		}).Info("request info")

		return resp, err
	}
}
