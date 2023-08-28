package MysqlDao

import (
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Util"
	"errors"
	"log"
)

// 实现相应的增删查改需求
// 需要对gorm进行初始化吗？
// Comment
// 直接调用定义中的数据结构体
type CommentDao interface {

	// 发表评论
	InsertComment(content string, createDate string, vedioid int64, userid int64) (RequestEntity.Comment, error)
	// 删除评论，传入评论id
	DeleteComment(CommentId int64) error
	// 获取视频下的评论列表
	GetCommentList(videoId int64) ([]RequestEntity.Comment, error)
	// Comment类型转换
	// Transfer(TcommentList []TableEntity.Comment) ([]RequestEntity.Comment, error)
}

type CommentDaoImpl struct {
}

// InsertComment
// 2、发表评论
func (c *CommentDaoImpl) InsertComment(content string, createdate string, vedioid int64, userid int64) (RequestEntity.Comment, error) {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//数据库中插入一条评论信息

	//根据雪花算法生成评论id
	id, _ := Util.MakeUid(Util.Snowflake.DataCenterId, Util.Snowflake.MachineId)
	//根据id获取user
	var user = &RequestEntity.User{}
	result := mysqldb.Model(&TableEntity.UserInfo{}).Where("id=?", id).First(&user)
	if result.RowsAffected == 0 { //查询到y
		log.Println("CommentDao-insertComment: return user is not existed") //函数返回提示错误信息
	}
	var Tcomment = TableEntity.Comment{
		Content:    content,
		CreateDate: createdate,
		ID:         id,
		VedioId:    vedioid,
		UserID:     userid,
	}
	//err := MysqlDb.Model(model.CommentModel{}).Create(&comment).Error
	var Rcomment = RequestEntity.Comment{
		Content:    content,
		CreateDate: createdate,
		ID:         id,
		User:       *user,
	}
	err := mysqldb.Create(Tcomment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return Rcomment, err
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return Rcomment, nil
}

// DeleteComment
// 3、删除评论，传入评论id
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

	result := mysqldb.Model(TableEntity.Comment{}).Where("video_id=?", videoId).Order("create_date desc").Find(&TcommentList)
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
