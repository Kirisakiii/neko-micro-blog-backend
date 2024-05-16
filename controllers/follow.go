/*
Package controllers - NekoBlog backend server controllers.
This file is for follow controller, which is used to create handlee follow related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/
package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
)

// FollowController 博文控制器结构体
type FollowController struct {
	followService *services.FollowService
}

// NewFollowController 创建博文控制器实例
//
// 返回：
//   - *FollowController: 返回一个新的评论控制器实例。
func (factory *Factory) NewFollowController() *FollowController {
	return &FollowController{
		followService: factory.serviceFactory.NewFollowService(),
	}
}

// NewCreateFollowHandler 返回一个用于处关注用户请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的关注用户函数
func (controller *FollowController) NewCreateFollowHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取followedID
		body := struct {
			UserID uint64 `json:"user_id" form:"user_id"`
		}{}
		err := ctx.BodyParser(&body)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"))
		}
		followedID := body.UserID

		// 执行关注操作
		if err := controller.followService.FollowUser(claims.UID, followedID); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewCancelFollowHandler 返回一个用于处理取消关注用户请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的取消关注用户函数
func (controller *FollowController) NewCancelFollowHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取followedID
		body := struct {
			UserID uint64 `json:"user_id" form:"user_id"`
		}{}
		err := ctx.BodyParser(&body)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"))
		}
		followedID := body.UserID

		// 执行取消关注操作
		if err := controller.followService.CancelFollowUser(claims.UID, followedID); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewFollowListHandler() 返回一个用于处理获取关注列表请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的获取关注列表函数
func (controller *FollowController) NewFollowListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取查询UID
		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		// 执行获取关注列表操作
		follows, err := controller.followService.GetFollowList(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewFollowListResponse(follows)),
		)
	}
}

// NewFollowCountHandler 返回一个用于处理获取关注人数请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的获取关注人数函数
func (controller *FollowController) NewFollowCountHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取查询UID
		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		// 执行获取关注人数操作
		count, err := controller.followService.GetFollowCountByUID(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", struct {Count int64}{count}),
		)
	}
}

// NewFollowerListHandler 返回一个用于处理获取粉丝列表请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的获取粉丝列表函数
func (controller *FollowController) NewFollowerListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取查询UID
		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		// 执行获取粉丝列表操作
		followers, err := controller.followService.GetFollowerList(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewFollowerListResponse(followers)),
		)
	}
}

// NewFollowerCountHandler 返回一个用于处理获取粉丝人数请求的 Fiber 处理函数
//
// 返回：
//   - fiber.Handler: 新的获取粉丝人数函数
func (controller *FollowController) NewFollowerCountHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取查询UID
		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		// 执行获取关注人数操作
		count, err := controller.followService.GetFollowerCountByUID(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", struct {Count int64}{count}),
		)
	}
}
