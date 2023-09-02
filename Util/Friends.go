package Util

import (
	"douyin/Log"
	"fmt"
	"sort"
	"strconv"
)

// FindFriends 做朋友查找
func FindFriends(followsId, fansId []string) []int64 {
	var (
		friendsList []int64
		err         error
	)
	// 将string类型的UID转换为int64类型的UID
	followList := make([]int64, len(followsId))
	fansList := make([]int64, len(fansId))

	for i := 0; i < len(followsId); i++ {
		followList[i], err = strconv.ParseInt(followsId[i], 10, 64)
		if err != nil {
			Log.ErrorLogWithoutPanic("Parse UID Error!", err)
			return nil
		}
	}
	for i := 0; i < len(fansList); i++ {
		fansList[i], err = strconv.ParseInt(fansId[i], 10, 64)
		if err != nil {
			Log.ErrorLogWithoutPanic("Parse UID Error!", err)
			return nil
		}
	}

	// 排序
	sort.Slice(followList, func(i, j int) bool {
		return followList[i] < followList[j]
	})
	sort.Slice(fansList, func(i, j int) bool {
		return fansList[i] < fansList[j]
	})
	fmt.Println(followList, fansList)
	// 获取共同元素
	p1 := 0
	p2 := 0
	for p1 < len(followList) && p2 < len(fansList) {
		if followList[p1] < fansList[p2] {
			p1++
			continue
		}
		if followList[p1] > fansList[p2] {
			p2++
			continue
		}
		if followList[p1] == fansList[p2] {
			fmt.Println(followList[p1])
			friendsList = append(friendsList, followList[p1])
			p1++
			p2++
			continue
		}
	}
	return friendsList
}
