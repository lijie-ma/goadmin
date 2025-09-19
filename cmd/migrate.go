package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var (
	migrationsDir string
	driverName    string
	databaseDSN   string

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "数据库迁移命令",
		Long:  `执行数据库迁移操作，包括创建、升级、回滚等`,
	}

	upCmd = &cobra.Command{
		Use:   "up",
		Short: "执行所有待执行的迁移",
		Long:  `执行所有尚未应用的数据库迁移`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("up", 0)
		},
	}

	upByOneCmd = &cobra.Command{
		Use:   "up-by-one",
		Short: "执行一个迁移",
		Long:  `执行下一个尚未应用的数据库迁移`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("up-by-one", 0)
		},
	}

	downCmd = &cobra.Command{
		Use:   "down",
		Short: "回滚最后一个迁移",
		Long:  `回滚最后应用的数据库迁移`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("down", 0)
		},
	}

	resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "回滚所有迁移",
		Long:  `回滚所有已应用的数据库迁移`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("reset", 0)
		},
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "显示迁移状态",
		Long:  `显示所有迁移的当前状态`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("status", 0)
		},
	}

	createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "创建新的迁移文件",
		Long:  `创建一个新的空白SQL迁移文件`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return createMigration(name, false)
		},
	}

	createGoCmd = &cobra.Command{
		Use:   "create-go [name]",
		Short: "创建新的Go迁移文件",
		Long:  `创建一个新的Go迁移文件`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return createMigration(name, true)
		},
	}
)

func init() {
	// 获取当前执行路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前路径失败: %v", err)
	}

	// 设置默认迁移目录
	defaultMigrationsDir := filepath.Join(currentDir, "migrations")

	// 添加命令行参数
	migrateCmd.PersistentFlags().StringVarP(&migrationsDir, "dir", "d", defaultMigrationsDir, "迁移文件目录")
	migrateCmd.PersistentFlags().StringVarP(&driverName, "driver", "r", "mysql", "数据库驱动名称")
	migrateCmd.PersistentFlags().StringVarP(&databaseDSN, "dsn", "s", "", "数据库连接字符串")

	// 添加子命令
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(upByOneCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(resetCmd)
	migrateCmd.AddCommand(statusCmd)
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(createGoCmd)
}

// runMigrate 执行迁移操作
func runMigrate(command string, steps int64) error {
	// 确保迁移目录存在
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("创建迁移目录失败: %w", err)
	}

	// 如果没有提供DSN，则从配置中获取
	if databaseDSN == "" && cfg != nil {
		master := cfg.Database.Master
		databaseDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			master.Username, master.Password, master.Host, master.Port, master.Database)
	}

	if databaseDSN == "" {
		return fmt.Errorf("未提供数据库连接字符串，请通过--dsn参数指定或在配置文件中配置")
	}

	// 连接数据库
	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 设置迁移目录
	goose.SetBaseFS(nil)
	if err := goose.SetDialect(driverName); err != nil {
		return fmt.Errorf("设置数据库方言失败: %w", err)
	}

	// 执行迁移命令
	switch command {
	case "up":
		err = goose.Up(db, migrationsDir)
	case "up-by-one":
		err = goose.UpByOne(db, migrationsDir)
	case "down":
		err = goose.Down(db, migrationsDir)
	case "reset":
		err = goose.Reset(db, migrationsDir)
	case "status":
		err = goose.Status(db, migrationsDir)
	default:
		return fmt.Errorf("未知的迁移命令: %s", command)
	}

	if err != nil {
		return fmt.Errorf("迁移操作失败: %w", err)
	}

	return nil
}

// createMigration 创建迁移文件
func createMigration(name string, goMigration bool) error {
	// 确保迁移目录存在
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("创建迁移目录失败: %w", err)
	}

	// 如果没有提供DSN，则从配置中获取
	if databaseDSN == "" && cfg != nil {
		master := cfg.Database.Master
		databaseDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			master.Username, master.Password, master.Host, master.Port, master.Database)
	}

	// 连接数据库
	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	// 创建迁移文件
	fileType := "sql"
	if goMigration {
		fileType = "go"
	}

	if err := goose.Create(db, migrationsDir, name, fileType); err != nil {
		return fmt.Errorf("创建迁移文件失败: %w", err)
	}

	fmt.Printf("迁移文件创建成功: %s\n", name)
	return nil
}
