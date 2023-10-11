package encrypt

import (
	"crypto"

	"github.com/golang-jwt/jwt/v4"
)

type JWTEncrypt struct {
	jwtToken    *jwt.Token
	token       interface{}
	error       error
	signedToken string
}

type JwtConfig struct {
	CommonOptions
	Token      string
	PrivateKey crypto.PrivateKey
}
type JwtOpt func(config *JwtConfig)

func (o JwtOpt) Apply(options *JwtConfig) {
	o(options)
}

func NewJwt(c JwtConfig) *JWTEncrypt {
	return &JWTEncrypt{
		token:    getEncryptKey(c.Token, c.PrivateKey),
		jwtToken: jwt.New(jwt.SigningMethodHS256),
	}
}

func getEncryptKey(token string, privateKey crypto.PrivateKey) interface{} {
	if token != "" {
		return []byte(token)
	}
	if privateKey != nil {
		return privateKey
	}
	return ""
}

// WithJwtToken 替换加密方式
func (e *JWTEncrypt) WithJwtToken(method jwt.SigningMethod, header map[string]interface{}) *JWTEncrypt {
	e.jwtToken = jwt.New(method)
	for k, v := range header {
		e.jwtToken.Header[k] = v
	}
	return e
}

func (e *JWTEncrypt) Encode(data map[string]interface{}) *JWTEncrypt {
	e.jwtToken.Claims = jwt.MapClaims(data)
	e.signedToken, e.error = e.jwtToken.SignedString(e.token)
	return e
}

func (e *JWTEncrypt) String() string {
	return e.signedToken
}

func (e *JWTEncrypt) Error() error {
	return e.error
}

func (e *JWTEncrypt) Decode(signedToken string) *JWTEncrypt {
	e.jwtToken, e.error = jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return e.token, nil
	})
	return e
}

func (e *JWTEncrypt) Verify() bool {
	if e.jwtToken == nil {
		return false
	}
	return e.jwtToken.Valid
}

func (e *JWTEncrypt) Data() map[string]interface{} {
	res := make(map[string]interface{})
	if e.jwtToken == nil {
		return res
	}
	data, ok := e.jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return res
	}
	return data
}
