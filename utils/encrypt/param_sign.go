package encrypt

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/130-133/go-common/utils/help"
)

type IParamSign interface {
	SetSignKey(key string) IParamSign
	SetSignTimeKey(key string) IParamSign
	SetLowerSign() IParamSign
	SetUpperSign() IParamSign
	SetFormatter(fn KVFormat) IParamSign
	SetSeparator(seq string) IParamSign
	SetSignFunc(fn SignFunc) IParamSign
	EnableValidTime(key bool) IParamSign
	Sign(req interface{}) IParamSignResult
	MiniWorldSign(req interface{}) IParamSignResult
}

type IParamSignResult interface {
	Debug() IParamSignResult
	Verify() error
	SignedVal() string
}

type KVFormat func(k, v string) string
type SignFunc func(signing string) string
type ParamSignEncrypt struct {
	signedStr     string        //已签名值
	signKey       string        //签名字段名
	signStr       string        //传入的验签值
	signTimeKey   string        //签名时间字段名
	signTimeStr   int64         //传入的签名时间
	isValidTime   bool          //是否验证签名时间
	signValidTime time.Duration //签名有效时长
	signCase      SignCase      //签名过程大小写
	formatter     KVFormat      //拼接方式
	signFunc      SignFunc      //加密逻辑
	separator     string        //参数见分隔符
	token         string        //加密密钥

}

type ParamSignResult struct {
	p          *ParamSignEncrypt
	signingStr string //签名原始参数
}

func NewParamSign(token string) IParamSign {
	return &ParamSignEncrypt{
		token:         token,
		signKey:       "sign",
		signTimeKey:   "signtime",
		signValidTime: 5 * time.Minute,
		separator:     "&",
		isValidTime:   true,
	}
}

// SetSignKey 加密字段名
func (e *ParamSignEncrypt) SetSignKey(key string) IParamSign {
	e.signKey = key
	return e
}

// SetSignTimeKey 签名时间字段名
func (e *ParamSignEncrypt) SetSignTimeKey(key string) IParamSign {
	e.signTimeKey = key
	return e
}

// EnableValidTime 开启校验时间
func (e *ParamSignEncrypt) EnableValidTime(key bool) IParamSign {
	e.isValidTime = key
	return e
}

// SetLowerSign 签名过程统一小写
func (e *ParamSignEncrypt) SetLowerSign() IParamSign {
	e.signCase = Lower
	return e
}

// SetUpperSign 签名过程统一大写
func (e *ParamSignEncrypt) SetUpperSign() IParamSign {
	e.signCase = Upper
	return e
}

// SetFormatter 设置待验签参数格式化方法
func (e *ParamSignEncrypt) SetFormatter(fn KVFormat) IParamSign {
	e.formatter = fn
	return e
}

// SetSeparator 参数间分隔符
func (e *ParamSignEncrypt) SetSeparator(sep string) IParamSign {
	e.separator = sep
	return e
}

func (e *ParamSignEncrypt) SetSignFunc(fc SignFunc) IParamSign {
	e.signFunc = fc
	return e
}

