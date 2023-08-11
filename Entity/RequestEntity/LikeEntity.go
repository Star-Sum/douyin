package RequestEntity

// LikeRequest 请求点赞实体
type LikeRequest struct {
	ActionType string `json:"action_type"` // 1-点赞，2-取消点赞
	Token      string `json:"token"`       // 用户鉴权token
	VideoID    string `json:"video_id"`    // 视频id
}

// LikeBack 请求点赞结果返回实体
type LikeBack struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// LikeListRequest 请求喜欢列表结构体
type LikeListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}
