package wxapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"tools/internal/pkg/network"
)

func (c *WxApiClient) OpenidList(openid string) ([]string, error) {

	tk, err := c.getToken()
	if err != nil {
		return nil, err
	}

	apiData, err := network.Get("https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + tk + "&next_openid=" + openid, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("next_openid:%s", openid))
	}

	var bodyJson struct {
		Errcode     int		`json:"errcode"`
		Errmsg      string	`json:"errmsg"`

		Data 		struct{
			Openid []string `json:"openid"`
		}		`json:"data"`
		NextOpenid string `json:"next_openid"`
		Total int `json:"total"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(apiData, &bodyJson); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("apiData:%s", string(apiData)))
	}

	if bodyJson.Errcode != 0 {
		return nil, errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.Data.Openid, nil
}

func (c *WxApiClient) UserInfo(openid string) (string, string, error) {

	tk, err := c.getToken()
	if err != nil {
		return "", "", err
	}

	apiData, err := network.Get("https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + tk + "&openid=" + openid, nil)
	if err != nil {
		return "", "", errors.WithMessage(err, fmt.Sprintf("openid:%s", openid))
	}

	var bodyJson struct {
		Errcode     int		`json:"errcode"`
		Errmsg      string	`json:"errmsg"`

		Nickname 	string	`json:"nickname"`
		Headimgurl	string	`json:"headimgurl"`
	}
	if err := json.Unmarshal(apiData, &bodyJson); err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("apiData:%s", string(apiData)))
	}

	if bodyJson.Errcode != 0 {
		return "", "", errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.Nickname, bodyJson.Headimgurl, nil
}


type getTagListRtnRowType struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Count int `json:"count"`
}
func (c *WxApiClient) GetUserTagList() ([]*getTagListRtnRowType, error) {

	tk, err := c.getToken()
	if err != nil {
		return nil, err
	}

	apiData, err := network.Get("https://api.weixin.qq.com/cgi-bin/tags/get?access_token=" + tk, nil)
	if err != nil {
			return nil, err
	}

	var bodyJson struct {
		Errcode     int		`json:"errcode"`
		Errmsg      string	`json:"errmsg"`

		Tags 		[]*getTagListRtnRowType `json:"tags"`
	}
	if err := json.Unmarshal(apiData, &bodyJson); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("apiData:%s", string(apiData)))
	}

	if bodyJson.Errcode != 0 {
		return nil, errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.Tags, nil
}