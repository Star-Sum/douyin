package Service

// import "douyin/Entity/RequestEntity"
import (
	"douyin/Entity/RequestEntity"
	"douyin/Entity/TableEntity"
	"douyin/Log"
	"douyin/Util"
	"strconv"
)

// 进行点赞或者取消点赞的过程
func LikeProcess(request RequestEntity.LikeRequest) RequestEntity.LikeBack {
	var (
		likeBack RequestEntity.LikeBack
	)
	// 判断token是否合法，并提取fromUserID
	fromUserID, err := Util.ParserToken(request.Token)
	if err != nil {
		likeBack.StatusCode = 1
		likeBack.StatusMsg = "Token Parse Error!"
		Log.ErrorLogWithoutPanic("Token Parse Error!", err)
		return likeBack
	}

	// 传入的ActionType参数是string，先解析为int64——如果actiontype不是1（点赞）或者2（取消点赞），则
	if (request.ActionType != "1") && (request.ActionType != "2") {
		likeBack.StatusCode = 1
		likeBack.StatusMsg = "Invalid action_type Param!"
		Log.ErrorLogWithoutPanic("Invalid action_type Param!", nil)
		return likeBack
	}

	// 传入的VideoID参数是string，先解析为int64
	likeVideoIDString := request.VideoID
	likeVideoID, err := strconv.ParseInt(likeVideoIDString, 10, 64)
	if err != nil {
		likeBack.StatusCode = 1
		likeBack.StatusMsg = "like_video_id Parse Error!"
		Log.ErrorLogWithoutPanic("like_video_id Parse Error!", err)
		return likeBack
	}

	// 判断数据库中是否存在likeVideoID的视频
	videoInfo := MFeedDaoImpl.FindVedioById(likeVideoID)
	videoId := videoInfo.ID
	if videoId == 0 {
		likeBack.StatusCode = 1
		likeBack.StatusMsg = "Video with the specified ID not found!"
		// 添加成功log，todo
		Log.ErrorLogWithoutPanic("Video with the specified ID not found!", nil)
		return likeBack
	}

	// 处理点赞/取消点赞数据库
	// 如果请求为点赞
	if request.ActionType == "1" {

		// 如果已经点赞过了
		if MLikeDaoImpl.IsFavorite(fromUserID, videoId) {
			likeBack.StatusCode = 1
			likeBack.StatusMsg = "request Error! can't like a video already liked before"
			return likeBack
		}

		// 数据库添加点赞信息
		err := MLikeDaoImpl.Add(fromUserID, videoId)

		// 如果数据库返回错误
		if err != nil {
			likeBack.StatusCode = 1
			likeBack.StatusMsg = "Database add like Error!"
			Log.ErrorLogWithoutPanic("Database add like Error!", err)
			return likeBack
		}

		likeBack.StatusCode = 0
		likeBack.StatusMsg = "Database add like success"

	} else if request.ActionType == "2" {

		// 如果还没有点赞过
		if !MLikeDaoImpl.IsFavorite(fromUserID, videoId) {
			likeBack.StatusCode = 1
			likeBack.StatusMsg = "request Error! can't dislike a video didn't like before"
			return likeBack
		}

		// 数据库删除点赞信息
		err := MLikeDaoImpl.Sub(fromUserID, videoId)

		// 如果数据库返回错误
		if err != nil {
			likeBack.StatusCode = 1
			likeBack.StatusMsg = "Database add like Error!"
			Log.ErrorLogWithoutPanic("Database add like Error!", err)
			return likeBack
		}

		likeBack.StatusCode = 0
		likeBack.StatusMsg = "Database add like success"

	}
	return likeBack

}

// LikeListProcess 请求点赞列表
func LikeListProcess(request RequestEntity.LikeListRequest) []RequestEntity.VedioRequest {
	var (
		likeListBack      []RequestEntity.VedioRequest
		likeVideoInfoList []TableEntity.VedioInfo
	)

	// 判断token是否合法，并提取fromUserID
	UserID, err := Util.ParserToken(request.Token)
	if err != nil {
		Log.ErrorLogWithoutPanic("Token Parse Error!", err)
		return likeListBack
	}

	// 根据id获取点赞过的视频

	likeVideoList, err := MLikeDaoImpl.FavoriteListByUserID(UserID)
	if err != nil {
		Log.ErrorLogWithoutPanic("database looking for user liked video err!", err)
		return likeListBack
	}
	for i := 0; i < len(likeVideoList); i++ {
		videoId := likeVideoList[i]
		videoInfo := MFeedDaoImpl.FindVedioById(videoId)
		likeVideoInfoList = append(likeVideoInfoList, videoInfo)
	}

	likeListBack = VedioListToRequest(likeVideoInfoList, -2)

	return likeListBack
}
