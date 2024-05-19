/*
Package controllers - NekoBlog backend server controllers.
This file is for topic controller, which is used to create handlee topic related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
)

// TopicController 用户控制器
type TopicController struct {
	topicService *services.TopicService
}

// NewTopicController 返回一个新的 TopicController 实例。
//
// 返回值：
//   - *TopicController：新的 TopicController 实例。
func (factory *Factory) NewTopicController() *TopicController {
	return &TopicController{
		topicService: factory.serviceFactory.NewTopicService(),
	}
}

// NewCreateTopicHandler 返回获取用户资料的处理函数。
//
// 参数：
//   - postStore *stores.PostStore：用于获取帖子信息的存储实例。
//
// 返回值：
//   - fiber.Handler：新的获取用户资料的处理函数。
func (controller *TopicController) NewCreateTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析请求
		reqBody := new(types.TopicCreateBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.PostID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post_id is required"),
			)
		}
		if reqBody.Description == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "description is required"),
			)
		}

		// 创建目标
		tagID, err := controller.topicService.CreateTopic(claims.UID, reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic created successfully", serializers.NewCreateTopicResponse(tagID)),
		)
	}
}

// DeleteTopicHandler 返回删除目标的处理函数
//
// 返回值：
//   - fiber.Handler：删除目标的处理函数。
func (controller *TopicController) DeleteTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析删除目标请求
		reqBody := new(types.TopicDeleteBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 删除目标
		err := controller.topicService.DeleteTarget(reqBody.TopicID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic deleted successfully"),
		)
	}
}

// GetTopicdetailHandler 获取目标列表的处理函数
//
// 返回值：
//   - fiber.Handler：获取目标列表的处理函数。
func (controller *TopicController) GetTopicdetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//解析请求
		reqBody := new(types.TopicListBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}
		//调用服务层获取目标列表
		tagList, err := controller.topicService.GetTopicList(reqBody.TopicID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		//返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic details successfully", serializers.NewTopicListResponse(tagList)),
		)

	}
}

// NewLikeTopicHandler 返回点赞目标的处理函数
//
// 返回值：
//   - fiber.Handler：点赞目标的处理函数
func (controller *TopicController) NewLikeTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析点赞目标请求
		reqBody := new(types.TopicLikeBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 点赞目标
		err := controller.topicService.NewLikeTopicHandler(reqBody.TopicID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic liked successfully"),
		)
	}
}

// NewCancelLikeTopicHandler 返回取消点赞目标的处理函数
//
// 返回值：
//   - fiber.Handler：取消点赞目标的处理函数
func (controller *TopicController) NewCancelLikeTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析取消点赞目标请求
		reqBody := new(types.TopicLikeBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 取消点赞目标
		err := controller.topicService.NewCancelLikeTopicHandler(reqBody.TopicID)

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic cance like successfully"),
		)
	}
}

// NewDislikeTopicHandler 返回点踩数的处理函数
//
// 返回值：
//   - fiber.Handler：点踩数的处理函数
func (controller *TopicController) NewDislikeTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析点踩目标请求
		reqBody := new(types.TopicDisLikeBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 点踩目标
		err := controller.topicService.NewDislikeTopicHandler(reqBody.TopicID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic disliked successfully"),
		)
	}
}

// GetTopicListHandler 返回话题列表
//
// 返回值：
//   - fiber.Handler：话题列表的处理函数
func (controller *TopicController) GetHotTopicsHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		limit := 10 // 设置默认的获取热门话题的数量

		topics, err := controller.topicService.GetHotTopics(limit)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(fiber.StatusOK).JSON(
			serializers.NewResponse(consts.SUCCESS, "success", topics),
		)
	}
}
