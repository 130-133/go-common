package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"

	"github.com/130-133/go-common/utils/logger"
)

// 接入Asynq实现延时任务，以方便任务的可视化监控

type AsynqServer struct {
	Server   *asynq.Server
	ServeMux *asynq.ServeMux
	rdsOpt   asynq.RedisClientOpt
	config   asynq.Config
	logger   asynq.Logger
	tasks    []string
}

func (r *MRedis) NewAsynqServer() *AsynqServer {
	rdsOpt := asynq.RedisClientOpt{
		Addr:      r.Options().Addr,
		Password:  r.Options().Password,
		DB:        r.Options().DB,
		PoolSize:  r.Options().PoolSize,
		TLSConfig: r.Options().TLSConfig,
	}
	config := asynq.Config{
		Concurrency:  1,
		ErrorHandler: asynq.ErrorHandlerFunc(r.reportError),
		Queues:       make(map[string]int),
	}
	return &AsynqServer{
		ServeMux: asynq.NewServeMux(),
		rdsOpt:   rdsOpt,
		config:   config,
		logger:   r.Logger,
	}
}

func (r *MRedis) reportError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		r.Logger.Error("重试次数超过最大值，", "，Queue：", task.Type(), "，Error：", err.Error())
	}
}

func (q *AsynqServer) Logger() logger.ILogger {
	return q.logger
}

// Queue 设置队列
func (q *AsynqServer) Queue(name string, priority int) *AsynqServer {
	q.config.Queues[name] = priority
	return q
}

// Queues 设置队列
func (q *AsynqServer) Queues(queues map[string]int) *AsynqServer {
	q.config.Queues = queues
	return q
}

// Concurrency 设置并发数
func (q *AsynqServer) Concurrency(concurrency int) *AsynqServer {
	q.config.Concurrency = concurrency
	return q
}

// Add 添加任务
func (q *AsynqServer) Add(typeName string, handler func(context.Context, *asynq.Task) error) *AsynqServer {
	q.tasks = append(q.tasks, typeName)
	q.ServeMux.HandleFunc(typeName, handler)
	return q
}

// Start 启动服务
func (q *AsynqServer) Start() error {
	if len(q.tasks) == 0 {
		return nil
	}
	q.Server = asynq.NewServer(q.rdsOpt, q.config)
	q.Server.Run(q.ServeMux)
	return nil
}

// Stop 停止服务
func (q *AsynqServer) Stop() {
	if q.Server == nil {
		return
	}
	q.Server.Stop()
}

type AsynqClient struct {
	client  *asynq.Client
	rdsOpt  asynq.RedisClientOpt
	options []asynq.Option
}

func (r *MRedis) NewAsyncClient() *AsynqClient {
	rdsOpt := asynq.RedisClientOpt{
		Addr:      r.Options().Addr,
		Password:  r.Options().Password,
		DB:        r.Options().DB,
		PoolSize:  r.Options().PoolSize,
		TLSConfig: r.Options().TLSConfig,
	}
	c := asynq.NewClient(rdsOpt)
	return &AsynqClient{
		client: c,
		rdsOpt: rdsOpt,
	}
}

// Queue 设置队列
func (r *AsynqClient) Queue(name string) *AsynqClient {
	r.options = append(r.options, asynq.Queue(name))
	return r
}

// Publish 发布任务
func (r *AsynqClient) Publish(ctx context.Context, taskType string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = r.client.EnqueueContext(ctx, asynq.NewTask(taskType, body), r.options...)
	return err
}

// PublishDelay 发布延迟任务
func (r *AsynqClient) PublishDelay(ctx context.Context, taskType string, data any, delay time.Duration) error {
	r.options = append(r.options, asynq.ProcessIn(delay))
	return r.Publish(ctx, taskType, data)
}

// PublishDelayAt 发布指定时间执行的任务
func (r *AsynqClient) PublishDelayAt(ctx context.Context, taskType string, data any, at time.Time) error {
	r.options = append(r.options, asynq.ProcessAt(at))
	return r.Publish(ctx, taskType, data)
}
