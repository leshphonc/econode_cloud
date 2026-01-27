package server

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, app *Container) {

	api := r.Group("/api/v1")

	RegisterDeviceRoutes(api, app.Handlers.Device, app.Middleware.AuthDevice)
}
