package RedisDao

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/utils"
	"strconv"
)

type FeedDao interface {
	InsertFeedInfo(vedioInfo []TableEntity.VedioInfo)
	FeedSetSize(UID int64) int
	SliceFeedInfo(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo
}

type FeedDaoImpl struct {
	Ctx context.Context
}

var (
	MFeedDaoImpl MysqlDao.FeedDao = &MysqlDao.FeedDaoImpl{}
)

// InsertFeedInfo 将每条信息ID都推送到粉丝手上
func (u *FeedDaoImpl) InsertFeedInfo(vedioInfo []TableEntity.VedioInfo) {
	for i := 0; i < len(vedioInfo); i++ {
		stringInfo, err := Util.JsonMarshal(vedioInfo[i])
		// 对公共Feed列表进行插入操作
		stringKey := "feed:-2"
		err = redisHandler.ZAdd(u.Ctx, stringKey, redis.Z{
			// 以时间戳（时间先后顺序）作为评判标准
			Score: float64(vedioInfo[i].Timestamp.Unix()),
			// 传入ID
			Member: stringInfo,
		}).Err()
		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Feed Push Error", err)
		}
		// 拉取自己的粉丝列表，给粉丝的Feed推送数据
		fanUID := MysqlDao.FindFansByUID(vedioInfo[i].ID)
		for j := 0; j < len(fanUID); j++ {
			stringKey = "feed:" + strconv.FormatInt(fanUID[j], 10)
			err = redisHandler.ZAdd(u.Ctx, stringKey, redis.Z{
				// 以时间戳（时间先后顺序）作为评判标准
				Score: float64(vedioInfo[i].Timestamp.Unix()),
				// 传入ID
				Member: stringInfo,
			}).Err()
			if err != nil {
				Log.ErrorLogWithoutPanic("Vedio Feed Push Error", err)
			}
		}
	}
}

// SliceFeedInfo 切片，选取若干数量的视频信息
func (u *FeedDaoImpl) SliceFeedInfo(timeStamp int64, num int, UID int64) []TableEntity.VedioInfo {
	var (
		vedioInfo []TableEntity.VedioInfo
		vedioList TableEntity.VedioInfo
	)
	stringKey := "feed:" + utils.ToString(UID)
	scoreList, err := redisHandler.ZRevRangeByScoreWithScores(u.Ctx, stringKey, &redis.ZRangeBy{
		Min:   strconv.FormatInt(0, 10),
		Max:   strconv.FormatInt(timeStamp, 10),
		Count: int64(num),
	}).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("Redis Get Error!", err)
		return nil
	}

	var count int
	count = 0
	for i := 0; i < len(scoreList); i++ {
		count++
		byteString := utils.ToString(scoreList[i].Member)
		err = Util.JsonUnmarshal([]byte(byteString), &vedioList)
		if err != nil {
			return nil
		}
		vedioInfo = append(vedioInfo, MFeedDaoImpl.FindVedioById(vedioList.ID))
		// 数据量大于30后就退出
		if count > 30 {
			break
		}
	}
	return vedioInfo
}

// FeedSetSize 查询Feed列表中元素的个数
func (u *FeedDaoImpl) FeedSetSize(UID int64) int {
	stringKey := "feed:" + strconv.FormatInt(UID, 10)
	num, err := redisHandler.ZCard(u.Ctx, stringKey).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("Redis Read Error!", err)
	}
	return int(num)
}
