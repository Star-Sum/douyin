package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
	"mime/multipart"
)

func PublishImp(request RequestEntity.PublishRequest, fileHeader *multipart.FileHeader) RequestEntity.PublishBack {
	var (
		publishBack RequestEntity.PublishBack
	)
	publishBack = Service.PublishProcess(request, fileHeader)
	return publishBack
}

func PublishListImp(request RequestEntity.UserInfoRequest) RequestEntity.PublishListBack {
	var (
		publishListBack RequestEntity.PublishListBack
	)
	publishListBack = Service.PublishListProcess(request)
	return publishListBack
}
