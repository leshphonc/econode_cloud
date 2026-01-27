package device

import (
	"time"

	"gorm.io/datatypes"
)

// RegisterRequest 注册
type RegisterRequest struct {
	SerialNo string         `json:"serial_no" binding:"required"`
	Meta     map[string]any `json:"meta"`
}

type RegisterResponse struct {
	SerialNo  string `json:"serial_no"`
	ClaimCode string `json:"claim_code"`
}

type ActivateRequest struct {
	SerialNo  string         `json:"serial_no" binding:"required"`
	Model     string         `json:"model" binding:"required"`
	PowerMode int16          `json:"power_mode" binding:"required"`
	HWVersion string         `json:"hw_version" binding:"required"`
	FWVersion string         `json:"fw_version" binding:"required"`
	ClaimCode string         `json:"claim_code" binding:"required"`
	Meta      map[string]any `json:"meta"`
}

type ActivateResponse struct {
	DeviceUID    string         `json:"device_uid"`
	Name         string         `json:"name"`
	Model        string         `json:"model"`
	Status       int16          `json:"status"`
	PowerMode    int16          `json:"power_mode"`
	HWVersion    string         `json:"hw_version"`
	FWVersion    string         `json:"fw_version"`
	ClaimAt      int64          `json:"claim_at"`
	ActiveErrors []string       `json:"active_errors"`
	Meta         map[string]any `json:"meta"`
}

// HeartbeatRequest 心跳
type HeartbeatRequest struct {
	SerialNo string `json:"serial_no"`

	// 可选运行态
	DoorOpen       *bool  `json:"door_open"`
	SignalStrength *int16 `json:"signal_strength"`
	BatteryLevel   *int16 `json:"battery_level"`

	Weight *string `json:"weight"`

	// 扩展
	Payload datatypes.JSONMap `json:"payload"`
}

type HeartbeatResponse struct {
	ServerTime int64 `json:"server_time"`
}

// CreateEventRequest 事件上传
type CreateEventRequest struct {
	SerialNo   string         `json:"serial_no"`
	TraceID    string         `json:"trace_id" binding:"required,uuid"`
	OccurredAt time.Time      `json:"occurred_at" binding:"required"`
	EventType  int16          `json:"event_type" binding:"required"`
	Severity   int16          `json:"severity,omitempty"` // nil => default 1
	Payload    map[string]any `json:"payload,omitempty"`  // nil => {}
}

type CreateEventResponse struct {
	EventID    int64 `json:"event_id"`
	Idempotent bool  `json:"idempotent"` // true 表示重复上报命中已有记录
}
