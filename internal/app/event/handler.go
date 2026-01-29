package event

import (
	"econode-cloud/api/v1/event"
	"econode-cloud/internal/app/device/ctxdev"
	"econode-cloud/internal/infra/http/resp"
	"econode-cloud/internal/pkg/bizerr"
	"econode-cloud/internal/pkg/ctxx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	eventService *Service
}

func NewHandler(es *Service) *Handler {
	return &Handler{
		eventService: es,
	}
}

func (h *Handler) Report(c *gin.Context) {
	var req event.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, bizerr.ErrParamInvalid)
		return
	}

	err := h.eventService.Report(c.Request.Context(), ReportParams{
		DeviceID:     ctxdev.DeviceID(c.Request.Context()),
		Type:         req.Type,
		Action:       req.Action,
		Code:         req.Code,
		Severity:     req.Severity,
		EventUID:     req.EventUID,
		ReportedAtMs: req.ReportedAtMs,
		Meta:         req.Meta,
	})
	if err != nil {
		ctxx.Logger(c.Request.Context()).Error(err.Error())
		resp.Fail(c, err)
		return
	}

	resp.OK(c, event.ReportResponse{})
}
