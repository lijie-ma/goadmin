package server

import (
	"context"
	"fmt"
	"goadmin/config"
	"goadmin/internal/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebServer struct {
	httpServer *http.Server
}

// RegisterServices 用于注册路由需要的 services
type RegisterServices = api.Services

func NewWebServer(cfg *config.Config, r *gin.Engine, services RegisterServices) *WebServer {
	// 设置Gin模式
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 配置路由
	api.RegisterRouter(r, services)

	return &WebServer{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.App.Port),
			Handler: r,
		},
	}
}

func (s *WebServer) Name() string {
	return "WebServer"
}

func (s *WebServer) Start(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		log.Println("[Web] starting on", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (s *WebServer) Stop(ctx context.Context) error {
	log.Println("[Web] shutting down...")
	return s.httpServer.Shutdown(ctx)
}