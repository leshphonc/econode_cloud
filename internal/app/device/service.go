package device

import (
	"context"
	"econode-cloud/internal/model"
	"econode-cloud/internal/pkg/nullable"
	"econode-cloud/internal/pkg/txm"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	txm        *txm.TxManager
	deviceRepo Repo
}

func NewService(txm *txm.TxManager, deviceRepo Repo) *Service {
	return &Service{
		txm,
		deviceRepo,
	}
}

type AuthService interface {
	AuthByDeviceUID(ctx context.Context, deviceUID string) (*IdentityResult, error)
}

type IdentityResult struct {
	DeviceID  int64  // 内部 ID，后续 service/repo 用它更快
	DeviceUID string // 外部 UUID（原 public_id）
}

func (s *Service) AuthByDeviceUID(ctx context.Context, deviceUID string) (*IdentityResult, error) {
	dev, err := s.deviceRepo.GetByDeviceUID(ctx, deviceUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
	}

	if dev.DisabledAt != nil || dev.RetiredAt != nil {
		return nil, ErrDeviceDisabled
	}

	return &IdentityResult{
		DeviceID:  dev.ID,
		DeviceUID: dev.DeviceUID.String(),
	}, err
}

type RegisterParams struct {
	SerialNo string
	Meta     map[string]any
}

type RegisterResult struct {
	SerialNo  string
	ClaimCode string
}

func (s *Service) Register(ctx context.Context, p RegisterParams) (RegisterResult, error) {
	dev, err := s.deviceRepo.Register(ctx, p.SerialNo, p.Meta)
	if err != nil {
		return RegisterResult{}, ErrDeviceRegisterFailed
	}

	return RegisterResult{
		SerialNo:  dev.SerialNo,
		ClaimCode: dev.ClaimCode,
	}, nil
}

type ActivateParams struct {
	SerialNo  string
	Model     string
	PowerMode int16
	HWVersion string
	FWVersion string
	ClaimCode string
	Meta      map[string]any
}

type ActivateResult struct {
	DeviceUID    string
	Name         string
	Model        string
	Status       int16
	PowerMode    int16
	HWVersion    string
	FWVersion    string
	ClaimAt      int64
	ActiveErrors []string
	Meta         map[string]any
}

func mapPowerMode(i int16) (v model.DevicePowerMode, err error) {
	switch i {
	case 1:
		v = model.PowerModeGrid
	case 2:
		v = model.PowerModeSolar
	case 3:
		v = model.PowerModeHybrid
	default:
		err = ErrDevicePowerModeUnknow
	}

	return
}

func (s *Service) Activate(ctx context.Context, p ActivateParams) (ActivateResult, error) {
	pw, err := mapPowerMode(p.PowerMode)
	if err != nil {
		return ActivateResult{}, err
	}
	p.PowerMode = int16(pw)

	dev, err := s.deviceRepo.Activate(ctx, p)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ActivateResult{}, ErrDeviceActivateFailed
		}
		return ActivateResult{}, err
	}

	return ActivateResult{
		DeviceUID:    dev.DeviceUID.String(),
		Name:         nullable.StrOrEmpty(dev.Name),
		Model:        nullable.StrOrEmpty(dev.Model),
		Status:       dev.Status,
		PowerMode:    nullable.Int16OrZero(dev.PowerMode),
		HWVersion:    nullable.StrOrEmpty(dev.HWVersion),
		FWVersion:    nullable.StrOrEmpty(dev.FWVersion),
		ClaimAt:      dev.ClaimedAt.Unix(),
		ActiveErrors: dev.ActiveErrors,
		Meta:         dev.Meta,
	}, nil

}

type UpsertDeviceBySerialNoParams struct {
	SerialNo string
	Meta     map[string]any
}

type UpdateDeviceStateParams struct {
	DeviceID        int64
	LastSeenAt      *time.Time
	LastHeartbeatAt *time.Time
	DoorOpen        *bool
	SignalStrength  *int16
	BatteryLevel    *int16

	// numeric(12,3) 输入暂用 string
	Weight decimal.Decimal

	Payload datatypes.JSONMap
}

type HeartbeatParams struct {
	SerialNo string

	DoorOpen       *bool
	SignalStrength *int16
	BatteryLevel   *int16
	Weight         *string // numeric(12,3) 先用 string 输入
	Payload        datatypes.JSONMap
}

type HeartbeatResult struct {
	DeviceID int64
}

var weightRe = regexp.MustCompile(`^\d+(\.\d{1,3})?$`)

func ParseWeightKgDecimal(input string) (decimal.Decimal, error) {
	s := strings.TrimSpace(input)

	// 1. 协议校验（严格）
	if !weightRe.MatchString(s) {
		return decimal.Zero, errors.New("invalid weight format")
	}

	// 2. parse 成 decimal（不会有精度问题）
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, err
	}

	// 3. 兜底校验（可选）
	if d.IsNegative() {
		return decimal.Zero, errors.New("weight must be >= 0")
	}

	return d, nil
}

func (s *Service) Heartbeat(ctx context.Context, p HeartbeatParams) (HeartbeatResult, error) {
	if p.SerialNo == "" {
		return HeartbeatResult{}, errors.New("serial_no required")
	}

	now := time.Now()

	var res HeartbeatResult
	err := s.txm.Transaction(ctx, func(tx *gorm.DB) error {
		deviceRepo := s.deviceRepo.WithDB(tx)

		// 1) 查设备（按 serial_no）
		dev, err := deviceRepo.GetBySerialNo(ctx, p.SerialNo)
		if err != nil {
			return err
		}
		if dev == nil {
			return errors.New("device not found")
		}
		res.DeviceID = dev.ID

		// 2) 确保 device_state 存在
		if err = s.deviceRepo.EnsureDeviceState(ctx, dev.ID, &now); err != nil {
			return err
		}

		// 3) 更新 device.last_seen_at
		if err = s.deviceRepo.TouchLastSeen(ctx, dev.ID, &now); err != nil {
			return err
		}

		// 4) 更新 device_state
		kgDecimal, err := ParseWeightKgDecimal(*p.Weight)
		if err != nil {
			return err
		}
		if err = s.deviceRepo.UpdateDeviceStateByHeartbeat(ctx, UpdateDeviceStateParams{
			DeviceID:        dev.ID,
			LastSeenAt:      &now,
			LastHeartbeatAt: &now,
			DoorOpen:        p.DoorOpen,
			SignalStrength:  p.SignalStrength,
			BatteryLevel:    p.BatteryLevel,
			Weight:          kgDecimal,
			Payload:         p.Payload,
		}); err != nil {
			return err
		}

		// 5) (可选) 落 heartbeat 表：你如果 heartbeat.sql 已建，我下一步可以补
		return nil
	})
	if err != nil {
		return HeartbeatResult{}, err
	}
	return res, nil
}
