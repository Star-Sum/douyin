package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"strconv"
	"time"
)

var (
	MLikeDaoImpl MysqlDao.LikeDao = &MysqlDao.LikeDaos{}
	MFeedDaoImpl MysqlDao.FeedDao = &MysqlDao.FeedDaoImpl{}
	RFeedDaoImpl RedisDao.FeedDao = &RedisDao.FeedDaoImpl{
		Ctx: context.Background()}
)

func VedioFeedProcess(request RequestEntity.FeedRequest) RequestEntity.FeedBack {
	var (
		vedioFeedBack RequestEntity.FeedBack
		UID           int64
		nextTime      *int64
		vedioRequest  []RequestEntity.VedioRequest
		finalTime     string
	)

	// 预先设定返回的正确与错误信息
	SuccessMsg := "Vedio Feed Success!"
	FailedMsg := "Vedio Feed Failed!"

	// 传入前端的token
	token := request.Token

	// 获取输入时间，截取时间戳前十位做解析
	inputTime := request.LatestTime
	finalTime = *inputTime
	finalTime = finalTime[0:10]

	// 解析传入的时间，并将其转换为int64位的时间戳
	midTime, err := strconv.Atoi(finalTime)
	latestTime := int64(midTime)
	// 如果出错就返回出错信息
	if err != nil {
		nowTime := time.Now().Unix()
		vedioFeedBack.StatusCode = 1
		vedioFeedBack.VideoList = nil
		vedioFeedBack.StatusMsg = &FailedMsg
		vedioFeedBack.NextTime = &nowTime
		Log.ErrorLogWithoutPanic("TimeStamp Parse Error!", err)
		return vedioFeedBack
	}

	// 解析UID,-2表示本处未指定UID
	UID = -2
	// 如果前端token不为空（处于登录状态）就对token进行解析，并提取UID
	if *token != "" {
		UID, err = Util.ParserToken(*token)
		if err != nil {
			nowTime := time.Now().Unix()
			vedioFeedBack.StatusCode = 1
			vedioFeedBack.VideoList = nil
			vedioFeedBack.StatusMsg = &FailedMsg
			vedioFeedBack.NextTime = &nowTime
			Log.ErrorLogWithoutPanic("Token Parse Error!", err)
			return vedioFeedBack
		}
	}

	// 通过UID来确定其关注者的UID
	//followUID := MysqlDao.FindFocusByUID(UID)
	//err = RedisDao.BuildFeedList(UID, followUID)
	if err != nil {
		nowTime := time.Now().Unix()
		vedioFeedBack.StatusCode = 1
		vedioFeedBack.VideoList = nil
		vedioFeedBack.StatusMsg = &FailedMsg
		vedioFeedBack.NextTime = &nowTime
		Log.ErrorLogWithoutPanic("Feed List Build Error!", err)
		return vedioFeedBack
	}

	// 根据UID和时间戳抽取数据，首先从redis中抽取
	vedioList := RFeedDaoImpl.SliceFeedInfo(latestTime, 30, UID)
	// 如果redis缓存中数据不够就从mysql中抽取
	if len(vedioList) < 30 {
		vedioList = MFeedDaoImpl.GetVedioList(latestTime, 30, UID)
		RFeedDaoImpl.InsertFeedInfo(vedioList)
	}

	// 可能会返回空列表，此时直接按照错误状况退出
	if len(vedioList) == 0 {
		nowTime := time.Now().Unix()
		vedioFeedBack.StatusCode = 1
		vedioFeedBack.VideoList = nil
		vedioFeedBack.StatusMsg = &FailedMsg
		vedioFeedBack.NextTime = &nowTime
		Log.ErrorLogWithoutPanic("Vedio Find Error!", err)
		return vedioFeedBack
	}

	// 从数据库中抽取的是原始的vedio信息，这里进行一个转换
	vedioRequest = VedioListToRequest(vedioList, UID)

	// 返回最后一条信息的时间，为下次传入做准备
	timeStamp := vedioList[len(vedioList)-1].Timestamp.Unix()
	nextTime = &timeStamp
	vedioFeedBack.StatusCode = 0
	vedioFeedBack.NextTime = nextTime
	vedioFeedBack.StatusMsg = &SuccessMsg
	vedioFeedBack.VideoList = vedioRequest
	return vedioFeedBack
}

// VedioListToRequest 将vedio列表转换成前端需要的部分
func VedioListToRequest(vedioList []TableEntity.VedioInfo, UID int64) []RequestEntity.VedioRequest {
	var (
		vedioId int64
	)
	length := len(vedioList)
	vedioRequest := make([]RequestEntity.VedioRequest, len(vedioList))
	//UID未指定时默认不喜欢该视频
	if UID == -2 {
		for i := 0; i < length; i++ {
			vedioRequest[i].IsFavorite = false
			vedioRequest[i].Title = vedioList[i].Title
			vedioRequest[i].ID = vedioList[i].ID
			vedioRequest[i].CommentCount = vedioList[i].CommentCount
			vedioRequest[i].FavoriteCount = vedioList[i].FavoriteCount
			vedioRequest[i].PlayURL = vedioList[i].PlayURL
			vedioRequest[i].CoverURL = vedioList[i].CoverURL
			user, err := MUserDaoImpl.GetUserById(vedioList[i].AuthorID)
			if err != nil {
				vedioRequest[i].Author = RequestEntity.AuthorUser{}
			}
			// 未登录自然是未关注
			vedioRequest[i].Author = RequestEntity.AuthorUser{
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				FavoriteCount:   user.FavoriteCount,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				ID:              user.ID,
				IsFollow:        false,
				Name:            user.Name,
				Signature:       user.Signature,
				TotalFavorited:  strconv.Itoa(int(user.TotalFavorited)),
				WorkCount:       user.WorkCount,
			}
		}
	} else {
		for i := 0; i < length; i++ {
			vedioId = vedioList[i].ID
			vedioRequest[i].IsFavorite = MLikeDaoImpl.IsFavorite(UID, vedioId)
			vedioRequest[i].Title = vedioList[i].Title
			vedioRequest[i].ID = vedioList[i].ID
			vedioRequest[i].CommentCount = vedioList[i].CommentCount
			vedioRequest[i].FavoriteCount = vedioList[i].FavoriteCount
			vedioRequest[i].PlayURL = vedioList[i].PlayURL
			vedioRequest[i].CoverURL = vedioList[i].CoverURL
			user, err := MUserDaoImpl.GetUserById(vedioList[i].AuthorID)
			if err != nil {
				vedioRequest[i].Author = RequestEntity.AuthorUser{}
			}
			isFollow, err := MUserDaoImpl.GetIsFollowByUserIdAndToUserId(UID, user.ID)
			if err != nil {
				isFollow = false
			}
			vedioRequest[i].Author = RequestEntity.AuthorUser{
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				FavoriteCount:   user.FavoriteCount,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				ID:              user.ID,
				IsFollow:        isFollow,
				Name:            user.Name,
				Signature:       user.Signature,
				TotalFavorited:  strconv.Itoa(int(user.TotalFavorited)),
				WorkCount:       user.WorkCount,
			}
		}
	}
	return vedioRequest
}
