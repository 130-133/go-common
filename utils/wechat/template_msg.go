package wechat

import (
	"context"
	"github.com/130-133/go-common/utils/request"
)

func SendTemplateMsg(ctx context.Context, token AccessToken, data interface{}) error {
	resp := request.NewRestyReq(ctx).Post(
		host,
		"/cgi-bin/message/template/send?access_token="+string(token),
		data,
	)
	if resp.GetError() != nil {
		return resp.GetError()
	}
	return nil
}
