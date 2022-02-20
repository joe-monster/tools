package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
	"tools/apps/wechat/internal/web/biz"
	"tools/internal/pkg/logger"
	"tools/internal/pkg/resp"
)

var Constructor = wire.NewSet(NewWechatHttpService)


type WechatHttpService struct {
	logClient *logger.LogClient
	userBiz *biz.UserBiz
	messageBiz *biz.MessageBiz
	massBiz *biz.MassBiz
}

func NewWechatHttpService(
	logClient *logger.LogClient,
	user *biz.UserBiz,
	message *biz.MessageBiz,
	mass *biz.MassBiz,
) *WechatHttpService {
	return &WechatHttpService{
		logClient: logClient,

		userBiz: user,
		messageBiz: message,
		massBiz: mass,
	}
}

//初始化一个gin引擎，并设置好路由等
func (s *WechatHttpService) InitEngine() *gin.Engine {

	e := gin.Default()

	e.Use(
		gin.Recovery(),
	)

	e.POST("/user/list", s.Handler(s.UserList))
	e.POST("/user/tag/list", s.Handler(s.TagList))
	e.POST("/message/template/send/publishid", s.Handler(s.SendPublishIdType))
	e.POST("/message/mass/send", s.Handler(s.SendMass))

	return e
}

//定义各restful api服务方法签名
type handlerFunc func(c *gin.Context) (interface{}, error)

//统一请求执行器，方便请求前后做其他操作，代替gin中间件
func (s *WechatHttpService) Handler(handler handlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {

		//执行前做些事。。。



		data, err := handler(c)
		if err != nil {
			s.logClient.Error(err)
			c.JSON(http.StatusOK, resp.Fail(err))
		} else {
			c.JSON(http.StatusOK, resp.Success(data))
		}

		//执行后做些事。。。




	}
}
