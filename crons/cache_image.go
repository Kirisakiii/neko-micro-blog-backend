package crons

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// CachedImageCleanJob 头像清理任务
type CachedImageCleanJob struct {
	logger *logrus.Logger // 日志记录器
	rds    *redis.Client  // 数据库连接
}

// NewCachedImageCleanJob 博文图片缓存清理任务
//
// 参数：
//   - logger：日志记录器
//   - redisClient：数据库连接
//
// 返回值：
//   - *CachedImageCleanJob：创建一个新的博文图片缓存清理任务。
func NewCachedImageCleanJob(logger *logrus.Logger, rds *redis.Client) *CachedImageCleanJob {
	return &CachedImageCleanJob{
		logger: logger,
		rds:    rds,
	}
}

// Run 执行头像清理任务。
func (job *CachedImageCleanJob) Run() {
	job.logger.Debugln("正在执行缓存图片清理任务...")

	ctx := context.Background()
	// 遍历缓存表并将过期的图片添加至清除队列

	keys, err := job.rds.Keys(ctx, consts.CACHE_IMAGE_LIST+":*").Result()
	if err != nil {
		job.logger.Errorln("获取缓存图片列表失败:", err)
		return
	}

	// 遍历缓存表
	for _, key := range keys {
		// 获取过期时间
		timestamp, err := job.rds.HGet(ctx, key, "expire").Int64()
		if err != nil {
			job.logger.Errorln("获取缓存图片过期时间失败:", err)
			continue
		}
		// 获取文件名
		filename, err := job.rds.HGet(ctx, key, "filename").Result()
		if err != nil {
			job.logger.Errorln("获取缓存图片文件名失败:", err)
			continue
		}

		// 如果过期时间小于当前时间，则添加至清除队列并删除缓存
		if timestamp < time.Now().Unix() {
			tx := job.rds.TxPipeline()
			_, err := tx.XAdd(ctx, &redis.XAddArgs{
				Stream: consts.CACHE_IMG_CLEAN_STREAM,
				Values: map[string]interface{}{"filename": filename},
			}).Result()
			if err != nil {
				tx.Discard()
				job.logger.Errorln("添加缓存图片清理队列失败:", err)
				return
			}

			_, err = tx.Del(ctx, key).Result()
			if err != nil {
				tx.Discard()
				job.logger.Errorln("删除缓存图片失败:", err)
				return
			}

			_, err = tx.Exec(ctx)
			if err != nil {
				tx.Discard()
				job.logger.Errorln("删除缓存图片失败:", err)
				continue
			}
		}
	}

	// 获取队列长度
	length, err := job.rds.XLen(ctx, consts.CACHE_IMG_CLEAN_STREAM).Result()
	if err != nil {
		job.logger.Errorln("获取清理队列长度失败:", err)
		return
	}
	if length == 0 {
		job.logger.Debugln("缓存图片清理任务执行完毕")
		return
	}

	// 获取清除队列
	messages, err := job.rds.XRead(ctx, &redis.XReadArgs{
		Streams: []string{consts.CACHE_IMG_CLEAN_STREAM, "0"},
		Count:   0,
		Block:   0,
	}).Result()
	if err != nil {
		job.logger.Errorln("获取清理队列失败:", err)
		return
	}

	// 清理缓存图片
	for _, item := range messages[0].Messages {
		filename := item.Values["filename"].(string)
		err := os.Remove(filepath.Join(consts.POST_IMAGE_CACHE_PATH, filename))
		// 头像文件不存在
		if errors.Is(err, os.ErrNotExist) {
			job.logger.Warningln("缓存图片文件不存在:", filename)
			_, err := job.rds.XDel(ctx, consts.CACHE_IMG_CLEAN_STREAM, item.ID).Result()
			if err != nil {
				job.logger.Errorln("清理缓存图片失败:", err)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("清理缓存图片失败:", err)
			continue
		}
		// 删除数据库记录
		_, err = job.rds.XDel(ctx, consts.CACHE_IMG_CLEAN_STREAM, item.ID).Result()
		if err != nil {
			job.logger.Errorln("清理缓存图片失败:", err)
		}
	}

	job.logger.Debugln("缓存图片清理任务执行完毕")
}
