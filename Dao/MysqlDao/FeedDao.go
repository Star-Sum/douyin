package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"fmt"
)

type FeedDao interface {
	GetVedioList(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo
}

type FeedDaoImpl struct {
}

// GetVedioList 获取vedio信息列表
func (u *FeedDaoImpl) GetVedioList(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo {
	var vedioTable []TableEntity.VedioInfo

	stringTime := Util.TimeStampToString(timeStamp)
	fmt.Println(stringTime)
	if stringTime == "" {
		return nil
	}

	if UID == -2 {
		// 不指定UID的情况，查询最近30条
		err := mysqldb.Model(&TableEntity.VedioInfo{}).
			Where("timestamp <= ?", stringTime).Limit(num).Order("timestamp desc").Find(&vedioTable).Error

		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Get Error!", err)
			return nil
		}

	} else {
		// 指定UID的情况下查询最近30条
		err := mysqldb.Model(TableEntity.VedioInfo{}).
			Where("timestamp <= ? AND author_id = ?", stringTime, UID).Limit(num).
			Order("timestamp desc").Find(&vedioTable).Error

		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Get Error!", err)
			return nil
		}
	}

	return vedioTable
}
