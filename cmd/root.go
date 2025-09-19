package main

import (
	"fmt"
	"goadmin/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	configPath string
	cfg        *config.Config
	rootCmd    = &cobra.Command{
		Use:   "goadmin",
		Short: "GoAdmin 是一个后台管理系统",
		Long:  `GoAdmin 是基于 Go 语言开发的现代化后台管理系统，提供用户管理、权限管理等功能。`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 跳过migrate命令的配置加载
			if cmd.Name() == "migrate" && cmd.Parent().Name() == "goadmin" {
				return nil
			}

			// 加载配置
			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("加载配置文件失败: %w", err)
			}
			return nil
		},
	}
)

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 获取可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		execDir = "."
	} else {
		execDir = filepath.Dir(execDir)
	}

	// 设置默认配置文件路径
	defaultConfigPath := filepath.Join(execDir, "config", "config.yaml")

	// 添加命令行参数
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", defaultConfigPath, "配置文件路径")

	// 添加子命令
	rootCmd.AddCommand(controlCmd)
	rootCmd.AddCommand(migrateCmd)
}
