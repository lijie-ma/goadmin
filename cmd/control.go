package cmd

import (
	"goadmin/cmd/server"
	"goadmin/internal/i18n"
	"goadmin/pkg/db"
	"goadmin/pkg/logger"
	"goadmin/pkg/redisx"
	"goadmin/pkg/task"

	"github.com/spf13/cobra"
)

var (
	controlCmd = &cobra.Command{
		Use:   "control",
		Short: "控制API服务器",
		Long:  `控制API服务器的启动、停止等操作`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer()
		},
	}
)

// runServer 启动HTTP服务器
func runServer() error {
	// 配置全局日志实例
	logger.SetGlobal(logger.New(
		logger.WithConfig(&cfg.Logger),
	))
	err := db.Init(&cfg.Database)
	if err != nil {
		return err
	}
	err = redisx.Init(&cfg.Redis)
	if err != nil {
		return err
	}
	i18n.Init()

	services := task.NewServiceManager()
	services.AddService(
		server.NewCronManager(),
		server.NewWebServer(cfg),
		server.NewHookServer(),
	)

	services.Run()
	return nil
}
