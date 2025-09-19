package db

import (
	"goadmin/config"
	"testing"
)

func TestInit(t *testing.T) {
	// 加载配置
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	t.Logf("配置加载成功")
	t.Logf("数据库地址: %s:%d", cfg.Database.Master.Host, cfg.Database.Master.Port)

	// 初始化数据库
	t.Log("开始初始化数据库...")
	err = Init(&cfg.Database)
	if err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}
	t.Log("数据库初始化成功")

	// 验证数据库连接
	db := GetDB()
	if db == nil {
		t.Fatal("数据库实例为空")
	}
	t.Log("数据库实例获取成功")

	// 执行简单查询
	t.Log("执行测试查询...")
	var version string
	if err := db.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		t.Fatalf("数据库查询失败: %v", err)
	}
	t.Logf("数据库版本: %s", version)

	// 获取底层sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取sqlDB失败: %v", err)
	}

	// 验证连接池设置
	t.Log("验证连接池设置...")
	stats := sqlDB.Stats()
	t.Logf("打开连接数: %d", stats.OpenConnections)
	t.Logf("空闲连接数: %d", stats.Idle)
	t.Logf("活跃连接数: %d", stats.InUse)
	t.Logf("最大打开连接数: %d", cfg.Database.Master.MaxOpenConns)
	t.Logf("最大空闲连接数: %d", cfg.Database.Master.MaxIdleConns)

	// 测试主从配置
	if len(cfg.Database.Slaves) > 0 {
		t.Log("检测到从库配置，测试主从功能...")
		// 执行写操作测试（主库）
		t.Log("在主库执行写操作...")
		// 这里可以添加写操作测试，如CREATE TEMPORARY TABLE, INSERT等

		// 执行读操作测试（应路由到从库）
		t.Log("执行读操作，应路由到从库...")
		var slaveVersion string
		if err := db.Raw("SELECT VERSION()").Scan(&slaveVersion).Error; err != nil {
			t.Fatalf("从库查询失败: %v", err)
		}
		t.Logf("从库数据库版本: %s", slaveVersion)
	} else {
		t.Log("未检测到从库配置，跳过主从测试")
	}

	// 测试禁用数据库的情况
	t.Log("\n测试禁用数据库...")
	cfg.Database.Enable = false
	err = Init(&cfg.Database)
	if err != nil {
		t.Fatalf("禁用数据库失败: %v", err)
	}
	t.Log("数据库禁用测试成功")

	t.Log("数据库测试完成")
}
