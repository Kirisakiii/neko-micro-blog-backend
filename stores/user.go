/*
Package stores - NekoBlog backend server data access objects.
This file is for user storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/functools"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// UserStore 用户信息数据库
type UserStore struct {
	db  *gorm.DB
	rds *redis.Client
}

// NewUserStore 返回一个新的 UserStore 实例。
//
// 返回值：
//   - *UserStore：新的 UserStore 实例。
func (factory *Factory) NewUserStore() *UserStore {
	return &UserStore{factory.db, factory.rds}
}

// RegisterUserByUsername 注册用户将提供的用户名、盐和哈希密码注册到数据库中。
//
// 参数：
//   - username：用户名
//   - salt：盐值
//   - hashedPassword：哈希密码
//
// 返回值：
//   - error：如果在注册过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) RegisterUserByUsername(username string, salt string, hashedPassword string) error {
	tx := store.db.Begin()

	user := models.UserInfo{
		UserName: username,
		NickName: &username,
	}
	result := tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	uid := user.ID
	userAuthInfo := models.UserAuthInfo{
		UID:          uint64(uid),
		UserName:     username,
		Salt:         salt,
		PasswordHash: hashedPassword,
	}
	result = tx.Create(&userAuthInfo)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	userPostStatus := models.UserPostStatus{
		UID:        uint64(uid),
		Viewed:     pq.Int64Array{},
		Liked:      pq.Int64Array{},
		Favourited: pq.Int64Array{},
		Commented:  pq.Int64Array{},
	}
	result = tx.Create(&userPostStatus)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	userCommentStatus := models.UserCommentStatus{
		UID:       uint64(uid),
		Liked:     pq.Int64Array{},
		Disliked:  pq.Int64Array{},
		Commented: pq.Int64Array{},
	}
	result = tx.Create(&userCommentStatus)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// GetUserByUID 通过用户ID获取用户信息。
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - *models.UserInfo：如果找到了相应的用户信息，则返回该用户信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserByUID(uid uint64) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户信息。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserInfo：如果找到了相应的用户信息，则返回该用户信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserByUsername(username string) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserAuthInfoByUsername 通过用户名获取用户的认证信息。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserAuthInfo：如果找到了相应的用户认证信息，则返回该用户认证信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserAuthInfoByUsername(username string) (*models.UserAuthInfo, error) {
	userAuthInfo := new(models.UserAuthInfo)
	result := store.db.Where("username = ?", username).First(userAuthInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return userAuthInfo, nil
}

// InsertUserLoginLog 插入用户登录日志。
//
// 参数：
//   - userLoginLogInfo：用户登录日志信息
//
// 返回值：
//   - error：如果在插入过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) CreateUserLoginLog(userLoginLogInfo *models.UserLoginLog) error {
	result := store.db.Create(userLoginLogInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateUserAvaliableToken 创建一个可用的 Token。
//
// 参数：
//   - token：Token
//   - claims：Token 的声明
//
// 返回值：
//   - error：如果在创建过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) CreateUserAvaliableToken(token string, claims *types.BearerTokenClaims) error {
	var sb strings.Builder
	sb.WriteString(consts.REDIS_AVAILABLE_USER_TOKEN_LIST)
	sb.WriteRune(':')
	sb.WriteString(strconv.FormatUint(claims.UID, 10))
	key := sb.String()

	ctx := context.Background()

	// 获取当前 Token 数量
	length, err := store.rds.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	fmt.Println(length)

	tx := store.rds.TxPipeline()
	// 如果 Token 数量超过限制，则移除最早的 Token
	if length >= consts.MAX_TOKENS_PER_USER {
		// 移除超出限制的 Token
		_, err = tx.LTrim(ctx, key, length-4, -1).Result()
		if err != nil {
			tx.Discard()
			return err
		}
	}

	// 添加新 Token
	_, err = tx.RPush(ctx, key, token).Result()
	if err != nil {
		tx.Discard()
		return err
	}

	// 执行事务
	_, err = tx.Exec(ctx)
	if err != nil {
		tx.Discard()
		return err
	}

	return nil
}

// BanUserToken 将 Token 禁用。
//
// 参数：
//   - token：Token
//
// 返回值：
//   - error：如果在禁用过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) BanUserToken(uid uint64, token string) error {
	var sb strings.Builder
	sb.WriteString(consts.REDIS_AVAILABLE_USER_TOKEN_LIST)
	sb.WriteRune(':')
	sb.WriteString(strconv.FormatUint(uid, 10))
	key := sb.String()

	ctx := context.Background()
	tx := store.rds.TxPipeline()

	// 移除 Token
	_, err := tx.LRem(ctx, key, 0, token).Result()
	if err != nil {
		tx.Discard()
		return err
	}

	// 执行事务
	_, err = tx.Exec(ctx)
	if err != nil {
		tx.Discard()
		return err
	}

	return nil
}

// IsUserTokenAvaliable 检查 Token 是否可用。
//
// 参数：
//   - token：Token
//
// 返回值：
//   - bool：如果 Token 可用，则返回 true，否则返回 false。
//   - error：如果在检查过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) IsUserTokenAvaliable(token string) (bool, error) {
	ctx := context.Background()

	// 获取所有用户的令牌列表
	keys, err := store.rds.Keys(ctx, consts.REDIS_AVAILABLE_USER_TOKEN_LIST+":*").Result()
	if err != nil {
		return false, err
	}

	// 遍历所有用户的 Token
	for _, key := range keys {
		// 获取用户的 Token
		tokens, err := store.rds.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return false, err
		}

		// 检查 Token 是否存在
		for _, t := range tokens {
			if t == token {
				return true, nil
			}
		}
	}

	return false, nil
}

// SaveUserAvatarByUID 保存用户头像。
//
// 参数：
//   - fileName：文件名
//   - data：文件数据
//
// 返回值：
//   - error：如果在保存过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) SaveUserAvatarByUID(uid uint64, fileName string, data []byte) error {
	savePath := filepath.Join(consts.AVATAR_IMAGE_PATH, fileName)

	// 创建目标文件
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 io.Copy 将数据写入文件
	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return err
	}

	// 用户信息记录
	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return result.Error
	}

	// 将旧头像文件加入清理队列
	if user.Avatar != "vanilla.webp" {
		ctx := context.Background()
		store.rds.XAdd(ctx, &redis.XAddArgs{
			Stream: consts.AVATAR_CLEAN_STREAM,
			Values: map[string]interface{}{
				"filename": user.Avatar,
			},
		})
	}

	// 更新头像文件名
	user.Avatar = fileName
	result = store.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateUserPasswordByUsername 更新用户密码。
//
// 参数：
//   - username：用户名
//   - hashedNewPassword：经过哈希处理的新密码
//
// 返回值：
//   - error：如果在更新过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) UpdateUserPasswordByUsername(username string, hashedNewPassword string) error {
	userAuthInfo := new(models.UserAuthInfo)
	result := store.db.Where("username = ?", username).First(userAuthInfo)
	if result.Error != nil {
		return result.Error
	}

	userAuthInfo.PasswordHash = hashedNewPassword
	result = store.db.Save(userAuthInfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateUserInfoByUID 更新用户信息。
//
// 参数：
//   - uid：用户ID
//   - updatedProfile：更新后的用户信息
//
// 返回值：
//   - error：如果在更新过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) UpdateUserInfoByUID(uid uint64, updatedProfile *models.UserInfo) error {
	var userProfile models.UserInfo
	result := store.db.Where("id = ?", uid).First(&userProfile)
	if result.Error != nil {
		return result.Error
	}

	userProfile.UpdatedAt = time.Now()
	userProfile.NickName = updatedProfile.NickName
	userProfile.Birth = updatedProfile.Birth
	userProfile.Gender = updatedProfile.Gender

	return store.db.Save(&userProfile).Error
}

// GetUserLikedRecord 获取用户点赞记录。
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - pq.Int64Array：用户点赞记录。
func (store *UserStore) GetUserLikedRecord(uid string) ([]int64, error) {
	userPostStatus := models.UserPostStatus{}

	result := store.db.Where("uid = ?", uid).First(&userPostStatus)
	if result.Error != nil {
		return nil, result.Error
	}

	return functools.Reverse(userPostStatus.Liked), nil
}

// GetUserFavoriteRecord 获取用户收藏记录。
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - pq.Int64Array：用户收藏记录。
func (store *UserStore) GetUserFavoriteRecord(uid string) ([]int64, error) {
	userPostStatus := models.UserPostStatus{}

	result := store.db.Where("uid = ?", uid).First(&userPostStatus)
	if result.Error != nil {
		return nil, result.Error
	}

	return functools.Reverse(userPostStatus.Favourited), nil
}

// AddUserLikedRecord 添加用户点赞记录。
//
// 参数：
//   - uid：用户ID
//   - postID：帖子ID
//   - tx：事务 当为nil时自动创建事务
//
// 返回值：
//   - error：如果在添加过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) AddUserLikedRecord(uid, postID int64, tx *gorm.DB) error {
	if tx == nil {
		tx = store.db.Begin()
		defer tx.Commit()
	}

	result := tx.Model(&models.UserPostStatus{}).
		Where("uid = ? AND NOT ARRAY[?::bigint] <@ \"liked\"", uid, postID).
		Update("liked", gorm.Expr("array_append(\"liked\", ?)", postID))

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user has liked this post")
	}

	return nil
}

// RemoveUserLikedRecord 移除用户点赞记录。
//
// 参数：
//   - uid：用户ID
//   - postID：帖子ID
//   - tx：事务 当为nil时自动创建事务
//
// 返回值：
//   - error：如果在移除过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) RemoveUserLikedRecord(uid, postID int64, tx *gorm.DB) error {
	if tx == nil {
		tx = store.db.Begin()
		defer tx.Commit()
	}

	result := tx.Model(&models.UserPostStatus{}).
		Where("uid = ? AND ARRAY[?::bigint] <@ \"liked\"", uid, postID).
		Update("liked", gorm.Expr("array_remove(\"liked\", ?)", postID))

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user has not liked this post")
	}

	return nil
}

// AddUserFavoriteRecord 添加用户收藏记录。
//
// 参数：
//   - uid：用户ID
//   - postID：帖子ID
//   - tx：事务 当为nil时自动创建事务
//
// 返回值：
//   - error：如果在添加过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) AddUserFavoriteRecord(uid, postID int64, tx *gorm.DB) error {
	if tx == nil {
		tx = store.db.Begin()
		defer tx.Commit()
	}

	result := tx.Model(&models.UserPostStatus{}).
		Where("uid = ? AND NOT ARRAY[?::bigint] <@ \"favourite\"", uid, postID).
		Update("favourite", gorm.Expr("array_append(\"favourite\", ?)", postID))

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user has favourited this post")
	}

	return nil
}

// RemoveUserFavoriteRecord 移除用户收藏记录。
//
// 参数：
//   - uid：用户ID
//   - postID：帖子ID
//   - tx：事务 当为nil时自动创建事务
//
// 返回值：
//   - error：如果在移除过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) RemoveUserFavoriteRecord(uid, postID int64, tx *gorm.DB) error {
	if tx == nil {
		tx = store.db.Begin()
		defer tx.Commit()
	}

	result := tx.Model(&models.UserPostStatus{}).
		Where("uid = ? AND ARRAY[?::bigint] <@ \"favourite\"", uid, postID).
		Update("favourite", gorm.Expr("array_remove(\"favourite\", ?)", postID))

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user has not favourited this post")
	}

	return nil
}