// Sign 参数签名
// 签名方式：
// 1、所有参数除（sign）都参与签名
// 2、按参数按键名ascii排序
// 3、通过键值对用 = 方式组合字符串，例如 key=value
// 4、多个参数之间使用 & 拼接，例如  key1=value1&key2=value2
// 5、拼接完整体根据要求转大小写
// 6、md5整个字符串
// 7、后面直接拼接token值，例如 e10adc3949ba59abbe56e057f20f883ee10adc3949ba59abbe56e057f20f883e
// 8、md5拼接后的值，生成最终验签密钥
func (e *ParamSignEncrypt) Sign(req interface{}) IParamSignResult {
	rType := reflect.TypeOf(req)
	rValue := reflect.ValueOf(req)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	rValue = reflect.Indirect(rValue)

	var (
		waitSign []string
		sorts    sort.StringSlice
	)
	extract := func(field interface{}, kind reflect.Kind) bool {
		var (
			key string
			tmp interface{}
		)
		keyValue, ok := field.(reflect.Value)
		if ok {
			key = keyValue.String()
		} else {
			key, _ = field.(string)
		}

		//提取验签值
		if strings.ToLower(key) == e.signKey {
			if kind == reflect.Struct {
				tmp = rValue.FieldByName(key).Interface()
			} else {
				tmp = rValue.MapIndex(keyValue).Interface()
			}
			e.signStr = help.ToString(tmp)
			return false
		}
		//提取签名时间值
		if strings.ToLower(key) == e.signTimeKey {
			if kind == reflect.Struct {
				tmp = rValue.FieldByName(key).Interface()
			} else {
				tmp = rValue.MapIndex(keyValue).Interface()
			}
			tmpStr := help.ToString(tmp)
			e.signTimeStr, _ = strconv.ParseInt(tmpStr, 10, 64)
		}
		return true
	}
	switch rType.Kind() {
	case reflect.Struct:
		for i := 0; i < rType.NumField(); i++ {
			key := rType.Field(i).Name
			if ok := extract(key, reflect.Struct); !ok {
				continue
			}
			sorts = append(sorts, key)
		}
		sorts.Sort()
		for _, k := range sorts {
			val := rValue.FieldByName(k).Interface()
			waitSign = append(waitSign, e.waitSignFormat(k, help.ToString(val)))
		}
	case reflect.Map:
		reqMap := make(map[string]interface{})
		m := rValue.MapRange()
		for m.Next() {
			key := m.Key()
			if ok := extract(key, reflect.Map); !ok {
				continue
			}
			sorts = append(sorts, key.String())
			reqMap[key.String()] = m.Value().Interface()
		}
		sorts.Sort()
		for _, k := range sorts {
			val, _ := reqMap[k]
			waitSign = append(waitSign, e.waitSignFormat(k, help.ToString(val)))
		}
	}
	signingStr := turnCase(strings.Join(waitSign, e.separator), e.signCase)
	e.signedStr = e.toSign(signingStr)
	return &ParamSignResult{
		p:          e,
		signingStr: signingStr,
	}
}

// MiniWorldSign 迷你世界签名
func (e *ParamSignEncrypt) MiniWorldSign(req interface{}) IParamSignResult {
	e.formatter = func(k, v string) string {
		return v
	}
	e.signFunc = func(signing string) string {
		return NewMD5(e.token).EncodeToken(signing)
	}
	e.separator = ""
	return e.Sign(req)
}

func (e *ParamSignEncrypt) waitSignFormat(k, v string) string {
	if e.formatter != nil {
		return e.formatter(k, v)
	}
	return fmt.Sprintf("%s=%v", k, v)
}

func (e *ParamSignEncrypt) toSign(signingStr string) string {
	if e.signFunc != nil {
		return e.signFunc(signingStr)
	}
	md5 := NewMD5(e.token)
	return md5.EncodeToken(md5.Encode(signingStr))
}

func (e *ParamSignResult) Debug() IParamSignResult {
	fmt.Printf("signingStr: 《%s》，signedStr:《%s》\n", e.signingStr, e.SignedVal())
	return e
}

// SignedVal 获取签名数据
func (e *ParamSignResult) SignedVal() string {
	return e.p.signedStr
}

// Verify 验签
func (e *ParamSignResult) Verify() (err error) {
	if res := e.VerifyTime(); !res {
		return errors.New("签名已过期")
	}
	if e.p.signedStr != e.p.signStr {
		return errors.New("验签失败")
	}
	return
}

// VerifyTime 验证签名时间
func (e *ParamSignResult) VerifyTime() bool {
	if !e.p.isValidTime {
		return true
	}
	signTime := time.Unix(e.p.signTimeStr, 0)
	now := time.Now().Local() //有效时间
	// 容忍5秒的服务间时间误差
	return signTime.Add(-5*time.Second).Before(now) && signTime.Add(e.p.signValidTime).After(now)
}
