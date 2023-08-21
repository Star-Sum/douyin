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
		err := mysqldb.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(follow).Error; err != nil {
				// 返回任何错误都会回滚事务
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
		// 取消关注关系
		err := mysqldb.Model(&TableEntity.Follow{}).Where("user_id = ? AND to_user_id = ?", userID, toUserId).Delete(&TableEntity.Follow{}).Error
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
