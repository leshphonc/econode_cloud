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
	AuthByDeviceUID(ctx context.Context, deviceUID string) (IdentityResult, error)
}

type IdentityResult struct {
	DeviceID  int64  // 内部 ID，后续 service/repo 用它更快
	DeviceUID string // 外部 UUID（原 public_id）
}

func (s *Service) AuthByDeviceUID(ctx context.Context, deviceUID string) (IdentityResult, error) {
	dev, err := s.deviceRepo.GetByDeviceUID(ctx, deviceUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return IdentityResult{}, ErrDeviceNotFound
		}
		return IdentityResult{}, err
	}

	if dev.DisabledAt != nil || dev.RetiredAt != nil {
		return IdentityResult{}, ErrDeviceDisabled
	}

	return IdentityResult{
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

type HeartbeatParams struct {
	DeviceID     int64
	ReportedAtMs *int64
	Meta         map[string]any
}

func (s *Service) Heartbeat(ctx context.Context, p HeartbeatParams) error {

	err := s.txm.Transaction(ctx, func(tx *gorm.DB) error {
		devRepo := s.deviceRepo.WithDB(tx)

		// 1) 落 heartbeat 表
		now := time.Now()
		var reportedAt *time.Time
		if p.ReportedAtMs != nil {
			t := time.UnixMilli(*p.ReportedAtMs)
			// 可选：做时间合理性校验，不合理就丢弃
			if t.After(now.Add(5*time.Minute)) || t.Before(now.Add(-365*24*time.Hour)) {
				reportedAt = nil
			} else {
				reportedAt = &t
			}
			reportedAt = &t
		}
		hb := model.Heartbeat{
			DeviceID:   p.DeviceID,
			ReportedAt: reportedAt,
			Meta:       p.Meta,
		}
		err := devRepo.InsertHeartbeat(ctx, &hb)
		if err != nil {
			return err
		}

		// 2) 更新 device.last_seen_at
		if err = devRepo.TouchLastSeen(ctx, p.DeviceID, &hb.CreatedAt); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
