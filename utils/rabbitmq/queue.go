package rabbitmq

import "github.com/streadway/amqp"

type Queue struct {
	Error  error
	Option *QueueOption
	ch     *amqp.Channel
	qu     *amqp.Queue
}

type QueueOption struct {
	QueueName  string // queue 名
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func NewQueueOption() *QueueOption {
	return &QueueOption{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
}

func (q *Queue) WithQueueName(name string) *Queue {
	q.Option.QueueName = name
	return q
}

func (q *Queue) WithArgs(args map[string]interface{}) *Queue {
	q.Option.Args = args
	return q
}

func (q *Queue) Declare() *Queue {
	if q.Error != nil {
		return q
	}
	queue, err := q.ch.QueueDeclare(
		q.Option.QueueName,  //队列名称
		q.Option.Durable,    //是否持久化
		q.Option.AutoDelete, //是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		q.Option.Exclusive,  //是否具有排他性
		q.Option.NoWait,     //是否阻塞
		q.Option.Args,       //额外信息
	)
	q.qu = &queue
	q.Error = err
	return q
}

func (q *Queue) Bind(exchangeName, routingKey string) *Queue {
	if q.Error != nil {
		return q
	}
	q.Error = q.ch.QueueBind(q.Option.QueueName, routingKey, exchangeName, false, nil)
	return q
}
