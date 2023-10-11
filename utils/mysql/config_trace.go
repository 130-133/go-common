package mysql

import (
	"fmt"

	ztrace "github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (o Conf) reportRawSql(db *gorm.DB) (err error) {
	if !o.TraceOn {
		return
	}
	if o.TraceLevel|Select == o.TraceLevel {
		err = db.Callback().Row().Before("gorm:row").Register("before:row:tracing", o.tracingBefore("db.query"))
		err = db.Callback().Row().After("gorm:row").Register("after:row:tracing", o.tracingAfter())
		err = db.Callback().Query().Before("gorm:query").Register("before:query:tracing", o.tracingBefore("db.query"))
		err = db.Callback().Query().After("gorm:query").Register("after:query:tracing", o.tracingAfter())
	}
	if o.TraceLevel|Insert == o.TraceLevel {
		err = db.Callback().Create().Before("gorm:create").Register("before:create:tracing", o.tracingBefore("db.create"))
		err = db.Callback().Create().After("gorm:create").Register("after:create:tracing", o.tracingAfter())
	}
	if o.TraceLevel|Update == o.TraceLevel {
		err = db.Callback().Update().Before("gorm:update").Register("before:update:tracing", o.tracingBefore("db.update"))
		err = db.Callback().Update().After("gorm:update").Register("after:update:tracing", o.tracingAfter())
	}
	if o.TraceLevel|Delete == o.TraceLevel {
		err = db.Callback().Delete().Before("gorm:delete").Register("before:delete:tracing", o.tracingBefore("db.delete"))
		err = db.Callback().Delete().After("gorm:delete").Register("after:delete:tracing", o.tracingAfter())
	}
	if o.TraceLevel|Raw == o.TraceLevel {
		err = db.Callback().Raw().Before("gorm:raw").Register("before:raw:tracing", o.tracingBefore("db.raw"))
		err = db.Callback().Raw().After("gorm:raw").Register("after:raw:tracing", o.tracingAfter())
	}
	return
}

// tracingBefore 执行前
func (o Conf) tracingBefore(cmd string) func(*gorm.DB) {
	tracer := otel.Tracer(ztrace.TraceName)
	return func(db *gorm.DB) {
		spanName := fmt.Sprintf("%s.%s", cmd, db.Statement.Table)
		ctx, _ := tracer.Start(db.Statement.Context, spanName)
		db.Statement.Context = ctx
	}
}

// tracingAfter 执行后
func (o Conf) tracingAfter() func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context == nil {
			return
		}
		span := trace.SpanFromContext(db.Statement.Context)
		sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		span.SetAttributes(attribute.String("db.statement", sql))
		if db.Error != nil {
			span.SetStatus(codes.Error, db.Error.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}
