package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func VedioFeedImp(request RequestEntity.FeedRequest) RequestEntity.FeedBack {
	var vedioFeedBack RequestEntity.FeedBack
	vedioFeedBack = Service.VedioFeedProcess(request)
	return vedioFeedBack
}
