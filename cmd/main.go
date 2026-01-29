package main

import (
	"context"
	"econode-cloud/internal/app/server"
	"econode-cloud/internal/infra/config"
	"econode-cloud/internal/infra/db"
	"econode-cloud/internal/infra/log"
	"econode-cloud/internal/infra/redis"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println(time.Now().UnixMilli())

	// 加载配置文件
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 创建 logger
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 连接 pg
	pg, err := db.NewPostgres(cfg.DB, logger)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	// 连接 redis
	rds, err := redis.New(cfg.Redis, logger)
	if err != nil {
		panic(err)
	}
	defer rds.Close()

	// 组装 repo/service/handler
	app := server.BuildContainer(pg.Gorm, rds.Rdb)

	// 启动服务
	srv := server.New(cfg.Server.Port, app, logger)
	go func() {
		if err = srv.Run(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
