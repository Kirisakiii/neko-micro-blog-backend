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
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/functools"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
)

// PostController 博文控制器结构体
type PostController struct {
	postService *services.PostService
}

// NewPostController 博文控制器工厂函数。
//
// 返回值：
//   - *PostController 博文控制器指针
func (factory *Factory) NewPostController(searchServiceClient search.SearchEngineClient) *PostController {
	return &PostController{
		postService: factory.serviceFactory.NewPostService(searchServiceClient),
	}
}

// NewPostListHandler 博文列表函数
//
// 返回值：
//   - fiber.Handle：新的博文列表函数
func (controller *PostController) NewPostListHandler(userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取请求参数
		reqType := ctx.Query("type")
		uid := ctx.Query("uid")
		length := ctx.Query("len")
		from := ctx.Query("from")
		if reqType == "user" || reqType == "liked" || reqType == "favourited" {
			_, err := strconv.ParseUint(uid, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid uid"),
				)
			}
		}
		if length != "" {
			_, err := strconv.ParseUint(length, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid length"),
				)
			}
		}
		if from != "" {
			_, err := strconv.ParseUint(from, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid from id"),
				)
			}
		}

		// 获取帖子列表
		var (
			posts []int64
			err   error
		)
		switch reqType {
		case "":
			posts, err = controller.postService.GetPostList("all", "", length, from, userStore)
		case "all":
			posts, err = controller.postService.GetPostList("all", "", length, from, userStore)
		case "user":
			posts, err = controller.postService.GetPostList("user", uid, length, from, userStore)
		case "liked":
			posts, err = controller.postService.GetPostList("liked", uid, length, from, userStore)
			posts = functools.Reverse(posts)
		case "favourited":
			posts, err = controller.postService.GetPostList("favourited", uid, length, from, userStore)
			posts = functools.Reverse(posts)
		case "topic":
			topicIDStr := ctx.Query("topic_id")
			topicID, err := primitive.ObjectIDFromHex(topicIDStr)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid topic id"),
				)
			}
			posts, err = controller.postService.GetPostListByTopic(topicID, from, length)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
				)
			}
		default:
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "invalid type"))
		}
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostListResponse(posts)),
		)
	}
}

// NewFollowPostListHandler 获取关注用户的帖子列表
//
// 返回值：
//   - fiber.Handler：新的获取关注用户的帖子列表函数
func (controller *PostController) NewFollowPostListHandler(followStore *stores.FollowStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取请求用户的 UID
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取请求参数
		length := ctx.Query("len")
		from := ctx.Query("from")
		var lengthUint uint64 = 10
		if length != "" {
			var err error
			lengthUint, err = strconv.ParseUint(length, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid length"),
				)
			}
		}
		if from != "" {
			_, err := strconv.ParseUint(from, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid from id"),
				)
			}
		}

		// 获取帖子列表
		posts, err := controller.postService.GetFollowPostList(claims.UID, from, lengthUint, followStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostListResponse(posts)),
		)
	}
}

// NewDetailHandler 获取文章信息的函数
//
// 返回值：
//   - fiber.Handler：新的获取文章信息的函数
func (controller *PostController) NewPostDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 根据 UID 获取用户信息
		postIDString := ctx.Params("post")
		if postIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post id is required"),
			)
		}

		//获取帖子的唯一标识符
		postID, err := strconv.ParseUint(postIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取帖子的详细信息
		post, likeCount, favouriteCount, err := controller.postService.GetPostInfo(postID)
		// 若post不存在则返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post does not exist"),
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
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostDetailResponse(post, likeCount, favouriteCount)),
		)
	}
}

