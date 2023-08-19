package Service

import "douyin/Entity/RequestEntity"

func LikeProcess(request RequestEntity.LikeRequest) RequestEntity.LikeBack {
	var (
		likeBack RequestEntity.LikeBack
	)
	return likeBack
}

func LikeListProcess(request RequestEntity.LikeListRequest) RequestEntity.VedioRequest {
	var (
		likeListBack RequestEntity.VedioRequest
	)
	return likeListBack
}
