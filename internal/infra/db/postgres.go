package db

import (
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type DB struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func NewPostgres(cfg Config, log *zap.Logger) (*DB, error) {
	if cfg.DSN == "" {
		return nil, errors.New("db dsn is empty")
	}

	// gorm logger（不想引入太复杂的话先用 gorm 自带 logger）
	gormLogger := newGormLogger(cfg, log)

	gdb, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqldb, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	// 连接池参数（给默认值，避免你 config 没配就乱）
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 50
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = 30 * time.Minute
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = 5 * time.Minute
	}

	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqldb.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// 启动时做一次 ping，快速失败
	if err = sqldb.Ping(); err != nil {
		_ = sqldb.Close()
		return nil, err
	}

	log.Info("postgres initialized",
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)

	return &DB{Gorm: gdb, SQL: sqldb}, nil
}

func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}

func newGormLogger(cfg Config, log *zap.Logger) glogger.Interface {
	level := glogger.Warn
	switch cfg.LogLevel {
	case "silent":
		level = glogger.Silent
	case "error":
		level = glogger.Error
	case "warn":
		level = glogger.Warn
	case "info":
		level = glogger.Info
	}

	slow := cfg.SlowThreshold
	if slow <= 0 {
		slow = 300 * time.Millisecond
	}

	return glogger.New(
		zap.NewStdLog(log.With(zap.String("component", "gorm"))),
		glogger.Config{
			SlowThreshold:             slow,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
