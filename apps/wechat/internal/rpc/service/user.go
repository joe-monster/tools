package service

import (
	"context"
	v1 "tools/apps/wechat/api/rpc/v1"
)

func (s *WechatRpcService) RpcGetUserList(ctx context.Context, in *v1.GetUserListRequest) (*v1.GetUserListResponse, error) {

	rtn, err := s.userBiz.GetUserList()
	if err != nil {
		s.logClient.Error(err)
		return nil, err
	}

	//重新整理返回值
	resp := new(v1.GetUserListResponse)
	for _,v := range rtn {
		resp.Users = append(resp.Users, &v1.GetUserListResponse_User{
			Nickname: v.Nickname,
			Avatar: v.Avatar,
			Openid: v.Openid,
		})
	}

	return resp, nil
}



