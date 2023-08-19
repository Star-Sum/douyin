package Util

import (
	"douyin/Log"
	"fmt"
	"time"
)

func StringToTimeStamp(stringTime string) (int64, error) {
	var Loc, _ = time.LoadLocation("Asia/Shanghai")
	//时间未设置就返回当前时间
	if stringTime == "" {
		return time.Now().Unix(), nil
	}
	timeTemplate := "2006-01-02 15:04:05"
	inputTime, err := time.ParseInLocation(timeTemplate, stringTime, Loc)
	// 返回值为-1表示时间解析错误
	if err != nil {
		Log.ErrorLogWithoutPanic("Time Parse Error!", err)
		return -1, err
	}
	return inputTime.Unix(), nil
}

func TimeStampToString(intTime int64) string {
	// 基准字符串，表示时间转换的基准
	baseString := "2006-01-02 15:04:05"
	//时间戳有错误就返回空
	if intTime == -1 {
		return ""
	}
	midTime := time.Unix(intTime, 0)
	fmt.Println(intTime)
	stringTime := midTime.Format(baseString)
	return stringTime
}
