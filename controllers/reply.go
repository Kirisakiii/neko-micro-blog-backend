/*
Package controllers - NekoBlog backend server controllers.
This file is for post controller, which is used to create handlee post related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/
package controllers

import (
	"strconv"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
)

// ReplyController 博文控制器结构体
type ReplyController struct {
	replyService *services.ReplyService
}

// NewReplyController 博文控制器工厂函数。
//
// 返回值：
//   - *ReplyController 博文控制器指针
func (factory *Factory) NewReplyController() *ReplyController {
	return &ReplyController{
		replyService: factory.serviceFactory.NewReplyService(),
	}
}

// NewCreateReplyHandler 是博文控制器的工厂函数，用于创建新的回复处理程序。
//
// 返回值：
//   - fiber.Handler：新增回复的请求handler
func (controller *ReplyController) NewCreateReplyHandler(commentStore *stores.CommentStore, userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.ReplyCreateBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 验证必要参数
		if reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "content is required"),
			)
		}

		if reqBody.CommentID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		// 获取 Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法创建回复
		err := controller.replyService.CreateReply(claims.UID, reqBody.CommentID, reqBody.ParentReplyID, reqBody.Content, commentStore, userStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "reply created successfully"),
		)
	}
}

// DeleteReply 删除回复的请求。
//
// 返回：
//   - fiber.Handler：删除回复的请求handler
func (controller *ReplyController) DeleteReplyHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体中的数据
		reqBody := new(types.UserReplyDeleteBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 检查回复ID是否为空
		if reqBody.ReplyID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "reply id is required"),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 执行删除操作
		if err := controller.replyService.DeleteReply(claims.UID, reqBody.ReplyID); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewUpdateReplyHandler 处理修改回复的请求。
//
// 返回：
//   - fiber.Handler：修改回复的请求handler
func (controller *ReplyController) NewUpdateReplyHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserReplyUpdateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		//校检参数
		if reqBody.Content == "" || reqBody.ReplyID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, " content or reply id is required"),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 调用服务方法修改回复
		err = controller.replyService.UpdateReply(claims.UID, reqBody.ReplyID, reqBody.Content)
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

// NewGetReplyListHandler 处理获取回复列表的请求。
//
// 返回：
//   - fiber.Handler：获取回复列表的请求handler
func (controller *ReplyController) NewGetReplyListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint64, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 调用服务方法获取回复列表
		replyList, err := controller.replyService.GetReplyList(commentIDUint64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewReplyListResponse(replyList)),
		)
	}
}

// NewGetReplyDetailHandler 处理获取回复的请求。
//
// 返回：
//   - fiber.Handler：获取回复的请求handler
func (controller *ReplyController) NewGetReplyDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		replyID := ctx.Query("reply-id")
		if replyID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "reply id is required"),
			)
		}

		replyIDUint64, err := strconv.ParseUint(replyID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 调用服务方法获取回复
		reply, err := controller.replyService.GetReplyDetail(replyIDUint64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewReplyDetailResponse(reply)),
		)
	}
}
