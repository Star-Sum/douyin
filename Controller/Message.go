package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func MessageSendImp(request RequestEntity.MessageSendRequest) RequestEntity.MessageSendBack {
	var (
		messageBack RequestEntity.MessageSendBack
	)
	messageBack = Service.MessageSendProcess(request)
	return messageBack
}

func MessageImp(request RequestEntity.MessageRequest) RequestEntity.MessageBack {
	var (
		messageBack RequestEntity.MessageBack
	)
	messageBack = Service.MessageProcess(request)
	return messageBack
}
