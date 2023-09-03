package RedisDao

import (
	"context"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"encoding/json"
	"errors"
	"strconv"
)

type MessageDao interface {
	// GetMessageList 从redis缓存中获取user的聊天记录
	GetMessageList(fromUserID int64, toUserID int64) ([]TableEntity.Message, error)
	SetMessageList(fromUserID int64, toUserID int64, messages []TableEntity.Message) error
}
type MessageDaoImpl struct {
	Ctx context.Context
}

func (u *MessageDaoImpl) GetMessageList(fromUserID int64, toUserID int64) ([]TableEntity.Message, error) {
	key := "messages:user_" + strconv.FormatInt(toUserID, 10) + "_from_" + strconv.FormatInt(fromUserID, 10)
	messagesMap, err := redisHandler.HGetAll(u.Ctx, key).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("GetMessagesFromRedis Error!", err)
		return nil, err
	}
	// Redis中没找到该message
	if len(messagesMap) == 0 {
		Log.NormalLog("GetMessagesFromRedis: No messages found in Redis.", nil)
		return nil, errors.New("no messages found in Redis")
	}

	var messages []TableEntity.Message

	// 转换为 TableEntity.Message 切片
	for _, messageData := range messagesMap {
		var messagesSlice []TableEntity.Message
		if err := json.Unmarshal([]byte(messageData), &messagesSlice); err != nil {
			Log.ErrorLogWithoutPanic("Message Unmarshal Error for messageData: "+messageData, err)
		} else {
			messages = append(messages, messagesSlice...)
		}
	}

	return messages, nil
}

func (u *MessageDaoImpl) SetMessageList(fromUserID int64, toUserID int64, messages []TableEntity.Message) error {
	key := "messages:user_" + strconv.FormatInt(toUserID, 10) + "_from_" + strconv.FormatInt(fromUserID, 10)

	// 将 messages 转换为 JSON 字符串
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		Log.ErrorLogWithoutPanic("Message Marshal Error!", err)
		return err
	}

	// 使用 HSet 设置键值对
	_, err = redisHandler.HSet(u.Ctx, key, "messages", messagesJSON).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("SetMessagesInRedis Error!", err)
		return err
	}

	return nil
}
