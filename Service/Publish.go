package Service

import (
	"douyin/Dao/MysqlDao"
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

var MPublishDaoImpl MysqlDao.PublishDao = &MysqlDao.PublishDaoImpl{}

func PublishProcess(request RequestEntity.PublishRequest,
	fileHeader *multipart.FileHeader) RequestEntity.PublishBack {
	var (
		publishBack RequestEntity.PublishBack
		fileType    string
		title       string
	)
	successMsg := "Vedio Publish Success!"
	failMsg := "Vedio Publish Fail"
	// 获取投稿视频的文件类型和文件名
	fileType = "mp4"
	// 获取投稿者ID
	token := request.Token
	UID, err := Util.ParserToken(token)
	fmt.Println(UID)
	if err != nil {
		Log.ErrorLogWithoutPanic("Parse Error!", err)
		publishBack.StatusCode = 1
		publishBack.StatusMsg = &failMsg
		return publishBack
	}

	// 获取标题描述
	title = request.Title

	// 获取用户信息描述
	user, err := MUserDaoImpl.GetUserById(UID)
	if err != nil {
		Log.ErrorLogWithoutPanic("User Info Read Error!", err)
		publishBack.StatusCode = 1
		publishBack.StatusMsg = &failMsg
		return publishBack
	}
	userName := user.Name

	// 保存文件
	vedioName := Util.VedioName(userName, fileType)
	file, err := os.Create("./Resource/Vedio/" + vedioName)
	if err != nil {
		Log.ErrorLogWithoutPanic("File Create Error!", err)
		publishBack.StatusCode = 1
		publishBack.StatusMsg = &failMsg
		return publishBack
	}
	defer file.Close()
	_, err = io.Copy(file, request.Data)
	if err != nil {
		Log.ErrorLogWithoutPanic("File Transfer Error!", err)
		publishBack.StatusCode = 1
		publishBack.StatusMsg = &failMsg
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

	// 返回成功信息
	publishBack.StatusCode = 0
	publishBack.StatusMsg = &successMsg

	return publishBack
}

func PublishListProcess(request RequestEntity.UserInfoRequest) RequestEntity.PublishListBack {
	var (
		publishListBack RequestEntity.PublishListBack
	)
	return publishListBack
}
