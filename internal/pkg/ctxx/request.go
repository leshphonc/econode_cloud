package ctxx

import (
	"context"
	"strings"
)

const ctxRequestIDKey ctxKey = "X-Request-ID"

func WithRequestID(c context.Context, requestID string) context.Context {
	return context.WithValue(c, ctxRequestIDKey, requestID)
}

func RequestID(c context.Context) string {
	if v, ok := c.Value(ctxRequestIDKey).(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}
