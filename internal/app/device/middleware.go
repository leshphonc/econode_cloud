package device

import (
	"econode-cloud/internal/app/device/ctxdev"
	"econode-cloud/internal/infra/http/resp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const headerDeviceUID ctxKey = "X-Device-UID"

func AuthDevice(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 取设备UID，校验格式
		deviceUID := strings.TrimSpace(c.GetHeader(string(headerDeviceUID)))
		if deviceUID == "" {
			c.Abort()
			resp.Fail(c, ErrDeviceInvalidXUID)
			return
		}
		if _, err := uuid.Parse(deviceUID); err != nil {
			c.Abort()
			resp.Fail(c, ErrDeviceInvalidXUID)
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
		ctx := ctxdev.WithDeviceID(c.Request.Context(), ident.DeviceID)
		ctx = ctxdev.WithDeviceUID(ctx, ident.DeviceUID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
