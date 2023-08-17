package RequestEntity

import "encoding/json"

// RegisterRequest 注册请求输入结构体
type RegisterRequest struct {
	Password string `json:"password"` // 密码，最长32个字符
	Username string `json:"username"` // 注册用户名，最长32个字符
}

// RegisterBack 注册请求输出结构体
type RegisterBack struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// UserInfoRequest 获取用户信息请求结构体
type UserInfoRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

// UserInfoBack 获取用户信息返回结构体
type UserInfoBack struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	User       *User   `json:"user"`        // 用户信息
}

// User 具体用户信息
type User struct {
	Avatar          string `redis:"avatar" json:"avatar"`                     // 用户头像
	BackgroundImage string `redis:"background_image" json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `redis:"favorite_count" json:"favorite_count"`     // 喜欢数
	FollowCount     int64  `redis:"follow_count" json:"follow_count"`         // 关注总数
	FollowerCount   int64  `redis:"follower_count" json:"follower_count"`     // 粉丝总数
	ID              int64  `redis:"id" json:"id"`                             // 用户id
	IsFollow        bool   `redis:"is_follow" json:"is_follow"`               // true-已关注，false-未关注
	Name            string `redis:"name" json:"name"`                         // 用户名称
	Signature       string `redis:"signature" json:"signature"`               // 个人简介
	TotalFavorited  int64  `redis:"total_favorited" json:"total_favorited"`   // 获赞数量
	WorkCount       int64  `redis:"work_count" json:"work_count"`             // 作品数
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)

}

// LoginRequest 登录请求输入结构体
type LoginRequest struct {
	Password string `json:"password"` // 登录密码
	Username string `json:"username"` // 登录用户名
}

// LoginBack 登录请求返回结构体
type LoginBack struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	Token      *string `json:"token"`       // 用户鉴权token
	UserID     *int64  `json:"user_id"`     // 用户id
}

// AuthorUser 视频作者信息
type AuthorUser struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
}
