package ctxx

import (
	"context"

	"go.uber.org/zap"
)

const ctxKeyLogger ctxKey = "logger"

func WithLogger(c context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(c, ctxKeyLogger, logger)
}

func Logger(c context.Context) *zap.Logger {
	if v := c.Value(ctxKeyLogger); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return zap.NewNop()
}
