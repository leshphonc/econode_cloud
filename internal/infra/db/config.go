package db

import "time"

type Config struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`

	// 可选：gorm 相关
	SlowThreshold time.Duration `mapstructure:"slow_threshold"` // 慢查询阈值
	LogLevel      string        `mapstructure:"log_level"`      // silent/error/warn/info
}
