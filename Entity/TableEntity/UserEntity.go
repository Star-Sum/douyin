package TableEntity

type UserInfo struct {
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	FavoriteCount   int64  // 喜欢数
	FollowCount     int64  // 关注总数
	FollowerCount   int64  // 粉丝总数
	ID              int64  // 用户id
	IsFollow        bool   // true-已关注，false-未关注
	Name            string // 用户名称
	Signature       string // 个人简介
	TotalFavorited  string // 获赞数量
	WorkCount       int64  // 作品数
}

type UserAccountInfo struct {
	UserID   int64  // 对应的用户编号
	Password string // 密码，最长32个字符
	Username string `gorm:"column:username"` // 注册用户名，最长32个字符
}
