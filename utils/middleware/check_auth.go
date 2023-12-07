package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"gitea.com/llm-PhotoMagic/go-common/utils/context/header"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
	"time"

	"gitea.com/llm-PhotoMagic/go-common/utils/context/auth"
	"gitea.com/llm-PhotoMagic/go-common/utils/encrypt"
	"gitea.com/llm-PhotoMagic/go-common/utils/errorx"
	"gitea.com/llm-PhotoMagic/go-common/utils/help"
)

const UNAUTHORIZED = "server.unauthorized"
const UserInfoKey = "ctx-user"

type IAuth interface {
	GetCheckAuthFun(next http.HandlerFunc) http.HandlerFunc
}

type UserInfo struct {
	ID     int64
	Name   string
	AreaID int64
}

type AuthKeys struct {
	JwtSecret string
	//AuthSecret  string
	OnlyExtract bool //纯提取可能存在的token，不校验
	DefaultLang string
}

type AuthOpt func(*AuthKeys)

func (a AuthOpt) Apply(keys *AuthKeys) {
	a(keys)
}

// Auth 鉴权提供新旧密钥解析中间件
func Auth(opts ...AuthOpt) IAuth {
	//默认测试环境密钥
	a := AuthKeys{
		JwtSecret:   "ecb0babf0687c6a7427e225f5a29b2ef",
		DefaultLang: "en",
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

func WithLang(lang string) AuthOpt {
	return func(keys *AuthKeys) {
		keys.DefaultLang = lang
	}
}

func OnlyExtract() AuthOpt {
	return func(keys *AuthKeys) {
		keys.OnlyExtract = true
	}
}

//func (a AuthKeys) GetCheckNormalFun(next http.HandlerFunc, expire time.Duration) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		//先获取head中是否存在Authorization
//		hAuthorization := r.Header.Get("Authorization")
//		split := strings.SplitN(hAuthorization, "Bearer ", 2)
//		authorization := hAuthorization
//		if len(split) > 1 {
//			authorization = split[1]
//		}
//		en := encrypt.NewJwt(encrypt.JwtConfig{
//			Token: a.JwtSecret,
//		}).Decode(authorization)
//		if !a.OnlyExtract && !en.Verify() {
//			a.errOutput(w)
//			return
//		}
//		data := en.Data()
//
//		ctx := auth.InjectJwt(r.Context(), data)
//		//ctx, _ = help.SetUinToMetadataCtx(ctx, id)
//		if a.CheckAuthExpire(ctx, expire) {
//			a.errOutput(w)
//			return
//		}
//		next(w, r.WithContext(ctx))
//	}
//}

func (a AuthKeys) GetCheckAuthFun(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			//id   int64
			//name string
			err  error
			lang string
		)
		//先获取head中是否存在Authorization
		hAuthorization := r.Header.Get("Authorization")
		lang = header.GetLangFromCtx(r.Context())
		if lang == "" {
			lang = a.DefaultLang
		}
		user := UserInfo{}
		if hAuthorization != "" {
			err = authorizationFun(a, hAuthorization, &user)
		}
		if hAuthorization == "" || err != nil || user.ID == 0 || user.Name == "" {
			data := errorx.NewSystemError(r.Context(), UNAUTHORIZED, 0).(*errorx.TyyCodeError).Data()
			body, _ := json.Marshal(data)
			w.Header().Set(httpx.ContentType, httpx.JsonContentType)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(body)
			return
		}
		ctx := context.WithValue(r.Context(), UserInfoKey, user)
		//ctx, err = help.SetIDNameToMetadataCtx(ctx, strconv.FormatInt(id, 10), name)
		next(w, r.WithContext(ctx))
	}
}

func GetUserFromCtx(ctx context.Context) (u UserInfo, err error) {
	user := ctx.Value(UserInfoKey)
	if user == nil {
		err = errors.New("ctx user is nil")
		return
	}
	u = user.(UserInfo)
	if u.ID == 0 {
		err = errors.New("user id is nil")
		return
	}
	return u, nil
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
func authorizationFun(authKeys AuthKeys, authorization string, user *UserInfo) (err error) {
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
		fid, _ := data["id"].(float64)
		name, _ := data["name"].(string)
		areaID, _ := data["areaID"].(float64)
		user.ID = int64(fid)
		user.Name = name
		user.AreaID = int64(areaID)
	default:
		err = errors.New("鉴权失败")
		return
	}
	return
}

//func (a AuthKeys) errOutput(w http.ResponseWriter) {
//	data := errorx.NewSystemError("system.unauthorized", 0, "en").(*errorx.TyyCodeError).Data()
//	body, _ := json.Marshal(data)
//	w.Header().Set(httpx.ContentType, httpx.JsonContentType)
//	w.WriteHeader(http.StatusUnauthorized)
//	w.Write(body)
//}
