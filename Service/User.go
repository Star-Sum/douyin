package Service

import (
	"douyin/Entity/RequestEntity"
)

func UserRegisterProcess(request RequestEntity.RegisterRequest) RequestEntity.RegisterBack {
	var (
		registerBack RequestEntity.RegisterBack
	)
	return registerBack
}

func UserLoginProcess(request RequestEntity.LoginRequest) RequestEntity.LoginBack {
	var (
		loginBack RequestEntity.LoginBack
	)
	return loginBack
}

func UserInfoProcess(request RequestEntity.UserInfoRequest) RequestEntity.UserInfoBack {
	var (
		userInfoBack RequestEntity.UserInfoBack
	)
	return userInfoBack
}
