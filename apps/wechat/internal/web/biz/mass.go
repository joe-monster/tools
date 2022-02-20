package biz

import (
	"tools/internal/pkg/wxapi"
)


type MassBiz struct {
	wxApiClient *wxapi.WxApiClient
}
func NewMassBiz(c *wxapi.WxApiClient) *MassBiz {
	return &MassBiz{
		wxApiClient: c,
	}
}

func (b *MassBiz) SendMassByTag(tagId int, msgType string, data map[string]string) error {

	var p wxapi.MassParamsType
	p.Filter.TagId = tagId

	b.sendMassMakeContentParams(&p, msgType, data)

	if _, err := b.wxApiClient.SendMassMessageByTag(&p);err != nil {
		return err
	}

	return nil
}

func (b *MassBiz) SendMassByUser(openids []string, msgType string, data map[string]string) error {

	var p wxapi.MassParamsType
	p.Touser = openids

	b.sendMassMakeContentParams(&p, msgType, data)

	if _, err := b.wxApiClient.SendMassMessageByUser(&p);err != nil {
		return err
	}

	return nil
}

func (b *MassBiz) sendMassMakeContentParams(p *wxapi.MassParamsType,msgType string, data map[string]string) {

	p.Msgtype = msgType

	switch msgType {
	//case wechat_api.MSG_TYPE_MPNEWS:
	//	p.Mpnews.MediaId = data["media_id"]
	//	p.SendIgnoreReprint = 0
	//	break
	case wxapi.MSG_TYPE_TEXT:
		p.Text.Content = data["text"]
		break
	//case wechat_api.MSG_TYPE_VOICE:
	//	p.Voice.MediaId = data["media_id"]
	//	break
	//case wechat_api.MSG_TYPE_IMAGE:
	//	p.Images.MediaIds = strings.Split(data["media_ids"], ",")
	//	p.Images.Recommend = ""
	//	p.Images.NeedOpenComment = 1
	//	p.Images.OnlyFansCanComment = 0
	//	break
	//case wechat_api.MSG_TYPE_VIDEO:
	//	p.Mpvideo.MediaId = data["media_id"]
	//	p.Mpvideo.Title = ""
	//	p.Mpvideo.Description = ""
	//	break
	default:
		break
	}

}