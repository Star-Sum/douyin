package Util

import (
	"douyin/Log"
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
)

// MakeUid 雪花算法实现uid生成
func MakeUid(dataCenterId int64, machineId int64) (int64, error) {
	snowId, err := snowflake.NewSnowflake(dataCenterId, machineId)
	if err != nil {
		Log.ErrorLogWithoutPanic("UID generation failed!", err)
		return -1, err
	}
	Log.NormalLog("UID generation success!", err)
	id := snowId.NextVal()
	return id, err
}
