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

func (tm *TokenManager) GetToken(ctx context.Context) (token string, err error) {
	key := fmt.Sprintf(wxAccessTokenKey, tm.AppId)
	token, err = tm.Redis.Get(key).Result()
	fmt.Println(err)
	if err != nil && err != red.Nil {
		return "", err
	}
	lockKey := fmt.Sprintf(lockToken, tm.AppId)
	defer tm.Redis.Del(lockKey)
	lock, err := tm.Redis.SetNX(lockKey, 1, time.Second*60).Result()
	if err != nil {
		return "", err
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
			return "", err
		}
		data := gjson.GetBytes(resp.GetBody(), "access_token").Raw
		if data != "" {
			tm.Redis.Set(key, data, time.Second*7000)
			return data, nil
		} else {
			return "", nil
		}
	} else {
		return "", nil
	}
}
