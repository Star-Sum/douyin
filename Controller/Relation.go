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

func FocusListImp(request RequestEntity.UserInfoRequest) RequestEntity.FocusListBack {
	var (
		focusListBack RequestEntity.FocusListBack
	)
	focusListBack = Service.FocusListProcess(request)
	return focusListBack
}

func FanListImp(request RequestEntity.UserInfoRequest) RequestEntity.FansListBack {
	var (
		fanListBack RequestEntity.FansListBack
	)
	fanListBack = Service.FanListProcess(request)
	return fanListBack
}

func FriendListImp(request RequestEntity.UserInfoRequest) RequestEntity.FriendsListBack {
	var (
		fanListBack RequestEntity.FriendsListBack
	)
	fanListBack = Service.FriendListProcess(request)
	return fanListBack
}
