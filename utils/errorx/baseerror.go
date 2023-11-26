package errorx

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitea.com/llm-PhotoMagic/go-common/utils/help"
)

var global *tyyError

type TyyErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type tyyError struct {
	SystemCode AppCode
	ErrMsgMap  map[int]string
	UnknownMsg string
	I18n       *I18n
}

type ErrOpt func(*tyyError)

func (o ErrOpt) Apply(e *tyyError) {
	o(e)
}

func Init(appCode AppCode, opts ...ErrOpt) {
	global = &tyyError{
		SystemCode: appCode,
		UnknownMsg: "服务器开小差了，请稍后再试",
	}
	for _, opt := range opts {
		opt(global)
	}
}

// WithErrMsgMap 设置全局映射错误码信息
func WithErrMsgMap(data map[int]string) ErrOpt {
	return func(e *tyyError) {
		e.ErrMsgMap = data
	}
}

// WithLocalize 设置国际化
func WithLocalize(data map[int]*i18n.Message, i18nFile []string, lang string) ErrOpt {
	return func(e *tyyError) {
		e.I18n = NewI18n(data, i18nFile, lang)
		e.I18n.UnknownMsg = e.UnknownMsg
	}
}

// WithUnknownMsg 设置默认未定义错误信息
func WithUnknownMsg(msg string) ErrOpt {
	return func(e *tyyError) {
		e.UnknownMsg = msg
		if e.I18n != nil {
			e.I18n.UnknownMsg = msg
		}
	}
}

func (e tyyError) GetMsg(code int) string {
	if e.I18n != nil {
		return e.I18n.Msg(code)
	}
	msg, ok := e.ErrMsgMap[code]
	if !ok {
		msg = e.UnknownMsg
	}
	return msg
}

type TyyCodeError struct {
	GrpcStatus  *status.Status //grpc状态码
	ErrMessage  string         //错误信息
	ErrCategory CategoryCode   //分类
	ErrCode     int            //7位错误码
}

func GetGlobal() *tyyError {
	if global == nil {
		Init(99)
	}
	return global
}

func formatCodeMessage(msg string, code int) string {
	return fmt.Sprintf("%s code:%d", msg, code)
}

func NewError(category CategoryCode, code int, msg string) error {
	if code < 1000 {
		codeStr := fmt.Sprintf("%02d%02d%03d", GetGlobal().SystemCode, category, code)
		code, _ = strconv.Atoi(codeStr)
	}
	statusCode := ToStatusCode(category)
	return &TyyCodeError{
		GrpcStatus:  status.New(statusCode, formatCodeMessage(msg, code)),
		ErrMessage:  msg,
		ErrCategory: category,
		ErrCode:     code,
	}
}
func NewSystemCodeError(code int) error {
	return NewError(SystemError, code, global.GetMsg(code))
}
func NewParamCodeError(code int) error {
	return NewError(ParamError, code, global.GetMsg(code))
}
func NewGetDataCodeError(code int) error {
	return NewError(GetDataError, code, global.GetMsg(code))
}
func NewCacheCodeError(code int) error {
	return NewError(CacheError, code, global.GetMsg(code))
}
func NewDbCodeError(code int) error {
	return NewError(DbError, code, global.GetMsg(code))
}
func NewMqCodeError(code int) error {
	return NewError(MqError, code, global.GetMsg(code))
}
func NewHttpCodeError(code int) error {
	return NewError(HttpError, code, global.GetMsg(code))
}
func NewRpcCodeError(code int) error {
	return NewError(RpcError, code, global.GetMsg(code))
}

func NewSystemError(msg string, code int) error {
	return NewError(SystemError, code, msg)
}
func NewParamError(msg string, code int) error {
	return NewError(ParamError, code, msg)
}
func NewGetDataError(msg string, code int) error {
	return NewError(GetDataError, code, msg)
}
func NewCacheError(msg string, code int) error {
	return NewError(CacheError, code, msg)
}
func NewDbError(msg string, code int) error {
	return NewError(DbError, code, msg)
}
func NewMqError(msg string, code int) error {
	return NewError(MqError, code, msg)
}
func NewHttpError(msg string, code int) error {
	return NewError(HttpError, code, msg)
}
func NewRpcError(msg string, code int) error {
	return NewError(RpcError, code, msg)
}

