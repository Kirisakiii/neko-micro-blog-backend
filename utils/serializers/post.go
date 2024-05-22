package serializers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

type PostListResponse struct {
	IDs []int64 `json:"ids"`
}

func NewPostListResponse(posts []int64) *PostListResponse {
	return &PostListResponse{IDs: posts}
}

// PostDetailResponse 文章信息响应结构
type PostDetailResponse struct {
	PostID       uint64   `json:"post_id"`        // 文章ID
	UID          uint64   `json:"uid"`            // 用户ID
	Timestamp    int64    `json:"timestamp"`      // 时间戳
	Title        string   `json:"title"`          // 标题
	Content      string   `json:"content"`        // 内容
	ParentPostID *uint64  `json:"parent_post_id"` // 转发自文章ID
	TopicID      *string  `json:"topic_id"`       // 所属话题ID
	Images       []string `json:"images"`         // 图片
	Like         int64    `json:"like"`           // 点赞数
	Favourite    int64    `json:"favourite"`      // 收藏数
}

// NewPostDetailResponse 创建新的文章信息响应
//
// 参数：
//   - model：文章信息模型
//
// 返回值：
//   - *PostProfileData：新的文章信息响应结构
func NewPostDetailResponse(post models.PostInfo, likeCount, favouriteCount int64) *PostDetailResponse {
	// 创建一个新的 PostProfileData 实例
	profileData := &PostDetailResponse{
		PostID:       uint64(post.ID),
		TopicID:      post.TopicID,
		UID:          post.UID,
		Timestamp:    post.CreatedAt.Unix(),
		Title:        post.Title,
		Content:      post.Content,
		ParentPostID: post.ParentPostID,
		Like:         likeCount,
		Favourite:    favouriteCount,
	}
	for _, image := range post.Images {
		profileData.Images = append(profileData.Images, "/resources/image/"+image)
	}

	return profileData
}

// CreatePostResponse 用于将 PostInfo 转换为 JSON 格式的结构体
type CreatePostResponse struct {
	ID uint64 `json:"id"`
}

// NewPostResponse 用于创建 PostResponse 实例
func NewCreatePostResponse(postInfo models.PostInfo) CreatePostResponse {
	var resp = CreatePostResponse{
		ID: uint64(postInfo.ID),
	}
	return resp
}

// UploadPostImageResponse 上传博文图片响应结构
type UploadPostImageResponse struct {
	UUID string `json:"uuid"` // 图片UUID
}

// NewUploadPostImageResponse 创建新的上传博文图片响应
func NewUploadPostImageResponse(uuid string) UploadPostImageResponse {
	return UploadPostImageResponse{UUID: uuid}
}

// PostUserStatus 用户文章状态响应结构
type PostUserStatus struct {
	PostID    uint64 `json:"post_id"`   // 文章ID
	UID       uint64 `json:"uid"`       // 用户ID
	Like      bool   `json:"like"`      // 是否点赞
	Favourite bool   `json:"favourite"` // 是否收藏
}

// NewPostUserStatus 创建新的用户文章状态响应
func NewPostUserStatus(postID uint64, uid uint64, like bool, favourite bool) PostUserStatus {
	return PostUserStatus{
		PostID:    postID,
		UID:       uid,
		Like:      like,
		Favourite: favourite,
	}
}
