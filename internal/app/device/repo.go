package device

import (
	"context"
	"econode-cloud/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repo interface {
	WithDB(db *gorm.DB) Repo

	// GetOrCreateBySerialNo 注册设备
	GetOrCreateBySerialNo(ctx context.Context, serialNo string, meta map[string]any) (*model.Device, error)
	// ClaimBySerialNo 激活设备
	ClaimBySerialNo(ctx context.Context, params ClaimParams) (*model.Device, error)

	GetBySerialNo(ctx context.Context, serialNo string) (*model.Device, error)
	GetByDeviceUID(ctx context.Context, deviceUID string) (*model.Device, error)

	CreateHeartbeat(ctx context.Context, hb *model.Heartbeat) error
	UpdateLastSeenAt(ctx context.Context, deviceID int64, t time.Time) error

	EnsureDeviceState(ctx context.Context, deviceID int64, lastSeenAt *time.Time) error
}
