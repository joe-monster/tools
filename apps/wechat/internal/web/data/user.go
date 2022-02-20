package data

import (
	"tools/apps/wechat/internal/web/biz"
)

type UserRepo struct {
	data *Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		data: data,
	}
}

func (r *UserRepo) GetTagList() ([]*biz.UserTagData, error){

	list, err := r.data.client.GetUserTagList()
	if err != nil {
		return nil, err
	}

	rtn := []*biz.UserTagData{}

	for _,v := range list {
		rtn = append(rtn, &biz.UserTagData{
			Id: v.Id,
			Name: v.Name,
			Count: v.Count,
		})
	}

	return rtn, nil
}

func (r *UserRepo) OpenidList(openid string) ([]string, error) {
	return r.data.client.OpenidList(openid)
}

func (r *UserRepo) UserInfo(openid string) (string, string, error) {
	return r.data.client.UserInfo(openid)
}