package Util

// 启动数据库
import (
	"context"
	"database/sql"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// 这里的下划线表示只使用mysql的init驱动
	_ "github.com/go-sql-driver/mysql"
	"os"
)

// MysqlConfig Mysql数据库配置信息，结构体内变量必须大写，yaml中变量必须小写
type MysqlConfig struct {
	Url        string `yaml:"url"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Dbname     string `yaml:"dbname"`
	Protocol   string `yaml:"protocol"`
	Charset    string `yaml:"charset"`
	Parsertime string `yaml:"parsertime"`
	Loc        string `yaml:"loc"`
	Env        string `yaml:"env"`
}

type RedisConfig struct {
	Url  string `yaml:"url"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

// DbConfig 配置信息
type DbConfig struct {
	Mysql MysqlConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
}

// ConfigGet 获取配置信息
func ConfigGet(path string) (DbConfig, error) {
	Content, err := os.ReadFile(path)
	if err != nil {
		Log.ErrorLogWithPanic("Failed to read database configuration file!", err)
	}
	Log.NormalLog("Successfully read database configuration file!", err)
	var config DbConfig
	// yaml反序列化进行信息读取
	err = yaml.Unmarshal(Content, &config)
	return config, err
}

// InitMysqlDb InitDb 初始化数据库
func InitMysqlDb() *gorm.DB {
	var (
		config          DbConfig
		mysqlConf       MysqlConfig
		mysqlInfoString string
	)
	config, err := ConfigGet("./Config/config.yaml")
	if err != nil {
		Log.ErrorLogWithPanic("Mysql Database settings acquisition failed!", err)
	}
	Log.NormalLog("Successfully obtained Mysql database settings!", err)
	mysqlConf = config.Mysql
	if mysqlConf.Env == "1" {
		mysqlConf.Url = EnvTransfer(mysqlConf.Url)
		mysqlConf.Port = EnvTransfer(mysqlConf.Port)
	}
	// 将数据库信息封装成可以识别的模式字符串
	mysqlInfoString = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true",
		mysqlConf.Username, mysqlConf.Password, mysqlConf.Protocol, mysqlConf.Url, mysqlConf.Port,
		mysqlConf.Dbname)
	mysqlDb, err := sql.Open("mysql", mysqlInfoString)
	if err != nil {
		Log.ErrorLogWithPanic("Database opening failed!", err)
	}
	Log.NormalLog("Database opened successfully!", err)
	// 测试数据库
	err = mysqlDb.Ping()
	if err != nil {
		Log.ErrorLogWithPanic("Database:"+mysqlConf.Dbname+" connected failed！", err)
	}
	Log.NormalLog("Database:"+mysqlConf.Dbname+" connected successfully！", err)
	//使用gorm框架
	GormDb, err := gorm.Open(mysql.Open(mysqlInfoString), &gorm.Config{})
	if err != nil {
		Log.ErrorLogWithPanic("Gorm framework startup failed!", err)
	}
	Log.NormalLog("Gorm framework startup successfully!", err)
	return GormDb

}

var ctx = context.Background()

func InitRedisDb() *redis.Client {
	var (
		config    DbConfig
		redisConf RedisConfig
	)
	config, err := ConfigGet("./Config/config.yaml")
	if err != nil {
		Log.ErrorLogWithPanic("Redis Database settings acquisition failed!", err)
	}
	Log.NormalLog("Successfully obtained Redis database settings!", err)
	redisConf = config.Redis
	if redisConf.Env == "1" {
		redisConf.Url = EnvTransfer(redisConf.Url)
		redisConf.Port = EnvTransfer(redisConf.Port)
	}
	addrString := fmt.Sprintf("%s:%s", redisConf.Url, redisConf.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addrString,
		Password: "",
		DB:       0,
	})
	_, err1 := rdb.Ping(ctx).Result()
	if err1 != nil {
		Log.ErrorLogWithPanic("Redis Connect Failed!", nil)
	}
	Log.NormalLog("Redis Connect Success!", nil)
	return rdb
}

func TableCreate(db *gorm.DB) {
	err := db.AutoMigrate(TableEntity.UserInfo{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create UserInfo Table Error!", err)
	}
	Log.NormalLog("Create UserInfo Table Success!", err)
	err = db.AutoMigrate(TableEntity.UserAccountInfo{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create UserAccount Table Error!", err)
	}
	Log.NormalLog("Create UserAccount Table Success!", err)
	err = db.AutoMigrate(TableEntity.VedioInfo{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create VedioInfo Table Error!", err)
	}
	Log.NormalLog("Create VedioInfo Table Success!", err)
	err = db.AutoMigrate(TableEntity.Comment{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create Comment Table Error!", err)
	}
	Log.NormalLog("Create Comment Table Success!", err)
	err = db.AutoMigrate(TableEntity.Message{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create Message Table Error!", err)
	}
	Log.NormalLog("Create Message Table Success!", err)
	err = db.AutoMigrate(TableEntity.Follow{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create Follow Table Error!", err)
	}
	Log.NormalLog("Create Follow Table Success!", err)
	err = db.AutoMigrate(TableEntity.LikeInfo{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create LikeInfo Table Error!", err)
	}
	Log.NormalLog("Create LikeInfo Table Success!", err)
	err = db.AutoMigrate(TableEntity.PublishInfo{})
	if err != nil {
		Log.ErrorLogWithoutPanic("Create PublishInfo Table Error!", err)
	}
	Log.NormalLog("Create PublishInfo Table Success!", err)
}
