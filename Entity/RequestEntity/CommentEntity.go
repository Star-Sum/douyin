package RequestEntity

// CommentRequest 请求评论结构体
type CommentRequest struct {
	ActionType  string  `json:"action_type"`            // 1-发布评论，2-删除评论
	CommentID   *string `json:"comment_id,omitempty"`   // 要删除的评论id，在action_type=2的时候使用
	CommentText *string `json:"comment_text,omitempty"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string  `json:"token"`                  // 用户鉴权token
	VideoID     string  `json:"video_id"`               // 视频id
}

// CommentBack 返回评论信息
type CommentBack struct {
	Comment    *Comment `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
	StatusCode int64    `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string  `json:"status_msg"`  // 返回状态描述
}

// Comment Comment结构体
type Comment struct {
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64  `json:"id"`          // 评论id
	User       User   `json:"user"`        // 评论用户信息
}

type CommentListRequest struct {
	Token   string `json:"token"`    // 用户鉴权token
	VideoID string `json:"video_id"` // 视频id
}

type CommentListBack struct {
	CommentList []Comment `json:"comment_list"` // 评论列表
	StatusCode  int64     `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   *string   `json:"status_msg"`   // 返回状态描述
}
