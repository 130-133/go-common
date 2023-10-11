package redis

import (
	"testing"
	"time"
)

func TestMPipeline_PushList(t *testing.T) {
	r := NewRedisConn(WithPwd("123456"))
	r.PushList("clint", "xxxx")

	go func() {
		pp := r.Pipeline()
		pp.PushList("clint_pp", []int{1})
		pp.Exec()
	}()
	go func() {
		pp := r.Pipeline()
		pp.PushList("clint_pp", []int{2})
	}()
	go func() {
		pp := r.Pipeline()
		pp.PushList("clint_pp", []int{3})
		pp.Exec()
	}()
	time.Sleep(time.Second)
}

func TestMRedis_BPopList(t *testing.T) {
	//r := NewRedisConn(WithPwd("123456"))
	//a := r.BPopList(0, "clint")
	//type A struct {
	//	A int
	//}
	//b := A{}
	////a.Unmarshal(&b)
	//t.Log(b)
}
