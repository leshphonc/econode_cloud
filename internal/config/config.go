package config

import (
	"econode-cloud/internal/infra/db"
	"econode-cloud/internal/infra/log"
	"econode-cloud/internal/infra/redis"
	"errors"
	"fmt"
	"time"
)

type Config struct {
	Env    string       `mapstructure:"env"` // dev/prod
	Server ServerConfig `mapstructure:"server"`

	DB    db.Config    `mapstructure:"db"`
	Redis redis.Config `mapstructure:"redis"`
	Log   log.Config   `mapstructure:"log"`

	// Consul consul.Config `mapstructure:"consul"` // 以后加
	// Device   DeviceConfig  `mapstructure:"device"`
	// Features FeatureConfig `mapstructure:"features"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"` // "8080"
}

type DeviceConfig struct {
	OfflineAfter time.Duration `mapstructure:"offline_after"`
}
type FeatureConfig struct {
	EnableMQTT bool `mapstructure:"enable_mqtt"`
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return errors.New("server.port is empty")
	}

	// db / redis / log 这些你可以按阶段逐步加严
	if c.DB.DSN == "" {
		return errors.New("db.dsn is empty")
	}
	if c.Redis.Addr == "" {
		// 如果你不是强依赖 Redis，可以先不校验，或给默认值
		return errors.New("redis.addr is empty")
	}

	if c.Log.Level == "" {
		// 给一个更友好的提示
		return fmt.Errorf("log.level is empty (suggest: debug/info/warn/error)")
	}
	return nil
}
