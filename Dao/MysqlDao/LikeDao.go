package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

// LikeDaos 建立单例模式
type LikeDaos struct {
}

var (
	likeDaosInstance *LikeDaos
	// 同步原语
	likeDaosOnce sync.Once
)

func NewLikeDao() *LikeDaos {
	// 结构体不会多次初始化，可以节省资源
	likeDaosOnce.Do(func() {
		likeDaosInstance = &LikeDaos{}
	})
	return likeDaosInstance
}

type LikeDao interface {
	// QueryCountOfFavorite 从点赞记录里面查询某一条视频的点赞数
	QueryCountOfFavorite(conditions map[string]interface{}) (int64, error)

	//IsFavorite 查询userID的用户是否对videoID的视频点赞
	//如果点赞返回true 否则返回false
	//如果中间发生错误，返回false
	IsFavorite(userID int64, videoID int64) bool

	// Add 增加 一条点赞记录
	Add(userID, videoID int64) error

	// Sub 删除一条点赞记录
	Sub(userID, videoID int64) error

	// FavoriteListByUserID 获取某用户点赞的所有视频的 ID 列表
	// 对于 favorite 实体来说，只有对应的 userID 或者 videoID 才是比较重要的，所以这里就只返回一个int64
	// 类型的切片，而不是 favorite 类型的切片
	FavoriteListByUserID(userID int64) ([]int64, error)
}

// QueryCountOfFavorite 从点赞记录里面查询某一条视频的点赞数
func (likeDaosInstance *LikeDaos) QueryCountOfFavorite(conditions map[string]interface{}) (int64, error) {
	var count int64
	err := mysqldb.Model(&TableEntity.LikeInfo{}).
		Select("id").
		Where(conditions).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IsFavorite 查询userID的用户是否对videoID的视频点赞
// 如果点赞返回true 否则返回false
// 如果中间发生错误，返回false
func (likeDaosInstance *LikeDaos) IsFavorite(userID int64, videoID int64) bool {
	var count int64
	err := mysqldb.Model(&TableEntity.LikeInfo{}).
		Select("UserID").
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count != 0
}

// Add 增加 一条点赞记录
func (likeDaosInstance *LikeDaos) Add(userID, videoID int64) error {
	f := TableEntity.LikeInfo{
		UserID:   userID,
		VideoID:  videoID,
		LikeTime: time.Now(),
	}
	err := mysqldb.Create(&f).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("Add Like Info Error!", err)
		return err
	}
	return nil
}

// Sub 删除一条点赞记录
func (likeDaosInstance *LikeDaos) Sub(userID, videoID int64) error {
	f := TableEntity.LikeInfo{
		UserID:  userID,
		VideoID: videoID,
	}
	return mysqldb.Model(&TableEntity.LikeInfo{}).
		Where("user_id =? AND video_id =?", userID, videoID).Delete(&f).Error
}

// FavoriteListByUserID 获取某用户点赞的所有视频的 ID 列表
// 对于 favorite 实体来说，只有对应的 userID 或者 videoID 才是比较重要的，所以这里就只返回一个int64
// 类型的切片，而不是 favorite 类型的切片
func (likeDaosInstance *LikeDaos) FavoriteListByUserID(userID int64) ([]int64, error) {
	var f []TableEntity.LikeInfo
	err := mysqldb.Model(&TableEntity.LikeInfo{}).
		Select("video_id").
		Where("user_id").
		Find(&f).Error

	// 除了判断err不为空，也要处理 RecordNotFound 问题
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var res []int64
	for _, i := range f {
		res = append(res, i.VideoID)
	}
	return res, nil
}
