package wxapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
	"tools/apps/wechat/configs"
	"tools/internal/pkg/logger"
	"tools/internal/pkg/network"
)

type WxApiClient struct {
	logClient *logger.LogClient
	appid string
	secret string
}
func NewWxApiClient(c *configs.Config, logClient *logger.LogClient) *WxApiClient {
	return &WxApiClient{
		logClient: logClient,
		appid: c.Wechat.Appid,
		secret: c.Wechat.Secret,
	}
}





var accessToken string
var expireTime int64
func (c *WxApiClient) getToken() (string, error) {

	if expireTime - time.Now().Unix() <= 600 {
		tk, expireIn, err := getAccessToken(c.appid, c.secret)
		if err != nil {
			//记录一下失败日志
			c.logClient.Error(err)
			return "", err
		}

		accessToken = tk
		expireTime = time.Now().Unix() + expireIn

		//记录一下成功日志
		c.logClient.Info(tk, expireIn)
	}

	return accessToken, nil
}
func getAccessToken(appid, secret string) (string, int64, error) {

	data, err := network.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret, nil)
	if err != nil {
		return "", 0, errors.WithMessage(err, fmt.Sprintf("appid:%s,secret:%s", appid, secret))
	}

	var bodyJson struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64    `json:"expires_in"`
	}
	if err := json.Unmarshal(data, &bodyJson); err != nil {
		return "", 0, errors.Wrap(err, fmt.Sprintf("data:%s", string(data)))
	}

	if bodyJson.Errcode != 0 {
		return "", 0, errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.AccessToken, bodyJson.ExpiresIn, nil
}