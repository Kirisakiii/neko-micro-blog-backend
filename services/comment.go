/*
Package services - NekoBlog backend server services.
This file is for comment related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package services

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

// CommentService 评论服务
type CommentService struct {
	commentStore *stores.CommentStore
}

// NewCommentService 返回一个新的评论服务实例。
//
// 返回：
//   - *CommentService: 返回一个指向新的评论服务实例的指针。
func (factory *Factory) NewCommentService() *CommentService {
	return &CommentService{
		commentStore: factory.storeFactory.NewCommentStore(),
	}
}

// NewCommentService 创建评论
//
// 参数：
//   - uid：用户ID
//   - postID: 博文编号
//   - content: 博文内容
//   - postStore，userStore：绑定post和user层来调用方法
//
// 返回值：
//
//	-error 创建失败返回创建失败时候的具体信息
func (service *CommentService) CreateComment(uid uint64, postID uint64, content string, postStore *stores.PostStore, userStore *stores.UserStore) (uint64, error) {
	// 校验评论是否存在
	existance, err := postStore.ValidatePostExistence(postID)
	if err != nil {
		return 0, err
	}
	if !existance {
		return 0, errors.New("post does not exist")
	}

	// 根据 UID 获取 Username
	user, err := userStore.GetUserByUID(uid)
	if err != nil {
		return 0, err
	}

	// 调用存储层的方法存储评论
	commentID, err := service.commentStore.CreateComment(uid, user.UserName, postID, content)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

// UpdateComment 修改评论
//
// 参数：
//   - comment：评论ID
//   - content: 博文内容
//
// 返回值：
//
//	-error 如果评论存在返回修改评论时候的信息
func (service *CommentService) UpdateComment(commentID uint64, content string) error {
	//检查评论是否存在
	// TODO: 改为控制器层判断
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用数据库或其他存储方法更新评论内容
	err = service.commentStore.UpdateComment(commentID, content)
	if err != nil {
		return err
	}

	// 如果更新成功，返回nil
	return nil
}

// DeleteComment 删除评论
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理删除的信息
func (service *CommentService) DeleteComment(commentID uint64) error {
	// TODO: 改为控制器层判断
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}
	// 调用评论存储中的删除评论方法
	err = service.commentStore.DeleteComment(commentID)
	if err != nil {
		// 如果发生错误，则返回错误
		return err
	}

	// 如果没有发生错误，则返回 nil
	return nil
}

// GetCommentList 获取评论列表
//
// 返回值：
//   - 成功则返回评论列表
//   - 失败返回nil
func (service *CommentService) GetCommentList(postID uint64) ([]models.CommentInfo, error) {
	return service.commentStore.GetCommentList(postID)
}

// GetCommentInfo 获取评论信息
//
// 返回值：
//   - 成功返回评论体
//   - 失败返回nil
func (service *CommentService) GetCommentInfo(commentID uint64) (models.CommentInfo, int64, error) {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return models.CommentInfo{}, 0, err
	}
	if !exists {
		return models.CommentInfo{}, 0, errors.New("comment does not exist")
	}

	return service.commentStore.GetCommentInfo(commentID)
}

// GetCommentUserStatus
//
// 参数：
//   - uid: 用户ID
//   - commentID: 评论ID
//
// 返回值：
//   - bool: 是否点赞
//   - bool: 是否点踩
//   - error: 错误信息
func (service *CommentService) GetCommentUserStatus(uid, commentID uint64) (bool, bool, error) {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return false, false, err
	}
	if !exists {
		return false, false, errors.New("comment does not exist")
	}

	// 调用存储层的方法获取用户对评论的点赞和点踩状态
	isLiked, isDisliked, err := service.commentStore.GetCommentUserStatus(uid, commentID)
	if err != nil {
		return false, false, err
	}

	// 返回用户对评论的点赞和点踩状态
	return isLiked, isDisliked, nil
}

// LikeComment 点赞评论
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理点赞的信息
func (service *CommentService) LikeComment(uid, commentID uint64) error {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用存储层的方法点赞评论
	err = service.commentStore.LikeComment(uid, commentID)
	if err != nil {
		return err
	}

	// 如果点赞成功，返回nil
	return nil
}

// CancelLikeComment
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理取消点赞的信息
func (service *CommentService) CancelLikeComment(uid, commentID uint64) error {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用存储层的方法取消点赞评论
	err = service.commentStore.CancelLikeComment(uid, commentID)
	if err != nil {
		return err
	}

	// 如果取消点赞成功，返回nil
	return nil
}

// DislikeComment 点踩评论
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理点踩的信息
func (service *CommentService) DislikeComment(uid, commentID uint64) error {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用存储层的方法点踩评论
	err = service.commentStore.DislikeComment(uid, commentID)
	if err != nil {
		return err
	}

	// 如果点踩成功，返回nil
	return nil
}

// CancelDislikeComment 取消点踩评论
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理取消点踩的信息
func (service *CommentService) CancelDislikeComment(uid, commentID uint64) error {
	// 检查评论是否存在
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用存储层的方法取消点踩评论
	err = service.commentStore.CancelDislikeComment(uid, commentID)
	if err != nil {
		return err
	}

	// 如果取消点踩成功，返回nil
	return nil
}
