package rabbitmq

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/atomic"
)

type Pool struct {
	mq     *RabbitMQ
	max    int
	locker sync.Mutex
	pool   []*ChannelPool
	freeze atomic.Bool
}

type ChannelPool struct {
	channel        *amqp.Channel
	channelCloseCh chan *amqp.Error
}

func NewChannelPool(mq *RabbitMQ) *Pool {
	return &Pool{
		max: 3,
		mq:  mq,
	}
}

func (p *Pool) Init() *Pool {
	for n := 0; n < p.max; n++ {
		p.locker.Lock()
		p.pool = append(p.pool, p.newChannel())
		p.locker.Unlock()
	}
	go p.watch()
	return p
}

func (p *Pool) Adjust(max int) *Pool {
	p.max = max
	p.CloseAll()
	p.Init()
	return p
}

func (p *Pool) CloseAll() {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.freeze.Swap(true)
	newPool := make([]*ChannelPool, 0)
	for _, pool := range p.pool {
		pool.channel.Close()
	}
	p.pool = newPool
}

func (p *Pool) Acquire() *amqp.Channel {
	l := len(p.pool)
	if l == 0 {
		return nil
	}
	hash := time.Now().UnixMilli() % int64(l)
	return p.pool[hash].channel
}

func (p *Pool) newChannel() *ChannelPool {
	p.mq.WaitReConnect()
	c := p.mq.NewChannel()
	ch := make(chan *amqp.Error)
	c.NotifyClose(ch)
	return &ChannelPool{
		channel:        c,
		channelCloseCh: ch,
	}
}

func (p *Pool) watch() {
	for {
		for k, pool := range p.pool {
			if p.freeze.Load() {
				return
			}
			select {
			case <-pool.channelCloseCh:
				p.pool[k] = p.newChannel()
			default:
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
