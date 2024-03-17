package rontines

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CachedImageCleanerJob 头像清理任务
type CachedImageCleanerJob struct {
	logger *logrus.Logger // 日志记录器
	db     *gorm.DB       // 数据库连接
}

// NewAvatarCleanerJob 创建一个新的头像清理任务。
//
// 参数：
//   - db：数据库连接
//   - logger：日志记录器
//
// 返回值：
//   - *AvatarCleanerJob：新的头像清理任务。
func NewCachedImageCleanerJob(logger *logrus.Logger, db *gorm.DB) *CachedImageCleanerJob {
	return &CachedImageCleanerJob{
		logger: logger,
		db:     db,
	}
}

// Run 执行头像清理任务。
func (job *CachedImageCleanerJob) Run() {
	job.logger.Debugln("正在执行缓存图片清理任务...")
	// 获取清除队列
	var waitList []models.DeletedCachedImage
	result := job.db.Find(&waitList)
	if result.Error != nil {
		job.logger.Errorln("获取缓存图片清理队列失败:", result.Error)
		return
	}

	// 清理头像
	for _, item := range waitList {
		err := os.Remove(filepath.Join(consts.POST_IMAGE_CACHE_PATH, item.FileName))
		// 头像文件不存在
		if errors.Is(err, os.ErrNotExist) {
			job.logger.Warnln("缓存图片文件不存在:", item.FileName)
			result = job.db.Unscoped().Delete(&item)
			if result.Error != nil {
				job.logger.Errorln("清理缓存图片失败:", result.Error)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("清理缓存图片失败:", err)
			continue
		}
		// 删除数据库记录
		result = job.db.Unscoped().Delete(&item)
		if result.Error != nil {
			job.logger.Errorln("清理缓存图片失败:", result.Error)
			continue
		}
	}
	job.logger.Debugln("缓存图片清理任务执行完毕")
}
