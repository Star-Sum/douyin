package RedisDao

import (
	"context"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/utils"
	"strconv"
	"time"
)

type PublishDao interface {
	GetVedioByAuthorId(UID int64) []TableEntity.VedioInfo
	InsertVedioInfo(UID int64, info []TableEntity.VedioInfo) error
}

type PublishDaoImpl struct {
	Ctx context.Context
}

// InsertVedioInfo 投稿后在作者栏内插入信息
func (p *PublishDaoImpl) InsertVedioInfo(UID int64, info []TableEntity.VedioInfo) error {
	stringKey := "Publish:_" + strconv.FormatInt(UID, 10)

	for i := 0; i < len(info); i++ {
		stringInfo, err := Util.JsonMarshal(info[i])
		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Info Marshal Error!", err)
			return err
		}
		err = redisHandler.ZAdd(p.Ctx, stringKey, redis.Z{
			Score:  float64(info[i].Timestamp.Unix()),
			Member: stringInfo,
		}).Err()
		if err != nil {
			Log.ErrorLogWithoutPanic("Publish Record Error!", err)
			return err
		}
	}
	Log.NormalLog("Publish Record Success!", nil)
	return nil
}

// GetVedioByAuthorId 根据作者ID查找数据
func (p *PublishDaoImpl) GetVedioByAuthorId(UID int64) []TableEntity.VedioInfo {
	var (
		vInfo TableEntity.VedioInfo
	)
	stringKey := "Publish:" + strconv.FormatInt(UID, 10)
	// 每次抽取100条数据
	vedioList, err := redisHandler.ZRevRangeByScoreWithScores(p.Ctx, stringKey, &redis.ZRangeBy{
		Min:   strconv.FormatInt(0, 10),
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: int64(100),
	}).Result()
	if err != nil {
		Log.ErrorLogWithoutPanic("Publish Info Get Error!", err)
		return nil
	}
	var vedioInfo []TableEntity.VedioInfo
	for i := 0; i < len(vedioList); i++ {
		byteString := utils.ToString(vedioList[i].Member)
		err = Util.JsonUnmarshal([]byte(byteString), &vInfo)
		vedioInfo = append(vedioInfo, vInfo)
	}
	return vedioInfo
}
