package queue

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/130-133/go-common/utils/locker"
	"github.com/130-133/go-common/utils/locker/client"
	"github.com/130-133/go-common/utils/logger"
	"github.com/130-133/go-common/utils/redis"
)

type BaseQueue struct {
	ExtendCtx *ExtendCtx
	Locker    *locker.ClientLocker
	logger    logger.ILogger
	plugins   map[PluginType]Plugin
}

type ExtendCtx struct {
	Redis *redis.MRedis
}

type IBaseQueue interface {
	Start()
	Stop()
	Add(taskName string, function any, pluginType PluginType)
	WithPlugin(...Plugin)
}

type PluginType int

const (
	NonePluginType     PluginType = iota - 1 //空插件
	LocalPluginType                          //默认本地
	AsynqPluginType                          //asynq包
	RabbitmqPluginType                       //rabbitmq包
)

func NewQueue(ctx *ExtendCtx) *BaseQueue {
	var Locker *locker.ClientLocker
	if ctx.Redis != nil {
		Locker = locker.NewLocker(client.NewRedisLocker(ctx.Redis))
	}
	queue := &BaseQueue{
		ExtendCtx: ctx,
		Locker:    Locker,
		logger:    ctx.Redis.Logger,
		plugins: map[PluginType]Plugin{
			LocalPluginType: NewLocalPlugin(WithLogger(ctx.Redis.Logger)),
		},
	}
	return queue
}

func (q *BaseQueue) WithPlugin(ps ...Plugin) {
	for _, plugin := range ps {
		q.plugins[plugin.GetType()] = plugin
	}
}

func (q *BaseQueue) Add(taskName string, function any, pluginType PluginType) {
	if _, ok := q.plugins[pluginType]; !ok {
		q.plugins[LocalPluginType].Add(taskName, function)
		return
	}
	q.plugins[pluginType].Add(taskName, function)
}

func (q *BaseQueue) AddProgress(queueName string, fn PluginFunc) {
	q.Add(queueName, fn, LocalPluginType)
}
func (q *BaseQueue) AddAsynq(typeName string, fn AsynqFunc) {
	q.Add(typeName, fn, AsynqPluginType)
}
func (q *BaseQueue) AddRabbitmq(queueName string, fn RabbitmqFunc) {
	q.Add(queueName, fn, RabbitmqPluginType)
}

func goRecover(logger logger.ILogger) {
	err := recover()
	switch err.(type) {
	case runtime.Error:
		logger.Error(fmt.Sprintf("%+v %s", err, debug.Stack()))
	}
}

func (q *BaseQueue) Start() {
	// 插件
	for _, plugin := range q.plugins {
		go func(p Plugin) {
			defer goRecover(q.logger)
			p.Start()
		}(plugin)
	}
}

func (q *BaseQueue) Stop() {
	// 插件
	for _, plugin := range q.plugins {
		plugin.Stop()
	}
}
