/*
Package stores - NekoBlog backend server data access objects.
This file is for comment storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"errors"
	"slices"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"gorm.io/gorm"
)

// Comment 评论信息数据库
type CommentStore struct {
	db *gorm.DB
}

// NewCommentStore 返回一个新的用户存储实例。
// 返回：
//   - *CommentStore: 返回一个指向新的用户存储实例的指针。
func (factory *Factory) NewCommentStore() *CommentStore {
	return &CommentStore{factory.db}
}

// NewCommentStore 存储comment
//
// 参数 ：- uid：用户id，- username: 用户名，- postID: 博文id，- content: 博文内容
//
// 返回：
//
//	-error 正确返回nil
func (store *CommentStore) CreateComment(uid uint64, username string, postID uint64, content string) (uint64, error) {
	newComment := models.CommentInfo{
		PostID:   postID,
		Username: username,
		Content:  content,
		UID:      uid,
		Like:     nil,
		Dislike:  nil,
		IsPublic: true,
	}

	result := store.db.Create(&newComment)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint64(newComment.ID), nil
}

//	ValidateCommentExistence 判断评论是否存在
//
//	参数：
//	- commentID: 评论ID
//
// 返回值：
//   - error：如果评论存在返回true，不存在判断具体的错误类型返回false
func (store *CommentStore) ValidateCommentExistence(commentID uint64) (bool, error) {
	var comment models.CommentInfo
	result := store.db.Where("id = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	// 返回错误类型
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// UpdateComment 修改评论
//
//	参数：
//	- commentID: 评论ID
//	- content: 修改内容
//
// 返回值：
//   - error：如果评论存在返回true，不存在判断具体的错误类型返回false
func (store *CommentStore) UpdateComment(commentID uint64, content string) error {
	commentInfo := new(models.CommentInfo)
	result := store.db.Where("id = ?", commentID).First(commentInfo)
	if result.Error != nil {
		return result.Error
	}

	commentInfo.Content = content
	result = store.db.Save(commentInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteComment 删除评论
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - error：返回删除处理的成功与否
func (store *CommentStore) DeleteComment(commentID uint64) error {
	return store.db.Where("id = ?", commentID).Unscoped().Delete(&models.CommentInfo{}).Error
}

// GetCommentList 获取评论列表
//
// 返回值：
//   - 成功则返回评论列表
//   - 失败返回nil
func (store *CommentStore) GetCommentList(postID uint64) ([]models.CommentInfo, error) {
	var commentInfos []models.CommentInfo
	result := store.db.Where("post_id = ?", postID).Order("id desc").Find(&commentInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	return commentInfos, nil
}

// GetCommentInfo 获取评论信息
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - models.CommentInfo：成功返回评论信息
//   - error：失败返回error
func (store *CommentStore) GetCommentInfo(commentID uint64) (models.CommentInfo, error) {
	comment := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&comment)
	return comment, result.Error
}

// GetCommentUserStatus 获取评论用户状态
//
// 参数：
//   - commentID：评论ID
//   - uid：用户ID
//
// 返回值：
//   - bool：返回用户状态
func (store *CommentStore) GetCommentUserStatus(uid, commentID uint64) (bool, bool, error) {
	userLikedRecord := models.UserLikedRecord{}
	result := store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return false, false, result.Error
	}
	userDislikeRecord := models.UserDislikeRecord{}
	result = store.db.Where("uid = ?", uid).First(&userDislikeRecord)
	if result.Error != nil {
		return false, false, result.Error
	}

	liked := slices.Index(userLikedRecord.LikedComment, int64(commentID)) != -1
	disliked := slices.Index(userDislikeRecord.DislikeComment, int64(commentID)) != -1
	return liked, disliked, nil
}

// LikeComment 点赞评论
//
// 参数：
//   - commentID：评论ID
//   - uid：用户ID
//
// 返回值：
//   - error：返回点赞处理的成功与否
func (store *CommentStore) LikeComment(uid, commentID uint64) error {
	// 获取评论信息
	commentInfo := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&commentInfo)
	if result.Error != nil {
		return result.Error
	}
	userLikedRecord := models.UserLikedRecord{}
	result = store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}
	userDislikeRecord := models.UserDislikeRecord{}
	result = store.db.Where("uid = ?", uid).First(&userDislikeRecord)
	if result.Error != nil {
		return result.Error
	}

	// 如果点踩列中存在 则先移除点踩记录
	index := slices.Index(userDislikeRecord.DislikeComment, int64(commentID))
	if index != -1 {
		userDislikeRecord.DislikeComment = slices.Delete(userDislikeRecord.DislikeComment, index, index+1)
		result = store.db.Save(&userDislikeRecord)
		if result.Error != nil {
			return result.Error
		}
	}
	if index := slices.Index(commentInfo.Dislike, int64(uid)); index != -1 {
		commentInfo.Dislike = slices.Delete(commentInfo.Dislike, index, index+1)
	}

	// 如果已经点赞过了
	if index := slices.Index(userLikedRecord.LikedComment, int64(uid)); index != -1 {
		return errors.New("you have liked this comment")
	}

	commentInfo.Like = append(commentInfo.Like, int64(uid))
	result = store.db.Save(&commentInfo)
	if result.Error != nil {
		return result.Error
	}

	userLikedRecord.LikedComment = append(userLikedRecord.LikedComment, int64(commentID))
	result = store.db.Save(&userLikedRecord)
	return result.Error
}

// CancelLikeComment 取消点赞评论
//
// 参数：
//   - commentID：评论ID
//   - uid：用户ID
//
// 返回值：
//   - error：返回取消点赞处理的成功与否
func (store *CommentStore) CancelLikeComment(uid, commentID uint64) error {
	commentInfo := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&commentInfo)
	if result.Error != nil {
		return result.Error
	}
	userLikedRecord := models.UserLikedRecord{}
	result = store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}

	index := slices.Index(userLikedRecord.LikedComment, int64(commentID))
	// 如果没有点赞过
	if index == -1 {
		return errors.New("you have not liked this comment")
	}
	userLikedRecord.LikedComment = slices.Delete(userLikedRecord.LikedComment, index, index+1)
	result = store.db.Save(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}

	index = slices.Index(commentInfo.Like, int64(uid))
	if index != -1 {
		commentInfo.Like = slices.Delete(commentInfo.Like, index, index+1)
		result = store.db.Save(&commentInfo)
	}

	return result.Error
}

// DislikeComment 点踩评论
//
// 参数：
//   - commentID：评论ID
//   - uid：用户ID
//
// 返回值：
//   - error：返回点踩处理的成功与否
func (store *CommentStore) DislikeComment(uid, commentID uint64) error {
	commentInfo := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&commentInfo)
	if result.Error != nil {
		return result.Error
	}
	userDislikeRecord := models.UserDislikeRecord{}
	result = store.db.Where("uid = ?", uid).First(&userDislikeRecord)
	if result.Error != nil {
		return result.Error
	}
	userLikedRecord := models.UserLikedRecord{}
	result = store.db.Where("uid = ?", uid).First(&userLikedRecord)
	if result.Error != nil {
		return result.Error
	}

	// 如果点赞列中存在 则先移除点赞记录
	index := slices.Index(userLikedRecord.LikedComment, int64(commentID))
	if index != -1 {
		userLikedRecord.LikedComment = slices.Delete(userLikedRecord.LikedComment, index, index+1)
		result = store.db.Save(&userLikedRecord)
		if result.Error != nil {
			return result.Error
		}
	}
	if index := slices.Index(commentInfo.Like, int64(uid)); index != -1 {
		commentInfo.Like = slices.Delete(commentInfo.Like, index, index+1)
	}

	// 如果已经点踩过了
	if index := slices.Index(userDislikeRecord.DislikeComment, int64(uid)); index != -1 {
		return errors.New("you have disliked this comment")
	}

	commentInfo.Dislike = append(commentInfo.Dislike, int64(uid))
	result = store.db.Save(&commentInfo)
	if result.Error != nil {
		return result.Error
	}

	userDislikeRecord.DislikeComment = append(userDislikeRecord.DislikeComment, int64(commentID))
	result = store.db.Save(&userDislikeRecord)
	return result.Error
}

// CancelDislikeComment 取消点踩评论
//
// 参数：
//   - commentID：评论ID
//   - uid：用户ID
//
// 返回值：
//   - error：返回取消点踩处理的成功与否
func (store *CommentStore) CancelDislikeComment(uid, commentID uint64) error {
	commentInfo := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&commentInfo)
	if result.Error != nil {
		return result.Error
	}
	userDislikeRecord := models.UserDislikeRecord{}
	result = store.db.Where("uid = ?", uid).First(&userDislikeRecord)
	if result.Error != nil {
		return result.Error
	}

	index := slices.Index(userDislikeRecord.DislikeComment, int64(commentID))
	// 如果没有点踩过
	if index == -1 {
		return errors.New("you have not disliked this comment")
	}
	userDislikeRecord.DislikeComment = slices.Delete(userDislikeRecord.DislikeComment, index, index+1)
	result = store.db.Save(&userDislikeRecord)
	if result.Error != nil {
		return result.Error
	}

	index = slices.Index(commentInfo.Dislike, int64(uid))
	if index != -1 {
		commentInfo.Dislike = slices.Delete(commentInfo.Dislike, index, index+1)
		result = store.db.Save(&commentInfo)
	}

	return result.Error
}
