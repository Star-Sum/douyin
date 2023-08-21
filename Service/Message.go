package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"sort"
	"strconv"
	"time"
)

var (
	MMessageDaoImpl MysqlDao.MessageDao = &MysqlDao.MessageDaoImpl{}
	RMessageDaoImpl RedisDao.MessageDao = &RedisDao.MessageDaoImpl{
		Ctx: context.Background(),
	}
)

func MessageSendProcess(request RequestEntity.MessageSendRequest) RequestEntity.MessageSendBack {
	var (
		messageSendBack RequestEntity.MessageSendBack
	)

	// 判断token是否合法，并提取fromUserID
	fromUserID, err := Util.ParserToken(request.Token)
	if err != nil {
		messageSendBack.StatusCode = 1
		messageSendBack.StatusMsg = "Token Parse Error!"
		Log.ErrorLogWithoutPanic("Token Parse Error!", err)
		return messageSendBack
	}

	// 传入的ActionType参数是string，先解析为int64
	if request.ActionType != "1" {
		messageSendBack.StatusCode = 1
		messageSendBack.StatusMsg = "Invalid action_type Param!"
		Log.ErrorLogWithoutPanic("Invalid action_type Param!", nil)
		return messageSendBack
	}

	// 传入的ToUserID参数是string，先解析为int64
	toUserIDString := request.ToUserID
	toUserID, err := strconv.ParseInt(toUserIDString, 10, 64)
	if err != nil {
		messageSendBack.StatusCode = 1
		messageSendBack.StatusMsg = "to_user_id Parse Error!"
		Log.ErrorLogWithoutPanic("to_user_id Parse Error!", err)
		return messageSendBack
	}

	// 判断数据库中是否存在toUserID的用户
	if _, err := MUserDaoImpl.GetUserById(toUserID); err != nil {
		messageSendBack.StatusCode = 1
		messageSendBack.StatusMsg = "User with the specified ID not found!"
		Log.ErrorLogWithoutPanic("User with the specified ID not found!", err)
		return messageSendBack
	}

	// 判断fromUserID和toUserID是否是同一个
	if fromUserID == toUserID {
		messageSendBack.StatusCode = 1
		messageSendBack.StatusMsg = "Same UserID!"
		Log.ErrorLogWithoutPanic("Same UserID!", err)
		return messageSendBack
	}

	// 获取消息内容
	content := request.Content

	// 获取当前时间戳
	createTime := time.Now().Unix()
	postStatusCode, err := MMessageDaoImpl.PostMessage(fromUserID, toUserID, content, createTime)

	// 数据库插入结果正确返回
	if postStatusCode == 0 && err == nil {
		messageSendBack.StatusCode = 0
		messageSendBack.StatusMsg = "success"
		return messageSendBack
	}

	// 数据库返回错误
	messageSendBack.StatusCode = 1
	messageSendBack.StatusMsg = "Database Insert Error!"
	return messageSendBack
}

func MessageProcess(request RequestEntity.MessageRequest) RequestEntity.MessageBack {
	var (
		messageBack RequestEntity.MessageBack
	)

	// 判断token是否合法，并提取fromUserID
	fromUserID, err := Util.ParserToken(request.Token)
	if err != nil {
		FailedMsg := "Token Parse Error"
		messageBack.StatusCode = "1"
		messageBack.MessageList = []RequestEntity.Message{}
		messageBack.StatusMsg = &FailedMsg
		Log.ErrorLogWithoutPanic("Token Parse Error!", err)
		return messageBack
	}

	// 传入的to_user_id参数是string，先解析为int64
	toUserIDString := request.ToUserID
	toUserID, err := strconv.ParseInt(toUserIDString, 10, 64)
	if err != nil {
		FailedMsg := "to_user_id Parse Error!"
		messageBack.StatusCode = "1"
		messageBack.MessageList = []RequestEntity.Message{}
		messageBack.StatusMsg = &FailedMsg
		Log.ErrorLogWithoutPanic("to_user_id Parse Error!", err)
		return messageBack
	}

	// 判断数据库中是否存在toUserID的用户
	if _, err := MUserDaoImpl.GetUserById(toUserID); err != nil {
		FailedMsg := "User with the specified ID not found!"
		messageBack.StatusCode = "1"
		messageBack.MessageList = []RequestEntity.Message{}
		messageBack.StatusMsg = &FailedMsg
		Log.ErrorLogWithoutPanic("User with the specified ID not found!", err)
		return messageBack
	}

	// 判断fromUserID和toUserID是否是同一个
	if fromUserID == toUserID {
		FailedMsg := "Same UserID!"
		messageBack.StatusCode = "1"
		messageBack.MessageList = []RequestEntity.Message{}
		messageBack.StatusMsg = &FailedMsg
		Log.ErrorLogWithoutPanic("Same UserID!", err)
		return messageBack
	}

	// 获取消息列表、合并并按时间戳排序
	if messageBack.MessageList, err = GetMessageList(fromUserID, toUserID); err != nil {
		FailedMsg := "Database error!"
		messageBack.StatusCode = "1"
		messageBack.MessageList = []RequestEntity.Message{}
		messageBack.StatusMsg = &FailedMsg
		Log.ErrorLogWithoutPanic("Database error!", err)
		return messageBack
	}

	// 正确返回
	SuccessMsg := "success"
	messageBack.StatusCode = "0"
	messageBack.StatusMsg = &SuccessMsg
	return messageBack
}

