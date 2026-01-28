package ctxx

import (
	"context"
)

type ctxKey string

const (
	ctxDeviceIDKey  ctxKey = "device_id"
	ctxDeviceUIDKey ctxKey = "device_uid"
)

func WithDeviceID(ctx context.Context, deviceID int64) context.Context {
	return context.WithValue(ctx, ctxDeviceIDKey, deviceID)
}
func DeviceID(c context.Context) int64 {
	if v, ok := c.Value(ctxDeviceIDKey).(int64); ok {
		return v
	}
	return 0
}

func WithDeviceUID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, ctxDeviceUIDKey, deviceID)
}

func DeviceUID(c context.Context) string {
	if v, ok := c.Value(ctxDeviceUIDKey).(string); ok {
		return v
	}
	return ""
}
