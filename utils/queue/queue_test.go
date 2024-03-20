package queue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/130-133/go-common/utils/redis"
)

func TestBaseQueue_Lock(t *testing.T) {
	q := NewQueue(&ExtendCtx{
		Redis: redis.NewRedisConn(redis.WithAddress("10.0.0.106:6379")),
	})
	t.Log(q.ExtendCtx.Redis.Ping().Err())
	go func() {
		t.Log("this one get lock")
		l := q.Locker.WithKey("test").Acquire()
		defer l.Release()
		t.Log("this one has lock")
		time.Sleep(10 * time.Second)
		t.Log("this one release lock")
	}()

	go func() {
		t.Log("this two get lock")
		l := q.Locker.WithKey("test").Acquire()
		defer l.Release()
		t.Log("this two has lock")
		time.Sleep(10 * time.Second)
		t.Log("this two release lock")
	}()
	time.Sleep(21 * time.Second)
}

func TestLocker_Acquire(t *testing.T) {
	q := BaseQueue{
		ExtendCtx: &ExtendCtx{
			Redis: redis.NewRedisConn(redis.WithAddress("10.0.0.106:6379"), redis.WithDb(4)),
		},
	}
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("get2-1")
		lock := q.Locker.WithKey(`sir@AlipaySigningCheck%s`, "1234").Acquire()
		fmt.Println("get2-2")
		defer lock.Release()
		time.Sleep(5 * time.Second)
	}()
	fmt.Println("get1-1")
	lock := q.Locker.WithKey(`sir@AlipaySigningCheck%s`, "1234").Acquire()
	fmt.Println("get1-2")
	time.Sleep(5 * time.Second)
	lock.Release()
	fmt.Println("get1-3")
	time.Sleep(5 * time.Second)
}

func TestNewQueue(t *testing.T) {
	q := NewQueue(&ExtendCtx{})
	q.AddProgress("a", func(ctx context.Context) { t.Log(1) })
	q.AddProgress("b", func(ctx context.Context) { t.Log(2) })

	q.Start()
	time.Sleep(time.Second)
}

func TestBaseQueue_Stop(t *testing.T) {
	q := NewQueue(&ExtendCtx{})
	q.AddProgress("a", func(ctx context.Context) {
		t.Log(1)
		time.Sleep(2 * time.Second)
		t.Log(2)
	})
	q.Start()
	time.Sleep(time.Second)
	q.Stop()
	time.Sleep(2 * time.Second)
}
