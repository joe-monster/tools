package wxapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"tools/internal/pkg/network"
)

const (
	MSG_TYPE_MPNEWS = "mpnews"
	MSG_TYPE_TEXT = "text"
	MSG_TYPE_VOICE = "voice"
	MSG_TYPE_IMAGE = "image"
	MSG_TYPE_VIDEO = "video"
)

//内容参数类型
type MassParamsType struct {
	Touser []string `json:"touser"`

	Filter struct{
		IsToAll bool `json:"is_to_all"`
		TagId int `json:"tag_id"`
	} `json:"filter"`

	Msgtype string `json:"msgtype"`

	//群发接口-图文
	//Mpnews struct{
	//	MediaId string `json:"media_id"`
	//} `json:"mpnews"`
	//SendIgnoreReprint int `json:"send_ignore_reprint"`

	//群发接口-文本
	Text struct{
		Content string `json:"content"`
	} `json:"text"`

	//群发接口-语音
	//Voice struct{
	//	MediaId string `json:"media_id"`
	//} `json:"voice"`
	//
	////群发接口-图片
	//Images struct{
	//	MediaIds []string `json:"media_ids"`
	//	Recommend string `json:"recommend"`
	//	NeedOpenComment int `json:"need_open_comment"`
	//	OnlyFansCanComment int `json:"only_fans_can_comment"`
	//} `json:"images"`
	//
	////群发接口-视频
	//Mpvideo struct{
	//	MediaId string `json:"media_id"`
	//	Title string `json:"title"`
	//	Description string `json:"description"`
	//} `json:"mpvideo"`
}

func (c *WxApiClient) SendMassMessageByUser(p *MassParamsType) (i int, err error) {
	tk, err := c.getToken()
	if err != nil {
		return 0, err
	}

	url := "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=" + tk
	return sendMassMessage(url, p)
}

func (c *WxApiClient) SendMassMessageByTag(p *MassParamsType) (i int, err error) {
	tk, err := c.getToken()
	if err != nil {
		return 0, err
	}

	url := "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=" + tk
	return sendMassMessage(url, p)
}

func sendMassMessage(url string, p interface{}) (int, error) {

	pJson, err := json.Marshal(p)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("p:%v", p))
	}

	apiData, err := network.PostJson(url, nil, pJson)
	if err != nil {
		return 0, errors.WithMessage(err, fmt.Sprintf("pJson:%s", string(pJson)))
	}

	var bodyJson struct {
		Errcode     int		`json:"errcode"`
		Errmsg      string	`json:"errmsg"`
		Msgid 		int		`json:"msg_id"`
		MsgDataId int		`json:"msg_data_id"`
	}
	if err = json.Unmarshal(apiData, &bodyJson); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("apiData:%s", string(apiData)))
	}

	if bodyJson.Errcode != 0 {
		return 0, errors.New(fmt.Sprintf("%d:%s", bodyJson.Errcode, bodyJson.Errmsg))
	}

	return bodyJson.Msgid, nil
}