package service

import (
	"github.com/gin-gonic/gin"
)

func (s *WechatHttpService) UserList(c *gin.Context) (interface{}, error) {

	rtn, err := s.userBiz.GetUserList()
	if err != nil {
		return nil, err
	}

	//这里需要一个dto对象做一次deep copy，暂时缺省




	return rtn, nil
}

func (s *WechatHttpService) TagList(c *gin.Context) (interface{}, error) {

	rtn, err := s.userBiz.GetTagList()
	if err != nil {
		return nil, err
	}

	//这里需要一个dto对象做一次deep copy，暂时缺省




	return rtn, nil
}
