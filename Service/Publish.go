package Service

import (
	"douyin/Entity/RequestEntity"
	"mime/multipart"
)

func PublishProcess(request RequestEntity.PublishRequest,
	fileHeader *multipart.FileHeader) RequestEntity.PublishBack {
	var (
		publishBack RequestEntity.PublishBack
	)
	return publishBack
}

func PublishListProcess(request RequestEntity.UserInfoRequest) RequestEntity.PublishListBack {
	var (
		publishListBack RequestEntity.PublishListBack
	)
	return publishListBack
}
