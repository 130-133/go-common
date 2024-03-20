package queue

import (
	"context"
	"fmt"

	"github.com/130-133/go-common/utils/logger"
	"github.com/130-133/go-common/utils/rabbitmq"
)

type rabbitmqPlugin struct {
	ctx      context.Context
	ctxClose context.CancelFunc
	logger   logger.ILogger
	mq       *rabbitmq.RabbitMQ
	queues   map[string]RabbitmqFunc
}

type RabbitmqFunc func(msg *rabbitmq.ConsumeMessage) error

func NewRabbitmqPlugin(rbmq *rabbitmq.RabbitMQ, opts ...pluginOpt) *rabbitmqPlugin {
	opt := &pluginOpts{
		ctx:    context.Background(),
		logger: logger.LocalLogger{},
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	return &rabbitmqPlugin{
		ctx:    opt.ctx,
		mq:     rbmq,
		logger: opt.logger,
		queues: make(map[string]RabbitmqFunc),
	}
}

func (q *rabbitmqPlugin) Add(queueName string, fn any) {
	fc, ok := fn.(RabbitmqFunc)
	if !ok {
		return
	}
	q.queues[queueName] = fc
}

func (q *rabbitmqPlugin) Start() error {
	q.ctx, q.ctxClose = context.WithCancel(q.ctx)
	for queueName, fc := range q.queues {
		q.logger.Info(fmt.Sprintf("queue start: %s", queueName))
		q.mq.WithContext(q.ctx).Queue(queueName).ConsumeFn(rabbitmq.Callback(fc))
	}
	return nil
}

func (q *rabbitmqPlugin) Stop() {
	for queueName := range q.queues {
		q.logger.Info(fmt.Sprintf("queue stop: %s", queueName))
	}
	q.ctxClose()
}

func (q *rabbitmqPlugin) GetType() PluginType {
	return RabbitmqPluginType
}
