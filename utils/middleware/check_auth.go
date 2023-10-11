package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"

	"git.tyy.com/llm-PhotoMagic/go-common/utils/context/auth"
	"git.tyy.com/llm-PhotoMagic/go-common/utils/encrypt"
	"git.tyy.com/llm-PhotoMagic/go-common/utils/errorx"
	"git.tyy.com/llm-PhotoMagic/go-common/utils/help"
)

type IAuth interface {
	GetCheckNormalFun(next http.HandlerFunc, expire time.Duration) http.HandlerFunc
	GetCheckAuthFun(next http.HandlerFunc) http.HandlerFunc
}

type AuthKeys struct {
	JwtSecret string
	//AuthSecret  string
	OnlyExtract bool //纯提取可能存在的token，不校验
}

type AuthOpt func(*AuthKeys)

func (a AuthOpt) Apply(keys *AuthKeys) {
	a(keys)
}

// Auth 鉴权提供新旧密钥解析中间件
func Auth(opts ...AuthOpt) IAuth {
	//默认测试环境密钥
	a := AuthKeys{
		JwtSecret: "ecb0babf0687c6a7427e225f5a29b2ef",
		//AuthSecret: "c8c93222583741bd828579b3d3efd43b_1",
	}
	for _, opt := range opts {
		opt.Apply(&a)
	}
	return a
}

//func WithAuthSecret(secret string) AuthOpt {
//	return func(keys *AuthKeys) {
//		keys.AuthSecret = secret
//	}
//}

func WithJwtSecret(secret string) AuthOpt {
	return func(keys *AuthKeys) {
		keys.JwtSecret = secret
	}
}

func OnlyExtract() AuthOpt {
	return func(keys *AuthKeys) {
		keys.OnlyExtract = true
	}
}

func (a AuthKeys) GetCheckNormalFun(next http.HandlerFunc, expire time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//先获取head中是否存在Authorization
		hAuthorization := r.Header.Get("Authorization")
		split := strings.SplitN(hAuthorization, "Bearer ", 2)
		authorization := hAuthorization
		if len(split) > 1 {
			authorization = split[1]
		}
		en := encrypt.NewJwt(encrypt.JwtConfig{
			Token: a.JwtSecret,
		}).Decode(authorization)
		if !a.OnlyExtract && !en.Verify() {
			a.errOutput(w)
			return
		}
		data := en.Data()

		ctx := auth.InjectJwt(r.Context(), data)
		//ctx, _ = help.SetUinToMetadataCtx(ctx, id)
		if a.CheckAuthExpire(ctx, expire) {
			a.errOutput(w)
			return
		}
		next(w, r.WithContext(ctx))
	}
}

func (a AuthKeys) GetCheckAuthFun(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			uid string
			err error
		)
		//先获取head中是否存在Authorization
		hAuthorization := r.Header.Get("Authorization")
		if hAuthorization != "" {
			uid, err = authorizationFun(a, hAuthorization)
		}
		if hAuthorization == "" || err != nil {
			data := errorx.NewSystemError("登录状态已失效、请重新登录", 0).(*errorx.TyyCodeError).Data()
			body, _ := json.Marshal(data)
			w.Header().Set(httpx.ContentType, httpx.ApplicationJson)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(body)
			return
		}

		ctx := context.WithValue(r.Context(), "Uid", uid)
		ctx, _ = help.SetUinToMetadataCtx(ctx, uid)

		next(w, r.WithContext(ctx))
	}
}

// CheckAuthExpire 校验业务过期
func (a AuthKeys) CheckAuthExpire(ctx context.Context, duration time.Duration) bool {
	signedUnix := auth.GetSignedUnix(ctx)
	if signedUnix > 0 && help.ParseUnix(signedUnix).Add(duration).Before(time.Now()) {
		return true
	}
	return false
}

// 四小强校验
//func fourAuthCheck(authKey string, uin string, time int64, s2t int64, sign string) error {
//	var hash1 = md5.Sum([]byte(uin + authKey + strconv.Itoa(int(s2t))))
//	var hash2 = md5.Sum([]byte(strconv.Itoa(int(time)) + hex.EncodeToString(hash1[:]) + uin))
//	//fmt.Println(hex.EncodeToString(hash2[:]))
//	if hex.EncodeToString(hash2[:]) != sign {
//		return errors.New("鉴权失败")
//	}
//	return nil
//}

// post 四小强鉴权 第一形态 （旧）
//func postAuthFun(authKey string, body string) (uin string, err error) {
//	UnBody := gjson.Parse(body)
//	uin = UnBody.Get("uin").String()
//	time := UnBody.Get("time").Int()
//	s2t := UnBody.Get("s2t").Int()
//	sign := UnBody.Get("sign").String()
//
//	if s2t == 0 || time == 0 || uin == "" || sign == "" {
//		err = errors.New("鉴权参数校验无效")
//		return
//	}
//	if err = fourAuthCheck(authKey, uin, time, s2t, sign); err != nil {
//		return
//	}
//	return
//}

// header 鉴权
func authorizationFun(authKeys AuthKeys, authorization string) (uid string, err error) {
	authSplit := strings.SplitN(authorization, " ", 2)
	if len(authSplit) != 2 {
		err = errors.New("鉴权参数无效")
		return
	}
	authType, authToken := authSplit[0], authSplit[1]

	switch authType {
	case "Bearer":
		en := encrypt.NewJwt(encrypt.JwtConfig{
			Token: authKeys.JwtSecret,
		}).Decode(authToken)
		if !en.Verify() {
			err = errors.New("鉴权失败")
			return
		}
		data := en.Data()
		uid, _ = data["uid"].(string)
	default:
		err = errors.New("鉴权失败")
		return
	}
	return
}

func (a AuthKeys) errOutput(w http.ResponseWriter) {
	data := errorx.NewSystemError("登录状态已失效、请重新登录", 0).(*errorx.TyyCodeError).Data()
	body, _ := json.Marshal(data)
	w.Header().Set(httpx.ContentType, httpx.ApplicationJson)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(body)
}
