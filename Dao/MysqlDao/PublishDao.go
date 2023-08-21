package MysqlDao

import "douyin/Entity/TableEntity"

func FindAuthorByVedioId(vedioId int64) int64 {
	var (
		UID int64
	)
	UID = 0
	err := mysqldb.Model(TableEntity.VedioInfo{}).Where("id =?", vedioId).Find(&UID).Error
	if err != nil {
		return -1
	}
	return UID
}
