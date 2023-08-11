package Controller

import (
	"douyin/Entity/RequestEntity"
	"douyin/Service"
)

func UserRegisterImp(request RequestEntity.RegisterRequest) RequestEntity.RegisterBack {
	var registerBack RequestEntity.RegisterBack
	registerBack = Service.UserRegisterProcess(request)
	return registerBack
}

func UserLoginImp(request RequestEntity.LoginRequest) RequestEntity.LoginBack {
	var loginBack RequestEntity.LoginBack
	loginBack = Service.UserLoginProcess(request)
	return loginBack
}

func UserInfoImp(request RequestEntity.UserInfoRequest) RequestEntity.UserInfoBack {
	var userInfoBack RequestEntity.UserInfoBack
	userInfoBack = Service.UserInfoProcess(request)
	return userInfoBack
}
