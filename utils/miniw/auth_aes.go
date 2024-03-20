package miniw

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/130-133/go-common/utils/encrypt"
	"github.com/tidwall/gjson"
)

type AuthInfo struct {
	Uin       int64  `json:"Uin,omitempty"`
	Token     string `json:"token,omitempty"`
	Sign      string `json:"sign,omitempty"`
	RegTime   int32  `json:"reg_time,omitempty"`
	IsNewUser bool   `json:"isnewuser,omitempty"`
}

type RoleInfo struct {
	Nickname string `json:"NickName,omitempty"`
	Model    int32  `json:"Model,omitempty"`
	SkinID   int64  `json:"SkinID,omitempty"`
}

type BaseInfo struct {
	Uin                 int64    `json:"Uin,omitempty"`
	AccountBindTime     int32    `json:"AccountBindTime,omitempty"`
	RoleInfo            RoleInfo `json:"RoleInfo,omitempty"`
	ProductionChannelID int32    `json:"production_channelid,omitempty"`
	CreateTime          int32    `json:"CreateTime,omitempty"`
	IsDeveloper         bool     `json:"isdeveloper,omitempty"`
	LastLoginTime       int32    `json:"LastLoginTime,omitempty"`
	Email               string   `json:"Email,omitempty"`
	IsFreeze            bool     `json:"isfreeze,omitempty"`
	AppID               int32    `json:"appid,omitempty"`
	UinFlag             int32    `json:"UinFlag,omitempty"`
	Leval               int32    `json:"level,omitempty"`
	LastLoginIP         int64    `json:"LastLoginIP,omitempty"`
	ProductionId        int32    `json:"production_id,omitempty"`
	Phone               string   `json:"Phone,omitempty"`
	IDCardAuthed        bool     `json:"idcard_authed,omitempty"`
	CtlVersion          int32    `json:"cltversion,omitempty"`
}

type MiniwBody struct {
	Code     int32     `json:"code,omitempty"`
	Msg      string    `json:"msg,omitempty"`
	AuthInfo *AuthInfo `json:"authinfo,omitempty"`
	BaseInfo *BaseInfo `json:"baseinfo,omitempty"`
}

type AESConf struct {
	Key string
	IV  string
}

type MiniwAuthByAES struct {
	aes               AESConf
	offset            int
	originBody        []byte
	encryptedAuthInfo string
	encryptedBaseInfo string
	decryptedAuthInfo string
	decryptedBaseInfo string
	body              *MiniwBody
	error             error
}

type MiniwAuthParam struct {
	Target    string
	Source    string
	Timestamp int64
	Key       string
}

func NewMiniwAuthByAES() *MiniwAuthByAES {
	return &MiniwAuthByAES{}
}

func (m *MiniwAuthByAES) WithAES(aes *AESConf) *MiniwAuthByAES {
	if len(aes.Key) == 0 || len(aes.IV) == 0 {
		m.error = errors.New("invalid aes key or iv")
	}
	m.aes.Key = aes.Key
	m.aes.IV = aes.IV
	return m
}

func (m *MiniwAuthByAES) WithBody(body []byte) *MiniwAuthByAES {
	if len(body) == 0 {
		m.error = errors.New("invalid body")
	}

	m.originBody = body
	m.body = &MiniwBody{}

	b := gjson.ParseBytes(body)
	if b.Get("code").Type == gjson.String {
		code := b.Get("code").String()
		if code != "OK" {
			m.body.Code = -1
		} else {
			m.body.Code = 0
		}
	} else {
		m.body.Code = int32(b.Get("code").Int())
	}
	m.body.Msg = b.Get("msg").String()

	if m.body.Code != 0 {
		m.error = errors.New(m.body.Msg)
		return m
	}

	// 计算offset
	iv := b.Get("iv").String()
	pattern := regexp.MustCompile(`\d+`)
	ns := pattern.FindAllString(iv, -1)
	n, _ := strconv.ParseInt(strings.Join(ns, ""), 10, 64)
	m.offset = int(n)

	// 计算加密数据
	encryptedAuthInfo, err := m.rebuild(b.Get("authinfo").String())
	if err != nil {
		return m
	}
	m.encryptedAuthInfo = encryptedAuthInfo

	encryptedBaseInfo, err := m.rebuild(b.Get("baseinfo").String())
	if err != nil {
		return m
	}
	m.encryptedBaseInfo = encryptedBaseInfo

	return m
}

