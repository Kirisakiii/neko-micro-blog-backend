/*
Package services - NekoBlog backend server services.
This file is for search related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package services

import (
	"context"

	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
)

type SearchService struct {
	searchServiceClient search.SearchEngineClient
}

func (factory *Factory) NewSearchService(searchServiceClient search.SearchEngineClient) *SearchService {
	return &SearchService{
		searchServiceClient: searchServiceClient,
	}
}

// SearchPost 搜索帖子
//
// 参数：
//   - queryString 搜索字符串
//
// 返回值：
//   - *search.SearchResponse 搜索结果
//   - error 错误
func (service *SearchService) SearchPost(queryString string) (*search.SearchResponse, error) {
	return service.searchServiceClient.Search(context.TODO(), &search.SearchRequest{
		Query: queryString,
	})
}
