package device

import (
	"context"
	"econode-cloud/internal/model"
	"econode-cloud/internal/pkg/ptr"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepo(db *gorm.DB) Repo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) WithDB(db *gorm.DB) Repo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Register(ctx context.Context, serialNo string, meta map[string]any) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	device := model.Device{}
	err := db.Where("serial_no = ?", serialNo).Attrs(model.Device{
		DeviceUID: uuid.New(),
		SerialNo:  serialNo,
		Status:    int16(model.DeviceStatusPreregistration),
		ClaimCode: uuid.New().String(),
	}).FirstOrCreate(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepo) Activate(ctx context.Context, params ActivateParams) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	device := model.Device{}
	// 查出待激活设备
	err := db.Where("serial_no = ? AND claim_code = ?", params.SerialNo, params.ClaimCode).First(&device).Error
	if err != nil {
		return nil, err
	}

	// 激活设备
	now := time.Now()
	err = db.Model(&device).Updates(map[string]any{
		"serial_no":  params.SerialNo,
		"model":      params.Model,
		"status":     int16(model.DeviceStatusActive),
		"power_mode": params.PowerMode,
		"hw_version": params.HWVersion,
		"fw_version": params.FWVersion,
		"claim_code": "",
		"claimed_at": now,
		"meta":       params.Meta,
	}).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// UpsertBySerialNo 按serialNo查询设备，如果存在返回，没有创建
func (r *deviceRepo) UpsertBySerialNo(ctx context.Context, params UpsertDeviceBySerialNoParams) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	device := model.Device{}
	err := db.Where("serial_no = ?", params.SerialNo).Attrs(model.Device{
		DeviceUID: uuid.New(),
		SerialNo:  params.SerialNo,
		Name:      ptr.String("测试"),
	}).FirstOrCreate(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// EnsureDeviceState 按device_id查询设备状态，如果存在更新，没有创建
func (r *deviceRepo) EnsureDeviceState(ctx context.Context, deviceID int64, lastSeenAt *time.Time) error {
	deviceState := model.DeviceState{
		DeviceID:   deviceID,
		LastSeenAt: lastSeenAt,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"last_seen_at": lastSeenAt,
			"updated_at":   gorm.Expr("now()"),
		}),
	}).Create(&deviceState).Error
}

// GetBySerialNo 根据SerialNo查询设备
func (r *deviceRepo) GetBySerialNo(ctx context.Context, serialNo string) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	var dev model.Device
	err := db.Where("serial_no = ?", serialNo).Take(&dev).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dev, nil
}

func (r *deviceRepo) TouchLastSeen(ctx context.Context, deviceID int64, t *time.Time) error {
	db := r.db.WithContext(ctx)
	return db.Model(&model.Device{}).
		Where("id = ?", deviceID).
		Updates(map[string]any{
			"last_seen_at": t,
			"updated_at":   gorm.Expr("now()"),
		}).Error
}

func (r *deviceRepo) UpdateDeviceStateByHeartbeat(ctx context.Context, p UpdateDeviceStateParams) error {
	db := r.db.WithContext(ctx)

	updates := map[string]any{
		"updated_at": gorm.Expr("now()"),
	}
	if p.LastSeenAt != nil {
		updates["last_seen_at"] = p.LastSeenAt
	}
	if p.LastHeartbeatAt != nil {
		updates["last_heartbeat_at"] = p.LastHeartbeatAt
	}
	if p.DoorOpen != nil {
		updates["door_open"] = *p.DoorOpen
	}
	if p.SignalStrength != nil {
		updates["signal_strength"] = *p.SignalStrength
	}
	if p.BatteryLevel != nil {
		updates["battery_level"] = *p.BatteryLevel
	}

	updates["weight"] = p.Weight

	// payload merge：如果你希望“合并而不是覆盖”，用 jsonb || 操作
	if p.Payload != nil {
		updates["payload"] = p.Payload
	}

	return db.Model(&model.DeviceState{}).
		Where("device_id = ?", p.DeviceID).
		Updates(updates).Error
}

func (r *deviceRepo) GetByDeviceUID(ctx context.Context, publicID string) (*model.Device, error) {
	return nil, nil
}
