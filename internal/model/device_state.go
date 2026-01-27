package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type DeviceState struct {
	// 一设备一行
	DeviceID int64 `gorm:"primaryKey;type:bigint;"`

	// 最近一次“有效活动时间”
	LastSeenAt      *time.Time `gorm:"type:timestamptz;index:device_state_last_seen_at_idx"`
	LastHeartbeatAt *time.Time `gorm:"type:timestamptz"` // 最近一次心跳记录（可选，用于快速关联）
	LastEventAt     *time.Time `gorm:"type:timestamptz"`

	// 当前关键状态
	DoorOpen       *bool  `gorm:"type:boolean"`
	SignalStrength *int16 `gorm:"type:smallint"`
	BatteryLevel   *int16 `gorm:"type:smallint"`

	// 当前称重相关（温漂 baseline 慢跟随）
	Weight   *decimal.Decimal `gorm:"type:numeric(12,3)"`
	Baseline *decimal.Decimal `gorm:"type:numeric(12,3)"`

	// 最近异常（可选）
	LastErrorCode *string    `gorm:"type:text"`
	LastErrorAt   *time.Time `gorm:"type:timestamptz"`

	// 扩展状态（兜底）
	Payload datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`

	// 审计字段（GORM 维护）
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (DeviceState) TableName() string {
	return "device_state"
}
