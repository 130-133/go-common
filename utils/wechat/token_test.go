package wechat

import (
	"context"
	"github.com/130-133/go-common/utils/redis"
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