// GetMessageList 获取fromUserID和toUserID互发的消息
func GetMessageList(fromUserID int64, toUserID int64) ([]RequestEntity.Message, error) {
	// 先获取自己发给对方的消息列表
	fromMessageList, err := RMessageDaoImpl.GetMessageList(fromUserID, toUserID)
	if err != nil {
		// 从 Redis 获取失败，从 MySQL 获取消息列表
		fromMessageList, err = MMessageDaoImpl.GetMessageList(fromUserID, toUserID)
		if err != nil {
			Log.ErrorLogWithoutPanic("GetMessageListFromMysql Error!", err)
			return []RequestEntity.Message{}, err
		}

		// 将获取的消息列表存储到 Redis
		err = RMessageDaoImpl.SetMessageList(fromUserID, toUserID, fromMessageList)
		if err != nil {
			Log.ErrorLogWithoutPanic("SetMessageList Error!", err)
		}
	}

	// 再获取对方发给自己的消息列表
	toMessageList, err := RMessageDaoImpl.GetMessageList(toUserID, fromUserID)
	if err != nil {
		// 从 Redis 获取失败，从 MySQL 获取消息列表
		toMessageList, err = MMessageDaoImpl.GetMessageList(toUserID, fromUserID)
		if err != nil {
			Log.ErrorLogWithoutPanic("GetMessageListFromMysql Error!", err)
			return []RequestEntity.Message{}, err
		}

		// 将获取的消息列表存储到 Redis
		err = RMessageDaoImpl.SetMessageList(toUserID, fromUserID, toMessageList)
		if err != nil {
			Log.ErrorLogWithoutPanic("SetMessageList Error!", err)
		}
	}
	mergedMessageList := mergeMessageLists(fromMessageList, toMessageList)
	return mergedMessageList, nil
}

// 合并消息列表并按时间戳排序
func mergeMessageLists(fromMessageList []TableEntity.Message, toMessageList []TableEntity.Message) []RequestEntity.Message {
	var requestMessages []RequestEntity.Message

	// TableEntity转为RequestEntity并合并
	for _, tableMessage := range fromMessageList {
		requestMessage := RequestEntity.Message{
			Content:    tableMessage.Content,
			CreateTime: tableMessage.CreateTime,
			FromUserID: tableMessage.FromUserID,
			ID:         tableMessage.ID,
			ToUserID:   tableMessage.ToUserID,
		}
		requestMessages = append(requestMessages, requestMessage)
	}
	for _, tableMessage := range toMessageList {
		requestMessage := RequestEntity.Message{
			Content:    tableMessage.Content,
			CreateTime: tableMessage.CreateTime,
			FromUserID: tableMessage.FromUserID,
			ID:         tableMessage.ID,
			ToUserID:   tableMessage.ToUserID,
		}
		requestMessages = append(requestMessages, requestMessage)
	}

	// 按时间戳排序
	sort.Slice(requestMessages, func(i, j int) bool {
		return requestMessages[i].CreateTime < requestMessages[j].CreateTime
	})
	return requestMessages
}
