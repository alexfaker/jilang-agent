package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alexfaker/jilang-agent/config"
	_ "github.com/go-sql-driver/mysql"
)

// DB 数据库连接实例封装
type DB struct {
	*sql.DB
}

// ConnectDB 连接到数据库
func ConnectDB(cfg config.DatabaseConfig) (*DB, error) {
	var dsn string

	switch cfg.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 验证连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库 Ping 失败: %v", err)
	}

	return &DB{db}, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	return db.DB.Close()
}

// Transaction 开始一个事务并执行回调函数
func (db *DB) Transaction(txFunc func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出 panic
		}
	}()

	if err := txFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