// 重组加密数据
func (m *MiniwAuthByAES) rebuild(originData string) (string, error) {
	l := len(originData)
	if l == 0 {
		m.error = errors.New(m.body.Msg)
		return "", m.error
	}
	left := l - (m.offset % l)
	return fmt.Sprintf("%s%s", originData[left:], originData[0:left]), nil
}

func (m *MiniwAuthByAES) GetOffset() int {
	return m.offset
}

// 获取code
func (m *MiniwAuthByAES) GetCode() int32 {
	return m.body.Code
}

// 获取msg
func (m *MiniwAuthByAES) GetMsg() string {
	return m.body.Msg
}

// 获取Uin
func (m *MiniwAuthByAES) GetUin() int64 {
	if m.body.AuthInfo.Uin > 0 {
		return m.body.AuthInfo.Uin
	}
	return m.body.BaseInfo.Uin
}

// 获取token
func (m *MiniwAuthByAES) GetToken() string {
	return m.body.AuthInfo.Token
}

// 获取sign
func (m *MiniwAuthByAES) GetSign() string {
	return m.body.AuthInfo.Sign
}

func (m *MiniwAuthByAES) GetAuthInfo() *AuthInfo {
	return m.body.AuthInfo
}

func (m *MiniwAuthByAES) GetBaseInfo() *BaseInfo {
	return m.body.BaseInfo
}

func (m *MiniwAuthByAES) GetRoleInfo() *RoleInfo {
	return &m.body.BaseInfo.RoleInfo
}

// 已设置密码
func (m *MiniwAuthByAES) HasPassword() bool {
	return m.body.BaseInfo.AccountBindTime > 0
}

func (m *MiniwAuthByAES) GetNickname() string {
	return m.body.BaseInfo.RoleInfo.Nickname
}

// 生成请求的auth参数
func (m *MiniwAuthByAES) GenRequestAuth(o *MiniwAuthParam) string {
	s := fmt.Sprintf("source=%s&target=%s&time=%d%s", o.Source, o.Target, o.Timestamp, o.Key)
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func (m *MiniwAuthByAES) Decrypt() *MiniwAuthByAES {
	if m.body.Code != 0 || m.error != nil {
		return m
	}

	key := m.aes.Key
	iv := m.aes.IV
	decodedauthinfo := encrypt.NewAes(key, iv, encrypt.OutBase64).Decode(m.encryptedAuthInfo)
	if decodedauthinfo.Error() != nil {
		m.error = decodedauthinfo.Error()
		return m
	}
	decodedbaseinfo := encrypt.NewAes(key, iv, encrypt.OutBase64).Decode(m.encryptedBaseInfo)
	if decodedbaseinfo.Error() != nil {
		m.error = decodedbaseinfo.Error()
		return m
	}

	authinfo := &AuthInfo{}
	if err := json.Unmarshal([]byte(decodedauthinfo.Data()), authinfo); err != nil {
		m.error = err
		return m
	}
	m.body.AuthInfo = authinfo
	if len(m.body.AuthInfo.Sign) > 0 {
		m.DecodeSign()
	}

	baseinfo := &BaseInfo{}
	if err := json.Unmarshal([]byte(decodedbaseinfo.Data()), baseinfo); err != nil {
		m.error = err
		return m
	}
	m.body.BaseInfo = baseinfo

	return m
}

func (m *MiniwAuthByAES) GetBody() *MiniwBody {
	return m.body
}

func (m *MiniwAuthByAES) Error() error {
	return m.error
}

// 解密sign
func (m *MiniwAuthByAES) DecodeSign() *MiniwAuthByAES {
	encodeSign := m.body.AuthInfo.Sign
	// fmt.Println("==== encodeSign", encodeSign)
	arr := strings.Split(encodeSign, "_")
	signA := arr[0]
	signB := arr[1]
	timestamp, _ := strconv.Atoi(arr[1])
	iv := timestamp % 32
	// fmt.Println("==== iv", iv)
	bp := len(signA) - iv
	signA = fmt.Sprintf("%s%s", signA[bp:], signA[0:bp])

	// 反转字符
	bytes := []rune(signA)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	signA = string(bytes)

	decodeSign := fmt.Sprintf("%s_%s", signA, signB)
	// fmt.Println("==== decodeSign", decodeSign)
	m.body.AuthInfo.Sign = decodeSign
	return m
}
