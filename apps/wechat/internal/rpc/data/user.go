package data

type UserRepo struct {
	data *Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		data: data,
	}
}

func (r *UserRepo) OpenidList(openid string) ([]string, error) {
	return r.data.client.OpenidList(openid)
}

func (r *UserRepo) UserInfo(openid string) (string, string, error) {
	return r.data.client.UserInfo(openid)
}