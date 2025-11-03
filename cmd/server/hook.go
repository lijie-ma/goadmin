// 服务启动需要初始的数据
// 跟随 web 一起运行， 在进程启动的后 运行一次 并退出
package server

import (
	"context"
	cusCtx "goadmin/internal/context"

	"golang.org/x/sync/errgroup"
)

type HookFunc func(ctx *cusCtx.CliContext) error

type HookServer struct {
	hooks []HookFunc
}

func NewHookServer() *HookServer {
	return &HookServer{}
}

func (s *HookServer) register() error {
	s.hooks = append(s.hooks,
		func(ctx *cusCtx.CliContext) error {
			ctx.Logger.Info("demo")
			return nil
		},
	)
	return nil
}

func (s *HookServer) Name() string {
	return "HookServer"
}

func (s *HookServer) Start(ctx context.Context) error {
	s.register()
	var g errgroup.Group
	for _, hook := range s.hooks {
		g.Go(func() error {
			cliCtx := cusCtx.NewCliContext(ctx)
			defer cliCtx.Close()
			return hook(cliCtx)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *HookServer) Stop(ctx context.Context) error {
	return nil
}
