package RedisDao

import (
	"context"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"encoding/json"
	"errors"
	"strconv"
)

type CommentDao interface {
	// 从redis缓存中获取视频的评论记录
	GetCommentList(videoId int64) ([]TableEntity.Comment, error)
	SetCommentList(videoId int64, comments []TableEntity.Comment) error
}

type CommentDaoImpl struct {
	Ctx context.Context
}

func (c *CommentDaoImpl) GetCommentList(videoId int64) ([]TableEntity.Comment, error) {
	key := "Comment:videoId_" + strconv.FormatInt(videoId, 10)
	commentMap, err := redisHandler.HGetAll(c.Ctx, key).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("GetCommentFromRedis Error!", err)
		return nil, err
	}
	// Redis中没找到该 comment
	if len(commentMap) == 0 {
		Log.NormalLog("GetCommentFromRedis: No comment found in Redis.", nil)
		return nil, errors.New("no Comment found in Redis")
	}
	var comments []TableEntity.Comment

	//转换为 TableEntity.Comment 切片
	for _, commentData := range commentMap {
		var commentSlice []TableEntity.Comment
		if err := json.Unmarshal([]byte(commentData), &commentSlice); err != nil {
			Log.ErrorLogWithoutPanic("Comment Unmarshal Error for commentDataData: "+commentData, err)
		} else {
			comments = append(comments, commentSlice...)
		}
	}
	return comments, nil
}

func (c *CommentDaoImpl) SetCommentList(videoId int64, comments []TableEntity.Comment) error {
	key := "comments:videoId_" + strconv.FormatInt(videoId, 10)

	// 将 comments 转换为 JSON 字符串
	commentsJSON, err := json.Marshal(comments)
	if err != nil {
		Log.ErrorLogWithoutPanic("Comments Marshal Error!", err)
		return err
	}

	// 使用 HSet 设置键值对
	_, err = redisHandler.HSet(c.Ctx, key, "comments", commentsJSON).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("SetCommentsInRedis Error!", err)
		return err
	}

	return nil
}
