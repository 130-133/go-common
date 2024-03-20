package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/130-133/go-common/utils/help"
	"github.com/130-133/go-common/utils/tracer"
)

func (r *MRedis) PushQueue(ctx context.Context, key string, data interface{}) *redis.StringCmd {
	value := map[string]interface{}{
		"data":        help.ToString(data),
		"traceparent": tracer.GetTraceParent(ctx),
		"retry_num":   0,
	}
	return r.XAdd(&redis.XAddArgs{
		Stream:       key,
		MaxLen:       300,
		MaxLenApprox: 300,
		ID:           "*",
		Values:       value,
	})
}

func (r *MRedis) RetryPushQueue(ctx context.Context, key string, data interface{}, retryNum int64) *redis.StringCmd {
	value := map[string]interface{}{
		"data":        help.ToString(data),
		"traceparent": tracer.GetTraceParent(ctx),
		"retry_num":   retryNum,
	}
	return r.XAdd(&redis.XAddArgs{
		Stream:       key,
		MaxLen:       300,
		MaxLenApprox: 300,
		ID:           "*",
		Values:       value,
	})
}

func (r *MRedis) WithPullGroup(key string) *MConsumer {
	// 初始化构建消费组
	groupName := fmt.Sprintf("%s_Group", key)
	r.XGroupCreateMkStream(key, groupName, "0")
	return &MConsumer{
		MRedis: r,
		group:  groupName,
		stream: key,
	}
}

type MConsumer struct {
	*MRedis
	group  string //消费组名
	stream string //stream，队列名
	i      int    //消费者起始编号
}

func (m *MConsumer) PullQueue(start string) *redis.XStreamSliceCmd {
	consumer := fmt.Sprintf("%s_%d", m.group, m.i)
	defer m.XGroupDelConsumer(m.stream, m.group, consumer)
	m.i++
	return m.XReadGroup(&redis.XReadGroupArgs{
		Group:    m.group,
		Consumer: consumer,
		Streams:  []string{m.stream, start},
		Count:    1,
		Block:    time.Millisecond,
	})
}

func (m *MConsumer) PullQueueBlock(fn func(msg *MyMessage) error) {
	hostName, _ := os.Hostname()
	consumer := fmt.Sprintf("%s_%d", hostName, m.i)
	defer m.XGroupDelConsumer(m.stream, m.group, consumer)
	m.i++
	resultChan := make(chan *redis.XMessage, 0)
	help.Go(func() {
		for {
			// 最新数据
			datas := m.XReadGroup(&redis.XReadGroupArgs{
				Group:    m.group,
				Consumer: consumer,
				Streams:  []string{m.stream, ">"},
				Count:    2,
				Block:    10 * time.Second,
			})
			if len(datas.Val()) == 0 {
				continue
			}
			data := datas.Val()[0]
			for _, item := range data.Messages {
				resultChan <- &item
			}
		}
	})
	help.Go(func() {
		for {
			// 未应答数据
			datas := m.XReadGroup(&redis.XReadGroupArgs{
				Group:    m.group,
				Consumer: consumer,
				Streams:  []string{m.stream, "0"},
				Count:    2,
				Block:    5 * time.Second,
			})
			if len(datas.Val()) == 0 {
				continue
			}
			data := datas.Val()[0]
			for _, item := range data.Messages {
				resultChan <- &item
			}
			time.Sleep(10 * time.Second)
		}
	})
	// 处理回调
	for {
		if message, ok := <-resultChan; ok {
			m.handler(message, fn)
		} else {
			break
		}
	}
}

func (m *MConsumer) handler(item *redis.XMessage, fn func(msg *MyMessage) error) {
	tracer := otel.Tracer(ztrace.TraceName)
	traceparent, _ := item.Values["traceparent"].(string)
	tracestate, _ := item.Values["tracestate"].(string)
	retryNumStr, _ := item.Values["retry_num"].(string)
	retryNum, _ := strconv.ParseInt(retryNumStr, 10, 64)
	ctx := m.ExtractTrace(context.Background(), traceparent, tracestate)
	ctx, span := tracer.Start(ctx, m.stream, trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()
	msg := MyMessage{
		redis:       m.MRedis,
		id:          item.ID,
		ctx:         ctx,
		key:         m.stream,
		cmd:         "stream",
		data:        item.Values["data"].(string),
		retryNum:    retryNum,
		maxRetryNum: 3,
	}
	if err := fn(&msg); err == nil {
		// 应答
		if err = m.XAck(m.stream, m.group, item.ID).Err(); err != nil {
			m.Logger.Error(err.Error())
		}
	}
}

func (m *MConsumer) ExtractTrace(ctx context.Context, traceparent, tracestate string) context.Context {
	propagator := otel.GetTextMapPropagator()
	carr := propagation.MapCarrier{}
	if traceparent != "" {
		carr.Set("traceparent", traceparent)
	}
	if tracestate != "" {
		carr.Set("tracestate", tracestate)
	}
	return propagator.Extract(ctx, &carr)
}
