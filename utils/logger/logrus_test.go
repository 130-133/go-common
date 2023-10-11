package logger

import (
	"context"
	"testing"
	"time"
)

func TestNewMLogger(t *testing.T) {
	m := NewMLogger("aaa")
	for i := 100000; i > 0; i-- {
		go m.WithCtx(context.Background()).Info("aaaa")
		t.Log(i)
	}
	time.Sleep(3 * time.Second)
}
