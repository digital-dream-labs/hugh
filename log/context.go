package log

import (
	"context"
)

type contextKey string

func (k contextKey) String() string {
	return "log context " + string(k)
}

var (
	middlewareContextKey = contextKey("grpc middleware")
)

// FromContext extracts a Logger from context returning base if not set
func FromContext(ctx context.Context) Logger {
	f, ok := ctx.Value(middlewareContextKey).(Fields)
	if !ok {
		return Base()
	}
	fields := Fields{}

	for k, v := range f {
		fields[k] = v
	}

	return WithFields(fields)
}

// AddToContext creates the context entry if it does not exist
func AddToContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(middlewareContextKey).(Fields); ok {
		return ctx
	}
	return context.WithValue(ctx, middlewareContextKey, Fields{})
}

// AddContextFields adds fields to be logged to the context
func AddContextFields(ctx context.Context, fields Fields) {
	f, ok := ctx.Value(middlewareContextKey).(Fields)
	if !ok {
		return
	}

	for k, v := range fields {
		f[k] = v
	}
}
