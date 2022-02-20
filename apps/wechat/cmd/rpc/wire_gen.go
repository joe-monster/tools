// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"tools/apps/wechat/configs"
	"tools/apps/wechat/internal/rpc/biz"
	"tools/apps/wechat/internal/rpc/data"
	"tools/apps/wechat/internal/rpc/server"
	"tools/apps/wechat/internal/rpc/service"
	"tools/internal/pkg/app"
	"tools/internal/pkg/logger"
	"tools/internal/pkg/wxapi"
)

// Injectors from wire.go:

func initApp(env *string, c *configs.Config, logClient *logger.LogClient) (*app.App, error) {
	wxApiClient := wxapi.NewWxApiClient(c, logClient)
	dataData := data.NewData(wxApiClient)
	userRepo := data.NewUserRepo(dataData)
	userBiz := biz.NewUserBiz(logClient, userRepo)
	wechatRpcService := service.NewWechatRpcService(logClient, userBiz)
	rpcServer := server.NewRpcServer(env, c, wechatRpcService)
	appApp := newApp(rpcServer)
	return appApp, nil
}