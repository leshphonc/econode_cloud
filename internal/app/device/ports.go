package device

import (
	"context"
	"econode-cloud/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repo interface {
	WithDB(db *gorm.DB) Repo
	Register(ctx context.Context, serialNo string, meta map[string]any) (*model.Device, error)
	Activate(ctx context.Context, params ActivateParams) (*model.Device, error)
	UpsertBySerialNo(ctx context.Context, params UpsertDeviceBySerialNoParams) (*model.Device, error)
	EnsureDeviceState(ctx context.Context, deviceID int64, lastSeenAt *time.Time) error
	GetBySerialNo(ctx context.Context, serialNo string) (*model.Device, error)
	TouchLastSeen(ctx context.Context, deviceID int64, t *time.Time) error
	UpdateDeviceStateByHeartbeat(ctx context.Context, p UpdateDeviceStateParams) error
	GetByDeviceUID(ctx context.Context, deviceUID string) (*model.Device, error)
}
