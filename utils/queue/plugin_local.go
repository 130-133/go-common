package queue

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/atomic"

	"github.com/130-133/go-common/utils/logger"
)

type localPlugin struct {
	taskList sync.Map
	mux      sync.Mutex
	Progress []*Progress
	logger   logger.ILogger
}

type Progress struct {
	Name     string
	ctx      context.Context
	ctxClose context.CancelFunc
	fn       Task
	state    atomic.Bool // 状态
	logger   logger.ILogger
}

type Task func(ctx context.Context)

func NewLocalPlugin(opts ...pluginOpt) *localPlugin {
	opt := &pluginOpts{
		logger: logger.LocalLogger{},
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	return &localPlugin{
		logger: opt.logger,
	}
}

func (p *localPlugin) GetType() PluginType {
	return LocalPluginType
}

func (p *localPlugin) Add(taskName string, fn any) {
	fc := fn.(PluginFunc)
	if _, ok := p.taskList.Load(taskName); ok {
		p.DelProgress(taskName)
	}
	progress := p.NewProgress(taskName, Task(fc))
	p.taskList.Store(taskName, progress)
	p.Progress = append(p.Progress, progress)
}

func (p *localPlugin) Start() error {
	for _, progress := range p.Progress {
		go func(pro *Progress) {
			defer goRecover(p.logger)
			pro.Start()
		}(progress)
	}
	return nil
}

func (p *localPlugin) Stop() {
	for _, progress := range p.Progress {
		progress.Stop()
	}
}

func (p *localPlugin) NewProgress(name string, fn Task) *Progress {
	ctx, closeFn := context.WithCancel(context.Background())
	ctx = trace.ContextWithSpanContext(ctx, trace.SpanContextFromContext(nil))
	return &Progress{
		ctx:      ctx,
		ctxClose: closeFn,
		Name:     name,
		fn:       fn,
		logger:   p.logger,
	}
}

func (p *localPlugin) StartProgress(name string) {
	if val, ok := p.taskList.Load(name); ok {
		progress := val.(*Progress)
		progress.ctx, progress.ctxClose = context.WithCancel(context.Background())
		progress.ctx = trace.ContextWithSpanContext(progress.ctx, trace.SpanContextFromContext(nil))
		progress.Start()
		p.taskList.Store(name, progress)
	}
}

func (p *localPlugin) StopProgress(name string) {
	if val, ok := p.taskList.Load(name); ok {
		progress := val.(*Progress)
		progress.Stop()
	}
}

func (p *localPlugin) DelProgress(name string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	var progressList []*Progress
	for _, v := range p.Progress {
		if v.Name == name {
			v.Stop()
			continue
		}
		progressList = append(progressList, v)
	}
	p.Progress = progressList
}

func (p *localPlugin) Show() []*Progress {
	return p.Progress
}

func (p *Progress) Start() {
	if p.state.Load() {
		return
	}
	p.state.Store(true)
	p.logger.Info(fmt.Sprintf("queue start: %s", p.Name))
	p.fn(p.ctx)
}

func (p *Progress) Stop() {
	if !p.state.Load() {
		return
	}
	p.ctxClose()
	p.state.Store(false)
	p.logger.Info(fmt.Sprintf("queue stop: %s", p.Name))
}
