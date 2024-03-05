package wechat

import (
	"context"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/request"
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
