package encrypt

type Encrypt struct {
	md5            IMD5
	sha256         ISha256
	paramSign      IParamSign
	miniAuthorized *MiniAuthorizedEncrypt
	jwt            *JWTEncrypt
	aes            *AesEncrypt

	secret          string // 加密密钥
	miniWorldSecret string // 迷你世界密钥
	aseSecret       string // aes加密密钥
	iv              string // aes向量
}

type Opt func(*Encrypt)

func (o Opt) Apply(e *Encrypt) {
	o(e)
}

func New(options ...Opt) *Encrypt {
	e := &Encrypt{
		secret:          "314EFCE872EAB321D92ABF49F732AF98",
		miniWorldSecret: "c8c93222583741bd828579b3d3efd43b_1",
		aseSecret:       "adFEdE2A10SO2022",
		iv:              "1234567891234567",
	}
	for _, o := range options {
		o.Apply(e)
	}
	return e
}

func WithSecret(secret string) Opt {
	return func(e *Encrypt) {
		e.secret = secret
	}
}

func WithMiniWorldSecret(secret string) Opt {
	return func(e *Encrypt) {
		e.miniWorldSecret = secret
	}
}

func WithAseSecret(secret string) Opt {
	return func(e *Encrypt) {
		e.aseSecret = secret
	}
}

func WithIV(iv string) Opt {
	return func(e *Encrypt) {
		e.iv = iv
	}
}

// MD5 加密方式
func (e *Encrypt) MD5() IMD5 {
	if e.md5 == nil {
		e.md5 = NewMD5(e.secret)
	}
	return e.md5
}

// Sha256 加密方式
func (e *Encrypt) Sha256() ISha256 {
	if e.sha256 == nil {
		e.sha256 = NewSha256(e.secret)
	}
	return e.sha256
}

// ParamSign 参数签名方式
func (e *Encrypt) ParamSign() IParamSign {
	if e.paramSign == nil {
		e.paramSign = NewParamSign(e.secret)
	}
	return e.paramSign
}

// MiniAuthorized 迷你玩鉴权Authorized
func (e *Encrypt) MiniAuthorized() *MiniAuthorizedEncrypt {
	if e.miniAuthorized == nil {
		e.miniAuthorized = NewMiniAuthorized(e.miniWorldSecret)
	}
	return e.miniAuthorized
}

// Jwt 加密方式
func (e *Encrypt) Jwt() *JWTEncrypt {
	if e.jwt == nil {
		e.jwt = NewJwt(JwtConfig{Token: e.secret})
	}
	return e.jwt
}

// Aes 加密方式
func (e *Encrypt) Aes() *AesEncrypt {
	if e.aes == nil {
		e.aes = NewAes(e.aseSecret, e.iv, OutBase64)
	}
	return e.aes
}
