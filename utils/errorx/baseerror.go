package errorx

import (
	"errors"
	"fmt"
	"gitea.com/llm-PhotoMagic/go-common/utils/i18n"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	//"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitea.com/llm-PhotoMagic/go-common/utils/help"
)

const UnknownMsgKey = "server.error"

var global *TyyError

type TyyErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type TyyError struct {
	SystemCode AppCode
	ErrMsgMap  map[int]string
	UnknownMsg string
}

type ErrOpt func(*TyyError)

func (o ErrOpt) Apply(e *TyyError) {
	o(e)
}

func Init(appCode AppCode, opts ...ErrOpt) {
	global = &TyyError{
		SystemCode: appCode,
		UnknownMsg: "server error please try later",
	}
	for _, opt := range opts {
		opt(global)
	}
}

// WithErrMsgMap 设置全局映射错误码信息
func WithErrMsgMap(data map[int]string) ErrOpt {
	return func(e *TyyError) {
		e.ErrMsgMap = data
	}
}

// WithLocalize 设置国际化
//func WithLocalize() ErrOpt {
//	return func(e *TyyError) {
//		//e.I18n = NewI18n(data, i18nFile, lang)
//		//e.I18n.UnknownMsg = e.UnknownMsg
//	}
//}

// WithUnknownMsg 设置默认未定义错误信息
func WithUnknownMsg(msg string) ErrOpt {
	return func(e *TyyError) {
		e.UnknownMsg = msg
		//if e.I18n != nil {
		//	e.I18n.UnknownMsg = msg
		//}
	}
}

func (e TyyError) GetMsg(code int) string {
	//if e.I18n != nil {
	//	return e.I18n.Tfd(code, nil)
	//}
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
	I18n        *i18n.I18n     // 国际化
}

func GetGlobal() *TyyError {
	if global == nil {
		Init(99)
	}
	return global
}

func formatCodeMessage(msg string, code int) string {
	return fmt.Sprintf("%s code:%d", msg, code)
}

func NewError(category CategoryCode, code int, msg, lang string) error {
	if code < 1000 {
		codeStr := fmt.Sprintf("%02d%02d%03d", GetGlobal().SystemCode, category, code)
		code, _ = strconv.Atoi(codeStr)
	}
	statusCode := ToStatusCode(category)
	in := i18n.NewI18n()
	in.SetLanguage(lang)
	return &TyyCodeError{
		GrpcStatus:  status.New(statusCode, formatCodeMessage(msg, code)),
		ErrMessage:  in.Tfd(msg, nil),
		ErrCategory: category,
		ErrCode:     code,
		I18n:        in,
	}
}

func NewSystemCodeError(code int, lang string) error {
	return NewError(SystemError, code, global.GetMsg(code), lang)
}
func NewParamCodeError(code int, lang string) error {
	return NewError(ParamError, code, global.GetMsg(code), lang)
}
func NewGetDataCodeError(code int, lang string) error {
	return NewError(GetDataError, code, global.GetMsg(code), lang)
}
func NewCacheCodeError(code int, lang string) error {
	return NewError(CacheError, code, global.GetMsg(code), lang)
}
func NewDbCodeError(code int, lang string) error {
	return NewError(DbError, code, global.GetMsg(code), lang)
}
func NewMqCodeError(code int, lang string) error {
	return NewError(MqError, code, global.GetMsg(code), lang)
}
func NewHttpCodeError(code int, lang string) error {
	return NewError(HttpError, code, global.GetMsg(code), lang)
}
func NewRpcCodeError(code int, lang string) error {
	return NewError(RpcError, code, global.GetMsg(code), lang)
}

func NewSystemError(msg string, code int, lang string) error {
	return NewError(SystemError, code, msg, lang)
}
func NewParamError(msg string, code int, lang string) error {
	return NewError(ParamError, code, msg, lang)
}
func NewGetDataError(msg string, code int, lang string) error {
	return NewError(GetDataError, code, msg, lang)
}
func NewCacheError(msg string, code int, lang string) error {
	return NewError(CacheError, code, msg, lang)
}
func NewDbError(msg string, code int, lang string) error {
	return NewError(DbError, code, msg, lang)
}
func NewMqError(msg string, code int, lang string) error {
	return NewError(MqError, code, msg, lang)
}
func NewHttpError(msg string, code int, lang string) error {
	return NewError(HttpError, code, msg, lang)
}
func NewRpcError(msg string, code int, lang string) error {
	return NewError(RpcError, code, msg, lang)
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
	var result *TyyCodeError
	if errors.As(err, &result) {
		return result
	}
	msg := err.Error()
	regex, _ := regexp.Compile(`([\s\S]*) code:(\d+)$`)
	if strings.HasPrefix(msg, "rpc error") {
		//eg:"rpc error: code = Unknown desc = 查询结果为空 code:2111007"
		regex, _ = regexp.Compile(`desc = ([\s\S]*) code:(\d+)$`)
	}
	match := regex.FindStringSubmatch(msg)
	//var result *TyyCodeError
	errors.As(NewSystemError(msg, 0, "en"), &result)
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
	return NewError(CategoryCode(categoryCode), errCode, sliceMsg, "en").(*TyyCodeError)
}

// HttpxHandler go-zero的http异常处理
func HttpxHandler(err error) (int, interface{}) {
	var e *TyyCodeError
	switch {
	case errors.As(err, &e):
		return http.StatusOK, HttpxErrMsgShow(e)
	default:
		tyyErr := ParseErr(err)
		if tyyErr != nil {
			return http.StatusOK, HttpxErrMsgShow(tyyErr)
		}
	}
	var initErr *TyyCodeError
	errors.As(NewSystemError(global.GetMsg(-1), 0, "en"), &initErr)
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
		if err.I18n != nil {
			result.Msg = err.I18n.Tfd(UnknownMsgKey, nil)
		}
	}
	return result
}
