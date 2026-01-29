package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// EventType: 1-状态 2-错误 3-告警 4-投递 5-维护 6-调试
type EventType int16

const (
	EventTypeStatus EventType = 1
	EventTypeError  EventType = 2
	EventTypeAlert  EventType = 3
	EventTypeDrop   EventType = 4
	EventTypeMaint  EventType = 5
	EventTypeDebug  EventType = 6
)

// EventAction: 1-发生/抛出(raise) 2-恢复/修复(fix) 3-更新(update)
type EventAction int16

const (
	EventActionRaise  EventAction = 1
	EventActionFix    EventAction = 2
	EventActionUpdate EventAction = 3
)

// EventSeverity: 1-debug 2-info 3-warn 4-error 5-critical
type EventSeverity int16

const (
	EventSeverityDebug    EventSeverity = 1
	EventSeverityInfo     EventSeverity = 2
	EventSeverityWarn     EventSeverity = 3
	EventSeverityError    EventSeverity = 4
	EventSeverityCritical EventSeverity = 5
)

type Event struct {
	ID       int64 `gorm:"primaryKey;column:id"`
	DeviceID int64 `gorm:"not null;column:device_id;index:event_device_created_at_idx,priority:1;index:event_device_code_created_at_idx,priority:1"`

	Type     int16  `gorm:"not null;column:type"`
	Action   int16  `gorm:"not null;column:action"`
	Code     string `gorm:"type:text;not null;column:code;index:event_device_code_created_at_idx,priority:2"`
	Severity int16  `gorm:"not null;column:severity"`

	// Optional idempotency key. Uniqueness enforced by a partial unique index:
	// (device_id, event_uid) where event_uid is not null
	EventUID *uuid.UUID `gorm:"type:uuid;column:event_uid;index:event_device_event_uid_uk,unique,priority:2"`

	Meta datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}';column:meta"`

	ReportedAt *time.Time `gorm:"type:timestamptz;column:reported_at"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now();column:created_at;index:event_device_created_at_idx,priority:2;index:event_device_code_created_at_idx,priority:3;index:event_created_at_idx"`

	// 注意：event 表是 append-only，所以没有 UpdatedAt 是合理的
}

func (Event) TableName() string { return "event" }
