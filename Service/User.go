package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"douyin/Util"
	"strconv"
)

var MUserDaoImpl MysqlDao.UserDao = &MysqlDao.UserDaoImpl{}

var RUserDaoImpl RedisDao.UserDao = &RedisDao.UserDaoImpl{
	Ctx: context.Background(),
}

func UserRegisterProcess(request RequestEntity.RegisterRequest) RequestEntity.RegisterBack {
	//var (
	//	registerBack RequestEntity.RegisterBack
	//)
	Log.NormalLog("UserRegisterProcess", nil)

	//获取使用该用户名的账号数量
	num, err := MUserDaoImpl.CountUserByUsername(request.Username)
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
	userId, err := MUserDaoImpl.AddUser(request)
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
	userId, err := MUserDaoImpl.GetUserIdByUsernameANDPassword(request.Username, Util.CalcMD5(request.Password))
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

	toUserId, _ := strconv.Atoi(request.UserID)

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

	var user *RequestEntity.User

	//判断缓存中是否存在用户信息
	exists, _ := RUserDaoImpl.Exists(int64(toUserId))

	if exists {

		user = RUserDaoImpl.GetUserInfo(int64(toUserId))

	} else {

		//获取用户信息
		user, err = MUserDaoImpl.GetUserById(int64(toUserId))
		if err != nil {
			statusMsg := "get user info failed"
			Log.ErrorLogWithoutPanic(statusMsg, err)
			return RequestEntity.UserInfoBack{
				StatusCode: 1,
				StatusMsg:  &statusMsg,
				User:       nil,
			}
		}

		//向缓存中新增user信息
		defer func() {
			err = RUserDaoImpl.AddUserInfo(user)
			if err != nil {
				Log.ErrorLogWithoutPanic("向redis中新增user信息失败", err)
			}
		}()

	}

	//判断是否关注
	isFollow, err := MUserDaoImpl.GetIsFollowByUserIdAndToUserId(userId, int64(toUserId))
	user.IsFollow = isFollow

	//返回用户信息
	statusMsg := "get user info success"
	return RequestEntity.UserInfoBack{
		StatusCode: 0,
		StatusMsg:  &statusMsg,
		User:       user,
	}
}
