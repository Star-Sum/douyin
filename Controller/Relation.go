package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func FocusImp(request RequestEntity.FocusRequest) RequestEntity.FocusBack {
	var (
		focusBack RequestEntity.FocusBack
	)
	focusBack = Service.FocusProcess(request)
	return focusBack
}

func FocusListImp(request RequestEntity.UserInfoRequest) RequestEntity.FocusBack {
	var (
		focusListBack RequestEntity.FocusBack
	)
	focusListBack = Service.FocusListProcess(request)
	return focusListBack
}

func FanListImp(request RequestEntity.UserInfoRequest) RequestEntity.FanListBack {
	var (
		fanListBack RequestEntity.FanListBack
	)
	fanListBack = Service.FanListProcess(request)
	return fanListBack
}

func FriendListImp(request RequestEntity.UserInfoRequest) RequestEntity.FanListBack {
	var (
		fanListBack RequestEntity.FanListBack
	)
	fanListBack = Service.FriendListProcess(request)
	return fanListBack
}
