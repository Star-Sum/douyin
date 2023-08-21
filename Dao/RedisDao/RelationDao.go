package RedisDao

import (
	"context"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"encoding/json"
	"errors"
)

type User struct {
	ID   int
	name string
}

type RelationDao struct {
	Ctx context.Context
}

func (dao *RelationDao) Follow(userID, followID, ActionType string) error {
	key := "follows:relation_" + userID
	err := redisHandler.HMSet(dao.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// 取消关注的操作
func (dao *RelationDao) Unfollow(userID, followID, ActionType string) error {
	key := "follows:relation_" + userID
	err := redisHandler.Del(dao.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// 查询关注列表
func (dao *RelationDao) GetFollowing(userID string) ([]TableEntity.Follow, error) {
	key := "Relation:user_" + userID
	getFollowingMap, err := redisHandler.HGetAll(dao.Ctx, key).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFollowingFromRedis Error!", err)
		return nil, err
	}
	// Redis中没找到该follow
	if len(getFollowingMap) == 0 {
		Log.NormalLog("GetFollowingFromRedis : No  found in Redis.", nil)
		return nil, errors.New("no focus found in Redis")
	}
	var follow []TableEntity.Follow
	// 转换为 TableEntity.Follow 切片
	for _, followData := range getFollowingMap {
		var followSlice []TableEntity.Follow
		if err := json.Unmarshal([]byte(followData), &followSlice); err != nil {
			Log.ErrorLogWithoutPanic("follow Unmarshal Error for followData: "+followData, err)
		} else {
			follow = append(follow, followSlice...)
		}
	}
	return follow, nil
}

// 查询粉丝列表
func (dao *RelationDao) GetFollowers(userID string) ([]TableEntity.Follow, error) {
	key := "Relation:user_" + userID
	getFollowingMap, err := redisHandler.HGetAll(dao.Ctx, key).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic(" GetFollowersFromRedis Error!", err)
		return nil, err
	}
	// Redis中没找到该follow
	if len(getFollowingMap) == 0 {
		Log.NormalLog(" GetFollowersFromRedis : No  found in Redis.", nil)
		return nil, errors.New("no focus found in Redis")
	}
	var follow []TableEntity.Follow
	// 转换为 TableEntity.Follow 切片
	for _, followData := range getFollowingMap {
		var followSlice []TableEntity.Follow
		if err := json.Unmarshal([]byte(followData), &followSlice); err != nil {
			Log.ErrorLogWithoutPanic("follow Unmarshal Error for followData: "+followData, err)
		} else {
			follow = append(follow, followSlice...)
		}
	}
	return follow, nil
}

// 查询好友列表
func (dao *RelationDao) GetFriends(userID string) ([]TableEntity.Follow, error) {
	key := "Relation:user_" + userID
	getFollowingMap, err := redisHandler.HGetAll(dao.Ctx, key).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("GetFriendsFromRedis Error!", err)
		return nil, err
	}
	// Redis中没找到该follow
	if len(getFollowingMap) == 0 {
		Log.NormalLog("GetFriendsFromRedis : No  found in Redis.", nil)
		return nil, errors.New("no friends found in Redis")
	}
	var follow []TableEntity.Follow
	// 转换为 TableEntity.Follow 切片
	for _, followData := range getFollowingMap {
		var followSlice []TableEntity.Follow
		if err := json.Unmarshal([]byte(followData), &followSlice); err != nil {
			Log.ErrorLogWithoutPanic("follow Unmarshal Error for followData: "+followData, err)
		} else {
			follow = append(follow, followSlice...)
		}
	}
	return follow, nil
}
