package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hibiken/asynq"
)

func TestAsynqClient_Publish(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	c := r.NewAsyncClient()
	c.Publish(context.TODO(), "222test", map[string]interface{}{
		"test2": "test22",
	})
}

func TestAsynqClient_PublishDelay(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	c := r.NewAsyncClient()
	c.PublishDelay(context.TODO(), "test111", map[string]interface{}{
		"test1": "test11",
	}, 10*time.Second)
}

func TestAsynqClient_PublishAt(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	c := r.NewAsyncClient()
	c.PublishDelayAt(context.TODO(), "test", map[string]interface{}{
		"test": "test",
	}, time.Now().Add(24*time.Hour))
}

func TestAsynqClient_PublishQueue(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	c := r.NewAsyncClient()
	c.Queue("test").Publish(context.TODO(), "test", map[string]interface{}{
		"test": "1",
	})
}

func TestAsynqServer_Start(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	s := r.NewAsynqServer()
	s.Add("test", func(ctx context.Context, task *asynq.Task) error {
		time.Sleep(2 * time.Second)
		t.Logf("%+v \t %+v", task.Type(), string(task.Payload()))
		b := struct {
			Test int `json:"test"`
		}{}
		_ = json.Unmarshal(task.Payload(), &b)
		t.Log(b)
		//if b.Test == 1 {
		//	t.Log("error")
		//	return errors.New("xxxxx")
		//}
		return nil
	})
	s.Start()
}

func TestAsynqClient_Queue(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	s := r.NewAsynqServer()
	s.Queue("queue", 1)
	s.Add("test", func(ctx context.Context, task *asynq.Task) error {
		t.Log(string(task.Payload()))
		return nil
	})
	s.Start()
}

func TestAsynqServer_Queues(t *testing.T) {
	r := NewRedisConn(WithAddress("127.0.0.1:6379"), WithPwd("123456"))
	s := r.NewAsynqServer()
	s.Queues(map[string]int{
		"queue": 1,
	})
	s.Add("test", func(ctx context.Context, task *asynq.Task) error {
		t.Log(string(task.Payload()))
		return nil
	})
	s.Start()
}
