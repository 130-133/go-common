package dao

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"gitea.com/llm-PhotoMagic/go-common/utils/help"
)

type IBaseDao interface {
	New()
	SetCtx(context.Context)
	GetCtx() context.Context
	Cancel(c context.Context)
	WithContext(ctx context.Context) TransactionCtx
	SetTimeOut(timeout time.Duration)
	GetMysql() *gorm.DB
	GetMongo() *mongo.Database
	TxMysql(ctx context.Context) *gorm.DB
	NewDB() *gorm.DB
	Transaction(fc func(context.Context) error)
}

type BaseDao struct {
	mysql   *gorm.DB
	mongo   *mongo.Database
	ctx     context.Context
	cancel  context.CancelFunc
	timeout time.Duration
	group   sync.Map
	txKey   string
}

func NewDao(mysql *gorm.DB, mongo *mongo.Database) *BaseDao {
	return &BaseDao{
		mysql: mysql,
		mongo: mongo,
		ctx:   context.Background(),
		txKey: "gorm-tx",
	}
}

// SetTimeOut set timeout
func (d *BaseDao) SetTimeOut(timeout time.Duration) {
	d.ctx, d.cancel = context.WithTimeout(context.Background(), timeout)
	d.timeout = timeout
}

// SetCtx set context
func (d *BaseDao) SetCtx(c context.Context) {
	if c != nil {
		d.ctx = c
	}
}

// GetCtx get context
func (d *BaseDao) GetCtx() context.Context {
	return d.ctx
}

// Cancel cancel context
func (d *BaseDao) Cancel(c context.Context) {
	d.cancel()
}

// GetMysql get gorm.DB info
func (d *BaseDao) GetMysql() *gorm.DB {
	return d.mysql
}

// GetMongo 返回Mongo客户端
func (d *BaseDao) GetMongo() *mongo.Database {
	return d.mongo
}

// UpdateDB update gorm.DB info
func (d *BaseDao) UpdateDB(db *gorm.DB) {
	d.mysql = db
}
func (d *BaseDao) UpdateMongo(db *mongo.Database) {
	d.mongo = db
}

// TxMysql 兼顾logic层事务
func (d *BaseDao) TxMysql(ctx context.Context) *gorm.DB {
	tx := ctx.Value(d.txKey)
	if tx != nil && tx.(*gorm.DB) != nil {
		return tx.(*gorm.DB).WithContext(ctx)
	}
	return d.mysql.WithContext(ctx)
}

// Transaction Mysql事务
func (d *BaseDao) Transaction(fc func(context.Context) error) {
	d.NewDB().Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(d.ctx, d.txKey, tx)
		return fc(ctx)
	})
}

type TransactionCtx struct {
	context.Context
	*BaseDao
}

func (d *BaseDao) WithContext(ctx context.Context) TransactionCtx {
	return TransactionCtx{
		Context: ctx,
		BaseDao: d,
	}
}

func (t TransactionCtx) Transaction(fc func(context.Context) error) error {
	return t.NewDB().Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(t.Context, t.txKey, tx)
		return fc(ctx)
	})
}

//
//func (d *BaseDao) Session(fc func(ctx mongo.SessionContext) error) {
//	d.GetMongo().Client().UseSession(d.ctx, func(sessionContext mongo.SessionContext) error {
//		return fc(sessionContext)
//	})
//}

// New new gorm.新gorm,重置条件
func (d *BaseDao) New() {
	d.mysql = d.NewDB()
}

// NewDB new gorm.新gorm
func (d *BaseDao) NewDB() *gorm.DB {
	return d.mysql.Session(&gorm.Session{NewDB: true, Context: d.ctx})
}

// ScopePage 加工分页参数
func (d *BaseDao) ScopePage(page IPage) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(page.GetOrderItemsString()) > 0 {
			db = db.Order(page.GetOrderItemsString())
		}
		return db.Limit(int(page.GetSize())).Offset(int(page.Offset()))
	}
}

// ScopeFilter 加工where参数，可支持后缀s字段进行in查询
func (d *BaseDao) ScopeFilter(filter map[string]interface{}, skipInFields []string) func(*gorm.DB) *gorm.DB {
	values := make([]interface{}, 0)
	fields := make([]string, 0)
	for key, val := range filter {
		if res, _ := help.InArray(key, skipInFields); !res && strings.HasSuffix(key, "s") {
			fields = append(fields, fmt.Sprintf("%s IN ?", strings.TrimRight(key, "s")))
		} else {
			fields = append(fields, fmt.Sprintf("%s = ?", key))
		}
		values = append(values, val)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(strings.Join(fields, " AND "), values...)
	}
}

// Condition 自定义sql查询
type Condition struct {
	list []*conditionInfo
}

func (c *Condition) AndWithCondition(condition bool, column string, cases string, value interface{}) *Condition {
	if condition {
		c.list = append(c.list, &conditionInfo{
			andor:  "and",
			column: column, // 列名
			case_:  cases,  // 条件(and,or,in,>=,<=)
			value:  value,
		})
	}
	return c
}

// And a Condition by and .and 一个条件
func (c *Condition) And(column string, cases string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, cases, value)
}

func (c *Condition) OrWithCondition(condition bool, column string, cases string, value interface{}) *Condition {
	if condition {
		c.list = append(c.list, &conditionInfo{
			andor:  "or",
			column: column, // 列名
			case_:  cases,  // 条件(and,or,in,>=,<=)
			value:  value,
		})
	}
	return c
}

// Or a Condition by or .or 一个条件
func (c *Condition) Or(column string, cases string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, cases, value)
}

func (c *Condition) Get() (where string, out []interface{}) {
	firstAnd := -1
	for i := 0; i < len(c.list); i++ { // 查找第一个and
		if c.list[i].andor == "and" {
			where = fmt.Sprintf("`%v` %v ?", c.list[i].column, c.list[i].case_)
			out = append(out, c.list[i].value)
			firstAnd = i
			break
		}
	}

	if firstAnd < 0 && len(c.list) > 0 { // 补刀
		where = fmt.Sprintf("`%v` %v ?", c.list[0].column, c.list[0].case_)
		out = append(out, c.list[0].value)
		firstAnd = 0
	}

	for i := 0; i < len(c.list); i++ { // 添加剩余的
		if firstAnd != i {
			where += fmt.Sprintf(" %v `%v` %v ?", c.list[i].andor, c.list[i].column, c.list[i].case_)
			out = append(out, c.list[i].value)
		}
	}

	return
}

type conditionInfo struct {
	andor  string
	column string // 列名
	case_  string // 条件(in,>=,<=)
	value  interface{}
}
