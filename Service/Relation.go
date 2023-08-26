package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"fmt"
	"strconv"
	"time"
)

var (
	MRelationDaoImpl MysqlDao.RelationDao = *&MysqlDao.RelationDao{}
	RRelationDaoImpl RedisDao.RelationDao = *&RedisDao.RelationDao{
		Ctx: context.Background(),
	}
)

func FocusProcess(request RequestEntity.FocusRequest) RequestEntity.FocusBack {
	userId, _ := strconv.ParseInt(request.UserId, 10, 64)
	ToUserId, _ := strconv.ParseInt(request.ToUserID, 10, 64)
	//如果user_id 和to_user_id一样的话报错
	if userId == ToUserId {
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed,user can't Focus themselves",
		}
	}
	//关注之前首先进行判断，查看该用户是否存在
	if !MRelationDaoImpl.IsExist(userId) || !MRelationDaoImpl.IsExist(ToUserId) {
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed,the user not exists",
		}
	}
	if !MRelationDaoImpl.FindRelation(userId, ToUserId, request.ActionType) {
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed,Already exists relationship",
		}
	}
	//每隔3秒将数据写入到mysql中
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		exists, _ := RRelationDaoImpl.ExistsFollow(request.UserId, request.ToUserID)
		if exists {
			fmt.Println("存在")
		} else {
			_ = RRelationDaoImpl.FollowUser(request.UserId, request.ToUserID)
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				// 从Redis中读取数据
				var followAll []RedisDao.Follow
				followAll, _ = RRelationDaoImpl.WriteFollowDataToDatabase(userId, ToUserId)
				for _, follow := range followAll {
					Log.NormalLog("FocusProcess", nil)
					now := time.Now()
					// 加载"Asia/Shanghai"时区
					loc, err := time.LoadLocation("Asia/Shanghai")
					timeNow := now.In(loc)
					_, err = MRelationDaoImpl.Follow(follow.UserID, follow.ToUserID, timeNow, request.ActionType)
					if err != nil {
						Log.ErrorLogWithoutPanic("get focus failed", err)
					}
				}
			}
		}
	}()
	return RequestEntity.FocusBack{
		StatusCode: 1,
		StatusMsg:  "get focus success",
	}
}

func FocusListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FocusListBack {

	userId, err := strconv.ParseInt(request.UserID, 10, 64)
	Log.NormalLog("FocusListProcess", nil)
	if err != nil {
		statusMsg := "get focusList failed"
		Log.ErrorLogWithoutPanic(statusMsg, err)
		return RequestEntity.FocusListBack{
			StatusCode: "1",
			StatusMsg:  "get focusList failed",
			UserList:   nil,
		}
	}

	var FocusList [][]RequestEntity.User

	exists, err := RRelationDaoImpl.ExistsFocus(userId)
	if exists {
		fmt.Println("存在")
		FocusList = RRelationDaoImpl.GetFocusRelation(userId)

	} else {

		//获取用户信息
		FocusList, err = MRelationDaoImpl.GetFollowing(userId)
		if err != nil {
			statusMsg := "get focusList failed"
			Log.ErrorLogWithoutPanic(statusMsg, err)
			return RequestEntity.FocusListBack{
				StatusCode: "1",
				StatusMsg:  "get focusList failed",
				UserList:   nil,
			}
		}
		defer func() {
			err = RRelationDaoImpl.AddFocusRelationInfo(userId, FocusList)
			if err != nil {
				fmt.Println(err)
				Log.ErrorLogWithoutPanic("向redis中新增relation信息失败", err)
			}
		}()
	}
	return RequestEntity.FocusListBack{
		StatusCode: "0",
		StatusMsg:  "get focusList success",
		UserList:   FocusList,
	}
}
func FanListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FansListBack {
	var FanList [][]RequestEntity.User
	Log.NormalLog("FanListProcess", nil)

	userId, err := strconv.ParseInt(request.UserID, 10, 64)
	exists, err := RRelationDaoImpl.ExistsFans(userId)
	if err != nil {
		statusMsg := "get fansList failed"
		Log.ErrorLogWithoutPanic(statusMsg, err)
		return RequestEntity.FansListBack{
			StatusCode: "1",
			StatusMsg:  "get fansList failed",
			UserList:   nil,
		}
	}
	if exists {
		fmt.Println("存在")
		FanList = RRelationDaoImpl.GetFansRelation(userId)
	} else {
		//获取用户信息
		FanList, err = MRelationDaoImpl.GetFollowers(userId)
		if err != nil {
			statusMsg := "get fansList failed"
			Log.ErrorLogWithoutPanic(statusMsg, err)
			return RequestEntity.FansListBack{
				StatusCode: "1",
				StatusMsg:  "get fansList failed",
				UserList:   nil,
			}
		}
		defer func() {
			err = RRelationDaoImpl.AddFansRelationInfo(userId, FanList)
			if err != nil {
				fmt.Println(err)
				Log.ErrorLogWithoutPanic("向redis中新增relation信息失败", err)
			}
		}()
	}

	return RequestEntity.FansListBack{
		StatusCode: "0",
		StatusMsg:  "get fansList success",
		UserList:   FanList,
	}
}

func FriendListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FriendsListBack {
	var FriendList [][]RequestEntity.User

	Log.NormalLog("FriendListProcess", nil)
	userId, err := strconv.ParseInt(request.UserID, 10, 64)
	exists, _ := RRelationDaoImpl.ExistsFriend(userId)

	if err != nil {
		statusMsg := "get friendsList failed"
		Log.ErrorLogWithoutPanic(statusMsg, err)
		return RequestEntity.FriendsListBack{
			StatusCode: "1",
			StatusMsg:  "get friendsList failed",
			UserList:   nil,
		}
	}
	if exists {
		fmt.Println("存在")
		FriendList = RRelationDaoImpl.GetFriendRelation(userId)
	} else {
		//获取用户信息
		FriendList, err = MRelationDaoImpl.GetFriends(userId)
		if err != nil {
			statusMsg := "get friendsList failed"
			Log.ErrorLogWithoutPanic(statusMsg, err)
			return RequestEntity.FriendsListBack{
				StatusCode: "1",
				StatusMsg:  "get friendsList failed",
				UserList:   nil,
			}
		}

		//向缓存中新增user信息
		defer func() {
			err = RRelationDaoImpl.AddFriendRelationInfo(userId, FriendList)
			if err != nil {
				Log.ErrorLogWithoutPanic("向redis中新增relation信息失败", err)
			}
		}()

	}

	return RequestEntity.FriendsListBack{
		StatusCode: "0",
		StatusMsg:  "get friendsList success",
		UserList:   FriendList,
	}

}
