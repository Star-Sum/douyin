package MysqlDao

import (
	"database/sql"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type RelationDao struct {
	db *sql.DB
}

// 关注操作进行时在测试中要加入user_id的参数名 ！！！

// 进行关注和取消关注的操作
func (dao *RelationDao) Follow(userID string, toUserId string, time time.Time, ActionType string) (int64, error) {
	id, _ := Util.MakeUid(Util.Snowflake.DataCenterId, Util.Snowflake.MachineId)
	var follow = &TableEntity.Follow{
		ID:       id,
		UserID:   userID,
		ToUserID: toUserId,
		Time:     time,
	}

	if ActionType == "1" {
		//开启事务
		err := mysqldb.Transaction(func(tx *gorm.DB) error {
			//创建新的关注记录
			if err := tx.Create(follow).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			//更新关注列表
			var FollowCount int64
			err := mysqldb.Model(&TableEntity.UserInfo{}).Select("follow_count").Where("id = ? ", userID).Scan(&FollowCount).Error
			if err != nil {
				return err
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", userID).
				Update("follow_count", FollowCount+1).Error; err != nil {
				return err
			}
			//更新粉丝列表
			var FollowerCount int64
			err2 := mysqldb.Model(&TableEntity.UserInfo{}).Select("follower_count").Where("id = ? ", toUserId).Scan(&FollowerCount).Error
			if err2 != nil {
				return err2
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", toUserId).
				Update("follower_count", FollowCount+1).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
		if err != nil {
			Log.ErrorLogWithoutPanic("Follow Error!", err)
			return 0, err
		}
		return 1, nil
	} else {
		//开启事务
		err := mysqldb.Transaction(func(tx *gorm.DB) error {
			if err := mysqldb.Model(&TableEntity.Follow{}).Where("user_id = ? AND to_user_id = ?", userID, toUserId).Delete(&TableEntity.Follow{}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			var FollowCount int64
			err := mysqldb.Model(&TableEntity.UserInfo{}).Select("follow_count").Where("id = ? ", userID).Scan(&FollowCount).Error
			if err != nil {
				return err
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", userID).
				Update("follow_count", FollowCount-1).Error; err != nil {
				return err
			}
			var FollowerCount int64
			err2 := mysqldb.Model(&TableEntity.UserInfo{}).Select("follower_count").Where("id = ? ", toUserId).Scan(&FollowerCount).Error
			if err != nil {
				return err2
			}
			if err := mysqldb.Model(&TableEntity.UserInfo{}).
				Where("id = ?", toUserId).
				Update("follower_count", FollowCount-1).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
			// 返回 nil 提交事务
			return nil

		})
		// 取消关注关系
		if err != nil {
			Log.ErrorLogWithoutPanic("UnFollow Error!", err)
			return 0, err
		}
		return 1, nil
	}
}

// 查询关注列表
func (dao *RelationDao) GetFollowing(userID string) ([][]RequestEntity.User, error) {
	var focus []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("to_user_id").
		Where("user_id = ?", userID).
		Find(&focus).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFollowing Error!", err)
		return nil, nil
	}

	var user [][]RequestEntity.User
	for i := 0; i < len(focus); i++ {
		var userInfo []RequestEntity.User
		err := mysqldb.Model(&TableEntity.UserInfo{}).Select("id , name").Where("id = ?", focus[i]).Find(&userInfo).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("GetFollowing Error!", err)
			return nil, nil
		}
		user = append(user, userInfo)
	}

	return user, nil
}

// 查询粉丝列表
func (dao *RelationDao) GetFollowers(userID string) ([][]RequestEntity.User, error) {
	var fans []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("user_id").
		Where("to_user_id = ?", userID).
		Find(&fans).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFollowers Error!", err)
		return nil, nil
	}

	var user [][]RequestEntity.User
	for i := 0; i < len(fans); i++ {
		var userInfo []RequestEntity.User
		err := mysqldb.Model(&TableEntity.UserInfo{}).Select("id , name").Where("id = ?", fans[i]).Find(&userInfo).Error
		if err != nil {
			Log.ErrorLogWithoutPanic("GetFollowers Error!", err)
			return nil, nil
		}
		user = append(user, userInfo)
	}
	return user, nil
}

// 查询好友列表
func (dao *RelationDao) GetFriends(userID string) ([][]RequestEntity.User, error) {

	var focus []string
	err := mysqldb.Model(&TableEntity.Follow{}).
		Select("to_user_id").
		Where("user_id = ?", userID).
		Find(&focus).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFriends Error!", err)
		return nil, nil
	}

	var fans []string
	err2 := mysqldb.Model(&TableEntity.Follow{}).
		Select("user_id").
		Where("to_user_id = ?", userID).
		Find(&fans).Error
	if err2 != nil {
		Log.ErrorLogWithoutPanic("GetFriends Error!", err)
		return nil, nil
	}

	mp := make(map[string]bool)

	for _, value := range focus {
		mp[value] = true
	}
	fmt.Println(focus)
	fmt.Println(fans)

	var user [][]RequestEntity.User
	for i := 0; i < len(fans); i++ {
		if mp[fans[i]] {
			var userInfo []RequestEntity.User
			err := mysqldb.Model(&TableEntity.UserInfo{}).Select("id , name").Where("id = ?", fans[i]).Find(&userInfo).Error
			if err != nil {
				Log.ErrorLogWithoutPanic("GetFriends Error!", err)
				return nil, nil
			}
			user = append(user, userInfo)
		}
	}
	return user, nil
}

// 查询该用户是否存在,0表示存在，1表示不存在
func (dao *RelationDao) IsExist(ID string) (res bool) {
	var count int64
	err := mysqldb.Model(&TableEntity.UserInfo{}).Select("Count(*)").Where("id = ?", ID).Scan(&count).Error
	//"SELECT COUNT(*) FROM users WHERE id = ?", ID
	if err != nil {
		Log.ErrorLogWithoutPanic("GetUserInfo failed!", err)
		return false
	}

	if count > 0 {
		return true
	}
	Log.ErrorLogWithoutPanic("The user does not exist!", err)
	return false
}

// 查询该用户与另一用户的关系 true:没关系 false:有关系
func (dao *RelationDao) FindRelation(userID, toUserID string) (res bool) {
	var result int64
	err := mysqldb.Model(&TableEntity.Follow{}).Select("Count(*)").Where("user_id = ? AND to_user_id = ?", userID, toUserID).Scan(&result).Error
	if err != nil {
		Log.ErrorLogWithoutPanic("FindRelation error!", err)
		return false
	}
	if result == 0 {
		return true
	}
	Log.ErrorLogWithoutPanic("adding failed!", err)
	return false
}

// FindFocusByUID 根据用户UID查询关注者的UID
func FindFocusByUID(UID int64) []int64 {
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
