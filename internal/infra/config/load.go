package config

import (
	conf "econode-cloud/internal/config"
	"fmt"

	"github.com/spf13/viper"
)

// Load 读取配置：默认 configs/config.yaml，可被环境变量覆盖。
// 约定：ENV 前缀 ECONODE_，嵌套用 __ 分隔，例如：ECONODE_DB__DSN
func Load() (*conf.Config, error) {
	v := viper.New()

	// 1) 默认值
	v.SetDefault("env", "dev")
	v.SetDefault("server.port", "8080")

	// 2) 配置文件
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".") // 兼容从项目根目录启动

	// 3) 读取配置
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// 4) 反序列化
	var cfg conf.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 5) 校验
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
