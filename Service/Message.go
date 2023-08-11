package Service

import "douyin/Entity/RequestEntity"

func MessageSendProcess(request RequestEntity.MessageSendRequest) RequestEntity.MessageSendBack {
	var (
		messageSendBack RequestEntity.MessageSendBack
	)
	return messageSendBack
}

func MessageProcess(request RequestEntity.MessageRequest) RequestEntity.MessageBack {
	var (
		messageBack RequestEntity.MessageBack
	)
	return messageBack
}
