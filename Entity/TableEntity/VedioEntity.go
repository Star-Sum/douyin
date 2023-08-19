package TableEntity

import (
	"time"
)

type VedioInfo struct {
	AuthorID      int64     // 视频作者ID
	CommentCount  int64     // 视频的评论总数
	CoverURL      string    // 视频封面地址
	FavoriteCount int64     // 视频的点赞总数
	ID            int64     `gorm:"primaryKey"` // 视频唯一标识
	PlayURL       string    // 视频播放地址
	Title         string    // 视频标题
	Timestamp     time.Time // 视频发布时间
}
