package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"time"
)

var (
	MRelationDaoImpl MysqlDao.RelationDao = *&MysqlDao.RelationDao{}
	RRelationDaoImpl RedisDao.RelationDao = *&RedisDao.RelationDao{
		Ctx: context.Background(),
	}
)

func FocusProcess(request RequestEntity.FocusRequest) RequestEntity.FocusBack {
	//关注之前首先进行判断，查看该用户是否存在
	if !MRelationDaoImpl.IsExist(request.UserId) || !MRelationDaoImpl.IsExist(request.ToUserID) {
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed,the user not exists",
		}
	}
	if !MRelationDaoImpl.FindRelation(request.UserId, request.ToUserID, request.ActionType) {
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed,Already exists relationship",
		}
	}
	Log.NormalLog("FocusProcess", nil)
	var now = time.Now()
	addFocus, err := MRelationDaoImpl.Follow(request.UserId, request.ToUserID, now, request.ActionType)
	if err != nil {
		Log.ErrorLogWithoutPanic("get focus failed", err)
		return RequestEntity.FocusBack{
			StatusCode: 0,
			StatusMsg:  "get focus failed",
		}
	}
	if addFocus != 0 {
		//将内容写入redis
		err = RRelationDaoImpl.Follow(request.UserId, request.ToUserID, request.ActionType)
		return RequestEntity.FocusBack{
			StatusCode: 1,
			StatusMsg:  "get focus success",
		}
	}
	//将内容写入redis
	err = RRelationDaoImpl.Unfollow(request.UserId, request.ToUserID, request.ActionType)
	return RequestEntity.FocusBack{
		StatusCode: 0,
		StatusMsg:  "out focus failed",
	}

}
func FocusListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FocusListBack {

	Log.NormalLog("FocusListProcess", nil)
	FocusList, err := MRelationDaoImpl.GetFollowing(request.UserID)
	if err != nil {
		Log.ErrorLogWithoutPanic("get focusList failed", err)
		return RequestEntity.FocusListBack{
			StatusCode: "1",
			StatusMsg:  "get focus failed",
			UserList:   nil,
		}
	}
	var _ []TableEntity.Follow
	_, err = RRelationDaoImpl.GetFollowing(request.UserID)
	return RequestEntity.FocusListBack{
		StatusCode: "0",
		StatusMsg:  "get focusList success",
		UserList:   FocusList,
	}
}
func FanListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FansListBack {
	Log.NormalLog("FanListProcess", nil)

	fanlist, err := MRelationDaoImpl.GetFollowers(request.UserID)
	if err != nil {
		Log.ErrorLogWithoutPanic("get fansList failed", err)
		return RequestEntity.FansListBack{
			StatusCode: "1",
			StatusMsg:  "get fansList failed",
			UserList:   nil,
		}
	}
	var _ []TableEntity.Follow
	_, err = RRelationDaoImpl.GetFollowers(request.UserID)
	return RequestEntity.FansListBack{
		StatusCode: "0",
		StatusMsg:  "get fansList success",
		UserList:   fanlist,
	}
}
func FriendListProcess(request RequestEntity.UserInfoRequest) RequestEntity.FriendsListBack {
	Log.NormalLog("FriendListProcess", nil)

	FriendList, err := MRelationDaoImpl.GetFriends(request.UserID)
	if err != nil {
		Log.ErrorLogWithoutPanic("get friendsList failed", err)
		return RequestEntity.FriendsListBack{
			StatusCode: "1",
			StatusMsg:  "get friendsList failed",
			UserList:   nil,
		}
	}
	var _ []TableEntity.Follow
	_, err = RRelationDaoImpl.GetFriends(request.UserID)
	return RequestEntity.FriendsListBack{
		StatusCode: "0",
		StatusMsg:  "get friendsList success",
		UserList:   FriendList,
	}
}
