package RedisDao

import (
	"context"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"encoding/json"
	"strings"

	"fmt"
	"strconv"
)

type RelationDao struct {
	Ctx context.Context
}

// ExistsFriend 判断该用户是否存在朋友缓存池
func (r *RelationDao) ExistsFriend(userId int64) (bool, error) {
	key := "Relation:_friend" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.Exists(r.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// ExistsFans 判断该用户是否存在粉丝缓存池
func (r *RelationDao) ExistsFans(userId int64) (bool, error) {
	key := "Relation:_fans" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.Exists(r.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// ExistsFocus 判断该用户是否存在关注缓存池
func (r *RelationDao) ExistsFocus(userId int64) (bool, error) {
	key := "Relation:_focus" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.Exists(r.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// AddFriendRelationInfo 添加朋友关系到redis中
func (r *RelationDao) AddFriendRelationInfo(userId int64, user []RequestEntity.User) error {
	key := "Relation:_friend" + strconv.FormatInt(userId, 10)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = redisHandler.LPush(r.Ctx, key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFriendRelation 从缓存中查询朋友关系
func (r *RelationDao) GetFriendRelation(userId int64) []RequestEntity.User {
	key := "Relation:_friend" + strconv.FormatInt(userId, 10)

	result, err := redisHandler.LRange(r.Ctx, key, 0, -1).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("Find Friend List Error!", err)
		return nil
	}
	var users []RequestEntity.User
	// 遍历获取的结果
	for _, data := range result {
		// 反序列化字符串为[][]User类型的数据
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			panic(err)
		}
	}
	return users
}

// AddFocusRelationInfo 添加关注信息到redis中
func (r *RelationDao) AddFocusRelationInfo(userId int64, user []RequestEntity.User) error {
	key := "Relation:_focus" + strconv.FormatInt(userId, 10)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = redisHandler.LPush(r.Ctx, key, data).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetFocusRelation 从缓存中查询关注关系
func (r *RelationDao) GetFocusRelation(userId int64) []RequestEntity.User {
	key := "Relation:_focus" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.LRange(r.Ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	var users []RequestEntity.User
	// 遍历获取的结果
	for _, data := range result {
		// 反序列化字符串为[][]User类型的数据
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			panic(err)
		}
	}
	return users
}

// AddFansRelationInfo 添加粉丝内容到redis中
func (r *RelationDao) AddFansRelationInfo(userId int64, user []RequestEntity.User) error {
	key := "Relation:_fans" + strconv.FormatInt(userId, 10)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = redisHandler.LPush(r.Ctx, key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFansRelation 从缓存中查询粉丝关系
func (r *RelationDao) GetFansRelation(userId int64) []RequestEntity.User {
	key := "Relation:_fans" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.LRange(r.Ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	var users []RequestEntity.User
	// 遍历获取的结果
	for _, data := range result {
		// 反序列化字符串为[][]User类型的数据
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			panic(err)
		}
	}
	return users
}

type Follow struct {
	UserID   int64
	ToUserID int64
}

// FollowUser 实现数据持久化
func (r *RelationDao) FollowUser(userID, toUserID string) error {
	key := "Relation:_follow" + userID + "to" + toUserID
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	toUserId, err := strconv.ParseInt(toUserID, 10, 64)
	if err != nil {
		return err
	}
	follow := Follow{
		UserID:   userId,
		ToUserID: toUserId,
	}

	// 将关注信息序列化为JSON字符串
	followJSON, err := json.Marshal(follow)
	if err != nil {
		return fmt.Errorf("failed to serialize follow data: %v", err)
	}

	// 将关注信息写入Redis列表中
	err = redisHandler.LPush(r.Ctx, key, followJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to write follow data to Redis: %v", err)
	}
	return nil
}

func (r *RelationDao) WriteFollowDataToDatabase(userId, toUserId int64) ([]Follow, error) {
	keys, err := redisHandler.Keys(context.Background(), "Relation:_follow*").Result()
	if err != nil {
		panic(err)
	}
	// 遍历所有键并输出
	var follows []Follow
	for _, key := range keys {
		var follow Follow
		fmt.Println(key)
		// 从Redis中读取关注数据
		followsJSON, err := redisHandler.LRange(r.Ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}

		for _, followJSON := range followsJSON {
			err = json.Unmarshal([]byte(followJSON), &follow)
			if err != nil {
				return nil, err
			}
			follows = append(follows, follow)
		}

	}
	keysAll, err := redisHandler.Keys(context.Background(), "Relation:_follow*").Result()
	// 清空Redis中的关注数据
	for _, key := range keysAll {
		fmt.Println(key)
		// 从Redis中读取关注数据
		err = redisHandler.Del(r.Ctx, key).Err()
		if err != nil {
			return nil, err
		}
	}

	var cursor uint64
	var keysOther []string

	//扫描所有的redis值
	for {
		var err error
		keysOther, cursor, err = redisHandler.Scan(r.Ctx, cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}
		// 遍历匹配到的key
		for _, key := range keysOther {
			// 根据特定的id值进行判断
			if strings.Contains(key, strconv.FormatInt(userId, 10)) {
				// 删除对应的key
				err := redisHandler.Del(r.Ctx, key).Err()
				if err != nil {
					panic(err)
				}
			}
			if strings.Contains(key, strconv.FormatInt(toUserId, 10)) {
				// 删除对应的key
				err := redisHandler.Del(r.Ctx, key).Err()
				if err != nil {
					panic(err)
				}
			}
		}
		// 如果cursor为0，表示遍历完成
		if cursor == 0 {
			break
		}
	}

	if err != nil {
		return nil, err
	}
	return follows, err
}

// ExistsFollow 判断是否存在该缓存
func (r *RelationDao) ExistsFollow(userID, touserID string) (bool, error) {
	key := "Relation:_follow" + userID + "to" + touserID
	result, err := redisHandler.Exists(r.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
