package crons

import (
	"context"
	"errors"
	
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/parsers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// TokenCleanJob 令牌清理任务
type TokenCleanJob struct {
	logger *logrus.Logger
	rds    *redis.Client
}

// NewTokenCleanJob 创建一个新的令牌清理任务。
//
// 参数：
//   - logger：日志记录器
//   - redisClient：数据库连接
//
// 返回值：
//   - *TokenCleanJob：新的令牌清理任务。
func NewTokenCleanJob(logger *logrus.Logger, redisClient *redis.Client) *TokenCleanJob {
	return &TokenCleanJob{
		logger: logger,
		rds:    redisClient,
	}
}

// Run 执行令牌清理任务。
func (job *TokenCleanJob) Run() {
	job.logger.Debugln("正在执行令牌清理任务...")

	ctx := context.Background()

	// 获取所有用户的令牌列表
	keys, err := job.rds.Keys(ctx, consts.REDIS_AVAILABLE_USER_TOKEN_LIST+":*").Result()
	if err != nil {
		job.logger.Errorln("获取用户令牌列表失败:", err)
		return
	}

	// 遍历所有用户的令牌列表
	for _, key := range keys {
		// 获取用户令牌列表
		tokens, err := job.rds.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			job.logger.Errorln("获取用户令牌列表失败:", err)
			continue
		}

		// 遍历用户令牌列表
		for _, token := range tokens {
			_, err := parsers.ParseToken(token)
			// 如果令牌过期，则删除令牌
			if errors.Is(err, jwt.ErrTokenExpired) {
				// 删除过期令牌
				_, err := job.rds.LRem(ctx, key, 0, token).Result()
				if err != nil {
					job.logger.Errorln("删除过期令牌失败:", err)
				}
				continue
			}
			// 如果解析令牌失败，则记录错误信息
			if err != nil {
				job.logger.Errorln("解析令牌失败:", err)
				continue
			}
		}
	}

	job.logger.Debugln("令牌清理任务执行完毕")
}
