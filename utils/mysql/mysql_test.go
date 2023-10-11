package mysql

import "testing"

func TestNewMysqlConn(t *testing.T) {
	NewMysqlConn(WithMaxIdleConn(10), WithConnMaxLifetime(10), WithMaxOpenConn(10))
}

func TestWithTraceLevel(t *testing.T) {
	NewMysqlConn(WithTraceOn(), WithTraceLevel(Select|Insert))
}
