package MysqlDao

import (
	"golang.org/x/crypto/openpgp"
	"sync"
	"gin_study/entity"
)

//建立单列模式
type like_Daos sturct{}

var(
	like_Daos_Instance *like_Daos
	like_Daos_Once sync.Once
)

func newlikedao() *like_Daos{
	like_Daos_Once.Do(func() {
		like_Daos_Once = &like_Daos_Instance{}
	})
	return like_Daos_Instance
}

type dianzandao interface {
	// 从点赞记录里面查询某一条视频的点赞数 QueryCountOfFavorite
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

// 从点赞记录里面查询某一条视频的点赞数 QueryCountOfFavorite
func (like_Daos_Instance *like_Daos) QueryCountOfFavorite(conditions map[string]interface{}) (int64, error) {
     var count int64
	 err:= db.Model(&tableentity.LikeInfo{}).
		 Select("id").
		 Where(conditions).Count(&count).Error
	 if err !=nil{
		 return 0, err
	 }
	 return count,nil
}
//IsFavorite 查询userID的用户是否对videoID的视频点赞
//如果点赞返回true 否则返回false
//如果中间发生错误，返回false
func (like_Daos_Instance *like_Daos) IsFavorite (userID int64, videoID int64) bool{
	var count int64
	err := db.Model(&tableentity.LikeInfo{}).
		Select("UserID").
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count != 0
}

// Add 增加 一条点赞记录
func (like_Daos_Instance *like_Daos) Add(userID, videoID int64) error{
	f := tableentity.LikeInfo{
		UserID:  userID,
		VideoID: videoID,
	}
	return db.Model(&tableentity.LikeInfo{}).Create(&f).Error
}

// Sub 删除一条点赞记录
func (like_Daos_Instance *like_Daos) Sub(userID, videoID int64) error{
	f := tableentity.LikeInfo{
		UserID:  userID,
		VideoID: videoID,
	}
	return db.Model(&tableentity.LikeInfo{}).Where("user_id =? AND video_id =?", userID, videoID).Delete(&f).Error
}

// FavoriteListByUserID 获取某用户点赞的所有视频的 ID 列表
// 对于 favorite 实体来说，只有对应的 userID 或者 videoID 才是比较重要的，所以这里就只返回一个int64
// 类型的切片，而不是 favorite 类型的切片
func (like_Daos_Instance *like_Daos) FavoriteListByUserID(userID int64) ([]int64, error){
	var f []tableentity.LikeInfo
	err := db.Model(&tableentity.LikeInfo{}).
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
