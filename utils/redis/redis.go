package redis

import (
	"context"

	red "github.com/go-redis/redis"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/errorx"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/logger"
)

type Option func(*Conf)
type Conf struct {
	Address string
	Pwd     string
	Db      int
	Logger  logger.ILogger
}

func (o Option) Apply(arg *Conf) {
	o(arg)
}

func WithAddress(data string) Option {
	return func(arg *Conf) {
		arg.Address = data
	}
}
func WithPwd(data string) Option {
	return func(arg *Conf) {
		arg.Pwd = data
	}
}
func WithDb(data int) Option {
	return func(arg *Conf) {
		arg.Db = data
	}
}
func WithLogger(logger logger.ILogger) Option {
	return func(arg *Conf) {
		arg.Logger = logger
	}
}

type MRedis struct {
	*red.Client
	Logger logger.ILogger
}
type MPipeline struct {
	red.Pipeliner
}
type MPop struct {
	*MRedis
	ctx context.Context
	key string
	ch  chan []string
}

func NewRedisConn(opts ...Option) *MRedis {
	opt := Conf{
		Address: "127.0.0.1:6379",
		Pwd:     "",
		Db:      0,
		Logger:  logger.LocalLogger{},
	}
	for _, o := range opts {
		o.Apply(&opt)
	}
	rdb := red.NewClient(&red.Options{
		Addr:     opt.Address,
		Password: opt.Pwd, // no password set
		DB:       opt.Db,  // use default DB
	})
	if rdb == nil && rdb.Ping().Err() != nil {
		panic(errorx.NewCacheError("连接Redis失败", 0))
	}
	return &MRedis{
		Client: rdb,
		Logger: opt.Logger,
	}
}
