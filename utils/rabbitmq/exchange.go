package rabbitmq

import "github.com/streadway/amqp"

type Exchange struct {
	Error  error
	Option *ExchangeOption
	ch     *amqp.Channel
}

type ExchangeOption struct {
	ExchangeName string // queue 名
	Kind         string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Args         amqp.Table
}

func NewExchangeOption() *ExchangeOption {
	return &ExchangeOption{
		Durable:    true,
		Kind:       amqp.ExchangeDirect,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
}

func (e *Exchange) WithExchangeName(name string) *Exchange {
	e.Option.ExchangeName = name
	return e
}
func (e *Exchange) WithKind(kind string) *Exchange {
	e.Option.Kind = kind
	return e
}

func (e *Exchange) Declare() *Exchange {
	if e.Error != nil {
		return e
	}
	e.Error = e.ch.ExchangeDeclare(
		e.Option.ExchangeName, //交换机名称
		e.Option.Kind,         //
		e.Option.Durable,      //是否持久化
		e.Option.AutoDelete,   //是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		e.Option.Internal,     //
		e.Option.NoWait,       //是否阻塞
		e.Option.Args,         //额外信息
	)
	return e
}

func (e *Exchange) Bind(destination, key string) *Exchange {
	if e.Error != nil {
		return e
	}
	e.Error = e.ch.ExchangeBind(destination, key, e.Option.ExchangeName, false, nil)
	return e
}
