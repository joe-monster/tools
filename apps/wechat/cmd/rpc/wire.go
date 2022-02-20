// +build wireinject

package main

import (
	"github.com/google/wire"
	"tools/apps/wechat/configs"
	"tools/apps/wechat/internal/rpc/biz"
	"tools/apps/wechat/internal/rpc/data"
	"tools/apps/wechat/internal/rpc/server"
	"tools/apps/wechat/internal/rpc/service"
	"tools/internal/pkg/app"
	"tools/internal/pkg/logger"
	"tools/internal/pkg/wxapi"
)

func initApp(env *string, c *configs.Config, logClient *logger.LogClient) (*app.App, error) {
	wire.Build(
		wxapi.NewWxApiClient,
		data.Constructor,
		biz.Constructor,
		service.Constructor,
		server.Constructor,
		newApp,
	)
	return &app.App{}, nil
}