package request

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"reflect"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"gitea.com/llm-PhotoMagic/go-common/config/cfginit"
	"gitea.com/llm-PhotoMagic/go-common/utils/help"
)

type IRequests interface {
	Get(host, path string, params map[string]interface{}, opts ...ReqOption) IResponses
	Post(host, path string, params interface{}, opts ...ReqOption) IResponses
}

type IResponses interface {
	IsOk() bool
	GetStatusCode() int
	GetBody() []byte
	GetError() error
	GetJson(data interface{}) error
	GetData(data interface{}) error
	GetRequest() *http.Request
	GetResponse() *http.Response
}

type Requests struct {
	client  *resty.Client
	ctx     context.Context
	header  map[string]string
	timeout time.Duration
}

type Responses struct {
	statusCode int
	body       []byte
	err        error
	request    *http.Request
	response   *http.Response
}

var globalResty *resty.Client

func getResty(global bool) *resty.Client {
	if global {
		if globalResty == nil {
			globalResty = resty.New().
				OnBeforeRequest(Before()).OnAfterResponse(After()).OnError(Error()).
				SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		}
		return globalResty
	}
	return resty.New()
}

func NewRestyReq(ctx context.Context) IRequests {
	return &Requests{
		client: getResty(true),
		ctx:    help.NewCtxFromTraceCtx(ctx),
		header: map[string]string{
			"Content-type": cfginit.ContentTypeJson,
		},
		timeout: 30 * time.Second,
	}
}

func NewRestyResp() *Responses {
	return &Responses{}
}

type ReqOption func(*Requests)

func (r ReqOption) apply(o *Requests) {
	r(o)
}

func WithCtx(ctx context.Context) ReqOption {
	return func(r *Requests) {
		r.ctx = ctx
	}
}

// WithHeader 设置单个header值
func WithHeader(key, val string) ReqOption {
	return func(r *Requests) {
		r.header[key] = val
	}
}

// WithHeaders 设置多个header值
func WithHeaders(val map[string]string) ReqOption {
	return func(r *Requests) {
		r.header = val
	}
}

// WithTimeout 设置请求超时时间
func WithTimeout(val time.Duration) ReqOption {
	return func(r *Requests) {
		r.timeout = val
	}
}

func (resp *Responses) SetBody(body []byte) *Responses {
	resp.body = body
	return resp
}

// GetBody 获取响应原始数据
func (resp *Responses) GetBody() []byte {
	return resp.body
}

func (resp *Responses) SetError(err error) *Responses {
	resp.err = err
	return resp
}

func (resp *Responses) GetError() error {
	return resp.err
}

func (resp *Responses) SetStatusCode(code int) *Responses {
	resp.statusCode = code
	return resp
}

func (resp *Responses) GetStatusCode() int {
	return resp.statusCode
}

// GetJson 获取响应数据并json解析到形参
func (resp *Responses) GetJson(data interface{}) error {
	return json.Unmarshal(resp.body, data)
}

// GetData 获取响应数据并根据协定规范code，message，data判断结果，直接提取data内容到形参字段
func (resp *Responses) GetData(data interface{}) error {
	body := gjson.ParseBytes(resp.body)
	code := body.Get("code").Int()
	message := body.Get("message").String()
	if code != 0 {
		return errors.New(fmt.Sprintf("http Request Failed 'code':%d, 'message':%s", code, message))
	}
	bodyValue := body.Get("data")
	if !bodyValue.Exists() {
		return errors.New(fmt.Sprintf("http Request Result Not Found 'data'"))
	}
	byteData, _ := json.Marshal(bodyValue.Value())
	json.Unmarshal(byteData, data)
	return nil
}

func (resp *Responses) IsOk() bool {
	if resp.statusCode == http.StatusOK {
		return true
	}
	return false
}

func (resp *Responses) SetRequest(request *http.Request) *Responses {
	resp.request = request
	return resp
}
func (resp *Responses) GetRequest() *http.Request {
	return resp.request
}
func (resp *Responses) SetResponse(response *http.Response) *Responses {
	resp.response = response
	return resp
}
func (resp *Responses) GetResponse() *http.Response {
	return resp.response
}

func MapInterfaceToString(params map[string]interface{}) (result map[string]string) {
	result = make(map[string]string)
	for key, value := range params {
		val := reflect.ValueOf(value)
		str := ""
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = strconv.Itoa(int(val.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str = strconv.Itoa(int(val.Uint()))
		case reflect.Float32, reflect.Float64:
			str = fmt.Sprintf("%.2f", val.Float())
		case reflect.String:
			str = val.String()
		}
		result[key] = str
	}
	return result
}

func Before() resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		var (
			ctx  context.Context
			span trace.Span
		)
		tracer := otel.Tracer(ztrace.TraceName)

		urlParse, _ := url.Parse(request.URL)
		urlStr := fmt.Sprintf("%s%s", urlParse.Host, urlParse.Path)
		ctx, span = tracer.Start(request.Context(), urlStr, trace.WithSpanKind(trace.SpanKindClient))
		span.SetAttributes(attribute.String("http.url", request.URL))
		span.SetAttributes(attribute.String("http.method", request.Method))
		span.SetAttributes(attribute.String("http.query", request.QueryParam.Encode()))
		span.SetAttributes(attribute.String("http.body", help.ToString(request.Body)))
		span.SetAttributes(attribute.String("http.header", help.ToString(request.Header)))
		request.SetContext(ctx)
		return nil
	}
}
func After() resty.ResponseMiddleware {
	return func(client *resty.Client, response *resty.Response) error {
		span := trace.SpanFromContext(response.Request.Context())
		span.SetStatus(codes.Ok, "")
		span.End()
		return nil
	}
}
func Error() resty.ErrorHook {
	return func(request *resty.Request, err error) {
		span := trace.SpanFromContext(request.Context())
		span.SetStatus(codes.Error, err.Error())
		span.End()
	}
}

func (r *Requests) Get(host, path string, params map[string]interface{}, opts ...ReqOption) IResponses {
	for _, opt := range opts {
		opt.apply(r)
	}
	c := r.client.
		SetTimeout(r.timeout).
		R().SetContext(r.ctx)
	if params != nil {
		c.SetQueryParams(MapInterfaceToString(params))
	}
	if r.header != nil {
		c.SetHeaders(r.header)
	}
	resp, err := c.Get(fmt.Sprintf("%s%s", host, path))
	return NewRestyResp().
		SetStatusCode(resp.StatusCode()).
		SetError(err).
		SetRequest(c.RawRequest).
		SetResponse(resp.RawResponse).
		SetBody(resp.Body())
}

func (r *Requests) Post(host, path string, params interface{}, opts ...ReqOption) IResponses {
	for _, opt := range opts {
		opt.apply(r)
	}
	c := r.client.
		SetTimeout(r.timeout).
		R().SetContext(r.ctx)
	if params != nil {
		c.SetBody(params)
	}
	if r.header != nil {
		c.SetHeaders(r.header)
	}
	resp, err := c.Post(fmt.Sprintf("%s%s", host, path))
	return NewRestyResp().
		SetStatusCode(resp.StatusCode()).
		SetError(err).
		SetRequest(c.RawRequest).
		SetResponse(resp.RawResponse).
		SetBody(resp.Body())
}
