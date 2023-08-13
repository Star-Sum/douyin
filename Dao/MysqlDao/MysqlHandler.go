package MysqlDao

import "gorm.io/gorm"

var mysqldb *gorm.DB

func GetMysqlDBHandler(db *gorm.DB) {
	mysqldb = db
}
