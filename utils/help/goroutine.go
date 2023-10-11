package help

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// Go 协程 带recover恢复
func Go(fn func()) {
	go func() {
		defer func() {
			err := recover()
			switch err.(type) {
			case runtime.Error:
				_ = fmt.Errorf("abnormal error %v", err)
			}
		}()
		fn()
	}()
}

// GoWait 协程 等待所有方法完成才结束
func GoWait(fn ...func()) {
	if len(fn) == 0 {
		return
	}
	var mux sync.WaitGroup
	mux.Add(len(fn))
	for _, call := range fn {
		Go(func() {
			defer mux.Done()
			call()
		})
	}
	mux.Wait()
}

// GoForeachLimit 并发带限流循环
func GoForeachLimit(data interface{}, limit int, fn func(index int)) {
	var (
		run int
		max int
		mux sync.WaitGroup
	)

	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		return
	}
	v := reflect.ValueOf(data)
	max = v.Len()
	for n := 0; n < max; n++ {
		run++
		mux.Add(1)
		go func(n int) {
			err := recover()
			if err != nil {
				fmt.Errorf("abnormal error %v", err)
			}
			defer mux.Done()
			fn(n)
		}(n)
		if run%limit == 0 {
			mux.Wait()
		}
	}
	mux.Wait()
}
