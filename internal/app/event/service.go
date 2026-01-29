package event

import (
	"context"
	"econode-cloud/internal/model"
	"econode-cloud/internal/pkg/timeutil"
	"econode-cloud/internal/pkg/txm"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	txm       *txm.TxManager
	eventRepo Repo
}

func NewService(txm *txm.TxManager, deviceRepo Repo) *Service {
	return &Service{
		txm,
		deviceRepo,
	}
}

type ReportParams struct {
	DeviceID     int64
	Type         int16
	Action       int16
	Code         string
	Severity     int16
	EventUID     *uuid.UUID
	ReportedAtMs *int64
	Meta         map[string]any
}

func (s *Service) Report(ctx context.Context, p ReportParams) error {
	err := validateEventTypeAction(p.Type, p.Action)
	if err != nil {
		return err
	}
	err = validateEventSeverity(p.Type, p.Severity)
	if err != nil {
		return err
	}
	reportedAt := timeutil.NormalizeReportedAt(p.ReportedAtMs, time.Now(), 5*time.Minute, 365*24*time.Hour)
	err = s.eventRepo.Create(ctx, &model.Event{
		DeviceID:   p.DeviceID,
		Type:       p.Type,
		Action:     p.Action,
		Code:       p.Code,
		Severity:   p.Severity,
		EventUID:   p.EventUID,
		Meta:       p.Meta,
		ReportedAt: reportedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
