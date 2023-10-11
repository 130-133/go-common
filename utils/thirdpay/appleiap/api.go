package appleiap

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"gitea.com/llm-PhotoMagic/go-common/utils/encrypt"
	"gitea.com/llm-PhotoMagic/go-common/utils/request"
)

const (
	AppleApi        = "https://api.storekit.itunes.apple.com"
	AppleApiSandbox = "https://api.storekit-sandbox.itunes.apple.com"
)

type Api struct {
	ctx   context.Context
	c     StoreAuthConfig
	p     *JwsParse
	cache keyCache
}

type StoreAuthConfig struct {
	Bid        string //app id
	Iss        string //app store令牌ID
	Kid        string //app store私钥ID
	Kp8        string //app store私钥文件
	PublicCert string //app store公钥
}

type keyCache struct {
	privateKey interface{}
	token      string
	expire     time.Time
}

type IApi interface {
	WithContext(ctx context.Context) IApi
	LookUp(orderId string) []*LookUpData
}

func NewApi(config StoreAuthConfig) IApi {
	return &Api{
		ctx:   context.Background(),
		c:     config,
		cache: keyCache{},
		p:     NewJws(config.PublicCert),
	}
}

func (a *Api) WithContext(ctx context.Context) IApi {
	a.ctx = ctx
	return a
}

func (a *Api) CheckAuthConfig() error {
	if a.c.Iss == "" {
		return errors.New("无效iss")
	}
	if a.c.Kid == "" {
		return errors.New("无效kid")
	}
	if a.c.Kp8 == "" {
		return errors.New("无效kp8")
	}
	if _, err := os.Stat(a.c.Kp8); err != nil && !os.IsExist(err) {
		return errors.New("找不到kp8文件")
	}
	return nil
}

// GetToken 获取请求加密token
func (a *Api) GetToken() string {
	var token string
	if a.cache.token != "" && a.cache.expire.After(time.Now()) {
		token = a.cache.token
	} else {
		//获取私钥
		buff, err := ioutil.ReadFile(a.c.Kp8)
		if err != nil {
			return ""
		}
		deBuff, _ := pem.Decode(buff)
		pKey, _ := x509.ParsePKCS8PrivateKey(deBuff.Bytes)

		//签名
		header := map[string]interface{}{
			"kid": a.c.Kid,
		}
		payload := map[string]interface{}{
			"iss": a.c.Iss,
			"iat": time.Now().Unix(),
			"exp": time.Now().Unix() + 3600,
			"aud": "appstoreconnect-v1",
			"bid": a.c.Bid,
		}
		token = encrypt.NewJwt(encrypt.JwtConfig{PrivateKey: pKey}).WithJwtToken(jwt.SigningMethodES256, header).Encode(payload).String()
		a.cache = keyCache{
			privateKey: pKey,
			token:      token,
			expire:     time.Now().Add(50 * time.Minute),
		}
	}
	return token
}

// LookUp 查询用户订单
func (a *Api) LookUp(orderId string) []*LookUpData {
	if err := a.CheckAuthConfig(); err != nil {
		return nil
	}
	token := a.GetToken()
	opt := request.WithHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	resp := request.NewRestyReq(a.ctx).Get(
		AppleApi,
		fmt.Sprintf("/inApps/v1/lookup/%s", orderId),
		nil,
		opt,
	)
	if !resp.IsOk() || resp.GetError() != nil {
		return nil
	}
	body := LookUpResp{}
	if err := resp.GetJson(&body); err != nil {
		return nil
	}
	if body.Status != 0 {
		//empty
		return nil
	}

	var result []*LookUpData
	for _, v := range body.SignedTransactions {
		data, err := a.p.IapParseLookUp(v)
		if err != nil {
			break
		}
		result = append(result, data)
	}
	return result
}
