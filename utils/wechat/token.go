package wechat

import (
	"context"
	"fmt"
	"gitea.com/llm-PhotoMagic/go-common/utils/redis"
	"gitea.com/llm-PhotoMagic/go-common/utils/request"
	red "github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"time"
)

type AccessToken string
type JSTicket string

type TokenManager struct {
	AppId  string
	Secret string
	Redis  *redis.MRedis
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
		data := gjson.GetBytes(resp.GetBody(), "access_token").Raw
		token = AccessToken(data)
		if data != "" {
			tm.Redis.Set(key, data, time.Second*7000)
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
				"type":         "jsapi_ticket",
				"access_token": string(token),
			})
		if err := resp.GetError(); err != nil {
			return "", err
		}
		data := gjson.GetBytes(resp.GetBody(), "ticket").Raw
		ticket = JSTicket(data)
		if data != "" {
			tm.Redis.Set(key, data, time.Second*7000)
			return ticket, nil
		} else {
			return
		}
	} else {
		return
	}
}
