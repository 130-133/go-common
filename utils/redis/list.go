package redis

import (
	"context"
	"runtime/debug"

	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"git.tyy.com/llm-PhotoMagic/go-common/utils/help"
)

func (r *MRedis) PushList(key string, data interface{}) *redis.IntCmd {
	return r.RPush(key, help.ToString(data))
}

func (r *MRedis) PullList(key string) *MPop {
	ch := make(chan []string, 1)
	defer help.Go(func() {
		for {
			slice := r.BLPop(0, key)
			if slice.Err() != nil {
				close(ch)
				return
			}
			ch <- slice.Val()
		}
	})
	return &MPop{
		MRedis: r,
		ctx:    context.Background(),
		key:    key,
		ch:     ch,
	}
}

func (r *MRedis) Pipeline() *MPipeline {
	return &MPipeline{
		Pipeliner: r.Client.Pipeline(),
	}
}

func (p *MPipeline) PushList(key string, data interface{}) {
	p.RPush(key, help.ToString(data))
}

func (p *MPop) Callback(fc func(msg *MyMessage) error) {
	tracer := otel.Tracer(ztrace.TraceName)
	for {
		msg, ok := <-p.ch
		//chan关闭退出
		if !ok {
			return
		}
		if msg == nil {
			continue
		}
		body := gjson.Parse(msg[1])
		traceparent := body.Get("traceparent").String()
		tracestate := body.Get("tracestate").String()
		retryNum := body.Get("retry_num").Int()
		ctx := p.ExtractTrace(traceparent, tracestate)
		ctx, _ = tracer.Start(ctx, p.key, trace.WithSpanKind(trace.SpanKindConsumer))
		args := &MyMessage{
			redis:       p.MRedis,
			ctx:         ctx,
			key:         p.key,
			cmd:         "list",
			data:        msg[1],
			retryNum:    retryNum,
			maxRetryNum: 5,
		}
		go func(args *MyMessage) {
			defer func() {
				err := recover()
				if err != nil {
					p.Logger.Error("abnormal error", err, string(debug.Stack()))
				}
			}()
			span := trace.SpanFromContext(args.ctx)
			defer span.End()
			span.SetStatus(codes.Ok, "")
			if err := fc(args); err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
		}(args)
	}
}

func (p *MPop) ExtractTrace(traceparent, tracestate string) context.Context {
	propagator := otel.GetTextMapPropagator()
	carr := propagation.MapCarrier{}
	if traceparent != "" {
		carr.Set("traceparent", traceparent)
	}
	if tracestate != "" {
		carr.Set("tracestate", tracestate)
	}
	ctx := propagator.Extract(p.ctx, &carr)
	return ctx
}
