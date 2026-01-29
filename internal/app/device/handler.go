package device

import (
	"econode-cloud/api/v1/device"
	"econode-cloud/internal/app/device/ctxdev"
	"econode-cloud/internal/infra/http/resp"
	"econode-cloud/internal/pkg/bizerr"
	"econode-cloud/internal/pkg/ctxx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req device.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, bizerr.ErrParamInvalid)
		return
	}

	dev, err := h.deviceService.Register(c.Request.Context(), RegisterParams{
		SerialNo: req.SerialNo,
		Meta:     req.Meta,
	})
	if err != nil {
		ctxx.Logger(c.Request.Context()).Error(err.Error())
		resp.Fail(c, err)
		return
	}

	resp.OK(c, device.RegisterResponse{
		SerialNo:  dev.SerialNo,
		ClaimCode: dev.ClaimCode,
	})
}

func (h *Handler) Activate(c *gin.Context) {
	var req device.ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, bizerr.ErrParamInvalid)
		return
	}

	dev, err := h.deviceService.Activate(c.Request.Context(), ClaimParams{
		SerialNo:  req.SerialNo,
		Model:     req.Model,
		PowerMode: req.PowerMode,
		HWVersion: req.HWVersion,
		FWVersion: req.FWVersion,
		ClaimCode: req.ClaimCode,
		Meta:      req.Meta,
	})
	if err != nil {
		ctxx.Logger(c.Request.Context()).Error(err.Error())
		resp.Fail(c, err)
		return
	}

	resp.OK(c, device.ActivateResponse{
		DeviceUID:    dev.DeviceUID,
		Name:         dev.Name,
		Model:        dev.Model,
		Status:       dev.Status,
		PowerMode:    dev.PowerMode,
		HWVersion:    dev.HWVersion,
		FWVersion:    dev.FWVersion,
		ClaimAt:      dev.ClaimAt,
		ActiveErrors: dev.ActiveErrors,
		Meta:         dev.Meta,
	})
}

func (h *Handler) Heartbeat(c *gin.Context) {
	var req device.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, bizerr.ErrParamInvalid)
		return
	}

	err := h.deviceService.Heartbeat(c.Request.Context(), HeartbeatParams{
		DeviceID:     ctxdev.DeviceID(c.Request.Context()),
		ReportedAtMs: req.ReportedAtMs,
		Meta:         req.Meta,
	})
	if err != nil {
		ctxx.Logger(c.Request.Context()).Error(err.Error())
		resp.Fail(c, err)
		return
	}

	resp.OK(c, nil)
}
