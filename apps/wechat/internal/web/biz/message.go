package biz

import (
	"math"
	"sync"
	"tools/internal/pkg/logger"
	"tools/internal/pkg/wxapi"
)

type MessageBiz struct {
	logClient *logger.LogClient
	wxApiClient *wxapi.WxApiClient
	userBiz *UserBiz
}
func NewMessageBiz(logClient *logger.LogClient, c *wxapi.WxApiClient, userBiz *UserBiz) *MessageBiz {
	return &MessageBiz{
		logClient: logClient,
		wxApiClient: c,
		userBiz: userBiz,
	}
}

func (b *MessageBiz) SendPublish(templateId, url string, openId string, data map[string]string) ([]string, error) {

	openidCh := make(chan string)
	ch := make(chan struct{})

	var failOpenids []string
	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		for {
			if openid,ok := <-openidCh;!ok {
				break
			} else {
				failOpenids = append(failOpenids, openid)
			}
		}
	}()

	//发通知
	if openId != "" {
		_, err := b.wxApiClient.SendPublishIdTemplateMessage(openId, templateId, url, data)
		if err != nil {
			b.logClient.Error(err)
			openidCh <- openId
		}
	} else {

		var totalOpenids [][]string
		var tmp []string

		//获取用户list
		n := 0
		var nextOpenid string
		for {
			openids, err := b.userBiz.userRepo.OpenidList(nextOpenid)
			if err != nil {
				return nil, err
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

		//开始发
		var wg sync.WaitGroup
		for _, openids := range totalOpenids {
			wg.Add(1)
			go func(openids []string) {
				defer func() {
					wg.Done()
				}()
				for _,openid := range openids {
					_, err := b.wxApiClient.SendPublishIdTemplateMessage(openid, templateId, url, data)
					if err != nil {
						b.logClient.Error(err)
						openidCh <- openid
					}
				}
			}(openids)
		}
		wg.Wait()

	}

	close(openidCh)

	<- ch

	return failOpenids, nil
}
