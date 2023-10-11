package queue

import (
	"context"

	"github.com/hibiken/asynq"

	"git.tyy.com/llm-PhotoMagic/go-common/utils/logger"
	"git.tyy.com/llm-PhotoMagic/go-common/utils/redis"
)

type asynqPlugin struct {
	asynqSer *redis.AsynqServer
	logger   logger.ILogger
}

func NewAsynqPlugin(rds *redis.MRedis, opts ...pluginOpt) *asynqPlugin {
	opt := &pluginOpts{
		logger: logger.LocalLogger{},
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	return &asynqPlugin{
		asynqSer: rds.NewAsynqServer(),
		logger:   opt.logger,
	}
}

type AsynqFunc func(ctx context.Context, task *asynq.Task) error

func (q *asynqPlugin) Add(taskName string, fn any) {
	fc, ok := fn.(AsynqFunc)
	if !ok {
		return
	}
	q.asynqSer.Add(taskName, fc)
}

func (q *asynqPlugin) Start() error {
	go func() {
		defer goRecover(q.asynqSer.Logger())
		q.asynqSer.Start()
	}()
	return nil
}

func (q *asynqPlugin) Stop() {
	q.asynqSer.Stop()
}

func (q *asynqPlugin) GetType() PluginType {
	return AsynqPluginType
}
