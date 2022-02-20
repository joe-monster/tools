package wxapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"tools/internal/pkg/network"
)

//整理数据发通知
func (c *WxApiClient) SendPublishIdTemplateMessage(openid, templateId, url string, data map[string]string) (int, error) {

	type pDataValueType struct{
		Value string `json:"value"`
	}
	type pDataType map[string]pDataValueType
	type paramType struct {
		Touser string `json:"touser"`
		TemplateID string `json:"template_id"`
		URL string `json:"url"`
		Miniprogram struct {
			Appid string `json:"appid"`
			Pagepath string `json:"pagepath"`
		} `json:"miniprogram"`
		Data pDataType `json:"data"`
	}

	var p paramType
	p.Touser = openid
	p.TemplateID = templateId
	p.URL = url
	p.Data = make(pDataType)
	for k,v := range data {
		p.Data[k] = pDataValueType{
			Value: v,
		}
	}

	pJson, err := json.Marshal(p)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("p:%v", p))
	}

	tk, err := c.getToken()
	if err != nil {
		return 0, err
	}

	apiData, err := network.PostJson("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + tk, nil, pJson)
	if err != nil {
		return 0, errors.WithMessage(err, fmt.Sprintf("pJson:%s", string(pJson)))
	}

	var bodyJson struct {
		Errcode     int		`json:"errcode"`
		Errmsg      string	`json:"errmsg"`
		Msgid 		int		`json:"msgid"`
	}
	if err := json.Unmarshal(apiData, &bodyJson); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("apiData:%s", string(apiData)))
	}

	if bodyJson.Errcode != 0 {
		return 0, errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.Msgid, nil
}