package crons

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/utils/jobs"
)

// InitJobs 初始化定时任务
//
// 参数：
//   - logger：日志记录器
//   - db：数据库连接
//   - redisClient：Redis 连接
func InitJobs(logger *logrus.Logger, db *gorm.DB, redisClient *redis.Client) {
	// 创建定时任务
	crontab := cron.New()

	// 头像清理任务
	_, err := jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewAvatarCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}
	// 缓存图片清理任务
	_, err = jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewCachedImageCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}
	// 令牌清理任务
	_, err = jobs.AddSkipIfStillRunningJob(crontab, "@every 1h", NewTokenCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}

	// 启动定时任务
	crontab.Start()
}
