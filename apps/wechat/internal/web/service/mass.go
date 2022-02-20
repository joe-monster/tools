package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func (s *WechatHttpService) SendMass(c *gin.Context) (interface{}, error) {
	type params struct {

		TagId 			int			`form:"tag_id"`
		Openids 		string		`form:"openids"`

		MsgType			string		`form:"msg_type"`

		//MediaId			string		`form:"media_id"`
		//MediaIds		string		`form:"media_ids"`
		Text			string		`form:"text"`

	}
	var param params
	if err := c.ShouldBind(&param); err != nil {
		return nil, err
	}

	paramData := make(map[string]string)
	//paramData["media_id"] = param.MediaId
	//paramData["media_ids"] = param.MediaIds
	paramData["text"] = param.Text

	if param.TagId != 0 {
		if err := s.massBiz.SendMassByTag(param.TagId, param.MsgType, paramData);err != nil {
			return nil, err
		}
	} else if param.Openids != "" {
		if err := s.massBiz.SendMassByUser(strings.Split(param.Openids, ","), param.MsgType, paramData);err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("发送对象未设置")
	}

	return nil, nil
}
