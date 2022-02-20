package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tools/apps/wechat/configs"
	"tools/apps/wechat/internal/web/service"
)

type HttpServer struct {
	server http.Server
}
func (h *HttpServer) Start() error {
	err := h.server.ListenAndServe()
	log.Println("http server stoped")
	return err
}
func (h *HttpServer) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("http server stopping")
	return h.server.Shutdown(ctx)
}

func NewHttpServer(env *string, c *configs.Config, svc *service.WechatHttpService) *HttpServer {

	if *env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := svc.InitEngine()

	return &HttpServer{
		server: http.Server{
			Addr:    ":" + c.HttpPort,
			Handler: e,
		},
	}
}
