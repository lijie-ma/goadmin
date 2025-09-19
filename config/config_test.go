package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// 加载配置文件
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 测试基本配置
	if cfg.App.Name != "goadmin" {
		t.Errorf("期望应用名称为 'goadmin', 实际为 '%s'", cfg.App.Name)
	}
	if cfg.App.Version != "1.0.0" {
		t.Errorf("期望版本号为 '1.0.0', 实际为 '%s'", cfg.App.Version)
	}

	// 测试数据库配置
	if cfg.Database.Master.Host != "localhost" {
		t.Errorf("期望数据库主机为 'localhost', 实际为 '%s'", cfg.Database.Master.Host)
	}
	if cfg.Database.Master.Port != 3306 {
		t.Errorf("期望数据库端口为 3306, 实际为 %d", cfg.Database.Master.Port)
	}
	if cfg.Database.Master.ConnMaxLifetime != time.Hour {
		t.Errorf("期望连接最大生命周期为 1h, 实际为 %v", cfg.Database.Master.ConnMaxLifetime)
	}

	// 测试从库配置
	if len(cfg.Database.Slaves) != 1 {
		t.Errorf("期望从库数量为 1, 实际为 %d", len(cfg.Database.Slaves))
	} else {
		if cfg.Database.Slaves[0].Port != 3307 {
			t.Errorf("期望从库端口为 3307, 实际为 %d", cfg.Database.Slaves[0].Port)
		}
	}

	// 测试Redis配置
	if cfg.Redis.Host != "localhost" {
		t.Errorf("期望Redis主机为 'localhost', 实际为 '%s'", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("期望Redis端口为 6379, 实际为 %d", cfg.Redis.Port)
	}
	if cfg.Redis.DialTimeout != 5*time.Second {
		t.Errorf("期望拨号超时为 5s, 实际为 %v", cfg.Redis.DialTimeout)
	}

	// 测试日志配置
	if cfg.Logger.Level != "info" {
		t.Errorf("期望日志级别为 'info', 实际为 '%s'", cfg.Logger.Level)
	}
	if cfg.Logger.Filename != "./logs/app.log" {
		t.Errorf("期望日志文件为 './logs/app.log', 实际为 '%s'", cfg.Logger.Filename)
	}
}
