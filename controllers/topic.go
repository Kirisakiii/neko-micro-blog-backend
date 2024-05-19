/*
Package controllers - NekoBlog backend server controllers.
This file is for topic controller, which is used to create handlee topic related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"strconv"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// TopicListHandler 话题列表处理函数
//
// 返回值：
//   - fiber.Handler：话题列表处理函数
func (controller *TopicController) TopicListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求
		from := ctx.Query("from")
		length := ctx.Query("length")

		var fromObejctID primitive.ObjectID
		if from == "" {
			fromObejctID = primitive.NewObjectID()
		} else {
			var err error
			fromObejctID, err = primitive.ObjectIDFromHex(from)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
				)
			}
		}

		lenghtUint, err := strconv.ParseUint(length, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if (lenghtUint > 20) || (lenghtUint == 0) {
			lenghtUint = 20
		}

		// 调用服务层获取话题列表
		topics, err := controller.topicService.GetTopicList(fromObejctID, lenghtUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "", serializers.NewTopicListResponse(topics)),
		)
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
		reqBody := types.TopicCreateBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 检查请求参数
		if reqBody.Name == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "name is required"),
			)
		}
		if reqBody.Description == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "description is required"),
			)
		}
		groupID, err := primitive.ObjectIDFromHex(reqBody.BundledGroupID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 创建目标
		tagID, err := controller.topicService.CreateTopic(claims.UID, reqBody, groupID)
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
		reqBody := types.TopicDeleteBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}
		topicID, err := primitive.ObjectIDFromHex(reqBody.TopicID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 删除目标
		err = controller.topicService.DeleteTopic(topicID)
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

// GetTopicDetailHandler 获取目标列表的处理函数
//
// 返回值：
//   - fiber.Handler：获取目标列表的处理函数。
func (controller *TopicController) GetTopicDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求
		objectIDString := ctx.Query("topic_id")
		objectID, err := primitive.ObjectIDFromHex(objectIDString)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		//调用服务层获取目标列表
		tagList, err := controller.topicService.GetTopicDetail(objectID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		//返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "", serializers.NewTopicDetailResponse(tagList)),
		)

	}
}

// NewLikeTopicHandler 返回点赞目标的处理函数
//
// 返回值：
//   - fiber.Handler：点赞目标的处理函数
func (controller *TopicController) NewLikeTopicHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析点赞目标请求
		reqBody := types.TopicLikeBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 点赞目标
		err := controller.topicService.NewLikeTopic(reqBody.TopicID, claims.UID)
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
		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析取消点赞目标请求
		reqBody := types.TopicLikeBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 取消点赞目标
		err := controller.topicService.NewCancelLikeTopic(reqBody.TopicID, claims.UID)

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
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析点踩目标请求
		reqBody := types.TopicDisLikeBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 点踩目标
		err := controller.topicService.NewDislikeTopic(reqBody.TopicID, claims.UID)
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

// NewCancelDislikeHandler 返回取消点踩数的处理函数
//
// 返回值：
//   - fiber.Handler：取消点踩数的处理函数
func (controller *TopicController) NewCancelDislikeHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析取消点踩目标请求
		reqBody := types.TopicDisLikeBody{}
		if err := ctx.BodyParser(&reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 取消点踩目标
		err := controller.topicService.CancelDislikeTopic(reqBody.TopicID, claims.UID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "topic cancel dislike successfully"),
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
