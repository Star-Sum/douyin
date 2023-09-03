package MysqlDao

import (
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"errors"
	"log"
)

type CommentDao interface {

	// InsertComment 发表评论
	InsertComment(content string, createDate string, vedioid int64, userid int64) (RequestEntity.Comment, error)
	// DeleteComment 删除评论，传入评论id
	DeleteComment(CommentId int64) error
	// GetCommentList 获取视频下的评论列表
	GetCommentList(videoId int64) ([]RequestEntity.Comment, error)
}

type CommentDaoImpl struct {
}

// InsertComment content:评论具体内容 createDate:创建时间 vedioId:视频ID userId:发表评论的用户ID
func (c *CommentDaoImpl) InsertComment(content string, createDate string, vedioId int64,
	userId int64) (RequestEntity.Comment, error) {
	// 标识方法已运行
	Log.NormalLog("CommentDao-InsertComment: running", nil)

	//根据id获取user

	var user TableEntity.UserInfo
	err := mysqldb.Model(&TableEntity.UserInfo{}).Where("id=?", userId).Find(&user).Error
	if err != nil {
		//函数返回提示错误信息
		Log.ErrorLogWithoutPanic("CommentDao-insertComment: return user is not existed", nil)
	}

	// 在表中插入评论信息
	var tComment = TableEntity.Comment{
		Content:    content,
		CreateDate: createDate,
		VedioId:    vedioId,
		UserID:     userId,
	}
	err = mysqldb.Create(&tComment).Error

	// 更新评论数
	var commentCount int64
	err = mysqldb.Model(TableEntity.VedioInfo{}).
		Select("comment_count").Where("id =?", vedioId).Find(&commentCount).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("Get Comment Count Error!", err)
	}
	err = mysqldb.Model(TableEntity.VedioInfo{}).
		Where("id =?", vedioId).Update("comment_count", commentCount+1).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("Set Comment Count Error!", err)
	}

	// 制作返回信息
	var rUser = RequestEntity.User{
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		ID:              user.ID,
		IsFollow:        true,
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
	var rComment = RequestEntity.Comment{
		Content:    content,
		CreateDate: createDate,
		User:       rUser,
	}

	if err != nil {
		//函数返回提示错误信息
		Log.ErrorLogWithoutPanic("CommentDao-InsertComment: return create comment failed", err)
		return rComment, err
	}
	Log.NormalLog("CommentDao-InsertComment: return success", nil) //函数执行成功，返回正确信息
	return rComment, nil
}

// DeleteComment 删除评论，传入评论id
func (c *CommentDaoImpl) DeleteComment(ID int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	var commentInfo TableEntity.Comment
	//先查询是否有此评论
	result := mysqldb.Model(TableEntity.Comment{}).Where("comment_id = ?", ID).First(&commentInfo)
	if result.RowsAffected == 0 { //查询到此评论数量为0则返回无此评论
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := mysqldb.Model(TableEntity.Comment{}).Where("comment_id = ?", ID).Delete(&TableEntity.Comment{}).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

// func (c *CommentDaoImpl) Transfer(TcommentList []TableEntity.Comment) ([]RequestEntity.Comment, error) {
// 	// for循环更改内部Comment类型
// 	RcommentList := []RequestEntity.Comment{}
// 	for _, v := range TcommentList {
// 		// 根据v查询user
// 		var user = &RequestEntity.User{}
// 		result := mysqldb.Model(&TableEntity.UserInfo{}).Where("id=?", v.ID).First(&user)
// 		if result.RowsAffected == 0 { //查询到y
// 			log.Println("CommentDao-insertComment: return user is not existed") //函数返回提示错误信息
// 		}
// 		Rcomment := RequestEntity.Comment{
// 			Content:    v.Content,
// 			CreateDate: v.CreateDate,
// 			ID:         v.ID,
// 			User:       *user,
// 		}
// 		RcommentList = append(RcommentList, Rcomment)
// 	}
// 	return RcommentList, nil
// }

// 4、获取评论列表
func (c *CommentDaoImpl) GetCommentList(videoId int64) ([]RequestEntity.Comment, error) {
	log.Println("CommentDao-GetCommentList: running") //函数已运行
	//数据库中查询评论信息list
	var TcommentList []TableEntity.Comment
	// var RcommentList []RequestEntity.Comment

	result := mysqldb.Model(TableEntity.Comment{}).
		Where("vedio_id=?", videoId).
		Order("create_date desc").Find(&TcommentList)
	// RcommentList, err := Transfer(TcommentList)
	// 做一个TcommentList -> RcommentList的转换
	RcommentList := []RequestEntity.Comment{}
	for _, v := range TcommentList {
		// 根据v查询user
		var user = &RequestEntity.User{}
		result := mysqldb.Model(&TableEntity.UserInfo{}).Where("id=?", v.ID).First(&user)
		if result.RowsAffected == 0 { //查询到y
			log.Println("CommentDao-insertComment: return user is not existed") //函数返回提示错误信息
		}
		Rcomment := RequestEntity.Comment{
			Content:    v.Content,
			CreateDate: v.CreateDate,
			ID:         v.ID,
			User:       *user,
		}
		RcommentList = append(RcommentList, Rcomment)
	}
	//若此视频没有评论信息
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") //函数返回提示无评论
		return RcommentList, errors.New("there are no comments")
	}
	//若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") //函数返回提示获取评论错误
		return RcommentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") //函数执行成功，返回正确信息
	return RcommentList, nil
}
