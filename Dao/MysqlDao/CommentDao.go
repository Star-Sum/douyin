package MysqlDao
<<<<<<< HEAD
=======

import (
	"douyin/Entity/RequestEntity"
	"douyin/Util"
	"errors"
	"log"
)

// 实现相应的增删查改需求
// 需要对gorm进行初始化吗？
// Comment
// 直接调用定义中的数据结构体
type CommentDao interface {

	// // Count 查询comment数量
	// Count(videoId int64) (int64, error)
	// 根据视频id获取评价id列表
	// CommentIdList(videoId int64) ([]string, error)
	// 发表评论
	InsertComment(content string, createDate string, vedioid int64, userid int64) (RequestEntity.Comment, error)
	// 删除评论，传入评论id
	DeleteComment(CommentId int64) error
	// 获取视频下的评论列表
	GetCommentList(videoId int64) ([]RequestEntity.Comment, error)
}

type CommentDaoImpl struct {
}

// type CommentListImp struct {
// }

// Count
// // 1、使用video id 查询Comment数量
// func (c *CommentDaoImpl) Count(videoId int64) (int64, error) {
// 	log.Println("CommentDao-Count: running") //函数已运行
// 	//Init()
// 	var count int64
// 	//数据库中查询评论数量
// 	err := mysqldb.Model(&TableEntity.Comment{}).Where("videoId= ?", videoId).Count(&count).Error
// 	if err != nil {
// 		log.Println("CommentDao-Count: return count failed") //函数返回提示错误信息
// 		return -1, errors.New("find comments count failed")
// 	}
// 	log.Println("CommentDao-Count: return count success") //函数执行成功，返回正确信息
// 	return count, nil
// }

// // CommentIdList 根据视频id获取评价id列表
// func (c *CommentDaoImpl) CommentIdList(videoId int64) ([]string, error) {
// 	var commentIdList []string
// 	err := mysqldb.Model(TableEntity.Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
// 	if err != nil {
// 		log.Println("CommentIdList:", err)
// 		return nil, err
// 	}
// 	return commentIdList, nil
// }

// InsertComment
// 2、发表评论
func (c *CommentDaoImpl) InsertComment(content string, createdate string, vedioid int64, userid int64) (RequestEntity.Comment, error) {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//数据库中插入一条评论信息

	//根据雪花算法生成评论id
	id, _ := Util.MakeUid(Util.Snowflake.DataCenterId, Util.Snowflake.MachineId)
	var user RequestEntity.User
	result := mysqldb.Model(RequestEntity.User{}).Where("Id = ?", userid).First(&user)
	if result.RowsAffected == 0 { //查询到y
		log.Println("CommentDao-insertComment: return user is not existed") //函数返回提示错误信息
	}
	var comment = RequestEntity.Comment{
		Content:    content,
		CreateDate: createdate,
		ID:         id,
		User:       user,
	}
	//err := MysqlDb.Model(model.CommentModel{}).Create(&comment).Error
	err := mysqldb.Create(comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return comment, nil
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return comment, nil
}

// DeleteComment
// 3、删除评论，传入评论id
func (c *CommentDaoImpl) DeleteComment(ID int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	var commentInfo RequestEntity.Comment
	//先查询是否有此评论
	result := mysqldb.Model(RequestEntity.Comment{}).Where("comment_id = ?", ID).First(&commentInfo)
	if result.RowsAffected == 0 { //查询到此评论数量为0则返回无此评论
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := mysqldb.Model(RequestEntity.Comment{}).Where("comment_id = ?", ID).Delete(&RequestEntity.Comment{}).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

// 4、获取评论列表
func (c *CommentDaoImpl) GetCommentList(videoId int64) ([]RequestEntity.Comment, error) {
	log.Println("CommentDao-GetCommentList: running") //函数已运行
	//数据库中查询评论信息list
	var commentList []RequestEntity.Comment
	result := mysqldb.Model(RequestEntity.Comment{}).Where("video_id=?", videoId).Order("create_date desc").Find(&commentList)
	//若此视频没有评论信息
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") //函数返回提示无评论
		return commentList, errors.New("there are no comments")
	}
	//若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") //函数返回提示获取评论错误
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") //函数执行成功，返回正确信息
	return commentList, nil
}
>>>>>>> 51633ea (我的第一次提交)
