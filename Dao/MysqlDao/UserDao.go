package MysqlDao

import (
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"gorm.io/gorm"
)

type UserDao interface {

	// AddUser 添加用户，同时操作user_infos和user_account_infos表，并且返回userId
	AddUser(rUser RequestEntity.RegisterRequest) (int64, error)

	// CountUserByUsername 根据用户名获取user的数量
	CountUserByUsername(username string) (int64, error)

	// GetUserById 根据id获取user
	GetUserById(id int64) (*RequestEntity.User, error)

	// GetUserIdByUsernameANDPassword 根据用户名和密码获取userId
	GetUserIdByUsernameANDPassword(username, password string) (*int64, error)

	IncrementFields(userId int64, fields string) error

	DecrementField(userId int64, fields string) error
}
type UserDaoImpl struct {
}

func (u *UserDaoImpl) AddUser(rUser RequestEntity.RegisterRequest) (int64, error) {

	//根据雪花算法生成id
	userId, _ := Util.MakeUid(Util.Snowflake.DataCenterId, Util.Snowflake.MachineId)

	//MD5加密密码
	password := Util.CalcMD5(rUser.Password)

	var userAccount = &TableEntity.UserAccountInfo{
		UserID:   userId,
		Password: password,
		Username: rUser.Username,
	}

	var user = &TableEntity.UserInfo{
		FavoriteCount:  0,
		FollowCount:    0,
		FollowerCount:  0,
		ID:             userId,
		IsFollow:       false,
		Name:           rUser.Username,
		Signature:      "大家好，我是" + rUser.Username,
		TotalFavorited: 0,
		WorkCount:      0,
	}

	//开启gorm事务
	err := mysqldb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userAccount).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (u *UserDaoImpl) CountUserByUsername(username string) (int64, error) {

	var num int64

	err := mysqldb.Model(&TableEntity.UserAccountInfo{}).Where("username =?", username).Count(&num).Error

	if err != nil {
		Log.ErrorLogWithoutPanic("CountUserByUsername出错", err)
		return -1, err
	}

	return num, nil
}

func (u *UserDaoImpl) GetUserById(id int64) (*RequestEntity.User, error) {

	//根据id获取user
	var user = &RequestEntity.User{}

	err := mysqldb.Model(&TableEntity.UserInfo{}).Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserDaoImpl) GetUserIdByUsernameANDPassword(username, password string) (*int64, error) {

	//根据UsernameANDPassword获取userId
	var users []TableEntity.UserAccountInfo

	err := mysqldb.Model(&TableEntity.UserAccountInfo{}).
		Where("username = ? AND password = ?", username, password).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	//没有查询到记录
	if len(users) == 0 {
		return nil, nil
	}

	return &users[0].UserID, nil
}

func (u *UserDaoImpl) IncrementFields(userId int64, fields string) error {

	if err := mysqldb.Model(&TableEntity.UserInfo{}).
		Where("id = ?", userId).
		Update(fields, gorm.Expr(fields+" + 1")).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserDaoImpl) DecrementField(userId int64, fields string) error {
	if err := mysqldb.Model(&TableEntity.UserInfo{}).
		Where("id = ?", userId).
		Update(fields, gorm.Expr(fields+" - 1")).Error; err != nil {
		return err
	}
	return nil
}
