package redis

import "time"

type Config struct {
	Addr         string        `mapstructure:"addr"`          // "127.0.0.1:6379"
	Username     string        `mapstructure:"username"`      // 可空
	Password     string        `mapstructure:"password"`      // 可空
	DB           int           `mapstructure:"db"`            // 默认 0
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`  // 默认 2s
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 默认 2s
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // 默认 2s
	PoolSize     int           `mapstructure:"pool_size"`     // 默认 20
	MinIdleConns int           `mapstructure:"min_idle_conns"`
}
