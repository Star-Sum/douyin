package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func CommentImp(request RequestEntity.CommentRequest) RequestEntity.CommentBack {
	var commentBack RequestEntity.CommentBack
	commentBack = Service.CommentProcess(request)
	return commentBack
}

func CommentListImp(request RequestEntity.CommentListRequest) RequestEntity.CommentListBack {
	var commentListBack RequestEntity.CommentListBack
	commentListBack = Service.CommentListProcess(request)
	return commentListBack
}
