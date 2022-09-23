package cron

import (
	"time"

	"git.code.oa.com/tmemesh.io/srf/srfs"
	"github.com/robfig/cron"
)

type CronJob struct {
	TTL      int32
	CronSpec string // 格式参考：https://segmentfault.com/a/1190000023029219
	Job      func(ctx *srfs.Context)
}

// AsyncRunOneJob: 运行一个任务
func AsyncRunOneJob(ctx *srfs.Context, one CronJob) error {
	ctx.AsyncDo(func(ctx *srfs.Context) {
		if one.CronSpec != "" {
			c := cron.New()
			c.AddFunc(one.CronSpec, func() {
				ctx.AsyncDo(one.Job)
			})
			c.Start()
			return
		}
		for range time.Tick(time.Duration(one.TTL) * time.Second) {
			ctx.AsyncDo(one.Job)
		}
	})
	return nil
}

// AsyncRunJobSet: 运行异步任务集
func AsyncRunJobSet(ctx *srfs.Context, set []CronJob) error {
	for _, one := range set {
		AsyncRunOneJob(ctx, one)
	}
	return nil
}
