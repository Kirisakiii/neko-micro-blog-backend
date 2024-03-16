package jobs

import (
	"github.com/robfig/cron/v3"
)

func AddSkipIfStillRunningJob(crontab *cron.Cron, spec string, job cron.Job) (cron.EntryID, error) {
	return crontab.AddJob(
		spec,
		cron.NewChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		).Then(
			job,
		),
	)
}
