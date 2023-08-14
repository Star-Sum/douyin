package Service

import (
	"douyin/Dao/MysqlDao"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"douyin/Util"
	"strconv"
)

var userDaoImpl MysqlDao.UserDao = &MysqlDao.UserDaoImpl{}

func UserRegisterProcess(request RequestEntity.RegisterRequest) RequestEntity.RegisterBack {
	//var (
	//	registerBack RequestEntity.RegisterBack
	//)
	Log.NormalLog("UserRegisterProcess", nil)

	//获取使用该用户名的账号数量
	num, err := userDaoImpl.CountUserByUsername(request.Username)
	if err != nil {
		Log.ErrorLogWithoutPanic("register failed", err)
		return RequestEntity.RegisterBack{
			StatusCode: 1,
			StatusMsg:  "register failed",
		}
	}

	//不为0，已经存在使用该用户名的用户，不允许注册
	if num != 0 {
		Log.ErrorLogWithoutPanic("user already registered", nil)
		return RequestEntity.RegisterBack{
			StatusCode: 1,
			StatusMsg:  "user already registered",
		}
	}

	//允许注册，进行持久化储存
	userId, err := userDaoImpl.AddUser(request)
	if err != nil {
		Log.ErrorLogWithoutPanic("register failed", err)
		return RequestEntity.RegisterBack{
			StatusCode: 1,
			StatusMsg:  "register failed",
		}
	}

	//生成token，并返回给客户端
	token := Util.MakeToken(userId)
	Log.NormalLog("user register success:"+strconv.FormatInt(userId, 10), nil)
	return RequestEntity.RegisterBack{
		StatusCode: 0,
		StatusMsg:  "register success",
		Token:      token,
		UserID:     userId,
	}
}

func UserLoginProcess(request RequestEntity.LoginRequest) RequestEntity.LoginBack {
	//var (
	//	loginBack RequestEntity.LoginBack
	//)

	//根据用户名和密码查询用户并返回useId
	Log.NormalLog("UserLoginProcess", nil)
	userId, err := userDaoImpl.GetUserIdByUsernameANDPassword(request.Username, Util.CalcMD5(request.Password))
	if err != nil || userId == nil {
		statusMsg := "username or password is incorrect"
		Log.NormalLog(statusMsg, err)
		return RequestEntity.LoginBack{
			StatusCode: 1,
			StatusMsg:  &statusMsg,
		}
	}

	//生成token，并返回给客户端
	statusMsg := "login success"
	token := Util.MakeToken(*userId)
	Log.NormalLog("user login success:"+strconv.FormatInt(*userId, 10), nil)
	return RequestEntity.LoginBack{
		StatusCode: 0,
		StatusMsg:  &statusMsg,
		Token:      &token,
		UserID:     userId,
	}
}

func UserInfoProcess(request RequestEntity.UserInfoRequest) RequestEntity.UserInfoBack {
	//var (
	//	userInfoBack RequestEntity.UserInfoBack
	//)

	Log.NormalLog("UserInfoProcess", nil)

	//解析token
	userId, err := Util.ParserToken(request.Token)

	//token解析错误，返回
	if err != nil {
		statusMsg := "token parse error, please log in and try again"
		Log.ErrorLogWithoutPanic(statusMsg, err)
		return RequestEntity.UserInfoBack{
			StatusCode: 1,
			StatusMsg:  &statusMsg,
			User:       nil,
		}
	}

	//获取用户信息
	user, err := userDaoImpl.GetUserById(userId)
	if err != nil {
		statusMsg := "get user info failed"
		Log.ErrorLogWithoutPanic(statusMsg, err)
		return RequestEntity.UserInfoBack{
			StatusCode: 1,
			StatusMsg:  &statusMsg,
			User:       nil,
		}
	}

	//返回用户信息
	statusMsg := "get user info success"
	userInfoBack := RequestEntity.UserInfoBack{
		StatusCode: 0,
		StatusMsg:  &statusMsg,
		User:       user,
	}
	return userInfoBack
}
