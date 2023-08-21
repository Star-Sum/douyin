package RequestEntity

// FocusRequest 关注请求结构体
type FocusRequest struct {
	UserId     string `json:"user_id"`
	ActionType string `json:"action_type"` // 1-关注，2-取消关注
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

// FocusBack 关注请求返回结构体
type FocusBack struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// FanListBack 返回粉丝列表
type FansListBack struct {
	StatusCode string   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`  // 返回状态描述
	UserList   [][]User `json:"user_list"`   // 用户列表
}

// FriendListBack 返回好友列表
type FriendsListBack struct {
	StatusCode string   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`  // 返回状态描述
	UserList   [][]User `json:"user_list"`   // 用户列表
}

// FocusListBack 返回关注列表
type FocusListBack struct {
	StatusCode string   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`  // 返回状态描述
	UserList   [][]User `json:"user_list"`   // 用户列表
}
