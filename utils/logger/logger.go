package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ILogger interface {
	Error(...interface{})
	Info(...interface{})
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

type LocalLogger struct{}

func (l LocalLogger) Info(msg ...interface{}) {
	fmt.Printf("%s.%03d [info] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Nanosecond()/1e6, l.TrimSlice(msg...))
}
func (l LocalLogger) Debug(msg ...interface{}) {
	fmt.Printf("%s.%03d [debug] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Nanosecond()/1e6, l.TrimSlice(msg...))
}

func (l LocalLogger) Warn(msg ...interface{}) {
	fmt.Printf("%s.%03d [warn] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Nanosecond()/1e6, l.TrimSlice(msg...))
}

func (l LocalLogger) Error(msg ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s.%03d [error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Nanosecond()/1e6, l.TrimSlice(msg...))
}

func (l LocalLogger) Fatal(msg ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s.%03d [fatal] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Nanosecond()/1e6, l.TrimSlice(msg...))
}
func (l LocalLogger) TrimSlice(msg ...any) string {
	return strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%v", msg), "["), "]")
}
