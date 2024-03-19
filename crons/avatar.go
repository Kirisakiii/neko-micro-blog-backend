package crons

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// AvatarCleanJob 头像清理任务
type AvatarCleanJob struct {
	logger *logrus.Logger // 日志记录器
	rds    *redis.Client  // 数据库连接
}

// NewAvatarCleanJob 创建一个新的头像清理任务。
//
// 参数：
//   - logger：日志记录器
//   - redisClient：数据库连接
//
// 返回值：
//   - *AvatarCleanJob：新的头像清理任务。
func NewAvatarCleanJob(logger *logrus.Logger, redisClient *redis.Client) *AvatarCleanJob {
	return &AvatarCleanJob{
		logger: logger,
		rds:    redisClient,
	}
}

// Run 执行头像清理任务。
func (job *AvatarCleanJob) Run() {
	job.logger.Debugln("正在执行头像清理任务...")

	ctx := context.Background()
	// 获取队列长度
	length, err := job.rds.XLen(ctx, consts.AVATAR_CLEAN_STREAM).Result()
	if err != nil {
		job.logger.Errorln("获取清理队列长度失败:", err)
		return
	}
	if length == 0 {
		job.logger.Debugln("头像清理任务执行完毕")
		return
	}

	// 获取清除队列
	messages, err := job.rds.XRead(ctx, &redis.XReadArgs{
		Streams: []string{consts.AVATAR_CLEAN_STREAM, "0"},
		Count:   0,
		Block:   0,
	}).Result()
	if err != nil {
		job.logger.Errorln("获取清理队列失败:", err)
		return
	}

	// 清理头像
	for _, item := range messages[0].Messages {
		filename := item.Values["filename"].(string)
		err := os.Remove(filepath.Join(consts.AVATAR_IMAGE_PATH, filename))
		// 头像文件不存在
		if errors.Is(err, os.ErrNotExist) {
			job.logger.Warnln("头像文件不存在:", filename)
			_, err := job.rds.XDel(ctx, consts.AVATAR_CLEAN_STREAM, item.ID).Result()
			if err != nil {
				job.logger.Errorln("清理头像失败:", err)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("清理头像失败:", err)
			continue
		}
		// 删除数据库记录
		_, err = job.rds.XDel(ctx, consts.AVATAR_CLEAN_STREAM, item.ID).Result()
		if err != nil {
			job.logger.Errorln("清理头像失败:", err)
			continue
		}
	}

	job.logger.Debugln("头像清理任务执行完毕")
}
