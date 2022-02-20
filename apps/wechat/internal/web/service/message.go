package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *WechatHttpService) SendPublishIdType(c *gin.Context) (interface{}, error) {
	type params struct {

		TemplateId 		string		`form:"template_id" binding:"required"`

		Url				string		`form:"url"`

		Data			string		`form:"data" binding:"required"`

		Openid			string		`form:"openid"`
	}
	var param params
	if err := c.ShouldBind(&param); err != nil {
		return nil, err
	}

	paramData := make(map[string]string)
	if err := json.Unmarshal([]byte(param.Data), &paramData);err != nil {
		return nil, err
	}

	rtn, err := s.messageBiz.SendPublish(param.TemplateId, param.Url, param.Openid, paramData)
	if err != nil {
		return nil, err
	}

	//这里需要一个dto对象做一次deep copy，暂时缺省



	return rtn, nil
}


