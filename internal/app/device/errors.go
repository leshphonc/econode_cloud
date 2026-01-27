package device

import (
	"econode-cloud/internal/pkg/bizerr"
	"net/http"
)

var (
	ErrDeviceInvalidXID      = bizerr.NewBizError(20000, http.StatusUnauthorized, "invalid X-Device-Id (must be UUID)")
	ErrDeviceNotFound        = bizerr.NewBizError(20001, http.StatusUnauthorized, "设备不存在")
	ErrDeviceDisabled        = bizerr.NewBizError(20002, http.StatusForbidden, "设备已禁用")
	ErrDeviceSvcNotSet       = bizerr.NewBizError(20003, http.StatusInternalServerError, "device svc 未设置")
	ErrDeviceRegisterFailed  = bizerr.NewBizError(20004, http.StatusInternalServerError, "设备注册失败")
	ErrDeviceActivateFailed  = bizerr.NewBizError(20005, http.StatusInternalServerError, "设备不存在或已激活")
	ErrDevicePowerModeUnknow = bizerr.NewBizError(20006, http.StatusBadRequest, "未知供电方式")
)