// NewCreatePostHandler 返回一个用于处理创建博文请求的 Fiber 处理函数
//
// 参数：
//   - topicStore：话题存储
//
// 返回值：
//   - fiber.Handler：新的创建博文函数
func (controller *PostController) NewCreatePostHandler(topicStore *stores.TopicStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析用户请求
		reqBody := types.PostCreateBody{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 验证参数
		if reqBody.Title == "" || reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post title or post content is required"),
			)
		}
		if len(reqBody.Images) > 9 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post images count exceeds the limit"),
			)
		}

		// 创建博文
		postInfo, err := controller.postService.CreatePost(claims.UID, ctx.IP(), reqBody, topicStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回成功响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(
				consts.SUCCESS,
				"post created successfully",
				serializers.NewCreatePostResponse(postInfo),
			),
		)
	}
}

// NewPostUserStatusHandler 返回一个用于处理获取用户对帖子的状态请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的获取用户对帖子的状态函数
func (controller *PostController) NewPostUserStatusHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Query("post-id")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 获取用户对帖子的状态
		isLiked, isFavourited, err := controller.postService.GetPostUserStatus(int64(claims.UID), int64(postIDUint))
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.Status(200).JSON(serializers.NewResponse(
			consts.SUCCESS,
			"succeed",
			serializers.NewPostUserStatus(
				postIDUint,
				claims.UID,
				isLiked,
				isFavourited,
			),
		),
		)
	}
}

// NewUploadPostImageFromFileHandler 返回一个用于处理上传博文图片请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的上传博文图片函数
func (controller *PostController) NewUploadPostImageFromFileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 接收文件
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验文件数量
		files := form.File["file"]
		if len(files) < 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "image is required"),
			)
		}
		if len(files) > 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "the number of image cannot exceed 1"),
			)
		}

		UUID, err := controller.postService.UploadPostImage(files[0])
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(serializers.NewResponse(
			consts.SUCCESS,
			"succeed",
			serializers.NewUploadPostImageResponse(UUID),
		))
	}
}

// NewUploadPostImageFromURLHandler
//
// 返回值：
//   - fiber.Handler：新的上传博文图片函数
func (controller *PostController) NewUploadPostImageFromURLHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求
		reqBody := types.PostURLImageBody{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 上传图片
		UUID, err := controller.postService.UploadPostImageFromURL(reqBody.URL)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(serializers.NewResponse(
			consts.SUCCESS,
			"succeed",
			serializers.NewUploadPostImageResponse(UUID),
		))
	}
}

// NewForwardPostHandler 返回一个用于处理转发博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的转发博文函数
func (controller *PostController) NewForwardPostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析请求
		reqBody := struct {
			Content string `json:"content"`
			PostID  uint64 `json:"post_id"`
		}{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 验证参数
		if reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post content is required"),
			)
		}

		// 执行转发操作
		err = controller.postService.ForwardPost(claims.UID, ctx.IP(), reqBody.PostID, reqBody.Content)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回成功响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(
				consts.SUCCESS,
				"",
			),
		)
	}
}

// NewLikePostHandler 返回一个用于处理点赞博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的博文点赞函数
func (controller *PostController) NewLikePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Query("post-id")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行点赞操作
		if err := controller.postService.LikePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewCancelLikePostHandler 返回一个用于处理取消点赞博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的取消点赞函数
func (controller *PostController) NewCancelLikePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Query("post-id")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行取消点赞操作
		if err := controller.postService.CancelLikePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewFavouritePostHandler 返回一个用于处理收藏博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的收藏博文函数
func (controller *PostController) NewFavouritePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Query("post-id")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行收藏操作
		if err := controller.postService.FavouritePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewCancelFavouritePostHandler 返回一个用于处理取消收藏博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的取消收藏博文函数
func (controller *PostController) NewCancelFavouritePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Query("post-id")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行取消收藏操作
		if err := controller.postService.CancelFavouritePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

// NewDeletePostHandler 返回一个用于处理删除博文请求的 Fiber 处理函数
//
// 返回值：
//   - fiber.Handler：新的博文删除函数
func (controller *PostController) NewDeletePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取PostID
		postID := ctx.Params("post")

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行删除操作
		if err := controller.postService.DeletePost(postIDUint, claims.UID); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}
