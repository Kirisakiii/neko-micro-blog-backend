/*
Package stores - NekoBlog backend server data access objects.
This file is for post storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
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
	"strings"
	"time"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// PostStore 博文信息数据库
type PostStore struct {
	db    *gorm.DB
	rds   *redis.Client
	mongo *mongo.Client
}

// NewPostStore 是一个工厂方法，用于创建 PostStore 的新实例。
//
// 参数
// - factory: 一个包含 gorm.DB 的 Factory 实例，用于初始化 PostStore 的数据库连接。
//
// 返回值
// 它初始化并返回一个 PostStore，并关联了相应的 gorm.DB。
func (factory *Factory) NewPostStore() *PostStore {
	return &PostStore{
		db:    factory.db,
		rds:   factory.rds,
		mongo: factory.mongo,
	}
}

// GetPostList 获取适用于用户查看的帖子信息列表。
//
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在检索过程中遇到的任何错误，如果有的话。
func (store *PostStore) GetPostList(from string, length int) ([]models.PostInfo, error) {
	var posts []models.PostInfo
	if from != "" {
		if result := store.db.Where("id < ?", from).Order("id desc").Limit(length).Find(&posts); result.Error != nil {
			return nil, result.Error
		}
		return posts, nil
	}
	if result := store.db.Order("id desc").Limit(length).Find(&posts); result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// GetPostListByUID 获取适用于用户查看的帖子信息列表。
//
// 参数：
// - uid：用户ID
//
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在检索过程中遇到的任何错误，如果有的话。
func (store *PostStore) GetPostListByUID(uid string) ([]models.PostInfo, error) {
	var userPosts []models.PostInfo
	if result := store.db.Where("uid = ?", uid).Order("id desc").Find(&userPosts); result.Error != nil {
		return nil, result.Error
	}
	return userPosts, nil
}

// ValidatePostExistence 用来检查是否存在Post博文
//
// 参数：postID：博文ID
//
// 返回值：
// - bool: 找到返回true ，找不到返回false
// - error: 返回的错误类型是否是post为空
func (store *PostStore) ValidatePostExistence(postID uint64) (bool, error) {
	var post models.PostInfo
	result := store.db.Where("id = ?", postID).First(&post)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	// 返回错误类型
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// GetPostByUID 通过用户UID获取用户信息。
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - *models.PostInfo：如果找到了相应的用户信息，则返回该用户信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *PostStore) GetPostInfo(postID uint64) (models.PostInfo, int64, int64, error) {
	post := models.PostInfo{}
	result := store.db.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return models.PostInfo{}, 0, 0, result.Error
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	likeCount, err := postLikeCollection.CountDocuments(context.Background(), bson.D{{Key: "post_id", Value: postID}})
	if err != nil {
		return models.PostInfo{}, 0, 0, err
	}
	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	favouriteCount, err := postFavouriteCollection.CountDocuments(context.Background(), bson.D{{Key: "post_id", Value: postID}})
	if err != nil {
		return models.PostInfo{}, 0, 0, err
	}

	return post, likeCount, favouriteCount, nil
}

// CreatePost 根据用户提交的帖子信息创建帖子。
//
// 参数：
//   - userID：用户ID，用于关联帖子与用户。
//   - ipAddr：IP地址
//   - postInfo：帖子信息，包含标题、内容等。
//   - images：帖子图片
//
// 返回值：
//   - error：如果在创建过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *PostStore) CreatePost(uid uint64, ipAddr string, postReqData types.PostCreateBody) (models.PostInfo, error) {
	var imageFileNames []string
	// 将文件复制出缓存
	for _, imageUUID := range postReqData.Images {
		srcImage, err := os.Open(filepath.Join(consts.POST_IMAGE_CACHE_PATH, imageUUID+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer srcImage.Close()
		dstImage, err := os.Create(filepath.Join(consts.POST_IMAGE_PATH, imageUUID+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer dstImage.Close()
		_, err = io.Copy(dstImage, srcImage)
		if err != nil {
			return models.PostInfo{}, err
		}
		imageFileNames = append(imageFileNames, imageUUID+".webp")

		// 删除缓存图片
		ctx := context.Background()
		tx := store.rds.TxPipeline()

		_, err = tx.XAdd(ctx, &redis.XAddArgs{
			Stream: consts.CACHE_IMG_CLEAN_STREAM,
			Values: map[string]interface{}{"filename": imageUUID + ".webp"},
		}).Result()
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}

		// 删除数据库记录
		var sb strings.Builder
		sb.WriteString(consts.CACHE_IMAGE_LIST)
		sb.WriteString(":")
		sb.WriteString(imageUUID)
		fmt.Println(sb.String())
		_, err = tx.Del(ctx, sb.String()).Result()
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}

		_, err = tx.Exec(ctx)
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}
	}

	// 将博文数据写入数据库
	postInfo := models.PostInfo{
		ParentPostID: nil,
		UID:          uid,
		IpAddrress:   &ipAddr,
		Title:        postReqData.Title,
		Content:      postReqData.Content,
		Images:       imageFileNames,
		Like:         pq.Int64Array{},
		Favourite:    pq.Int64Array{},
		Farward:      pq.Int64Array{},
		IsPublic:     true,
	}
	result := store.db.Create(&postInfo)
	return postInfo, result.Error
}

// CachePostImage 缓存博文图片
//
// 参数：
//   - image []byte：待缓存的图片
//
// 返回值：
//   - string：如果缓存成功，返回缓存图片的UUID；否则返回空字符串
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) CachePostImage(image []byte) (string, error) {
	// 生成文件名
	var (
		fileNameBuilder strings.Builder
		UUID            string
		savePath        string
	)
	for {
		UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
		fileNameBuilder.WriteString(UUID)
		fileNameBuilder.WriteString(".webp")
		savePath = filepath.Join(consts.POST_IMAGE_CACHE_PATH, fileNameBuilder.String())
		_, err := os.Stat(savePath)
		if os.IsExist(err) {
			fileNameBuilder.Reset()
			continue
		}
		break
	}

	// 保存图片
	file, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(image))
	if err != nil {
		return "", err
	}

	// 写入缓存列表
	ctx := context.Background()
	var sb strings.Builder
	sb.WriteString(consts.CACHE_IMAGE_LIST)
	sb.WriteString(":")
	sb.WriteString(UUID)

	_, err = store.rds.HSet(ctx, sb.String(), map[string]interface{}{
		"filename": fileNameBuilder.String(),
		"expire":   time.Now().Add(consts.CACHE_IMAGE_EXPIRE_TIME * time.Second).Unix(),
	}).Result()

	if err != nil {
		return "", err
	}

	return UUID, nil
}

// CheckCacheImageAvaliable 检查缓存图片是否存在
//
// 参数：
//   - uuid string：待检查的缓存图片UUID
//
// 返回值：
//   - bool：如果缓存图片存在，返回true；否则返回false
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) CheckCacheImageAvaliable(uuid string) (bool, error) {
	// 检查缓存图片是否存在
	ctx := context.Background()

	var sb strings.Builder
	sb.WriteString(consts.CACHE_IMAGE_LIST)
	sb.WriteString(":")
	sb.WriteString(uuid)

	// 遍历缓存列表
	flag := false
	keys, err := store.rds.Keys(ctx, consts.CACHE_IMAGE_LIST+":*").Result()
	if err != nil {
		return false, err
	}
	for _, key := range keys {
		if store.rds.HGet(ctx, key, "filename").Val() == uuid+".webp" {
			// 如果存在则检测是否过期
			expire, err := store.rds.HGet(ctx, key, "expire").Int64()
			if err != nil {
				return false, err
			}
			// 返回过期
			if time.Now().Unix() > expire {
				return false, nil
			}
			flag = true
			break
		}
	}
	// 不存在
	if !flag {
		return false, nil
	}

	_, err = os.Stat(filepath.Join(consts.POST_IMAGE_CACHE_PATH, uuid+".webp"))
	// 文件不存在
	if os.IsNotExist(err) {
		// 删除缓存记录
		tx := store.rds.TxPipeline()
		var sb strings.Builder
		sb.WriteString(consts.CACHE_IMAGE_LIST)
		sb.WriteString(":")
		sb.WriteString(uuid)
		_, err = tx.Del(ctx, sb.String()).Result()
		if err != nil {
			tx.Discard()
			return false, err
		}
		_, err = tx.Exec(ctx)
		if err != nil {
			tx.Discard()
			return false, err
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// LikePost 点赞博文
//
// 参数：
//   - postID uint64：待点赞博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) LikePost(uid, postID int64) error {
	// 查询博文信息
	postInfo := models.PostInfo{}
	result := store.db.Where("id = ?", postID).First(&postInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("post does not exist")
	}
	if result.Error != nil {
		return result.Error
	}

	// 构造查询条件
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
		{Key: "poster_uid", Value: postInfo.UID},
	}
	// 构造更新内容
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "liked_at", Value: time.Now()},
		}},
	}

	// 更新博文点赞记录
	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	_, err := postLikeCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	// 重复点赞
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has liked this post")
	}
	return err
}

// CancelLikePost 取消点赞博文
//
// 参数：
//   - uid：用户ID
//   - postID：待取消点赞博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) CancelLikePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	_, err := postLikeCollection.DeleteOne(context.Background(), filter)
	if mongo.ErrNoDocuments == err {
		return errors.New("user has not liked this post")
	}
	return err
}

// FavouritePost 收藏博文
//
// 参数：
//   - uid：用户ID
//   - postID：待收藏博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) FavouritePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "favourited_at", Value: time.Now()},
		}},
	}

	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	_, err := postFavouriteCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has favourited this post")
	}
	return err
}

// CancelFavouritePost 取消收藏博文
//
// 参数：
//   - uid：用户ID
//   - postID：待取消收藏博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) CancelFavouritePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	_, err := postFavouriteCollection.DeleteOne(context.Background(), filter)
	if mongo.ErrNoDocuments == err {
		return errors.New("user has not favourited this post")
	}
	return err
}

// GetPostUserStatus 获取用户对帖子的状态
//
// 参数：
//   - uid int64：用户ID
//   - postID int64：帖子ID
//
// 返回值：
//   - bool：用户是否点赞
//   - bool：用户是否收藏
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) GetPostUserStatus(uid, postID int64) (bool, bool, error) {
	// 构造查询条件
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	// 查询博文点赞记录
	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	count, err := postLikeCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, false, err
	}
	isLiked := count > 0

	// 查询博文收藏记录
	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	count, err = postFavouriteCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, false, err
	}
	isFavourited := count > 0

	return isLiked, isFavourited, nil
}

// DeletePost 通过博文ID删除博文的存储方法
//
// 参数：
// - postID uint64：待删除博文的ID
//
// 返回值：
// - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *PostStore) DeletePost(postID uint64) error {
	return store.db.Where("id = ?", postID).Unscoped().Delete(&models.PostInfo{}).Error
}
