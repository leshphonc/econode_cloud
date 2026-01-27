package device

import (
	"econode-cloud/api/v1/device"
	"econode-cloud/internal/infra/http/resp"
	"econode-cloud/internal/pkg/bizerr"
	"econode-cloud/internal/pkg/ctxx"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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

	dev, err := h.deviceService.Activate(c.Request.Context(), ActivateParams{
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
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}
	if req.SerialNo == "" {
		c.JSON(400, gin.H{"error": "serial_no required"})
		return
	}
	if req.Payload == nil {
		req.Payload = datatypes.JSONMap{}
	}

	// transport -> service params 映射
	_, err := h.deviceService.Heartbeat(c.Request.Context(), HeartbeatParams{
		SerialNo:       req.SerialNo,
		DoorOpen:       req.DoorOpen,
		SignalStrength: req.SignalStrength,
		BatteryLevel:   req.BatteryLevel,
		Weight:         req.Weight, // 先用 string，service 再 parse
		Payload:        req.Payload,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, device.HeartbeatResponse{ServerTime: time.Now().Unix()})
}
