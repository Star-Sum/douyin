package Service

import (
	"context"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

var (
	MPublishDaoImpl MysqlDao.PublishDao = &MysqlDao.PublishDaoImpl{}
	RPublishDaoImpl RedisDao.PublishDao = &RedisDao.PublishDaoImpl{
		Ctx: context.Background(),
	}
)

func PublishProcess(request RequestEntity.PublishRequest,
	fileHeader *multipart.FileHeader) RequestEntity.PublishBack {
	var (
		publishBack RequestEntity.PublishBack
		fileType    string
		title       string
	)
	successMsg := "Vedio Publish Success!"
	failMsg := "Vedio Publish Fail"
	publishBack.StatusCode = 1
	publishBack.StatusMsg = &failMsg
	// 获取投稿视频的文件类型和文件名
	fileType = "mp4"
	// 获取投稿者ID
	token := request.Token
	UID, err := Util.ParserToken(token)
	fmt.Println(UID)
	if err != nil {
		Log.ErrorLogWithoutPanic("Parse Error!", err)
		return publishBack
	}

	// 获取标题描述
	title = request.Title

	// 获取用户信息描述
	user, err := MUserDaoImpl.GetUserById(UID)
	if err != nil {
		Log.ErrorLogWithoutPanic("User Info Read Error!", err)
		return publishBack
	}
	userName := user.Name

	// 保存文件
	vedioName := Util.VedioName(userName, fileType)
	file, err := os.Create("./Resource/Vedio/" + vedioName)
	if err != nil {
		Log.ErrorLogWithoutPanic("File Create Error!", err)
		return publishBack
	}
	defer file.Close()
	_, err = io.Copy(file, request.Data)
	if err != nil {
		Log.ErrorLogWithoutPanic("File Transfer Error!", err)
		return publishBack
	}

	// 为视频制作截图
	coverName := Util.GetVedioCover(vedioName, userName)
	// 将视频信息写入数据库
	infoOne := TableEntity.PublishInfo{}
	infoTwo := TableEntity.VedioInfo{}

	infoOne.UseID = UID
	infoOne.Title = title
	infoOne.Time = time.Now()
	infoOne.OpID = 0
	infoOne.UrlRoot = Util.DirToUrl("vedio", vedioName)

	infoTwo.AuthorID = UID
	infoTwo.Title = title
	infoTwo.Timestamp = time.Now()
	infoTwo.CoverURL = Util.DirToUrl("image", coverName)
	infoTwo.PlayURL = Util.DirToUrl("vedio", vedioName)
	infoTwo.FavoriteCount = 0
	infoTwo.CommentCount = 0

	MPublishDaoImpl.PublishVedio(infoOne)
	MPublishDaoImpl.VedioInfoInsert(infoTwo)

	info := make([]TableEntity.VedioInfo, 1)
	info[0] = infoTwo
	err = RPublishDaoImpl.InsertVedioInfo(UID, info)

	if err != nil {
		Log.ErrorLogWithoutPanic("Insert Info Error!", err)
		return publishBack
	}

	// 返回成功信息
	publishBack.StatusCode = 0
	publishBack.StatusMsg = &successMsg

	return publishBack
}

func PublishListProcess(request RequestEntity.UserInfoRequest) RequestEntity.PublishListBack {
	var (
		publishListBack RequestEntity.PublishListBack
		vedioInfo       []TableEntity.VedioInfo
		vRequest        []RequestEntity.VedioRequest
		vRequestInfo    RequestEntity.VedioRequest
	)
	successMsg := "Vedio Publish Success!"
	failMsg := "Vedio Publish Fail"
	publishListBack.StatusCode = 1
	publishListBack.StatusMsg = &failMsg
	// 首先解析token
	uid, err := strconv.ParseInt(request.UserID, 10, 64)
	if err != nil {
		Log.ErrorLogWithoutPanic("UID Transform Error!", err)
		return publishListBack
	}
	judge, err := Util.TokenJudge(request.Token, uid)
	if err != nil {
		Log.ErrorLogWithoutPanic("UID Parse Error!", err)
		return publishListBack
	}
	if !judge {
		Log.ErrorLogWithoutPanic("UID Judge Error!", err)
		return publishListBack
	}

	// 拉取相关用户的投稿视频信息
	vedioInfo = RPublishDaoImpl.GetVedioByAuthorId(uid)
	if len(vedioInfo) < 30 {
		vInfo := MPublishDaoImpl.FindVedioByAuthorUid(uid)
		err = RPublishDaoImpl.InsertVedioInfo(uid, vInfo)
		if err != nil {
			Log.ErrorLogWithoutPanic("Vedio Info Insert Error!", err)
		}
		if len(vInfo) <= 0 {
			publishListBack.StatusCode = 1
			publishListBack.StatusMsg = &failMsg
		} else {
			publishListBack.StatusCode = 0
			publishListBack.StatusMsg = &failMsg
		}
		for i := 0; i < len(vInfo); i++ {
			vRequestInfo.ID = vInfo[i].ID
			vRequestInfo.FavoriteCount = vInfo[i].FavoriteCount
			vRequestInfo.Title = vInfo[i].Title
			vRequestInfo.PlayURL = vInfo[i].PlayURL
			vRequestInfo.CoverURL = vInfo[i].CoverURL
			vRequestInfo.CommentCount = vInfo[i].CommentCount
			vRequestInfo.Author = GetAuthorInfo(vInfo[i].AuthorID)
			vRequestInfo.IsFavorite = MLikeDaoImpl.IsFavorite(uid, vInfo[i].ID)
			vRequest = append(vRequest, vRequestInfo)
		}
	}
	publishListBack.VideoList = vRequest
	publishListBack.StatusCode = 0
	publishListBack.StatusMsg = &successMsg
	return publishListBack
}

func GetAuthorInfo(UID int64) RequestEntity.AuthorUser {
	user := RUserDaoImpl.GetUserInfo(UID)
	var err error
	if user.ID == -1 {
		user, err = MUserDaoImpl.GetUserById(UID)
		if err != nil {
			Log.ErrorLogWithoutPanic("No such Author, Please Retry", err)
			return RequestEntity.AuthorUser{ID: -1}
		}
	}
	var aUser RequestEntity.AuthorUser
	aUser.ID = user.ID
	aUser.Name = user.Name
	aUser.FavoriteCount = user.FavoriteCount
	aUser.WorkCount = user.WorkCount
	aUser.TotalFavorited = strconv.FormatInt(user.TotalFavorited, 10)
	aUser.Signature = user.Signature
	aUser.IsFollow = MRelationDaoImpl.FindRelation(UID, user.ID, "1")
	return aUser
}
