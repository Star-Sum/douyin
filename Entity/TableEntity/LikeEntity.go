package TableEntity

import "time"

type LikeInfo struct {
	UserID   int64     // 用户ID
	VideoID  int64     // 视频ID
	LikeID   int64     `gorm:"primaryKey"` // 点赞行为ID
	LikeTime time.Time // 点赞时间
}
