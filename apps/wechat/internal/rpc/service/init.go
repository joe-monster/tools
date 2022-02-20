package service

import (
	"github.com/google/wire"
	"tools/apps/wechat/internal/rpc/biz"
	"tools/internal/pkg/logger"
)

var Constructor = wire.NewSet(NewWechatRpcService)


type WechatRpcService struct {
	logClient *logger.LogClient

	userBiz *biz.UserBiz
}

func NewWechatRpcService(
	logClient *logger.LogClient,

	user *biz.UserBiz,
) *WechatRpcService {

	return &WechatRpcService{
		logClient: logClient,

		userBiz: user,
	}

}