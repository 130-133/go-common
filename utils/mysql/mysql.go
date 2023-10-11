package mysql

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlConn(opts ...Option) *gorm.DB {
	conf := Conf{
		Dsn:        "root:123456@tcp(127.0.0.1:3306)?charset=utf8mb4&parseTime=True&loc=Local",
		TraceLevel: Select | Insert | Update | Delete | Raw,
	}
	for _, o := range opts {
		o.Apply(&conf)
	}
	mysqlConf := mysql.Config{
		DSN: conf.Dsn,
	}
	db, err := gorm.Open(mysql.Open(mysqlConf.DSN), &gorm.Config{})
	if err != nil {
		panic(errors.New("连接Mysql数据库失败"))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(errors.New("获取Mysql数据库失败"))
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifetime)
	if err = conf.reportRawSql(db); err != nil {
		panic(errors.New("链路追踪设置失败"))
	}
	return db
}
