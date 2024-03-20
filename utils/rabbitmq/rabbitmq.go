package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"github.com/130-133/go-common/utils/logger"
)

type Option func(*RabbitMQ)
type RabbitMQ struct {
	conf         *ConnectConf     //连接配置
	conn         *amqp.Connection //连接
	channelPool  *Pool
	serverName   string
	rabbitUrl    string
	queueName    string //队列名称
	exchangeName string //交换机
	routingName  string //路由值
	mu           sync.Mutex
	logger       *logger.MLogger
}

type ConnectConf struct {
	Host        string
	User        string
	Pwd         string
	VirtualHost string
}

func (o Option) Apply(arg *RabbitMQ) {
	o(arg)
}

func WithHost(data string) Option {
	return func(arg *RabbitMQ) {
		arg.conf.Host = data
	}
}
func WithVirtualHost(data string) Option {
	return func(arg *RabbitMQ) {
		arg.conf.VirtualHost = data
	}
}
func WithUser(data string) Option {
	return func(arg *RabbitMQ) {
		arg.conf.User = data
	}
}
func WithPwd(data string) Option {
	return func(arg *RabbitMQ) {
		arg.conf.Pwd = data
	}
}

func WithLogger(mLogger *logger.MLogger) Option {
	return func(arg *RabbitMQ) {
		arg.logger = mLogger
	}
}

// 申请队列参数
type QueueDeclareParams struct {
	QueueName                              string // queue 名
	Durable, AutoDelete, Exclusive, NoWait bool
	Args                                   amqp.Table
}

// 消费者参数
type ConsumeSimpleParams struct {
	QueueName, Consumer                 string
	AutoAck, Exclusive, NoLocal, NoWait bool
	Args                                amqp.Table
	QueueDeclare                        QueueDeclareParams
}

// NewRabbitMQ 创建RabbitMQ结构体实例
func NewRabbitMQ(serverName string, opts ...Option) *RabbitMQ {
	rabbitmq := RabbitMQ{
		serverName: serverName,
		conf: &ConnectConf{
			Host: "127.0.0.1:5672",
			User: "root",
			Pwd:  "",
		},
	}
	for _, o := range opts {
		o.Apply(&rabbitmq)
	}
	if rabbitmq.logger == nil {
		rabbitmq.logger = logger.NewMLogger("default")
	}

	rabbitmq.conn = rabbitmq.NewConnect()
	rabbitmq.channelPool = NewChannelPool(&rabbitmq).Init()
	go rabbitmq.WatchConnect()
	return &rabbitmq
}

func (r *RabbitMQ) NewChannel() *amqp.Channel {
	r.mu.Lock()
	defer r.mu.Unlock()
	channel, err := r.conn.Channel()
	if err != nil {
		r.logger.WithError(err).Error("获取RabbitMQ失败", r.rabbitUrl)
		r.ReConnect()
	}
	return channel
}

// ChannelAcquire 从Channel池获取
func (r *RabbitMQ) ChannelAcquire() *amqp.Channel {
	return r.channelPool.Acquire()
}

// ChannelAcquire 从Channel池获取
func (r *RabbitMQ) ChannelAjust(max int) *RabbitMQ {
	r.channelPool.Adjust(max)
	return r
}

func (r *RabbitMQ) NewConnect() *amqp.Connection {
	//创建rabbitmq连接
	r.rabbitUrl = fmt.Sprintf("amqp://%s:%s@%s/%s", r.conf.User, r.conf.Pwd, r.conf.Host, r.conf.VirtualHost)
	conn, err := amqp.Dial(r.rabbitUrl)
	if err != nil {
		r.logger.WithError(err).Error("连接取RabbitMQ失败", r.rabbitUrl)
		r.ReConnect()
	}
	return conn
}

func (r *RabbitMQ) ReConnect() {
	tk := time.NewTicker(100 * time.Millisecond)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	for {
		if r.conn != nil {
			tk.Stop()
			break
		}
		select {
		case _ = <-tk.C:
			var err error
			if r.conn == nil {
				r.conn, err = amqp.Dial(r.rabbitUrl)
				if err != nil {
					break
				}
			}
		case _ = <-ctx.Done():
			panic(fmt.Sprintf("连接Rabbitmq失败：%s", r.rabbitUrl))
		}
	}
}

func (r *RabbitMQ) WatchConnect() {
	for {
		// 连接和通道重连
		if r.conn.IsClosed() {
			r.conn = r.NewConnect()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// WaitReConnect 阻塞等待重连
func (r *RabbitMQ) WaitReConnect() {
	for {
		if r.conn.IsClosed() {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
}

func (r *RabbitMQ) Channel() IEntry {
	return NewEntry(r).Channel(r.NewChannel())
}

func (r *RabbitMQ) WithContext(ctx context.Context) IEntry {
	return NewEntry(r).WithContext(ctx)
}

func (r *RabbitMQ) Exchange(data string) IEntry {
	return NewEntry(r).Exchange(data)
}

func (r *RabbitMQ) Queue(data string) IEntry {
	return NewEntry(r).Queue(data)
}

func (r *RabbitMQ) RoutingKey(data string) IEntry {
	return NewEntry(r).RoutingKey(data)
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() error {
	// 先关闭管道,再关闭链接
	r.channelPool.CloseAll()
	err := r.conn.Close()
	if err != nil {
		r.logger.WithError(err).Error("rmq链接关闭失败", r.rabbitUrl)
		return errors.New("rmq链接关闭失败")
	}
	return nil
}

// Publish 发布
func (r *RabbitMQ) Publish(ctx context.Context, message []byte, opts ...EntryOption) error {
	return NewEntry(r).WithContext(ctx).Publish(message, opts...)
}

// Consume 消费
func (r *RabbitMQ) Consume(ctx context.Context, opts ...EntryOption) (<-chan amqp.Delivery, error) {
	return NewEntry(r).WithContext(ctx).Consume(opts...)
}

// ConsumeFn 回调消费
func (r *RabbitMQ) ConsumeFn(ctx context.Context, fn Callback, opts ...EntryOption) error {
	return NewEntry(r).WithContext(ctx).ConsumeFn(fn, opts...)
}
