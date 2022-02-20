package biz

import (
	"math"
	"sync"
	"tools/internal/pkg/logger"
)

type UserInfoData struct{
	Openid string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar string `json:"avatar"`
}

type UserRepoInterface interface {
	OpenidList(openid string) ([]string, error)
	UserInfo(openid string) (string, string, error)
}

//user领域要求该领域的repo具有UserRepo规定的能力
type UserBiz struct {
	logClient *logger.LogClient
	userRepo UserRepoInterface
}
func NewUserBiz(logClient *logger.LogClient, userRepo UserRepoInterface) *UserBiz {
	return &UserBiz{
		logClient: logClient,
		userRepo: userRepo,
	}
}

func (b *UserBiz) GetUserList() ([]*UserInfoData, error) {

	//开始获取用户信息
	var rtn []*UserInfoData

	var totalOpenids [][]string
	var tmp []string

	//获取用户list
	n := 0
	var nextOpenid string
	for {
		openids, err := b.userRepo.OpenidList(nextOpenid)
		if err != nil {
			return rtn, err
		}

		if len(openids) == 0 {
			break
		}

		//分组
		for _,openid := range openids {
			tmp = append(tmp, openid)
			n += 1

			m := int(math.Mod(float64(n), 100))
			if m == 0 {
				totalOpenids = append(totalOpenids, tmp)
				tmp = []string{}
			}
		}

		nextOpenid = openids[len(openids) - 1]
	}
	if len(tmp) != 0 {
		totalOpenids = append(totalOpenids, tmp)
		tmp = []string{}
	}

	//开始获取用户信息
	openidCh := make(chan *UserInfoData)
	ch := make(chan struct{})

	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		for {
			if user,ok := <-openidCh;!ok {
				break
			} else {
				rtn = append(rtn, user)
			}
		}
	}()

	//开始发
	var wg sync.WaitGroup
	for _, openids := range totalOpenids {
		wg.Add(1)
		go func(openids []string) {
			defer func() {
				wg.Done()
			}()
			for _,openid := range openids {
				nickname, avatar, err := b.userRepo.UserInfo(openid)
				if err != nil {
					b.logClient.Error(err)
				}
				openidCh <- &UserInfoData{
					Openid: openid,
					Nickname: nickname,
					Avatar: avatar,
				}
			}
		}(openids)
	}
	wg.Wait()

	close(openidCh)

	<- ch

	return rtn, nil
}


