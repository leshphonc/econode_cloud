package ctxx

import (
	"context"
	"strings"
)

type ctxKey string

const ctxTraceIDKey ctxKey = "X-Trace-ID"

func WithTraceID(c context.Context, traceID string) context.Context {
	return context.WithValue(c, ctxTraceIDKey, traceID)
}

func TraceID(c context.Context) string {
	if v, ok := c.Value(ctxTraceIDKey).(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}
