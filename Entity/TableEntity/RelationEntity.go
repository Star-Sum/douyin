package TableEntity

import (
	"time"
)

type Follow struct {
	ID       int64     `gorm:"primaryKey;AUTO_INCREMENT"` // 关注操作ID
	UserID   string    // 关注者ID
	ToUserID string    // 被关注者ID
	Time     time.Time // 关注时间
}
