package data

import (
	"github.com/google/wire"
	"tools/apps/wechat/internal/rpc/biz"
	"tools/internal/pkg/wxapi"
)

var Constructor = wire.NewSet(
	NewData,
	NewUserRepo,
	wire.Bind(new(biz.UserRepoInterface), new(*UserRepo)),
)

type Data struct {
	client *wxapi.WxApiClient
}
func NewData(wxApiClient *wxapi.WxApiClient) *Data {
	return &Data{
		client: wxApiClient,
	}
}