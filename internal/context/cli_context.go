package context

import (
	"context"
	"goadmin/pkg/logger"
	"goadmin/pkg/trace"
	"sync"
)

var (
	cliCtxPool = sync.Pool{
		New: func() any {
			return &CliContext{}
		},
	}
)

type CliContext struct {
	context.Context
	CancelFunc context.CancelFunc
	Logger     logger.Logger
}

func NewCliContext(parent context.Context) *CliContext {
	// 从池中取出
	cliCtx := cliCtxPool.Get().(*CliContext)
	// 创建可取消的 context
	nctx, cancel := context.WithCancel(parent)
	// 初始化字段
	cliCtx.Context = nctx
	cliCtx.CancelFunc = cancel
	cliCtx.Logger = logger.Global().With(trace.NewTraceID())
	return cliCtx
}

func PutCliContext(c *CliContext) {
	if c == nil {
		return
	}
	// 清理资源，防止数据污染
	if c.CancelFunc != nil {
		c.CancelFunc()
	}
	c.Context = nil
	c.CancelFunc = nil
	c.Logger = nil
	// 放回池中
	cliCtxPool.Put(c)
}
