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

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/lib/pq"
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
		Like:     pq.Int64Array{},
		Dislike:  pq.Int64Array{},
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
	var count int64
	result := store.db.Model(&models.CommentInfo{}).Where("id = ? AND ? = ANY(\"like\")", commentID, uid).Count(&count)
	if result.Error != nil {
		return false, false, result.Error
	}
	isLiked := count > 0

	result = store.db.Model(&models.CommentInfo{}).Where("id = ? AND ? = ANY(\"dislike\")", commentID, uid).Count(&count)
	if result.Error != nil {
		return false, false, result.Error
	}
	isDisliked := count > 0

	return isLiked, isDisliked, nil
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
	tx := store.db.Begin()

	// 删除点踩记录
	result := tx.Model(&models.CommentInfo{}).Where("id = ? AND ARRAY[?::bigint] <@ comment_infos.\"dislike\"", commentID, uid).Update("dislike", gorm.Expr("array_remove(\"dislike\", ?)", uid))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// 更新博文点赞记录
	result = tx.Model(&models.CommentInfo{}).Where("id = ? AND NOT ARRAY[?::bigint] <@ comment_infos.\"like\"", commentID, uid).Update("like", gorm.Expr("array_append(\"like\", ?)", uid))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user has liked this comment")
	}

	return tx.Commit().Error
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
	// 更新博文点赞记录
	result := store.db.Model(&models.CommentInfo{}).Where("id = ? AND ARRAY[?::bigint] <@ comment_infos.\"like\"", commentID, uid).Update("like", gorm.Expr("array_remove(\"like\", ?)", uid))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user has not liked this comment")
	}

	return nil
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
	tx := store.db.Begin()

	// 删除点赞记录
	result := tx.Model(&models.CommentInfo{}).Where("id = ? AND ARRAY[?::bigint] <@ comment_infos.\"like\"", commentID, uid).Update("like", gorm.Expr("array_remove(\"like\", ?)", uid))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// 更新博文点踩记录
	result = tx.Model(&models.CommentInfo{}).Where("id = ? AND NOT ARRAY[?::bigint] <@ comment_infos.\"dislike\"", commentID, uid).Update("dislike", gorm.Expr("array_append(\"dislike\", ?)", uid))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user has disliked this comment")
	}

	return tx.Commit().Error
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
	// 更新博文点踩记录
	result := store.db.Model(&models.CommentInfo{}).Where("id = ? AND ARRAY[?::bigint] <@ comment_infos.\"dislike\"", commentID, uid).Update("dislike", gorm.Expr("array_remove(\"dislike\", ?)", uid))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user has not disliked this comment")
	}

	return nil
}
