package device

import (
	"context"
	"econode-cloud/internal/model"
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

func (r *deviceRepo) GetBySerialNo(ctx context.Context, serialNo string) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	var dev model.Device
	err := db.Where("serial_no = ?", serialNo).Take(&dev).Error
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (r *deviceRepo) GetByDeviceUID(ctx context.Context, DeviceUID string) (*model.Device, error) {
	db := r.db.WithContext(ctx)

	var dev model.Device
	err := db.Where("device_uid = ?", DeviceUID).Take(&dev).Error
	if err != nil {
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

func (r *deviceRepo) InsertHeartbeat(ctx context.Context, hb *model.Heartbeat) error {
	return r.db.WithContext(ctx).Create(hb).Error
}

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
