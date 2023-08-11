package RequestEntity

import (
	"mime/multipart"
)

// PublishRequest 发布视频实体
type PublishRequest struct {
	Data  multipart.File `json:"data"`
	Token string         `json:"token"`
	Title string         `json:"title"`
}

// PublishBack 返回状态信息
type PublishBack struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  *string `json:"status_msg"`
}

type PublishListBack struct {
	StatusCode int64         `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string       `json:"status_msg"`  // 返回状态描述
	VideoList  []VideRequest `json:"video_list"`  // 用户发布的视频列表
}
