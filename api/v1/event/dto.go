package event

import "github.com/google/uuid"

type ReportRequest struct {
	Type         int16          `json:"type" binding:"required,oneof=1 2 3 4 5 6"`
	Action       int16          `json:"action" binding:"required,oneof=1 2 3"`
	Code         string         `json:"code" binding:"required"`
	Severity     int16          `json:"severity" binding:"required,oneof=1 2 3 4 5"`
	EventUID     *uuid.UUID     `json:"event_uid"`
	ReportedAtMs *int64         `json:"reported_at_ms"`
	Meta         map[string]any `json:"meta"`
}
type ReportResponse struct {
}

// 典型事件案例
// 1.状态（type=1）
// 开门
//{
//	"type": 1,
//	"action": 1,
//	"code": "DOOR_OPEN",
//	"severity": 2,
//	"event_uid": "....",
//	"meta": { "slot_no": 3, "method": "screen" }
//}

// 状态快照
//{
//	"type": 1,
//	"action": 3,
//	"code": "STATE_SNAPSHOT",
//	"severity": 2,
//	"event_uid": "....",
//	"meta": {
//		"rssi": -63,
//		"temp_c": 28.4,
//		"door_open": false,
//		"slots": [{ "slot_no": 1, "full": false }]
//	}
//}

// 2.错误（type=2，仅 action=1/2）
// 故障发生
//{
//	"type": 2,
//	"action": 1,
//	"code": "TOF_TIMEOUT",
//	"severity": 4,
//	"event_uid": "....",
//	"meta": { "slot_no": 2, "timeout_ms": 1500, "retry": 1 }
//}

// 故障恢复
//{
//	"type": 2,
//	"action": 2,
//	"code": "TOF_TIMEOUT",
//	"severity": 2,
//	"event_uid": "....",
//	"meta": { "slot_no": 2, "recovered_by": "auto" }
//}

// 3.告警（type=3，仅 action=1/2）
// 满溢告警
//{
//	"type": 3,
//	"action": 1,
//	"code": "BIN_FULL",
//	"severity": 5,
//	"event_uid": "....",
//	"meta": { "slot_no": 1, "tof_mm": 32, "threshold_mm": 40 }
//}

// 4.投递（type=4）
// 投递完成
//{
//	"type": 4,
//	"action": 1,
//	"code": "DROP_COMPLETED",
//	"severity": 2,
//	"event_uid": "....",
//	"meta": {
//		"slot_no": 1,
//		"session_id": "S-20260128-0001",
//		"weight_g": 235,
//		"baseline_before_g": 12030,
//		"baseline_after_g": 12265
//	}
//}

// 5.维护（type=5）
// 清运完成
//{
//	"type": 5,
//	"action": 1,
//	"code": "MAINT_EMPTY_COMPLETED",
//	"severity": 2,
//	"event_uid": "....",
//	"meta": { "operator_id": "u123", "slot_no": 1, "emptied_weight_g": 5200 }
//}

// 6.调试（type=6，仅 action=3）
// 传感器采样
//{
//	"type": 6,
//	"action": 3,
//	"code": "SENSOR_SAMPLE",
//	"severity": 1,
//	"event_uid": "....",
//	"meta": { "slot_no": 2, "tof_mm": 36, "hx711_raw": 833120, "temp_c": 29.1 }
//}
