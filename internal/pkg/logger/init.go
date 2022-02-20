package logger

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

//配置相关
type configs struct {
	dir string
	suffix string
	cutTime time.Duration
	fileSaveTime time.Duration
}
type dialOptionSetter func(options *configs)

func SetDir(dir string) dialOptionSetter {
	return func(config *configs) {
		config.dir = dir
	}
}
func SetSuffix(suffix string) dialOptionSetter {
	return func(config *configs) {
		config.suffix = suffix
	}
}
func SetCutTime(t time.Duration) dialOptionSetter {
	return func(config *configs) {
		config.cutTime = t
	}
}
func SetFileSaveTime(t time.Duration) dialOptionSetter {
	return func(config *configs) {
		config.fileSaveTime = t
	}
}

//日志客户端封装
type LogClient struct {
	appName string
	client *logrus.Logger
}
func NewLogClient(appName string, options ...dialOptionSetter) *LogClient {

	//设置配置
	config := configs {
		dir: "./log/",	//默认存储目录
		suffix: ".log.%Y%m%d",	//日志后缀
		cutTime: 24 * time.Hour,	//默认一天一个文件
		fileSaveTime: 30 * 24 * time.Hour,	//默认文件保存存30天
	}
	for _, optionFunc := range options {
		optionFunc(&config)
	}

	//创建logrus客户端
	var client = logrus.New()

	// 目录不存在则创建
	if _, err := os.Stat(config.dir); os.IsNotExist(err) {
		os.MkdirAll(config.dir, os.ModePerm)
	}

	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	client.Out = src

	client.SetReportCaller(true)

	var logWriterfunc = func (t string) (*rotatelogs.RotateLogs, error) {
		return rotatelogs.New(
			path.Join(config.dir, t + config.suffix),
			rotatelogs.WithMaxAge(config.fileSaveTime),   // 文件最大保存时间
			rotatelogs.WithRotationTime(config.cutTime), // 日志切割时间间隔
		)
	}

	infoLogWriter, err := logWriterfunc("info")
	if err != nil {
		panic(err)
	}

	errorLogWriter, err := logWriterfunc("error")
	if err != nil {
		panic(err)
	}

	fatalLogWriter, err := logWriterfunc("fatal")
	if err != nil {
		panic(err)
	}

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  infoLogWriter,
		logrus.ErrorLevel: errorLogWriter,
		logrus.FatalLevel: fatalLogWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{})

	client.AddHook(lfHook)

	return &LogClient{
		appName: appName,
		client: client,
	}
}

func (c *LogClient) Info(args ...interface{}) {
	c.client.Infof("app:%s info:%+v \n", c.appName, args)
}

func (c *LogClient) Error(args ...interface{}) {
	c.client.Errorf("app:%s error:%+v \n", c.appName, args)
}

func (c *LogClient) Fatal(args ...interface{}) {
	c.client.Fatalf("app:%s fatal:%+v \n", c.appName, args)
}


