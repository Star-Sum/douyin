package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func LikeImp(request RequestEntity.LikeRequest) RequestEntity.LikeBack {
	var (
		likeBack RequestEntity.LikeBack
	)
	likeBack = Service.LikeProcess(request)
	return likeBack
}

func LikeListImp(request RequestEntity.LikeListRequest) RequestEntity.VideRequest {
	var (
		likeBack RequestEntity.VideRequest
	)
	likeBack = Service.LikeListProcess(request)
	return likeBack
}
