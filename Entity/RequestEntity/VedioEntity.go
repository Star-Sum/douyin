package RequestEntity

// FeedRequest Feed流请求结构体（内部参数非必选）
type FeedRequest struct {
	LatestTime *string `json:"latest_time,omitempty"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token      *string `json:"token,omitempty"`       // 用户登录状态下设置
}

// FeedBack Feed流返回信息
type FeedBack struct {
	NextTime   *int64        `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64         `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string       `json:"status_msg"`  // 返回状态描述
	VideoList  []VideRequest `json:"video_list"`  // 视频列表
}

// VideRequest 视频具体访问信息
type VideRequest struct {
	Author        AuthorUser `json:"author"`         // 视频作者信息
	CommentCount  int64      `json:"comment_count"`  // 视频的评论总数
	CoverURL      string     `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64      `json:"favorite_count"` // 视频的点赞总数
	ID            int64      `json:"id"`             // 视频唯一标识
	IsFavorite    bool       `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string     `json:"play_url"`       // 视频播放地址
	Title         string     `json:"title"`          // 视频标题
}
