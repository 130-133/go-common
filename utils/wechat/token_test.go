package wechat

import (
	"context"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/redis"
	"testing"
)

func TestNewTokenManager(t *testing.T) {
	rdb := redis.NewRedisConn(
		redis.WithAddress("localhost:6379"),
		redis.WithPwd(""),
		redis.WithDb(0),
	)
	tm := NewTokenManager("wxab1a8beb9c37dd47", "", rdb)
	tm.GetToken(context.Background())
}
