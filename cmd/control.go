package cmd

import (
	"goadmin/internal/wire"

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
	// 使用 Wire 初始化应用
	app, err := wire.InitializeApp()
	if err != nil {
		return err
	}

	// 运行服务管理器
	app.ServiceManager.Run()
	return nil
}
