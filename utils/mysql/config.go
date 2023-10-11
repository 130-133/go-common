package mysql

import "time"

type Option func(*Conf)
type Conf struct {
	Dsn             string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime time.Duration
	TraceOn         bool
	TraceLevel      Level
}

type Level int

const (
	Select Level = 1 << iota
	Insert
	Update
	Delete
	Raw
)

func (o Option) Apply(arg *Conf) {
	o(arg)
}

// WithDsn 设置dsn
func WithDsn(data string) Option {
	return func(conf *Conf) {
		conf.Dsn = data
	}
}

// WithMaxIdleConn 设置最大空闲连接数
func WithMaxIdleConn(data int) Option {
	return func(conf *Conf) {
		conf.MaxIdleConn = data
	}
}

// WithMaxOpenConn 设置最大连接数
func WithMaxOpenConn(data int) Option {
	return func(conf *Conf) {
		conf.MaxOpenConn = data
	}
}

// WithConnMaxLifetime 设置连接最大生命周期
func WithConnMaxLifetime(data time.Duration) Option {
	return func(conf *Conf) {
		conf.ConnMaxLifetime = data
	}
}

// WithTraceOn 设置链路追踪开关
func WithTraceOn() Option {
	return func(conf *Conf) {
		conf.TraceOn = true
	}
}

// WithTraceLevel 设置链路追踪语句类别
func WithTraceLevel(data Level) Option {
	return func(conf *Conf) {
		conf.TraceLevel = data
	}
}
