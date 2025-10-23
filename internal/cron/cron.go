package cron

import (
	"log"
	"time"
)

type Job struct {
	Name string
	Spec string
	Fn   func() error
}

func Register() []*Job {
	return []*Job{
		{
			Name: "示例任务",
			Spec: "*/30 * * * * *", // 每30秒执行一次  秒 分 时 日 月 周
			Fn: func() error {
				log.Println("[Cron] tick:", time.Now().Format(time.RFC3339))
				return nil
			},
		},
	}
}
