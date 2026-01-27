package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// DeviceStatus 对应 device.status: 1=active, 2=inactive, 3=maintenance, 4=decommissioned
type DeviceStatus int16

const (
	DeviceStatusPreregistration DeviceStatus = 0
	DeviceStatusActive          DeviceStatus = 1
	DeviceStatusInactive        DeviceStatus = 2
	DeviceStatusMaintenance     DeviceStatus = 3
	DeviceStatusDecommissioned  DeviceStatus = 4
)

type DevicePowerMode int16

const (
	PowerModeGrid   DevicePowerMode = 1
	PowerModeSolar  DevicePowerMode = 2
	PowerModeHybrid DevicePowerMode = 3
)

type Device struct {
	// 内部主键
	ID int64 `gorm:"primaryKey;autoIncrement"`

	// 对外公开 ID（不可变）
	DeviceUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:device_device_uid_uk"`

	// 设备出厂序列号
	SerialNo string `gorm:"type:text;not null;uniqueIndex:device_serial_no_uk"`

	// 设备名称
	Name *string `gorm:"type:text;not null"`

	// 型号
	Model *string `gorm:"type:text;not null"`

	// 设备状态: 0-预注册 1-激活 2-禁用 3-退役
	Status int16 `gorm:"type:smallint;not null;check:status in (0,1,2,3)"`

	// 供电方式：1-市电 2-太阳能 3-混合
	PowerMode *int16 `gorm:"type:smallint;not null;check:power_mode"`

	// 硬件 / 固件版本
	HWVersion *string `gorm:"type:text"`
	FWVersion *string `gorm:"type:text"`

	// 认领相关
	ClaimCode string     `gorm:"type:text;not null"`
	ClaimedAt *time.Time `gorm:"type:timestamptz"`

	// 禁用信息
	DisabledAt     *time.Time `gorm:"type:timestamptz"`
	DisabledReason *string    `gorm:"type:text"`

	// 退役信息
	RetiredAt     *time.Time `gorm:"type:timestamptz"`
	RetiredReason *string    `gorm:"type:text"`

	// 最近状态
	LastSeenAt  *time.Time `gorm:"type:timestamptz"`
	LastErrorAt *time.Time `gorm:"type:timestamptz"`

	// 当前活跃错误列表
	ActiveErrors pq.StringArray `gorm:"type:text[];not null;default:'{}'"`

	// 扩展元数据
	Meta datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`

	// 时间戳
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName 显式指定表名
func (Device) TableName() string { return "device" }
