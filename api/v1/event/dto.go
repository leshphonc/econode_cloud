package event

import "github.com/google/uuid"

type ReportRequest struct {
	Type         int16          `json:"type" binding:"required,oneof=1 2 3 4 5 6"`
	Action       int16          `json:"action" binding:"required,oneof=1 2 3"`
	Code         string         `json:"code" binding:"required"`
	Severity     int16          `json:"severity" binding:"required,oneof=1 2 3 4 5"`
	EventUID     *uuid.UUID     `json:"event_uid"`
	OccurredAtMs *int64         `json:"occurred_at_ms"`
	ReportedAtMs *int64         `json:"reported_at_ms"`
	Meta         map[string]any `json:"meta"`
}
