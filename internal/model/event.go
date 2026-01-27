package model

import (
	"time"
)

type Event struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	DeviceID   int64          `gorm:"column:device_id;not null"`
	OccurredAt time.Time      `gorm:"column:occurred_at;not null"`
	EventType  int16          `gorm:"column:event_type;not null"` // event_type: 1=heartbeat_gap, 2=door_open, 3=door_close, 10=deposit, 11=weigh, 20=error, 21=warning, 30=maintenance
	Severity   int16          `gorm:"column:severity;not null"`   // severity: 1=info, 2=warn, 3=error, 4=critical
	Status     int16          `gorm:"column:status;not null"`     //status: 1=new, 2=acked, 3=resolved, 4=ignored
	TraceID    string         `gorm:"column:trace_id;type:uuid;not null"`
	Payload    map[string]any `gorm:"column:payload;type:jsonb;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (Event) TableName() string { return "event" }
