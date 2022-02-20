package main

import (
	"flag"
	"time"
	"tools/apps/wechat/configs"
	"tools/apps/wechat/internal/web/server"
	"tools/internal/pkg/app"
	"tools/internal/pkg/logger"
)

var (
	Name = "tools.wechat.web"
	Env = flag.String("env", "", "...")
)

func newApp(s *server.HttpServer) *app.App {
	return app.New(s)
}

func main() {
	flag.Parse()

	//加载配置文件
	c, err := configs.MakeFromYaml("../../configs/config.yaml")
	if err != nil {
		panic(err)
	}

	//初始化log，用毛师傅教的支持默认值以及按需配置参数的方式，相当舒坦
	logClient := logger.NewLogClient(
		Name,
		logger.SetDir(c.Log.Dir),
		logger.SetSuffix(c.Log.Suffix),
		logger.SetCutTime(time.Hour*24*time.Duration(c.Log.CutDays)),
		logger.SetFileSaveTime(time.Hour*24*time.Duration(c.Log.FileSaveDays)),
	)

	//初始化应用，自动依赖注入
	app, err := initApp(Env, c, logClient)
	if err != nil {
		panic(err)
	}

	//启动相关应用服务
	if err := app.Run(); err != nil {
		panic(err)
	}
}
