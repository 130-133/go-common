package dao

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type Option struct {
	FirstFind bool
}

type Options func(*Page)

func (os Options) Apply(o *Page) {
	os(o)
}

func WithSort(sort OrderItem) Options {
	return func(p *Page) {
		p.AddOrderItem(sort)
	}
}

func WithSize(size int64) Options {
	return func(p *Page) {
		p.SetSize(size)
	}
}

func WithCurrent(current int64) Options {
	return func(p *Page) {
		p.SetCurrent(current)
	}
}

func WithFields(fields []string) Options {
	return func(p *Page) {
		p.SetFields(fields)
	}
}

func WithFirstFind() Options {
	return func(p *Page) {
		p.opt.init().FirstFind  = true
	}
}

func (o *Option) init() *Option {
	if o == nil {
		*o = Option{}
	}
	return o
}

// PageScope mysql附加便捷范围限制
func PageScope(page IPage) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := page.Offset()
		size := page.GetSize()
		order := page.GetOrderItemsString()
		fields := page.GetFields()
		if page.Offset() > 0 {
			db = db.Offset(int(offset))
		}
		if page.GetSize() > 0 {
			db = db.Limit(int(size))
		}
		if order != "" {
			db = db.Order(page.GetOrderItemsString())
		}
		if len(fields) > 0 {
			db = db.Select(fields)
		}
		return db
	}
}

// PageOptions mongo附加便捷范围限制
func PageOptions(page IPage) *options.FindOptions {
	opt := options.Find()
	offset := page.Offset()
	size := page.GetSize()
	order := page.GetOrderItemsString()
	if offset > 0 {
		opt.SetSkip(offset)
	}
	if size > 0 {
		opt.SetLimit(size)
	}
	if order != "" {
		opt.SetSort(page.GetBsonSort())
	}
	return opt
}

// ApplyOpts 便捷接收参数返回Page
func ApplyOpts(opts ...Options) IPage {
	p := Page{}
	for _, o := range opts {
		o.Apply(&p)
	}
	return &p
}
