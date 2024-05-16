/*
Package services - NekoBlog backend server services.
This file is for user related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/

package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/converters"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/validers"
	"github.com/lib/pq"
)

// PostService 博文服务
type PostService struct {
	postStore           *stores.PostStore
	searchServiceClient search.SearchEngineClient
}

// PostService 返回一个新的 PostService 实例
//
// 返回值：
//   - *PostService：新的 PostService 实力。
func (factory *Factory) NewPostService(searchServiceClient search.SearchEngineClient) *PostService {
	return &PostService{
		postStore:           factory.storeFactory.NewPostStore(),
		searchServiceClient: searchServiceClient,
	}
}

// GetPostList 获取适用于用户查看的帖子信息列表。
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在获取帖子信息过程中遇到的任何错误，如果有的话。
func (service *PostService) GetPostList(reqType, uid, length, from string, userStore *stores.UserStore) ([]int64, error) {
	var (
		postInfos  []models.PostInfo
		userRecord pq.Int64Array
		err        error
		uidInt64   int64
	)

	if reqType != "all" {
		uidInt64, err = strconv.ParseInt(uid, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	queryLenth := 10
	if length != "" {
		queryLenth, err = strconv.Atoi(length)
		if err != nil {
			return nil, err
		}
		if queryLenth > 10 {
			queryLenth = 10
		}
	}

	switch reqType {
	case "all":
		postInfos, err = service.postStore.GetPostList(from, queryLenth)
	case "user":
		postInfos, err = service.postStore.GetPostListByUID(uid)
	case "liked":
		userRecord, err = userStore.GetUserLikedRecord(uidInt64)
	case "favourited":
		userRecord, err = userStore.GetUserFavoriteRecord(uidInt64)
	}
	if err != nil {
		return nil, err
	}

	if reqType == "all" || reqType == "user" {
		postIDs := make([]int64, len(postInfos))
		for index, post := range postInfos {
			postIDs[index] = int64(post.ID)
		}
		return postIDs, nil
	}

	if userRecord == nil {
		return nil, nil
	}
	postIDs := make([]int64, len(userRecord))
	for index, id := range userRecord {
		postIDs[index] = int64(id)
	}
	return postIDs, nil
}

// GetPostInfoByUsername 根据用户名获取用户信息。
//
// 参数：
//   - UID：用户ID
//
// 返回值：
//   - *models.postInfo：用户信息模型。
func (service *PostService) GetPostInfo(postID uint64) (models.PostInfo, int64, int64, error) {
	return service.postStore.GetPostInfo(postID)
}

// CreatePost 根据用户提交的帖子信息创建帖子。
//
// 参数：
//   - userID：用户ID，用于关联帖子与用户。
//   - ipAddr：IP地址
//   - postReqInfo：帖子信息，包含标题、内容等。
//   - postImages:帖子图片
//
// 返回值：
//   - error：如果在创建过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *PostService) CreatePost(uid uint64, ipAddr string, postReqInfo types.PostCreateBody) (models.PostInfo, error) {
	// 校验图片是否可用
	for _, image := range postReqInfo.Images {
		existence, err := service.postStore.CheckCacheImageAvaliable(image)
		if err != nil {
			return models.PostInfo{}, err
		}
		if !existence {
			return models.PostInfo{}, errors.New("image does not exist")
		}
	}

	// 调用存储层的方法创建帖子
	postInfo, err := service.postStore.CreatePost(uid, ipAddr, postReqInfo)
	if err != nil {
		return models.PostInfo{}, err
	}

	// 写入搜索引擎索引库
	_, err = service.searchServiceClient.CreatePostIndex(context.TODO(), &search.CreatePostIndexRequest{
		Id:      int64(postInfo.ID),
		Title:   postReqInfo.Title,
		Content: postReqInfo.Content,
	})
	if err != nil {
		return models.PostInfo{}, err
	}

	return postInfo, nil
}

// UploadPostImage 上传博文图片
//
// 参数：
//   - postImage：待上传的博文图片
//
// 返回值：
//   - string：图片UUID
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) UploadPostImage(postImage *multipart.FileHeader) (string, error) {
	// 打开图像文件
	imageFile, err := postImage.Open()
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	// 验证图像文件的有效性，包括尺寸和文件类型等
	fileType, err := validers.ValidImageFile(
		postImage,
		&imageFile,
		consts.POST_IMAGE_MIN_WIDTH,
		consts.POST_IMAGE_MIN_HEIGHT,
		consts.POST_IMAGE_MAX_FILE_SIZE,
	)
	if err != nil {
		return "", err
	}

	// 调整图片大小
	convertedImage, err := converters.ResizePostImage(fileType, &imageFile)
	if err != nil {
		return "", err
	}

	// 保存在暂存区并返回UUID
	return service.postStore.CachePostImage(convertedImage)
}

// UploadPostImageFromURL 上传博文图片
//
// 参数：
//   - imageURL：待上传的博文图片URL
//
// 返回值：
//   - string：图片UUID
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) UploadPostImageFromURL(imageURL string) (string, error) {
	// 下载图片
	client := &http.Client{}
	if os.Getenv("HTTP_PROXY") != "" {
		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	}

	resp, err := client.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download image")
	}

	// 将reader的数据读取到内存中
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 验证图像文件的有效性，包括尺寸和文件类型等
	fileType, err := validers.ValidImageFileFromReader(
		data,
		resp.Header.Get("Content-Type"),
		consts.POST_IMAGE_MIN_WIDTH,
		consts.POST_IMAGE_MIN_HEIGHT,
		consts.POST_IMAGE_MAX_FILE_SIZE,
	)
	if err != nil {
		return "", err
	}

	// 调整图片大小
	convertedImage, err := converters.ResizeURLImage(fileType, data)
	if err != nil {
		return "", err
	}

	// 保存在暂存区并返回UUID
	return service.postStore.CachePostImage(convertedImage)
}

// LikePost 点赞博文
//
// 参数：
//   - postID：待点赞博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) LikePost(uid, postID int64) error {
	// 调用post存储中的点赞方法
	return service.postStore.LikePost(uid, postID)
}

// CancelLikePost 取消点赞博文
//
// 参数：
//   - uid：用户ID
//   - postID：待取消点赞博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) CancelLikePost(uid, postID int64) error {
	// 调用post存储中的取消点赞方法
	return service.postStore.CancelLikePost(uid, postID)
}

// FavouritePost 收藏博文
//
// 参数：
//   - uid：用户ID
//   - postID：待收藏博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) FavouritePost(uid, postID int64) error {
	// 调用post存储中的收藏方法
	return service.postStore.FavouritePost(uid, postID)
}

// CancelFavouritePost 取消收藏博文
//
// 参数：
//   - uid：用户ID
//   - postID：待取消收藏博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) CancelFavouritePost(uid, postID int64) error {
	// 调用post存储中的取消收藏方法
	return service.postStore.CancelFavouritePost(uid, postID)
}

// GetPostUserStatus 获取用户对帖子的状态
//
// 参数：
//   - uid：用户ID
//   - postID：帖子ID
//
// 返回值：
//   - bool：用户是否点赞
//   - bool：用户是否收藏
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) GetPostUserStatus(uid, postID int64) (bool, bool, error) {
	// 调用post存储中的获取用户帖子状态方法
	return service.postStore.GetPostUserStatus(uid, postID)
}

// DeletePost 删除博文
//
// 参数：
//   - postID uint64：待删除博文的ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) DeletePost(postID uint64) error {
	// 调用post存储中的删除post方法
	return service.postStore.DeletePost(postID)
}
