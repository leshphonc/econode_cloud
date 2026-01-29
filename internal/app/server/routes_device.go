package server

import (
	"econode-cloud/internal/app/device"
	"econode-cloud/internal/app/event"

	"github.com/gin-gonic/gin"
)

func RegisterDeviceRoutes(api *gin.RouterGroup, devH *device.Handler, evtH *event.Handler, auth gin.HandlerFunc) {

	dev := api.Group("/device")
	{
		// 注册
		dev.POST("/register", devH.Register)
		dev.POST("/activate", devH.Activate)
	}

	authDev := dev.Use(auth)
	{
		// 心跳
		authDev.POST("/heartbeat", devH.Heartbeat)

		// 事件
		authDev.POST("/event", evtH.Report)

		// 命令
		//authDev.GET("/commands", h.PullCommands)
		//authDev.POST("/commands/:command_id/ack", h.AckCommand)
	}

	// 资源面（后台/管理端/查询用）
	//devices := api.Group("/devices")
	{
		// 例如：后台查某设备心跳历史、事件历史、命令历史等
		// devices.GET("/:device_id/heartbeats", h.ListHeartbeats)
		// devices.GET("/:device_id/events", h.ListEvents)
		// devices.GET("/:device_id/commands", h.ListCommands)
	}
}
