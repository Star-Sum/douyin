package TableEntity

type Message struct {
	Content    string // 消息内容
	CreateTime int64  // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  // 消息发送者id
	ID         int64  `gorm:"primaryKey"` // 消息id
	ToUserID   int64  // 消息接收者id
}
