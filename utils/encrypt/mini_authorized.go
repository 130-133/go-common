package encrypt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/130-133/go-common/utils/help"
	"strconv"
	"strings"
)

type MiniAuthorizedEncrypt struct {
	secret      string
	error       error
	DecryptData *MiniTokenParam
}

type MiniTokenParam struct {
	Uin  string
	Time int64
	S2t  int64
	Sign string
}

func NewMiniAuthorized(secret string) *MiniAuthorizedEncrypt {
	return &MiniAuthorizedEncrypt{
		secret: secret,
	}
}

func (e *MiniAuthorizedEncrypt) Decrypt(authorization string) *MiniAuthorizedEncrypt {
	if authorization == "" {
		e.error = errors.New("authorization empty")
		return e
	}
	if !strings.HasPrefix(authorization, "Basic") {
		e.error = errors.New("authorization not basic type")
		return e
	}
	authorization = authorization[6:]
	sDec, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		e.error = errors.New("authorization base64 decode fail")
		return e
	}
	arr := strings.SplitN(string(sDec), ".", 4)
	times, _ := strconv.ParseInt(arr[1], 10, 64)
	s2t, _ := strconv.ParseInt(arr[2], 10, 64)
	e.DecryptData = &MiniTokenParam{
		Uin:  arr[0],
		Time: times,
		S2t:  s2t,
		Sign: arr[3],
	}
	return e
}

func (e *MiniAuthorizedEncrypt) Data() *MiniTokenParam {
	if e.DecryptData == nil {
		return &MiniTokenParam{}
	}
	return e.DecryptData
}

func (e *MiniAuthorizedEncrypt) Error() error {
	return e.error
}

func (e *MiniAuthorizedEncrypt) Verify() bool {
	str1 := fmt.Sprintf("%s%s%d", e.DecryptData.Uin, e.secret, e.DecryptData.S2t)
	strMd5 := strings.ToLower(help.MD5(str1))
	str2 := fmt.Sprintf("%d%s%s", e.DecryptData.Time, strMd5, e.DecryptData.Uin)
	token := strings.ToLower(help.MD5(str2))
	return e.DecryptData.Sign == token
}
