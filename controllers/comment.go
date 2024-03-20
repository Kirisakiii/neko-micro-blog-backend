/*
Package controllers - NekoBlog backend server controllers.
This file is for comment controller, which is used to create handlee comment related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"errors"
	"strconv"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController 创建一个新的评论控制器实例。
//
// 返回：
//   - *CommentController: 返回一个新的评论控制器实例。
func (factory *Factory) NewCommentController() *CommentController {
	return &CommentController{
		commentService: factory.serviceFactory.NewCommentService(),
	}
}

// NewCreateCommentHandler 处理创建评论的请求。
//
// 参数：- postStore，userStore：绑定post和user层 来调用其中方法
//
// 返回：
//   - 处理的成功和失败
func (controller *CommentController) NewCreateCommentHandler(postStore *stores.PostStore, userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserCommentCreateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.Content == "" || reqBody.PostID == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post_id or content is required"),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法创建评论
		commentID, err := controller.commentService.CreateComment(claims.UID, *reqBody.PostID, reqBody.Content, postStore, userStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "comment created successfully", serializers.NewCreateCommentResponse(commentID)),
		)
	}
}

// NewUpdateCommentHandler 处理修改评论的请求。
//
// 返回：
//   - 处理的成功和失败
// TODO: 未完成
func (controller *CommentController) NewUpdateCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserCommentUpdateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		//校检参数
		if reqBody.Content == "" || reqBody.CommentID == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "content or comment id is required"),
			)
		}

		// 获取Token Claims
		// claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法修改评论
		err = controller.commentService.UpdateComment(*reqBody.CommentID, reqBody.Content)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// DeleteCommentHandler 处理删除评论的请求
//
// 返回：
//   - 返回是否成功处理
func (controller *CommentController) DeleteCommentHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析请求体中的数据
		reqBody := new(types.UserCommentDeleteBody)
		if err := c.BodyParser(reqBody); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 检查评论ID是否为空
		if reqBody.CommentID == nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		// 执行删除操作
		if err := controller.commentService.DeleteComment(*reqBody.CommentID); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return c.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewCommentListHandler 下拉评论列表请求
//
// 返回值：
//   - 成功则返回评论列表
//   - 失败返回nil
func (controller *CommentController) NewCommentListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postID := c.Query("post-id")
		if postID == "" {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post id is required"),
			)
		}
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		comments, err := controller.commentService.GetCommentList(postIDUint)
		if err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return c.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentListResponse(comments)),
		)
	}
}

// NewCommentDetailHandler 获取文章信息的函数
//
// 返回值：
//   - fiber.Handler：新的获取文章信息的函数
func (controller *CommentController) NewCommentDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentIDString := ctx.Query("comment-id")
		if commentIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentID, err := strconv.ParseUint(commentIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		comment, err := controller.commentService.GetCommentInfo(commentID)
		// 若comment不存在则返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment does not exist"),
			)
		}
		// 返回其他错误
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentDetailResponse(comment)),
		)
	}
}

// NewCommentUserStatusHandler 获取用户对评论的状态
//
// 返回值：
//   - fiber.Handler：新的获取用户对评论的状态的函数
func (controller *CommentController) NewCommentUserStatusHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法获取用户对评论的状态
		liked, disliked, err := controller.commentService.GetCommentUserStatus(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentUserStatusResponse(liked, disliked)),
		)
	}
}


// NewLikeCommentHandler 点赞评论的函数
//
// 返回值：
//   - fiber.Handler：新的点赞评论的函数
func (controller *CommentController) NewLikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法点赞评论
		err = controller.commentService.LikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewCancelLikeCommentHandler 取消点赞评论的函数
//
// 返回值：
//   - fiber.Handler：新的取消点赞评论的函数
func (controller *CommentController) NewCancelLikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法取消点赞评论
		err = controller.commentService.CancelLikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewDislikeCommentHandler 点踩评论的函数
//
// 返回值：
//   - fiber.Handler：新的点踩评论的函数
func (controller *CommentController) NewDislikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法点踩评论
		err = controller.commentService.DislikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewCancelDislikeCommentHandler
//
// 返回值：
//   - fiber.Handler：新的取消点踩评论的函数
func (controller *CommentController) NewCancelDislikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法取消点踩评论
		err = controller.commentService.CancelDislikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}
