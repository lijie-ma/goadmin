package cmd

import (
	"context"
	"fmt"
	"goadmin/internal/api"
	"goadmin/pkg/db"
	"goadmin/pkg/logger"
	"goadmin/pkg/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
	err = redis.Init(&cfg.Redis)
	if err != nil {
		return err
	}

	// 设置Gin模式
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin实例
	r := gin.Default()

	// 配置路由
	api.RegisterRouter(r)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		fmt.Printf("服务器已启动，监听端口: %d\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务器启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("服务器关闭错误: %v\n", err)
		return err
	}

	fmt.Println("服务器已关闭")
	return nil
}
