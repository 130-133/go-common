package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Option func(*Conf)
type Conf struct {
	URI string
}

func (o Option) Apply(arg *Conf) {
	o(arg)
}

func WithURI(uri string) Option {
	return func(arg *Conf) {
		arg.URI = uri
	}
}

func NewMongoDb(opts ...Option) *mongo.Client {
	c := &Conf{}
	for _, o := range opts {
		o.Apply(c)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.URI))
	if err != nil {
		panic(err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return mongoClient
}
