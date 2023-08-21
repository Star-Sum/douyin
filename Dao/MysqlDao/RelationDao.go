package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
)

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
