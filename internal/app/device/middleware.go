package device

import (
	"econode-cloud/internal/app/device/ctxx"
	"econode-cloud/internal/infra/http/resp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthDevice(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 取设备UID，校验格式
		deviceUID := ctxx.DeviceUID(c)
		if _, err := uuid.Parse(deviceUID); err != nil {
			c.Abort()
			resp.Fail(c, ErrDeviceInvalidXID)
			return
		}

		if svc == nil {
			c.Abort()
			resp.Fail(c, ErrDeviceSvcNotSet)
			return
		}

		ident, err := svc.AuthByDeviceUID(c.Request.Context(), deviceUID)
		if err != nil {
			c.Abort()
			resp.Fail(c, err)
			return
		}

		// 3) ctx 放入 device_identity
		if ident != nil {
			ctx := ctxx.WithDeviceUID(c.Request.Context(), ident.DeviceUID)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}
