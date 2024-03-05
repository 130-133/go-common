package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/help"
	"gitlab.darmod.cn/llm-PhotoMagic/go-common/utils/logger"
)

type IHttpLog interface {
	Interceptor(next http.HandlerFunc) http.HandlerFunc
}

type log struct {
	*logger.MLogger
	ignore []string
}

// HttpLog API请求日志
func HttpLog(l *logger.MLogger) IHttpLog {
	return log{
		MLogger: l,
		ignore: []string{
			"/ping",
			"/checkhealth",
			"login",
		},
	}
}

func (l log) Interceptor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, _ := url.Parse(r.RequestURI)
		if l.ignorePath(urls.Path) {
			next(w, r)
			return
		}
		track := map[string]interface{}{
			"body":   help.GetRawBody(r),
			"query":  l.formatQuery(urls.Query()),
			"method": r.Method,
			"header": l.formatQuery(r.Header),
		}
		l.WithCtx(r.Context()).WithTrack(track).Info(fmt.Sprintf("请求参数 - %s", urls.Path))
		next(w, r)
	}
}

func (l log) ignorePath(path string) bool {
	for _, v := range l.ignore {
		if strings.HasSuffix(strings.ToLower(path), v) {
			return true
		}
	}
	return false
}

func (l log) formatQuery(query map[string][]string) map[string]string {
	var (
		data    []string
		mapData = make(map[string]string)
	)
	for k, v := range query {
		data = append(data, fmt.Sprintf("\"%s\":\"%s\"", k, v[0]))
	}
	str := fmt.Sprintf("{%s}", strings.Join(data, ","))
	_ = json.Unmarshal([]byte(str), &mapData)
	return mapData
}
