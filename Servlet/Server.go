package Servlet

import "C"
import (
	"context"
	"douyin/Controller"
	"douyin/Dao/MysqlDao"
	"douyin/Dao/RedisDao"
	"douyin/Entity/RequestEntity"
	"douyin/Log"
	"douyin/Util"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NetServerStartup 启动前端网络服务
func NetServerStartup() {
	r := gin.Default()

	//基础功能
	//视频Feed流，Query类
	r.GET("/douyin/feed/", func(c *gin.Context) {
		var (
			request    RequestEntity.FeedRequest
			latestTime string
			token      string
		)

		latestTime = c.Query("latest_time")
		token = c.Query("token")
		request.LatestTime = &latestTime
		request.Token = &token
		c.JSON(http.StatusOK, Controller.VedioFeedImp(request))
	})

	//注册功能，Query类型
	r.POST("/douyin/user/register/", func(c *gin.Context) {
		var request RequestEntity.RegisterRequest
		request.Username = c.Query("username")
		request.Password = c.Query("password")
		c.JSON(http.StatusOK, Controller.UserRegisterImp(request))
	})

	//用户登录，Query类型
	r.POST("/douyin/user/login/", func(c *gin.Context) {
		var request RequestEntity.LoginRequest
		request.Username = c.Query("username")
		request.Password = c.Query("password")
		c.JSON(http.StatusOK, Controller.UserLoginImp(request))
	})

	//获取用户信息，Query类型
	r.GET("/douyin/user/", func(c *gin.Context) {
		var request RequestEntity.UserInfoRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.UserInfoImp(request))
	})

	//投稿，Body类型
	r.POST("/douyin/publish/action/", func(c *gin.Context) {
		var (
			request    RequestEntity.PublishRequest
			data       multipart.File
			fileHeader *multipart.FileHeader
		)
		data, fileHeader, err := c.Request.FormFile("data")
		if err != nil {
			Log.ErrorLogWithoutPanic("Read File Error!", err)
		}
		Log.NormalLog("Read File Success!", err)
		request.Token = c.PostForm("token")
		request.Title = c.PostForm("title")
		request.Data = data
		c.JSON(http.StatusOK, Controller.PublishImp(request, fileHeader))
	})

	//列出用户投稿后的视频，Query类型
	r.GET("/douyin/publish/list/", func(c *gin.Context) {
		var request RequestEntity.UserInfoRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.PublishListImp(request))
	})

	//点赞功能，Query类型
	r.POST("/douyin/favorite/action/", func(c *gin.Context) {
		var request RequestEntity.LikeRequest
		request.ActionType = c.Query("action_type")
		request.Token = c.Query("token")
		request.VideoID = c.Query("video_id")
		c.JSON(http.StatusOK, Controller.LikeImp(request))
	})

	//列出喜欢列表，Query类型
	r.GET("/douyin/favorite/list/", func(c *gin.Context) {
		var request RequestEntity.LikeListRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.LikeListImp(request))
	})

	//评论功能，Query类型
	r.POST("/douyin/comment/action/", func(c *gin.Context) {
		var (
			request     RequestEntity.CommentRequest
			CommentID   string
			CommentText string
		)
		request.ActionType = c.Query("action_type")
		request.Token = c.Query("token")
		request.VideoID = c.Query("video_id")
		CommentID = c.Query("comment_id")
		CommentText = c.Query("comment_text")
		request.CommentID = &CommentID
		request.CommentText = &CommentText
		c.JSON(http.StatusOK, Controller.CommentImp(request))
	})

	//评论列表，Query类型
	r.GET("/douyin/comment/list/", func(c *gin.Context) {
		var request RequestEntity.CommentListRequest
		request.Token = c.Query("token")
		request.VideoID = c.Query("video_id")
		c.JSON(http.StatusOK, Controller.CommentListImp(request))
	})

	//关注操作，Query类型
	r.POST("/douyin/relation/action/", func(c *gin.Context) {
		var request RequestEntity.FocusRequest
		request.Token = c.Query("token")
		request.ToUserID = c.Query("to_user_id")
		request.ActionType = c.Query("action_type")
		c.JSON(http.StatusOK, Controller.FocusImp(request))
	})

	//关注列表，Query类型
	r.GET("/douyin/relation/follow/list/", func(c *gin.Context) {
		var request RequestEntity.UserInfoRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.FocusListImp(request))
	})

	//获取粉丝列表，Query类型
	r.GET("/douyin/relation/follower/list/", func(c *gin.Context) {
		var request RequestEntity.UserInfoRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.FanListImp(request))
	})

	//获取好友列表，Query类型
	r.GET("/douyin/relation/friend/list/", func(c *gin.Context) {
		var request RequestEntity.UserInfoRequest
		request.Token = c.Query("token")
		request.UserID = c.Query("user_id")
		c.JSON(http.StatusOK, Controller.FriendListImp(request))
	})

	//发送消息，Query类型
	r.POST("/douyin/message/action/", func(c *gin.Context) {
		var request RequestEntity.MessageSendRequest
		request.Token = c.Query("token")
		request.Content = c.Query("content")
		request.ToUserID = c.Query("to_user_id")
		request.ActionType = c.Query("action_type")
		c.JSON(http.StatusOK, Controller.MessageSendImp(request))
	})

	//调取聊天记录，Query类型
	r.GET("/douyin/message/chat/", func(c *gin.Context) {
		var request RequestEntity.MessageRequest
		request.Token = c.Query("token")
		request.ToUserID = c.Query("to_user_id")
		request.PreMsgTime = c.Query("pre_msg_time")
		c.JSON(http.StatusOK, Controller.MessageImp(request))
	})

	//加载静态资源，实现路径到URL的映射
	r.Static("/douyin/image", "Resource/Image")
	r.Static("/douyin/text", "Resource/text")
	r.Static("/douyin/vedio", "Resource/Vedio")

	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err1 := s.ListenAndServe(); err1 != nil && err1 != http.ErrServerClosed {
			Log.ErrorLogWithoutPanic("Listen: %s\n", err1)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	Log.NormalLog("Shutdown Server ...", nil)

	//创建超时上下文，Shutdown可以让未处理的连接在这个时间内关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//停止HTTP服务器
	if err := s.Shutdown(ctx); err != nil {
		Log.ErrorLogWithoutPanic("Server Shutdown:", err)
	}

	Log.NormalLog("Server Exiting", nil)

}

func DataBaseInit() (*gorm.DB, *redis.Client) {
	mdb := Util.InitMysqlDb()
	rdb := Util.InitRedisDb()
	Util.TableCreate(mdb)
	return mdb, rdb
}

// DbServerStartup 启动数据库服务，初始化数据库
func DbServerStartup() {
	mdb, rdb := DataBaseInit()
	MysqlDao.GetMysqlDBHandler(mdb)
	RedisDao.GetRedisDBHandler(rdb)
}
