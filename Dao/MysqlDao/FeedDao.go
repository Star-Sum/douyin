package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
)

type FeedDao interface {
	GetVedioList(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo
}

type FeedDaoImpl struct {
}

// GetVedioList 获取vedio信息列表
// 基本实现思路是通过UID来获取其关注者发布作品的列表，一次拉取30条
// 如果没有输入UID，表明现在处于未登录情况，就直接拉取最新的30条数据
func (u *FeedDaoImpl) GetVedioList(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo {
	var (
		vedioTable []TableEntity.VedioInfo
	)

	// 将int64形式的时间戳转换为string类型的时间戳，便于比较
	stringTime := Util.TimeStampToString(timeStamp)

	// 时间转换失效就返回空指针
	if stringTime == "" {
		return nil
	}

	// 对应未登录的情况，直接从mysql中查询30条
	if UID == -2 {
		// 不指定UID的情况，查询最近30条
		err := mysqldb.Model(&TableEntity.VedioInfo{}).
			Where("timestamp <= ?", stringTime).Limit(num).Order("timestamp desc").Find(&vedioTable).Error

		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Get Error!", err)
			return nil
		}

	} else {
		// 指定UID的情况下查询其关注者最近30条
		var midInfo []TableEntity.VedioInfo
		focusList := FindFocusByUID(UID)
		for i := 0; i < len(focusList); i++ {
			err := mysqldb.Model(TableEntity.VedioInfo{}).
				Where("timestamp <= ? AND author_id = ?", stringTime, focusList[i]).Limit(num).
				Order("timestamp desc").Find(&midInfo).Error

			if err != nil {
				Log.ErrorLogWithoutPanic("Vedio Get Error!", err)
				continue
			}
			for j := 0; j < len(midInfo); j++ {
				vedioTable = append(vedioTable, midInfo[j])
			}
		}
	}
	return vedioTable
}
