package model

import (
	"time"

	"gorm.io/datatypes"
)

type Heartbeat struct {
	ID         int64             `gorm:"primaryKey"`
	DeviceID   int64             `gorm:"not null;index:heartbeat_device_created_at_idx,priority:1"`
	ReportedAt *time.Time        `gorm:"type:timestamptz"`
	Meta       datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt  time.Time         `gorm:"type:timestamptz;not null;default:now();index:heartbeat_device_created_at_idx,priority:2,sort:desc"`
}

func (Heartbeat) TableName() string { return "heartbeat" }
