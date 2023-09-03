package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
)

type MessageDao interface {
	// 聊天记录
	GetMessageList(fromUserID int64, toUserID int64) ([]TableEntity.Message, error)
	// 发送消息
	PostMessage(fromUserID int64, toUserID int64, content string, createTime int64) (int64, error)
}

type MessageDaoImpl struct {
}

// GetMessageList 获取Message信息列表
func (u *MessageDaoImpl) GetMessageList(fromUserID int64, toUserID int64) ([]TableEntity.Message, error) {
	var messages []TableEntity.Message

	err := mysqldb.Model(&TableEntity.Message{}).
		Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).
		Find(&messages).Error

	if err != nil {
		Log.ErrorLogWithoutPanic("GetMessagesByToUserID Error!", err)
		return nil, err
	}

	return messages, nil
}

// PostMessage 发送Message
func (u *MessageDaoImpl) PostMessage(fromUserID int64, toUserID int64, content string, createTime int64) (int64, error) {
	var message = TableEntity.Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreateTime: createTime,
	}

	if err := mysqldb.Create(&message).Error; err != nil {
		return 1, err
	}

	return 0, nil
}
