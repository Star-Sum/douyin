package TableEntity

import (
	"time"
)

type PublishInfo struct {
	UrlRoot string    // URL路径
	UseID   int64     // 视频博主ID
	Time    time.Time // 视频发布时间
	Title   string    // 视频标题
	VedioID int64     // 视频ID
	OpID    int64     `gorm:"primaryKey;AUTO_INCREMENT"` // 发布操作ID
}
