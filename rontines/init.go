package rontines

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/utils/jobs"
)

func InitJobs(logger *logrus.Logger, db *gorm.DB) {
	// 创建定时任务
	crontab := cron.New()
	_, err := jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewAvatarCleanerJob(logger, db))
	if err != nil {
		logger.Panicln(err.Error())
	}
	_, err = jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewCachedImageCleanerJob(logger, db))
	if err != nil {
		logger.Panicln(err.Error())
	}
	crontab.Start()
}
