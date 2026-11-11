package logger

import (
	"context"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	return zap.NewProduction()
}

type contextKey struct{}

var nop = zap.NewNop()

func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

func FromContext(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return log
	}
	return nop
}
