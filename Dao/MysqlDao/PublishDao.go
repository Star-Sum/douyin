package MysqlDao

import (
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"gorm.io/gorm"
)

type PublishDao interface {
	FindAuthorByVedioId(vedioId int64) int64
	PublishVedio(info TableEntity.PublishInfo)
	VedioInfoInsert(info TableEntity.VedioInfo)
}
type PublishDaoImpl struct {
}

func (p *PublishDaoImpl) FindAuthorByVedioId(vedioId int64) int64 {
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

func (p *PublishDaoImpl) PublishVedio(info TableEntity.PublishInfo) {
	err := mysqldb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&info).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		Log.ErrorLogWithoutPanic("Publish Vedio Info Error!", err)
	}
	Log.NormalLog("Publish Vedio Info Success!", nil)
}

func (p *PublishDaoImpl) VedioInfoInsert(info TableEntity.VedioInfo) {
	err := mysqldb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&info).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		Log.ErrorLogWithoutPanic("Vedio Info Insert Error!", err)
	}
	Log.NormalLog("Vedio Info Insert Success!", nil)
}
