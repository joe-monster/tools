package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"tools/apps/wechat/api/rpc/v1"
)

func main() {

	auth := new(customCredential)
	auth.Token = "qiaoyixuan"

	// 连接
	conn, err := grpc.Dial("127.0.0.1:9001", grpc.WithInsecure(), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := v1.NewUserClient(conn)

	// 调用方法
	req := &v1.GetUserListRequest{}
	res, err := c.RpcGetUserList(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Users)
}


// customCredential 自定义认证
type customCredential struct{
	Token string
}
// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization":  c.Token,
	}, nil
}
// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return false
}