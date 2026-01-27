package ctxx

import (
	"context"
	"strings"
)

type ctxKey string

const ctxDeviceUIDKey ctxKey = "X-Device-UID"

func WithDeviceUID(c context.Context, deviceUID string) context.Context {
	return context.WithValue(c, ctxDeviceUIDKey, deviceUID)
}

func DeviceUID(c context.Context) string {
	if v, ok := c.Value(ctxDeviceUIDKey).(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}
