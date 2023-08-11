package Service

import "douyin/Entity/RequestEntity"

func FocusProcess(request RequestEntity.FocusRequest) RequestEntity.FocusBack {
	var (
		focusBack RequestEntity.FocusBack
	)
	return focusBack
}

func FocusListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FocusBack {
	var (
		focusListBack RequestEntity.FocusBack
	)
	return focusListBack
}

func FanListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FanListBack {
	var (
		fanListBack RequestEntity.FanListBack
	)
	return fanListBack
}

func FriendListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FanListBack {
	var (
		fanListBack RequestEntity.FanListBack
	)
	return fanListBack
}
