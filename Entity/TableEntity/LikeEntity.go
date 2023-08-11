package TableEntity

type LikeInfo struct {
	UserID   string // 用户ID
	VideoID  string // 视频ID
	LikeID   string `gorm:"primaryKey"` // 点赞行为ID
	LikeTime string // 点赞时间
}
