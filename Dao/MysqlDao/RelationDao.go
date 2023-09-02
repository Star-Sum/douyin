package MysqlDao

import (
	"database/sql"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RelationDao struct {
	db *sql.DB
}
type Follow struct {
	userID   int64
	toUserID int64
}

// Follow 进行关注和取消关注的操作
func (dao *RelationDao) Follow(userID int64, toUserId int64,
	time time.Time, ActionType string) error {

	// 制作关系信息
	var follow = &TableEntity.Follow{
		UserID:   strconv.FormatInt(userID, 10),
		ToUserID: strconv.FormatInt(toUserId, 10),
		Time:     time,
	}

	// 行为编码为1时进行关注
	if ActionType == "1" {
		//开启事务
		err := mysqldb.Transaction(func(tx *gorm.DB) error {
			//创建新的关注记录
			if err := tx.Create(follow).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}

			//获取关注数并更新
			var FollowCount int64
			err := mysqldb.Model(&TableEntity.UserInfo{}).Select("follow_count").
				Where("id = ? ", userID).Scan(&FollowCount).Error
			if err != nil {
				return err
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", userID).
				Update("follow_count", FollowCount+1).Error; err != nil {
				return err
			}

			//获取粉丝数并更新
			var FollowerCount int64
			err2 := mysqldb.Model(&TableEntity.UserInfo{}).Select("follower_count").
				Where("id = ? ", toUserId).Scan(&FollowerCount).Error
			if err2 != nil {
				return err2
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", toUserId).
				Update("follower_count", FollowerCount+1).Error; err != nil {
				return err
			}

			// 返回 nil 提交事务
			return nil
		})
		if err != nil {
			Log.ErrorLogWithoutPanic("Follow Operation Error!", err)
			return err
		}
		return nil
	} else if ActionType == "2" {
		// 此时操作为取消关注
		// 开启事务

		// 删除关注表信息
		err := mysqldb.Transaction(func(tx *gorm.DB) error {
			if err := mysqldb.Model(&TableEntity.Follow{}).
				Where("user_id = ? AND to_user_id = ?", userID, toUserId).
				Delete(&TableEntity.Follow{}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}

			// 更新关注数
			var FollowCount int64
			err := mysqldb.Model(&TableEntity.UserInfo{}).Select("follow_count").
				Where("id = ? ", userID).
				Scan(&FollowCount).Error
			// 当存在关注数为0时（非法请求）时对相关请求进行拦截
			if err != nil || FollowCount == 0 {
				return err
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", userID).
				Update("follow_count", FollowCount-1).Error; err != nil {
				return err
			}

			// 更新粉丝数
			var FollowerCount int64
			err2 := mysqldb.Model(&TableEntity.UserInfo{}).Select("follower_count").
				Where("id = ? ", toUserId).
				Scan(&FollowerCount).Error
			// 当存在粉丝数为0时（非法请求）时对相关请求进行拦截
			if err != nil || FollowerCount == 0 {
				return err2
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", toUserId).
				Update("follower_count", FollowerCount-1).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})

		if err != nil {
			Log.ErrorLogWithoutPanic("UnFollow Error!", err)
			return err
		}
		return nil
	}
	return nil
}

// GetFollowing 查询关注列表
func (dao *RelationDao) GetFollowing(userID int64) ([]RequestEntity.User, error) {

	// 查询关注者的ID列表
	var focus []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("to_user_id").
		Where("user_id = ?", userID).
		Find(&focus).Error

	if err != nil {
		Log.ErrorLogWithoutPanic("Get Following List Error!", err)
		return nil, nil
	}

	// 获取关注者的个人信息列表
	var user []RequestEntity.User
	for i := 0; i < len(focus); i++ {
		var userInfo TableEntity.UserInfo

		// 将表中存储的string类型UID转换成int64类型
		focusUserId, err := strconv.ParseInt(focus[i], 10, 64)
		if err != nil {
			Log.ErrorLogWithoutPanic("UID Transform Error!", err)
			return nil, err
		}

		// 查找表内user信息
		err = mysqldb.Model(&TableEntity.UserInfo{}).Where("id = ?", focusUserId).
			Find(&userInfo).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("GetFollowing Error!", err)
			return nil, nil
		}

		// 做信息转换并添加到列表
		var rUserInfo = RequestEntity.User{
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			FavoriteCount:   userInfo.FavoriteCount,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			ID:              userInfo.ID,
			IsFollow:        true,
			Name:            userInfo.Name,
			Signature:       userInfo.Signature,
			TotalFavorited:  userInfo.TotalFavorited,
			WorkCount:       userInfo.WorkCount,
		}
		user = append(user, rUserInfo)
	}

	return user, nil
}

// GetFollowers 查询粉丝列表
func (dao *RelationDao) GetFollowers(toUserID int64) ([]RequestEntity.User, error) {

	// 查询粉丝的ID列表
	var fans []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("user_id").
		Where("to_user_id = ?", toUserID).
		Find(&fans).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFollowers Error!", err)
		return nil, nil
	}

	// 获取粉丝的信息列表
	var user []RequestEntity.User
	for i := 0; i < len(fans); i++ {
		var userInfo TableEntity.UserInfo

		// 将表中存储的string类型UID转换成int64类型
		fanUserId, err := strconv.ParseInt(fans[i], 10, 64)
		if err != nil {
			Log.ErrorLogWithoutPanic("UID Transform Error!", err)
			return nil, err
		}

		// 查找user信息
		err = mysqldb.Model(&TableEntity.UserInfo{}).
			Where("id = ?", fanUserId).Find(&userInfo).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("GetFollowers Error!", err)
			return nil, nil
		}

		// 做信息转换
		var rUserInfo = RequestEntity.User{
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			FavoriteCount:   userInfo.FavoriteCount,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			ID:              userInfo.ID,
			// 这里调用方法进行判别
			IsFollow:       dao.FindRelation(toUserID, userInfo.ID, "1"),
			Name:           userInfo.Name,
			Signature:      userInfo.Signature,
			TotalFavorited: userInfo.TotalFavorited,
			WorkCount:      userInfo.WorkCount,
		}
		user = append(user, rUserInfo)
	}
	return user, nil
}

// GetFriends 查询好友列表
func (dao *RelationDao) GetFriends(userID int64) ([]RequestEntity.User, error) {

	// 先获取关注者列表
	var focus []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("to_user_id").
		Where("user_id = ?", userID).
		Find(&focus).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFriends Error!", err)
		return nil, nil
	}

	//再获取粉丝列表
	var fans []string
	err2 := mysqldb.Model(&TableEntity.Follow{}).
		Select("user_id").
		Where("to_user_id = ?", userID).
		Find(&fans).Error
	if err2 != nil {
		Log.ErrorLogWithoutPanic("GetFriends Error!", err)
		return nil, nil
	}

	friendsList := Util.FindFriends(focus, fans)
	var user []RequestEntity.User
	for i := 0; i < len(friendsList); i++ {
		var userInfo TableEntity.UserInfo
		err := mysqldb.Model(&TableEntity.UserInfo{}).
			Where("id = ?", friendsList[i]).Find(&userInfo).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("GetFriends Error!", err)
			return nil, err
		}
		var vUserInfo = RequestEntity.User{
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			FavoriteCount:   userInfo.FavoriteCount,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			ID:              userInfo.ID,
			IsFollow:        true,
			Name:            userInfo.Name,
			Signature:       userInfo.Signature,
			TotalFavorited:  userInfo.TotalFavorited,
			WorkCount:       userInfo.WorkCount,
		}
		user = append(user, vUserInfo)
	}
	return user, nil
}

