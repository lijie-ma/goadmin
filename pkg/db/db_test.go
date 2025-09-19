package db

import (
	"fmt"
	"goadmin/config"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// 打印当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前工作目录失败: %v", err)
	}
	t.Logf("当前工作目录: %s", pwd)

	configPath := "../../config/config.yaml"
	t.Logf("尝试加载配置文件: %s", configPath)

	// 检查配置文件是否存在
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Fatalf("配置文件不存在: %s", configPath)
	} else if err != nil {
		t.Fatalf("检查配置文件失败: %v", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	t.Logf("配置加载成功")
	t.Logf("数据库主机: %s:%d", cfg.Database.Master.Host, cfg.Database.Master.Port)

	// 初始化数据库
	t.Log("开始初始化数据库...")
	err = Init(cfg)
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

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取底层sqlDB失败: %v", err)
	}
	t.Log("获取底层sqlDB成功")

	// 尝试ping数据库
	t.Log("尝试Ping数据库...")
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库连接Ping失败: %v", err)
	}

	t.Log("数据库连接测试成功")
}

// 在main函数中手动运行测试，便于观察日志输出
func ExampleInit() {
	configPath := "../../config/config.yaml"
	fmt.Printf("尝试加载配置文件: %s\n", configPath)

	// 检查配置文件是否存在
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		fmt.Printf("配置文件不存在: %s\n", configPath)
		return
	} else if err != nil {
		fmt.Printf("检查配置文件失败: %v\n", err)
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}
	fmt.Println("配置加载成功")
	fmt.Printf("数据库主机: %s:%d\n", cfg.Database.Master.Host, cfg.Database.Master.Port)

	// 初始化数据库
	fmt.Println("开始初始化数据库...")
	err = Init(cfg)
	if err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}
	fmt.Println("数据库初始化成功")

	// 验证数据库连接
	db := GetDB()
	if db == nil {
		fmt.Println("数据库实例为空")
		return
	}
	fmt.Println("数据库实例获取成功")

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("获取底层sqlDB失败: %v\n", err)
		return
	}
	fmt.Println("获取底层sqlDB成功")

	// 尝试ping数据库
	fmt.Println("尝试Ping数据库...")
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("数据库连接Ping失败: %v\n", err)
		return
	}

	fmt.Println("数据库连接测试成功")
}
