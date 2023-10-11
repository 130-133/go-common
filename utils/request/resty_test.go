package request

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go func() {
		http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(request.Body)
			result := map[string]interface{}{
				"code":    0,
				"message": "success",
				"data": map[string]interface{}{
					"reqBody":   buffer.String(),
					"reqPath":   request.RequestURI,
					"reqHeader": request.Header,
				},
			}
			str, _ := json.Marshal(result)
			writer.Write(str)
		})
		http.ListenAndServe("127.0.0.1:23456", nil)
	}()
	m.Run()
}

func Get(opt ...ReqOption) IResponses {
	return NewRestyReq(context.Background()).Get("http://127.0.0.1:23456", "/test?b=2", map[string]interface{}{"a": 1}, opt...)
}
func Post(opt ...ReqOption) IResponses {
	return NewRestyReq(context.Background()).Post("http://127.0.0.1:23456", "/test", map[string]interface{}{"a": 1}, opt...)
}

func TestNewRestyReq(t *testing.T) {
	NewRestyReq(context.Background())
}

func TestWithOpts(t *testing.T) {
	resp := Get(
		WithCtx(context.Background()),
		WithTimeout(1*time.Second),
		WithHeader("X-AAA", "13"),
	)
	if resp.GetError() != nil {
		t.Error(resp.GetError())
		return
	}
	t.Logf("%+v", string(resp.GetBody()))
}

func TestRequests_Get(t *testing.T) {
	resp := Get()
	if resp.GetError() != nil {
		t.Error(resp.GetError())
		return
	}
	t.Logf("%+v", string(resp.GetBody()))
	t.Log(resp.GetStatusCode())
	t.Logf("%+v", resp.GetRequest())
	t.Logf("%+v", resp.GetResponse())
}

func TestRequests_Post(t *testing.T) {
	resp := Post()
	if resp.GetError() != nil {
		t.Error(resp.GetError())
		return
	}
	t.Logf("%+v", string(resp.GetBody()))
	t.Logf("%+v", resp.GetRequest())
	t.Logf("%+v", resp.GetResponse())
}

func TestResponses_GetData(t *testing.T) {
	resp := Get()
	if resp.GetError() != nil {
		t.Error(resp.GetError())
		return
	}
	t.Logf("%+v", string(resp.GetBody()))
	data := make(map[string]interface{})
	resp.GetData(&data)
	t.Log(data)
}

func TestResponses_GetJson(t *testing.T) {
	resp := Get()
	if resp.GetError() != nil {
		t.Error(resp.GetError())
		return
	}
	t.Logf("%+v", string(resp.GetBody()))
	data := make(map[string]interface{})
	resp.GetJson(&data)
	t.Log(data)
}
