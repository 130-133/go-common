package locker

import (
	"testing"
	"time"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/locker/client"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/redis"
)

func TestNewLocker(t *testing.T) {
	conn := redis.NewRedisConn(redis.WithAddress("10.0.0.106:6379"), redis.WithDb(4))
	l := NewLocker(client.NewRedisLocker(conn))
	a := l.WithKey("xxx").Acquire()
	t.Log(1)
	go l.WithKey("xxx").Acquire()
	t.Log(2)
	time.Sleep(10 * time.Second)
	a.Release()
	t.Log(3)
	time.Sleep(10 * time.Second)
}

func TestLocker_AcquireNoBlock(t *testing.T) {
	conn := redis.NewRedisConn(redis.WithAddress("10.0.0.106:6379"), redis.WithDb(4))
	l := NewLocker(client.NewRedisLocker(conn))
	go func() {
		a := l.WithKey("xxx").AcquireNoBlock()
		t.Logf("a: isLock:%v isGetLock:%v", a.IsLock(), a.IsGetLock())
	}()
	go func() {
		b := l.WithKey("xxx").AcquireNoBlock()
		t.Logf("b: isLock:%v isGetLock:%v", b.IsLock(), b.IsGetLock())
	}()
	time.Sleep(time.Second)
}