// Error 默认输出message带code
func (e *TyyCodeError) Error() string {
	return e.ErrMessage
}

func (e *TyyCodeError) GRPCStatus() *status.Status {
	return e.GrpcStatus
}

// Message 返回带code message
func (e *TyyCodeError) Message() string {
	return e.GRPCStatus().Message()
}

// Code 返回自定义错误码
func (e *TyyCodeError) Code() int {
	return e.ErrCode
}

func (e *TyyCodeError) Category() CategoryCode {
	return e.ErrCategory
}

// Data 返回http结构
func (e *TyyCodeError) Data() *TyyErrorResponse {
	return &TyyErrorResponse{
		Code: e.Code(),
		Msg:  e.Error(),
	}
}

// IsCode 判断两个错误是否一致（只针对code一致）
func (e *TyyCodeError) IsCode(err error) bool {
	errA := ParseErr(err)
	if e.Code() == errA.Code() {
		return true
	}
	return false
}

// IsErr 判断两个错误是否一致 （针对code和message一致）
func (e *TyyCodeError) IsErr(err error) bool {
	errA := ParseErr(err)
	if e.Code() == errA.Code() && e.Message() == errA.Message() {
		return true
	}
	return false
}

// CompCode 对比错误码
func CompCode(errA, errB error) bool {
	return ParseErr(errA).IsCode(errB)
}

// CompErr 对比错误
func CompErr(errA, errB error) bool {
	return ParseErr(errA).IsErr(errB)
}

func ToStatusCode(category CategoryCode) codes.Code {
	var statusCode codes.Code
	switch category {
	case SystemError:
		statusCode = codes.Unavailable
	case ParamError:
		statusCode = codes.InvalidArgument
	case GetDataError:
		statusCode = codes.DataLoss
	case CacheError, DbError, MqError, HttpError, RpcError:
		statusCode = codes.FailedPrecondition
	}
	return statusCode
}

// ParseErr 解析GRPC返回错误
func ParseErr(err error) *TyyCodeError {
	if err == nil {
		return nil
	}
	if result, ok := err.(*TyyCodeError); ok {
		return result
	}
	msg := err.Error()
	regex, _ := regexp.Compile(`([\s\S]*) code:(\d+)$`)
	if strings.HasPrefix(msg, "rpc error") {
		//eg:"rpc error: code = Unknown desc = 查询结果为空 code:2111007"
		regex, _ = regexp.Compile(`desc = ([\s\S]*) code:(\d+)$`)
	}
	match := regex.FindStringSubmatch(msg)
	result := NewSystemError(msg, 0).(*TyyCodeError)
	if len(match) != 3 {
		return result
	}
	sliceMsg := match[1]
	sliceCode := match[2]
	if len(sliceCode) != 7 {
		return result
	}
	errCode, cErr := strconv.Atoi(sliceCode)
	if cErr != nil {
		return result
	}
	categoryCode, cErr := strconv.Atoi(sliceCode[2:4])
	if cErr != nil {
		return nil
	}
	return NewError(CategoryCode(categoryCode), errCode, sliceMsg).(*TyyCodeError)
}

// HttpxHandler gozero的http异常处理
func HttpxHandler(err error) (int, interface{}) {
	switch e := err.(type) {
	case *TyyCodeError:
		return http.StatusOK, HttpxErrMsgShow(e)
	default:
		tyyErr := ParseErr(err)
		if tyyErr != nil {
			return http.StatusOK, HttpxErrMsgShow(tyyErr)
		}
	}

	fmt.Errorf("SetErrorHandler Err:%s Stack:%s", err.Error(), debug.Stack())

	initErr := NewSystemError(global.GetMsg(-1), 0).(*TyyCodeError)
	if global.I18n != nil {
		initErr = NewSystemError(global.I18n.Msg(-1), 0).(*TyyCodeError)
	}
	return http.StatusInternalServerError, initErr.Data()
}

func HttpxErrMsgShow(err *TyyCodeError) *TyyErrorResponse {
	result := err.Data()
	codeStr := help.ToString(err.Code())
	if len(codeStr) != 7 {
		return result
	}
	switch err.ErrCategory {
	case SystemError, DbError, MqError, HttpError, RpcError, GetDataError, CacheError:
		result.Msg = global.GetMsg(-1)
		if global.I18n != nil {
			result.Msg = global.I18n.Msg(-1)
		}
	}
	return result
}
