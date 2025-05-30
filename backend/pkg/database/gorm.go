package database

import (
	"fmt"
	"time"

	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectGormDB 连接到数据库并返回gorm.DB实例
func ConnectGormDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 配置GORM日志
	logLevel := logger.Silent
	if cfg.LogLevel == "info" {
		logLevel = logger.Info
	} else if cfg.LogLevel == "warn" {
		logLevel = logger.Warn
	} else if cfg.LogLevel == "error" {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// 根据驱动类型连接数据库
	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// 设置连接最大空闲时间
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接验证失败: %w", err)
	}

	return db, nil
}

// AutoMigrate 自动迁移数据库模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Workflow{},
		&models.WorkflowExecution{},
		&models.Agent{},
		&models.PointsTransaction{},
		&models.RechargeOrder{},
	)
}

// GormDB 封装GORM数据库实例（保持向后兼容）
type GormDB struct {
	DB *gorm.DB
}

// Close 关闭数据库连接
func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	return sqlDB.Close()
}
