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
//   - fiber.Handler：回复处理程序
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
				serializers.NewResponse(consts.PARAMETER_ERROR, " content is required"),
			)
		}

		if *reqBody.CommentID == 0&&*reqBody.ReplyToReplyID == 0 {
		    return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, " commentID or replyToReplyID is required"),
			)
		}


		// 获取 Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)
		if claims == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not available"),
			)
		}


		// 调用服务方法创建回复
		err := controller.replyService.CreateReply(claims.UID, reqBody.CommentID, reqBody.ReplyToReplyID, reqBody.Content, commentStore, userStore)
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
//   - 处理的成功和失败
func (controller *ReplyController) DeleteReplyHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 解析请求体中的数据
		reqBody := new(types.UserReplyDeleteBody)
		if err := c.BodyParser(reqBody); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 检查回复ID是否为空
		if reqBody.ReplyID == nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "reply id is required"),
			)
		}

		// 执行删除操作
		if err := controller.replyService.DeleteReply(*reqBody.ReplyID); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return c.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewUpdateReplyHandler 处理修改回复的请求。
//
// 返回：
//   - 处理的成功和失败
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
        if reqBody.Content == "" ||reqBody.ReplyID == nil {
            return ctx.Status(200).JSON(
                serializers.NewResponse(consts.PARAMETER_ERROR, " content or replyID is required"),
            )
        }

        // 获取Token Claims
        claims := ctx.Locals("claims").(*types.BearerTokenClaims)
        if claims == nil {
            return ctx.Status(200).JSON(
                serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not available"),
            )
        }

        // 调用服务方法修改回复
        err = controller.replyService.UpdateReply(*reqBody.ReplyID, reqBody.Content)
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