package server

import (
	"context"
	"log"
	"runtime/debug"

	bizCron "goadmin/internal/cron"
	"goadmin/pkg/logger"

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
	// 封装任务注册
	add := func(job *bizCron.Job) {
		_, _ = cm.c.AddFunc(job.Spec, func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("[Cron] job %s panic: %v %s", job.Name, r, debug.Stack())
				}
			}()
			if err := job.Fn(); err != nil {
				logger.Errorf("[Cron] job %s error: %v", job.Name, err)
			}
		})
	}

	tasks := bizCron.Register()
	for _, task := range tasks {
		add(task)
	}

	cm.c.Start()
	defer cm.c.Stop()

	<-ctx.Done()
	log.Println("[Cron] stop signal from context")
	return ctx.Err()
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
