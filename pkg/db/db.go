package db

import (
	"fmt"
	"goadmin/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// DB 全局数据库实例
var DB *gorm.DB

// 初始化数据库连接
func Init(cfg *config.Config) error {
	var err error

	// 初始化主库连接
	DB, err = initDB(cfg.Database.Master)
	if err != nil {
		return fmt.Errorf("初始化主库失败: %w", err)
	}

	// 如果配置了从库，添加数据库解析器
	if len(cfg.Database.Slaves) > 0 {
		resolverCfg := dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(buildDSN(cfg.Database.Master))},
			Replicas: make([]gorm.Dialector, len(cfg.Database.Slaves)),
			Policy:   dbresolver.RandomPolicy{}, // 随机策略
		}

		// 配置所有从库
		for i, slave := range cfg.Database.Slaves {
			resolverCfg.Replicas[i] = mysql.Open(buildDSN(slave))
		}

		// 注册数据库解析器
		err = DB.Use(dbresolver.Register(resolverCfg).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(cfg.Database.Master.ConnMaxLifetime).
			SetMaxIdleConns(cfg.Database.Master.MaxIdleConns).
			SetMaxOpenConns(cfg.Database.Master.MaxOpenConns))
		if err != nil {
			return fmt.Errorf("配置主从失败: %w", err)
		}
	}

	return nil
}

// 初始化单个数据库连接
func initDB(dbCfg config.DBConfig) (*gorm.DB, error) {
	// 配置GORM日志
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,                   // 慢SQL阈值
			LogLevel:                  parseLogLevel(dbCfg.LogLevel), // 日志级别
			IgnoreRecordNotFoundError: true,                          // 忽略记录未找到错误
			Colorful:                  true,                          // 彩色输出
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(buildDSN(dbCfg)), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// 获取底层的sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)

	return db, nil
}

// 构建数据库连接DSN
func buildDSN(cfg config.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)
}

// 解析日志级别
func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info // 默认info级别
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
