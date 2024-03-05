package wechat

import (
	"context"
	"fmt"
	red "github.com/go-redis/redis"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/redis"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/request"
	"time"
)

type AccessToken string
type JSTicket string

type TokenManager struct {
	AppId  string
	Secret string
	Redis  *redis.MRedis
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type jsTicketResp struct {
	ErrCode   int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func NewTokenManager(appId, secret string, rds *redis.MRedis) *TokenManager {
	return &TokenManager{
		AppId:  appId,
		Secret: secret,
		Redis:  rds,
	}
}

// GetToken 获取access_token
func (tm *TokenManager) GetToken(ctx context.Context) (token AccessToken, err error) {
	key := fmt.Sprintf(wxAccessTokenKey, tm.AppId)
	s, err := tm.Redis.Get(key).Result()
	if err != nil && err != red.Nil {
		return
	} else if s != "" {
		token = AccessToken(s)
		return
	}
	lockKey := fmt.Sprintf(lockAccessToken, tm.AppId)
	defer tm.Redis.Del(lockKey)
	lock, err := tm.Redis.SetNX(lockKey, 1, time.Second*60).Result()
	if err != nil {
		return
	}
	if lock { // 抢锁成功
		resp := request.NewRestyReq(ctx).Get(
			host,
			"/cgi-bin/token",
			map[string]interface{}{
				"grant_type": "client_credential",
				"appid":      tm.AppId,
				"secret":     tm.Secret,
			})
		if err := resp.GetError(); err != nil {
			return token, err
		}
		//data := gjson.GetBytes(resp.GetBody(), "access_token").Raw
		at := &accessTokenResp{}
		err = resp.GetJson(at)
		if err != nil {
			return "", err
		}
		token = AccessToken(at.AccessToken)
		if at.AccessToken != "" {
			tm.Redis.Set(key, at.AccessToken, time.Second*7000)
			return token, nil
		} else {
			return
		}
	} else {
		return
	}
}

// GetTicket 获取jsapi_ticket
func (tm *TokenManager) GetTicket(ctx context.Context, token AccessToken) (ticket JSTicket, err error) {
	if token == "" {
		return "", nil
	}
	key := fmt.Sprintf(wxJSTicket, tm.AppId)
	s, err := tm.Redis.Get(key).Result()
	if err != nil && err != red.Nil {
		return
	} else if s != "" {
		ticket = JSTicket(s)
		return
	}
	lockKey := fmt.Sprintf(lockJSTicket, tm.AppId)
	defer tm.Redis.Del(lockKey)
	lock, err := tm.Redis.SetNX(lockKey, 1, time.Second*60).Result()
	if err != nil {
		return
	}
	if lock { // 抢锁成功
		resp := request.NewRestyReq(ctx).Get(
			host,
			"/cgi-bin/ticket/getticket",
			map[string]interface{}{
				"type":         "jsapi",
				"access_token": string(token),
			})
		if err := resp.GetError(); err != nil {
			return "", err
		}
		//data := gjson.GetBytes(resp.GetBody(), "ticket").Raw
		tk := &jsTicketResp{}
		err = resp.GetJson(tk)
		if err != nil {
			return "", err
		}
		ticket = JSTicket(tk.Ticket)
		if tk.Ticket != "" {
			tm.Redis.Set(key, tk.Ticket, time.Second*7000)
			return ticket, nil
		} else {
			return
		}
	} else {
		return
	}
}
