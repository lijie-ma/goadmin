package server

import (
	"context"
	"fmt"
	"log"

	bizCron "goadmin/internal/cron"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	c *cron.Cron
}

func NewCronManager() *CronManager {
	return &CronManager{c: cron.New(cron.WithSeconds())}
}

func (cm *CronManager) Name() string { return "CronManager" }

func (cm *CronManager) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	// 封装任务注册
	add := func(spec string, fn func() error) {
		_, _ = cm.c.AddFunc(spec, func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in job: %v", r)
				}
			}()
			if err := fn(); err != nil {
				errChan <- err
			}
		})
	}

	tasks := bizCron.Register()
	for _, task := range tasks {
		add(task.Spec, task.Fn)
	}

	cm.c.Start()
	defer cm.c.Stop()

	select {
	case <-ctx.Done():
		log.Println("[Cron] stop signal from context")
		return ctx.Err()
	case err := <-errChan:
		return fmt.Errorf("cron job error: %w", err)
	}
}

func (cm *CronManager) Stop(ctx context.Context) error {
	ch := cm.c.Stop().Done()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
