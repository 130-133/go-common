package queue

import (
	"context"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/logger"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/redis"
)

type Plugin interface {
	Add(taskName string, function any)
	Start() error
	Stop()
	GetType() PluginType
}

type pluginOpt func(*pluginOpts)

func (o pluginOpt) Apply(arg *pluginOpts) {
	o(arg)
}

type pluginOpts struct {
	ctx    context.Context
	logger logger.ILogger
	redis  *redis.MRedis
}

func WithContext(ctx context.Context) pluginOpt {
	return func(opt *pluginOpts) {
		opt.ctx = ctx
	}
}

func WithLogger(logger logger.ILogger) pluginOpt {
	return func(opt *pluginOpts) {
		opt.logger = logger
	}
}
func WithRedis(redis *redis.MRedis) pluginOpt {
	return func(opt *pluginOpts) {
		opt.redis = redis
	}
}

// NewNonePlugin 空插件
func NewNonePlugin(opts *ExtendCtx) Plugin {
	return nonePlugin{}
}

type PluginFunc func(ctx context.Context)

// nonePlugin 空插件
type nonePlugin struct{}

func (p nonePlugin) Add(taskName string, fc any) {}
func (p nonePlugin) Start() error                { return nil }
func (p nonePlugin) Stop()                       {}
func (p nonePlugin) GetType() PluginType         { return NonePluginType }
