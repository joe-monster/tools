package server

import (
	"google.golang.org/grpc"
	"log"
	"net"
	v1 "tools/apps/wechat/api/rpc/v1"
	"tools/apps/wechat/configs"
	"tools/apps/wechat/internal/rpc/service"
)

type RpcServer struct {
	server *grpc.Server
	port   string
}
func (h *RpcServer) Start() error {
	listen, err := net.Listen("tcp", ":"+h.port)
	if err != nil {
		return err
	}
	err = h.server.Serve(listen)
	log.Println("rpc server stoped")
	return err
}
func (h *RpcServer) Stop() error {
	h.server.Stop()
	log.Println("rpc server stopping")
	return nil
}

func NewRpcServer(env *string, c *configs.Config, svc *service.WechatRpcService) *RpcServer {
	log.Println(*env)

	s := grpc.NewServer()

	v1.RegisterUserServer(s, svc)

	return &RpcServer{
		server: s,
		port: c.RpcPort,
	}
}





//// interceptor 拦截器
//func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//
//	//请求安全验证
//	if err := authMiddleware(ctx);err != nil {
//		return nil, err
//	}
//
//
//	// 继续处理请求
//	return handler(ctx, req)
//}
//func authMiddleware(ctx context.Context) error {
//
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return errors.New("metadata信息不存在")
//	}
//
//	val, ok := md["authorization"]
//	if !ok || len(val) == 0 {
//		return errors.New("认证token不存在或为空")
//	}
//
//	tk := val[0]
//	if tk == "" {
//		return errors.New("认证token为空")
//	}
//
//
//	return nil
//}