package RedisDao

import (
	"context"
	"douyin/Entity/RequestEntity"
	"strconv"
	"time"
)

type UserDao interface {

	// Exists 根据userId判断redis中是否缓存了对于user的信息
	Exists(userId int64) (bool, error)

	// AddUserInfo 向redis缓存中增加user的信息
	AddUserInfo(user *RequestEntity.User) error

	// RemoveUserInfo 向redis缓存中删除对应user的信息
	RemoveUserInfo(userId int64) error

	// GetUserInfo 从redis缓存中获取user的信息
	GetUserInfo(userId int64) *RequestEntity.User
}

type UserDaoImpl struct {
	Ctx context.Context
}

func (u *UserDaoImpl) Exists(userId int64) (bool, error) {
	key := "userInfo_" + strconv.FormatInt(userId, 10)
	result, err := redisHandler.Exists(u.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (u *UserDaoImpl) AddUserInfo(user *RequestEntity.User) error {
	key := "userInfo_" + strconv.FormatInt(user.ID, 10)
	err := redisHandler.Set(u.Ctx, key, user, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserDaoImpl) RemoveUserInfo(userId int64) error {
	key := "userInfo_" + strconv.FormatInt(userId, 10)
	err := redisHandler.Del(u.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
func (u *UserDaoImpl) GetUserInfo(userId int64) *RequestEntity.User {
	key := "userInfo_" + strconv.FormatInt(userId, 10)
	user := &RequestEntity.User{}
	err := redisHandler.Get(u.Ctx, key).Scan(user)
	if err != nil {
		return nil
	}
	return user
}
