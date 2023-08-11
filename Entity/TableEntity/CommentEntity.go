package TableEntity

type Comment struct {
	Content    string // 评论内容
	CreateDate string // 评论发布时间，格式 mm-dd
	ID         int64  `gorm:"primaryKey"` // 评论id
	VedioId    int64  // 视频id
	UserID     int64  // 评论者id
}
