package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"douyin/Util"
	"strconv"
	"time"
)

var MCommentDaoImpl MysqlDao.CommentDao = &MysqlDao.CommentDaoImpl{}

var RCommentDaoImpl RedisDao.CommentDao = &RedisDao.CommentDaoImpl{
	Ctx: context.Background(),
}

func CommentProcess(request RequestEntity.CommentRequest) RequestEntity.CommentBack {
	// var (
	// 	commentBack RequestEntity.CommentBack
	// )
	Log.NormalLog("CommentProcess", nil)

	//解析token
	userID, err := Util.ParserToken(request.Token)
	if err != nil {
		statusMsg := "Token Parse Error!"
		Log.ErrorLogWithoutPanic("Token Parse Error!", err)
		return RequestEntity.CommentBack{
			Comment:    nil,
			StatusCode: 1,
			StatusMsg:  &statusMsg,
		}
	}
	// var comment RequestEntity.Comment
	content := *request.CommentText
	createData := time.Now().Format("2006-01-02 15:04:05")
	videoid, _ := strconv.ParseInt(request.VideoID, 10, 64)
	if request.ActionType == "1" {
		// 添加评论,在这里创建commnet的表,同时还需要将TableEntity里的Comment转化成RequestEntity里的Comment

		Tcomment, err := MCommentDaoImpl.InsertComment(content, createData, videoid, userID)
		if err != nil {
			statusMsg := "Comment insert failed"
			Log.ErrorLogWithPanic("Comment insert failed", err)
			return RequestEntity.CommentBack{
				Comment:    nil,
				StatusCode: 1,
				StatusMsg:  &statusMsg,
			}
		}
		statusMsg := "comment success"
		return RequestEntity.CommentBack{
			Comment:    &comment,
			StatusCode: 0,
			StatusMsg:  &statusMsg,
		}
	} else if request.ActionType == "2" {
		//删除评论
		commentid, _ := strconv.ParseInt(*request.CommentID, 10, 64)

		err := MCommentDaoImpl.DeleteComment(commentid)
		if err != nil {
			statusMsg := "Comment delete failed"
			Log.ErrorLogWithPanic("Comment delete failed", err)
			return RequestEntity.CommentBack{
				StatusCode: 1,
				StatusMsg:  &statusMsg,
			}
		}
	}
	statusMsg := "Comment ActionType failed"
	Log.ErrorLogWithPanic("Comment ActionType failed", err)
	return RequestEntity.CommentBack{
		StatusCode: 1,
		StatusMsg:  &statusMsg,
	}

}

func CommentListProcess(request RequestEntity.CommentListRequest) RequestEntity.CommentListBack {
	//获取评论列表
	Log.NormalLog("CommentListProcess", nil)
	//解析token
	_, err := Util.ParserToken(request.Token)
	videoid, _ := strconv.ParseInt(request.VideoID, 10, 64)
	// 先尝试从redis中获取
	commentList, err := RCommentDaoImpl.GetCommentList(videoid)
	if err != nil {
		commentList, err := MCommentDaoImpl.GetCommentList(videoid)
		if err != nil {
			statusMsg := "GetCommentListFromMysql Error!"
			Log.ErrorLogWithoutPanic("GetCommentListFromMysql Error!", err)
			return RequestEntity.CommentListBack{
				StatusCode: 1,
				StatusMsg:  &statusMsg,
			}
		}
		statusMsg := "get CommentList success"
		return RequestEntity.CommentListBack{
			CommentList: commentList,
			StatusCode:  0,
			StatusMsg:   &statusMsg,
		}
	}
	statusMsg := "get CommentList success"
	return RequestEntity.CommentListBack{
		CommentList: commentList,
		StatusCode:  0,
		StatusMsg:   &statusMsg,
	}

}
