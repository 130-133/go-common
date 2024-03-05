package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/help"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type MLogger struct {
	Mode      string
	Rotate    *MRotate
	Body      *LogBody
	Level     logrus.Level
	Sync      sync.Mutex
	Logger    *logrus.Logger
	ErrLogger *logrus.Logger
}

type LogBody struct {
	Sender    string      `json:"sender"`
	TraceId   string      `json:"trace_id"`
	SpanId    string      `json:"span_id"`
	Level     string      `json:"level"`
	Code      int         `json:"code"`
	Line      string      `json:"line"`
	Msg       string      `json:"msg"`
	Time      string      `json:"time"`
	Uin       string      `json:"uin"`
	Req       interface{} `json:"req"`
	Resp      interface{} `json:"resp"`
	TrackData interface{} `json:"track_data"`
}

type Option func(*MLogger)

func (l Option) Apply(log *MLogger) {
	l(log)
}

func WithMode(mode string) Option {
	return func(m *MLogger) {
		m.Mode = mode
	}
}
func WithPathStr(data string) Option {
	return func(m *MLogger) {
		m.Rotate.PathStr = data
	}
}
func WithLevel(data logrus.Level) Option {
	return func(m *MLogger) {
		m.Level = data
	}
}
func WithRotationTime(data time.Duration) Option {
	return func(m *MLogger) {
		m.Rotate.RotationTime = data
	}
}
func WithKeepDays(data time.Duration) Option {
	return func(m *MLogger) {
		m.Rotate.KeepDays = data
	}
}

func NewMLogger(serverName string, opts ...Option) *MLogger {
	m := &MLogger{
		Rotate: &MRotate{
			PathStr:      "./logs/",
			PathNameStr:  fmt.Sprintf("%s.log", serverName),
			RotationTime: 24 * time.Hour,
			KeepDays:     24 * time.Hour,
		},
		Body: &LogBody{
			Sender:    serverName,
			TrackData: make(map[string]interface{}),
		},
		Level: logrus.InfoLevel,
	}
	for _, opt := range opts {
		opt.Apply(m)
	}
	m.NewInfoLogger(serverName)
	m.NewErrLogger(serverName)
	return m
}

func (m *MLogger) NewInfoLogger(serverName string) *MLogger {
	r := &MRotate{
		PathStr:      m.Rotate.PathStr,
		PathNameStr:  fmt.Sprintf("%s.log", serverName),
		RotationTime: m.Rotate.RotationTime,
		KeepDays:     m.Rotate.KeepDays,
	}
	// 实例化
	logger := logrus.New()

	if ok, _ := help.InArray(m.Mode, []string{"dev", "local"}); ok {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		writer, err := r.New()
		if err != nil {
			panic(err)
		}
		logger.SetReportCaller(true)
		logger.SetOutput(writer)
		logger.SetLevel(m.Level)
	}
	logger.SetFormatter(m)
	m.Logger = logger
	return m
}

func (m *MLogger) NewErrLogger(serverName string) *MLogger {
	r := &MRotate{
		PathStr:      m.Rotate.PathStr,
		PathNameStr:  fmt.Sprintf("%s-error.log", serverName),
		RotationTime: m.Rotate.RotationTime,
		KeepDays:     m.Rotate.KeepDays,
	}
	// 实例化
	logger := logrus.New()
	if ok, _ := help.InArray(m.Mode, []string{"dev", "local"}); ok {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		writer, err := r.New()
		if err != nil {
			panic(err)
		}
		logger.SetReportCaller(true)
		logger.SetOutput(writer)
		logger.SetLevel(logrus.ErrorLevel)
	}
	logger.SetFormatter(m)
	m.ErrLogger = logger
	return m
}

func (m *MLogger) Debug(data ...interface{}) {
	m.NewEntry().Debug(data...)
}
func (m *MLogger) Info(data ...interface{}) {
	m.NewEntry().Info(data...)
}
func (m *MLogger) Warn(data ...interface{}) {
	m.NewEntry().Warn(data...)
}
func (m *MLogger) Error(data ...interface{}) {
	m.NewEntry().Error(data...)
}
func (m MLogger) Fatal(data ...interface{}) {
	m.NewEntry().Fatal(data...)
}

func (m *MLogger) NewEntry() *MEntry {
	return NewMEntry(m.Logger, m.ErrLogger)
}
func (m *MLogger) WithCode(data int) *MEntry {
	return m.NewEntry().WithCode(data)
}
func (m *MLogger) WithReq(data interface{}) *MEntry {
	return m.NewEntry().WithReq(data)
}
func (m *MLogger) WithResp(data interface{}) *MEntry {
	return m.NewEntry().WithResp(data)
}
func (m *MLogger) WithTrack(data interface{}) *MEntry {
	return m.NewEntry().WithTrack(data)
}
func (m *MLogger) WithUin(data string) *MEntry {
	return m.NewEntry().WithUin(data)
}
func (m *MLogger) WithCtx(ctx context.Context) *MEntry {
	return m.WithContext(ctx)
}
func (m *MLogger) WithContext(ctx context.Context) *MEntry {
	return m.NewEntry().WithContext(ctx)
}
func (m *MLogger) WithError(err error) *MEntry {
	return m.NewEntry().WithError(err)
}

func (m *MLogger) Format(entry *logrus.Entry) ([]byte, error) {
	m.Sync.Lock()
	defer m.Sync.Unlock()
	data := *m.Body
	code, _ := entry.Data["code"].(int)
	req, _ := entry.Data["req"]
	resp, _ := entry.Data["resp"]
	errStr, _ := entry.Data["error"].(string)
	uin, _ := entry.Data["uin"].(string)
	lineStr, _ := entry.Data["line"].(string)
	trackData, ok := entry.Data["track_data"]
	if !ok {
		trackData = make(map[string]interface{})
	}
	if errStr != "" {
		errStr = fmt.Sprintf(" Error: %s", errStr)
	}
	if lineStr == "" {
		_, file, line := m.CallFrom(2)
		lineStr = fmt.Sprintf("%s/%s:%d", path.Base(path.Dir(file)), path.Base(file), line)
	}

	data.Code = code
	data.Req = req
	data.Resp = resp
	data.Uin = uin
	data.TrackData = trackData
	data.Level = entry.Level.String()
	data.Line = lineStr
	data.SpanId = spanIdFromContext(entry.Context)
	data.TraceId = traceIdFromContext(entry.Context)
	data.Msg = help.ToString(entry.Message) + errStr
	data.Time = fmt.Sprintf("%s.%03d", entry.Time.Format(help.FormatTime), entry.Time.Nanosecond()/1e6)
	msg, err := json.Marshal(data)
	msg = []byte(fmt.Sprintf("%s\n", msg))
	return msg, err
}

// CallFrom 获取调用方信息
func (m *MLogger) CallFrom(offset int) (funcName, file string, line int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(0, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line = f.FileLine(pc[i])
		if strings.Contains(f.Name(), "(*Entry).write") {
			i += 4 + offset
			f = runtime.FuncForPC(pc[i])
			file, line = f.FileLine(pc[i])
			funcName = f.Name()
			break
		}
	}
	return
}

// spanIdFromContext 提取Tracer SpanId
func spanIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

// spanIdFromContext 提取Tracer TraceId
func traceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
