package locker

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/atomic"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/help"
)

type Client interface {
	Set(key string, value interface{}, expiration time.Duration) (bool, error)
	Get(key string) string
	Del(key string) error
}

type ClientLocker struct {
	client  Client
	timeout time.Duration
}

type locker struct {
	client  Client
	key     string
	value   string
	timeout time.Duration
	error   error
	isLock  *atomic.Bool
}

var GLocker *ClientLocker

func InitGlobalLocker(c Client) *ClientLocker {
	if GLocker == nil {
		GLocker = NewLocker(c)
	}
	return GLocker
}

func NewLocker(c Client) *ClientLocker {
	return &ClientLocker{
		client:  c,
		timeout: 30 * time.Second,
	}
}

func (l *ClientLocker) copy() *locker {
	return &locker{
		client:  l.client,
		timeout: l.timeout,
		isLock:  atomic.NewBool(false),
	}
}

func (l *ClientLocker) WithKey(key string, match ...interface{}) *locker {
	c := l.copy()
	c.key = fmt.Sprintf(key, match...)
	c.value = help.GetRandstring(10)
	return c
}

// WithTimeout 加锁超时时长
func (l *locker) WithTimeout(value time.Duration) *locker {
	l.timeout = value
	return l
}

// Lock 加锁
func (l *locker) Lock() bool {
	ok, err := l.client.Set(l.key, l.value, l.timeout)
	if err != nil {
		return false
	}
	if ok {
		l.isLock.CAS(false, true)
	}
	return ok
}

// IsLock 判断全局是否有锁
func (l *locker) IsLock() bool {
	resValue := l.client.Get(l.key)
	if len(resValue) > 0 {
		return true
	}
	return false
}

// IsGetLock 判断是否当前逻辑进程加的锁
func (l *locker) IsGetLock() bool {
	return l.isLock.Load()
}

// AcquireNoBlock 获取锁，不阻塞
func (l *locker) AcquireNoBlock() *locker {
	if ok := l.Lock(); ok {
		return l
	}
	return l
}

// Acquire 获取锁，阻塞等待锁
func (l *locker) Acquire() *locker {
	ok := l.Lock()
	if ok {
		return l
	}
	t := time.NewTicker(10 * time.Millisecond)
	timeout := time.After(l.timeout)
	for {
		select {
		case <-t.C:
			if ok := l.Lock(); ok {
				t.Stop()
				l.error = nil
				return l
			}
		case <-timeout:
			l.error = errors.New("acquire lock timeout")
			t.Stop()
			return l
		}
	}
}

// Release 释放锁
func (l *locker) Release() error {
	if l.error != nil {
		return l.error
	}
	resValue := l.client.Get(l.key)
	if resValue == l.value {
		return l.client.Del(l.key)
	}
	return errors.New("can not release other lock, please wait a moment")
}
