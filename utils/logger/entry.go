package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/errorx"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/help"
)

type MEntry struct {
	logger    *logrus.Logger
	errLogger *logrus.Logger
	body      map[string]interface{}
	ctx       context.Context
	err       error
}

func NewMEntry(loggers ...*logrus.Logger) *MEntry {
	var (
		logger    *logrus.Logger
		errLogger *logrus.Logger
	)
	logger = loggers[0]
	errLogger = loggers[0]
	if len(loggers) > 1 {
		errLogger = loggers[1]
	}
	return &MEntry{
		logger:    logger,
		errLogger: errLogger,
		body:      make(map[string]interface{}),
		ctx:       context.Background(),
	}
}

func (e *MEntry) WithCode(data int) *MEntry {
	e.body["code"] = data
	return e
}

func (e *MEntry) WithReq(data interface{}) *MEntry {
	e.body["req"] = help.ToString(data)
	return e
}

func (e *MEntry) WithResp(data interface{}) *MEntry {
	e.body["resp"] = help.ToString(data)
	return e
}

func (e *MEntry) WithTrack(data interface{}) *MEntry {
	e.body["track_data"] = data
	return e
}
func (e *MEntry) WithTracks(kv ...string) *MEntry {
	track := make(map[string]string)
	length := len(kv)
	for n := 0; n < length; n += 2 {
		track[kv[n]] = kv[n+1]
	}
	if len(kv)%2 != 0 {
		k := kv[length-1]
		track[k] = ""
	}
	e.body["track_data"] = track
	return e
}
func (e *MEntry) WithUin(data string) *MEntry {
	e.body["uin"] = data
	return e
}

func (e *MEntry) WithError(err error) *MEntry {
	e.err = err
	return e
}

func (e *MEntry) WithCtx(ctx context.Context) *MEntry {
	return e.WithContext(ctx)
}
func (e *MEntry) WithContext(ctx context.Context) *MEntry {
	e.ctx = ctx
	return e
}
func (e *MEntry) WithField(key string, value interface{}) *MEntry {
	e.body[key] = value
	return e
}

// LinePrev 记录行向上偏移
func (e *MEntry) LinePrev(prev int) *MEntry {
	// 减去LinePrev方法自身1层
	_, file, line := e.CallFrom(prev - 2)
	e.body["line"] = fmt.Sprintf("%s/%s:%d", path.Base(path.Dir(file)), path.Base(file), line)
	return e
}

// BeforeOut 最后输出前处理
func (e *MEntry) beforeOut() *MEntry {
	// 空时自动提取context里的uin
	if _, ok := e.body["uin"]; !ok {
		uin, _ := help.GetUinFromCtx(e.ctx)
		e.body["uin"] = uin
	}
	if e.err != nil {
		err := errorx.ParseErr(e.err)
		e.body["code"] = err.Code()
		e.body["error"] = err.Message()
	}
	return e
}

func (e *MEntry) Debug(data ...interface{}) {
	e.beforeOut().logger.WithContext(e.ctx).WithFields(e.body).Debugln(data...)
}
func (e *MEntry) Info(data ...interface{}) {
	e.beforeOut().logger.WithContext(e.ctx).WithFields(e.body).Infoln(data...)
}
func (e *MEntry) Warn(data ...interface{}) {
	e.beforeOut().logger.WithContext(e.ctx).WithFields(e.body).Warnln(data...)
}
func (e *MEntry) Error(data ...interface{}) {
	e.beforeOut().errLogger.WithContext(e.ctx).WithError(e.err).WithFields(e.body).Errorln(data...)
}
func (e *MEntry) Fatal(data ...interface{}) {
	e.beforeOut().errLogger.WithContext(e.ctx).WithError(e.err).WithFields(e.body).Fatalln(data...)
}

// CallFrom 获取调用方信息
func (e *MEntry) CallFrom(offset int) (funcName, file string, line int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(0, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line = f.FileLine(pc[i])
		if strings.Contains(f.Name(), "(*MEntry).CallFrom") {
			i += 4 + offset
			f = runtime.FuncForPC(pc[i])
			file, line = f.FileLine(pc[i])
			funcName = f.Name()
			break
		}
	}
	return
}
