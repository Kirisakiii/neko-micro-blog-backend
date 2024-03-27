/*
Package controllers - NekoBlog backend server controllers.
This file is for search controller.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package controllers

import (
	"net/url"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
)

type SearchController struct {
	searchService *services.SearchService
}

func (factory *Factory) NewSearchController(searchServiceClient search.SearchEngineClient) *SearchController {
	return &SearchController{
		searchService: factory.serviceFactory.NewSearchService(searchServiceClient),
	}
}

// NewSearchPostHandler 创建一个新的搜索帖子的handler
func (controller *SearchController) NewSearchPostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求参数
		queryString := ctx.Query("q")
		if queryString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "query content is required"),
			)
		}

		// 解码URL
		decodedQueryString, err := url.QueryUnescape(queryString)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "query content is invalid"),
			)
		}

		result, err := controller.searchService.SearchPost(decodedQueryString)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "", serializers.NewPostListResponse(result.Ids)),
		)
	}
}