// IsExist 查询该用户是否存在,0表示不存在，1表示存在
func (dao *RelationDao) IsExist(ID int64) bool {
	var count int64
	err := mysqldb.Model(&TableEntity.UserInfo{}).
		Select("Count(*)").Where("id = ?", ID).
		Scan(&count).Error
	//"SELECT COUNT(*) FROM users WHERE id = ?", ID
	if err != nil {
		Log.ErrorLogWithoutPanic("GetUserInfo failed!", err)
		return false
	}

	if count > 0 {
		return true
	} else {
		Log.ErrorLogWithoutPanic("The user does not exist!", err)
		return false
	}
}

// FindRelation 查询该用户与另一用户的关系 true:有关注关系 false:没有关注关系
func (dao *RelationDao) FindRelation(userID int64, toUserID int64, ActionType string) bool {
	if ActionType == "1" {
		var result int64
		err := mysqldb.Model(&TableEntity.Follow{}).
			Select("Count(*)").
			Where("user_id = ? AND to_user_id = ?", userID, toUserID).
			Scan(&result).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("FindRelation error!", err)
			return false
		}
		if result <= 0 {
			return false
		} else {
			Log.ErrorLogWithoutPanic("adding failed!", err)
			return true
		}
	}
	return true
}

// FindFocusByUID 根据用户UID查询关注者的UID
func (dao *RelationDao) FindFocusByUID(UID int64) []int64 {
	var (
		followUid []int64
		err       error
	)
	if UID == -2 {
		err = mysqldb.Model(TableEntity.Follow{}).
			Select("to_user_id").Find(&followUid).Error
	} else {
		err = mysqldb.Model(TableEntity.Follow{}).
			Select("to_user_id").Where("user_id =?", UID).Find(&followUid).Error
	}
	if err != nil {
		Log.ErrorLogWithoutPanic("Follow Find Error!", err)
		return nil
	}
	return followUid
}

func FindFansByUID(UID int64) []int64 {
	var (
		followUid []int64
		err       error
	)
	if UID == -2 {
		err = mysqldb.Model(TableEntity.Follow{}).
			Select("user_id").Find(&followUid).Error
	} else {
		err = mysqldb.Model(TableEntity.Follow{}).
			Select("user_id").Where("to_user_id =?", UID).Find(&followUid).Error
	}
	if err != nil {
		Log.ErrorLogWithoutPanic("Follow Find Error!", err)
		return nil
	}
	return followUid
}
