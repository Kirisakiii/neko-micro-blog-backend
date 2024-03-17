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
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostStore 博文信息数据库
type PostStore struct {
	db *gorm.DB
}

// NewPostStore 是一个工厂方法，用于创建 PostStore 的新实例。
//
// 参数
// - factory: 一个包含 gorm.DB 的 Factory 实例，用于初始化 PostStore 的数据库连接。
//
// 返回值
// 它初始化并返回一个 PostStore，并关联了相应的 gorm.DB。
func (factory *Factory) NewPostStore() *PostStore {
	return &PostStore{factory.db}
}

// GetPostList 获取适用于用户查看的帖子信息列表。
//
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在检索过程中遇到的任何错误，如果有的话。
func (store *PostStore) GetPostList() ([]models.PostInfo, error) {
	var posts []models.PostInfo
	if result := store.db.Order("id desc").Find(&posts); result.Error != nil {
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
func (store *PostStore) GetPostInfo(postID uint64) (models.PostInfo, error) {
	post := models.PostInfo{}
	result := store.db.Where("id = ?", postID).First(&post)
	return post, result.Error
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
	for _, image := range postReqData.Images {
		srcImage, err := os.Open(filepath.Join(consts.POST_IMAGE_CACHE_PATH, image+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer srcImage.Close()
		dstImage, err := os.Create(filepath.Join(consts.POST_IMAGE_PATH, image+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer dstImage.Close()
		_, err = io.Copy(dstImage, srcImage)
		if err != nil {
			return models.PostInfo{}, err
		}
		imageFileNames = append(imageFileNames, image+".webp")

		// 删除缓存中的文件
		result := store.db.Create(&models.DeletedCachedImage{
			FileName: image + ".webp",
		})
		if result.Error != nil {
			return models.PostInfo{}, result.Error
		}

		// 删除数据库记录
		result = store.db.Where("file_name = ?", image+".webp").Unscoped().Delete(&models.CachedPostImage{})
		if result.Error != nil {
			return models.PostInfo{}, result.Error
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
		IsPublic:     true,
	}
	result := store.db.Create(&postInfo)
	return postInfo, result.Error
}

func (store *PostStore) CachePostIamge(image []byte) (string, error) {
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
	cachedPostImage := models.CachedPostImage{
		FileName:   fileNameBuilder.String(),
		ExpireTime: time.Now().Add(consts.CACHE_IMAGE_EXPIRE_TIME).Unix(),
	}
	result := store.db.Create(&cachedPostImage)
	return UUID, result.Error
}

func (store *PostStore) CheckCacheImageExistence(uuid string) (bool, error) {
	// 检查缓存图片是否存在
	var cachedPostImage models.CachedPostImage
	result := store.db.Where("file_name = ?", uuid+".webp").First(&cachedPostImage)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	_, err := os.Stat(filepath.Join(consts.POST_IMAGE_CACHE_PATH, cachedPostImage.FileName))
	if os.IsNotExist(err) {
		store.db.Unscoped().Delete(&cachedPostImage)
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
	// 检测目标博文是否存在
	var postInfo models.PostInfo
	result := store.db.Where("id = ?", postID).First(&postInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("post not found")
	}
	if result.Error != nil {
		return result.Error
	}

	// 检查是否已经点赞
	var userLikedRecord models.UserLikedRecord
	result = store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}
	if slices.Contains(userLikedRecord.LikedPost, postID) {
		return errors.New("user has liked this post")
	}

	// 更新用户点赞记录
	userLikedRecord.LikedPost = append(userLikedRecord.LikedPost, postID)
	result = store.db.Save(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}

	// 更新博文点赞记录
	postInfo.Like = append(postInfo.Like, uid)
	result = store.db.Save(&postInfo)
	return result.Error
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
	// 检测目标博文是否存在
	var postInfo models.PostInfo
	result := store.db.Where("id = ?", postID).First(&postInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("post not found")
	}
	if result.Error != nil {
		return result.Error
	}

	// 检查是否已经点赞
	var userLikedRecord models.UserLikedRecord
	result = store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}
	if !slices.Contains(userLikedRecord.LikedPost, postID) {
		return errors.New("user has not liked this post")
	}

	// 更新用户点赞记录
	index := slices.Index(userLikedRecord.LikedPost, postID)
	if index != -1 {
		userLikedRecord.LikedPost = slices.Delete(userLikedRecord.LikedPost, index, index+1)
		result = store.db.Save(&userLikedRecord)
		if result.Error != nil {
			return result.Error
		}
	}

	// 更新博文点赞记录
	index = slices.Index(postInfo.Like, uid)
	if index != -1 {
		postInfo.Like = slices.Delete(postInfo.Like, index, index+1)
		result = store.db.Save(&postInfo)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
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
	// 检测目标博文是否存在
	var postInfo models.PostInfo
	result := store.db.Where("id = ?", postID).First(&postInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("post not found")
	}
	if result.Error != nil {
		return result.Error
	}

	// 检查是否已经收藏
	var userFavouriteRecord models.UserFavouriteRecord
	result = store.db.Where("uid = ?", uid).First(&userFavouriteRecord)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if result.Error != nil {
		return result.Error
	}
	if slices.Contains(userFavouriteRecord.Favourite, postID) {
		return errors.New("user has favourite this post")
	}

	// 更新用户收藏记录
	userFavouriteRecord.Favourite = append(userFavouriteRecord.Favourite, postID)
	result = store.db.Save(&userFavouriteRecord)
	if result.Error != nil {
		return result.Error
	}

	// 更新博文收藏记录
	postInfo.Favourite = append(postInfo.Favourite, uid)
	result = store.db.Save(&postInfo)
	return result.Error
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
	// 检测目标博文是否存在
	var postInfo models.PostInfo
	result := store.db.Where("id = ?", postID).First(&postInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("post not found")
	}
	if result.Error != nil {
		return result.Error
	}

	// 检查是否已经收藏
	var userFavouriteRecord models.UserFavouriteRecord
	result = store.db.Where("uid = ?", uid).First(&userFavouriteRecord)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if result.Error != nil {
		return result.Error
	}
	if !slices.Contains(userFavouriteRecord.Favourite, postID) {
		return errors.New("user has not favourite this post")
	}

	// 更新用户收藏记录
	index := slices.Index(userFavouriteRecord.Favourite, postID)
	if index != -1 {
		userFavouriteRecord.Favourite = slices.Delete(userFavouriteRecord.Favourite, index, index+1)
		result = store.db.Save(&userFavouriteRecord)
		if result.Error != nil {
			return result.Error
		}
	}

	// 更新博文收藏记录
	index = slices.Index(postInfo.Favourite, uid)
	if index != -1 {
		postInfo.Favourite = slices.Delete(postInfo.Favourite, index, index+1)
		result = store.db.Save(&postInfo)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
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
	// 获取用户点赞记录
	var userLikedRecord models.UserLikedRecord
	result := store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return false, false, result.Error
	}
	liked := slices.Contains(userLikedRecord.LikedPost, postID)

	// 获取用户收藏记录
	var userFavouriteRecord models.UserFavouriteRecord
	result = store.db.Where("uid = ?", uid).First(&userFavouriteRecord)
	if result.Error != nil {
		return false, false, result.Error
	}
	favourite := slices.Contains(userFavouriteRecord.Favourite, postID)

	return liked, favourite, nil
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
