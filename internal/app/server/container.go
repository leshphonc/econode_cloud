package server

import (
	"econode-cloud/internal/app/device"
	"econode-cloud/internal/pkg/txm"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handlers struct {
	Device *device.Handler
}

type Middlewares struct {
	AuthDevice gin.HandlerFunc
}

type Container struct {
	Handlers   *Handlers
	Middleware *Middlewares
}

func BuildContainer(db *gorm.DB, rds *redis.Client) *Container {
	// 1) Repositories（数据访问）
	deviceRepo := device.NewDeviceRepo(db)

	// 2) Services（业务用例）
	tx := txm.NewTxManager(db)
	deviceSvc := device.NewService(tx, deviceRepo)

	// 3) Handlers（HTTP 层：bind DTO -> call service -> response）
	deviceHandler := device.NewHandler(deviceSvc)

	return &Container{
		Handlers: &Handlers{
			Device: deviceHandler,
		},
		Middleware: &Middlewares{
			AuthDevice: device.AuthDevice(deviceSvc),
		},
	}
}
