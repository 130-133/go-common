package redis

import (
	"context"
	"encoding/json"
	"reflect"
)

type MyMessage struct {
	redis       *MRedis
	ctx         context.Context
	id          string
	cmd         string
	key         string
	data        string
	retryNum    int64
	maxRetryNum int64
}

func (m *MyMessage) Context() context.Context {
	return m.ctx
}
func (m *MyMessage) ID() string {
	return m.id
}
func (m *MyMessage) Key() string {
	return m.key
}
func (m *MyMessage) SetMaxRetry(max int64) {
	m.maxRetryNum = max
}

func (m *MyMessage) Unmarshal(data interface{}) error {
	err := json.Unmarshal([]byte(m.data), data)
	if err == nil {
		return nil
	}
	tType := reflect.TypeOf(data)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	if tType.Kind() != reflect.String {
		return nil
	}
	tVal := reflect.Indirect(reflect.ValueOf(data))
	vVal := reflect.ValueOf(m.data)
	tVal.Set(vVal)
	return nil
}

func (m *MyMessage) String() string {
	return m.data
}

func (m *MyMessage) RetryNum() int64 {
	return m.retryNum
}

// Retry 重试 共可重试几次
func (m *MyMessage) Retry() error {
	if m.retryNum >= m.maxRetryNum {
		return nil
	}
	m.retryNum++
	switch m.cmd {
	case "list":
		tmpMap := make(map[string]any)
		_ = m.Unmarshal(&tmpMap)
		tmpMap["retry_num"] = m.retryNum
		if bytes, err := json.Marshal(tmpMap); err != nil {
			// 没有标记容易导致死循环
			return m.redis.PushList(m.key, m.data).Err()
		} else {
			return m.redis.PushList(m.key, string(bytes)).Err()
		}
	case "stream":
		return m.redis.RetryPushQueue(m.ctx, m.key, m.data, m.retryNum).Err()
	}
	return nil
}

type QueueCommon struct {
	TraceParent     string `json:"traceparent"`       //链路追踪
	TraceState      string `json:"tracestate"`        //链路追踪
	CreatedUnixMill int64  `json:"created_unix_mill"` //创建的毫秒时间戳
	RetryNum        int64  `json:"retry_num"`         //重试次数
}
