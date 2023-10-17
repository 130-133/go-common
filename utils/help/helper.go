package help

import (
	"bytes"
	"context"
	"crypto/md5"
	cRand "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tidwall/gjson"
	"google.golang.org/grpc/metadata"
)

func GetUinFromCtx(ctx context.Context) (uin string, err error) {
	uinInter := ctx.Value("Uin")
	if uinInter == nil {
		err = errors.New("ctx uin interface is nil")
		return
	}

	uin = uinInter.(string)
	if uin == "" {
		err = errors.New("uin is nil")
		return
	}

	return
}

func SetUinToMetadataCtx(ctx context.Context, uin string) (rsCtx context.Context, err error) {
	if uin == "" {
		err = errors.New("uin is nil")
		return
	}

	md := metadata.Pairs("rUin", uin)
	rsCtx = metadata.NewOutgoingContext(ctx, md)

	return
}

// GetRpcUinFromCtx 从metadata中获取Uin (需要上层ctx透传)
func GetRpcUinFromCtx(ctx context.Context) (uin string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("metadata fromIncomingContext is false")
		return
	}

	uinArr := md.Get("rUin")
	if len(uinArr) <= 0 {
		err = errors.New("uinArr is nil")
		return
	}

	uin = uinArr[0]
	if uin == "" {
		err = errors.New("uin is nil")
		return
	}

	return
}

func GetTodayTimeRemaining() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"

	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)

	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second

	return remainSecond
}

func GetPagingParam(pageIndex int64, pageSize int64) (limit int64, offset int64) {
	if pageSize <= 0 {
		limit = 10
	} else {
		limit = pageSize
	}

	if pageIndex <= 1 {
		offset = 0
	} else {
		offset = (pageIndex - 1) * pageSize
	}

	return
}

func GetRandString(length int) string {
	if length < 1 {
		return ""
	}
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))

	var rchar string = ""
	for i := 1; i <= length; i++ {
		rchar = rchar + charArr[ran.Intn(charlen)]
	}
	return rchar
}

//--------------------------以上为山本先人coding-----------------------------------//

// GetRandRange 获取指定大小范围内随机正整数
func GetRandRange(min, max int64) int64 {
	if min > max {
		return min
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := cRand.Int(cRand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := cRand.Int(cRand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

// StructToStruct 结构体同名引用复制 同名同类型字段
func StructToStruct(source, target interface{}) {
	if source == nil {
		return
	}
	aType := reflect.TypeOf(source)
	aValue := reflect.ValueOf(source)
	bType := reflect.TypeOf(target)
	bValue := reflect.ValueOf(target)
	aValue = reflect.Indirect(aValue)
	bValue = reflect.Indirect(bValue)
	if aType.Kind() == reflect.Ptr {
		aType = aType.Elem()
	}
	if bType.Kind() == reflect.Ptr {
		bType = bType.Elem()
	}
	for i := 0; i < aType.NumField(); i++ {
		aField := aType.Field(i).Name
		aChildType := aType.Field(i).Type
		aChildVal := aValue.FieldByName(aField)
		bChildTypeTmp, ok := bType.FieldByName(aField)
		bChildVal := bValue.FieldByName(aField)
		if !ok {
			continue
		}
		bChildType := bChildTypeTmp.Type
		if aChildType.Kind() == reflect.Ptr {
			aChildType = aChildType.Elem()
		}
		if bChildType.Kind() == reflect.Ptr {
			bChildType = bChildType.Elem()
		}
		if aChildType != bChildType {
			continue
		}
		if !bChildVal.CanSet() {
			continue
		}
		bChildVal.Set(aChildVal)
	}
}

func MD5(data string, salts ...string) string {
	str := fmt.Sprintf("%s%s", data, strings.Join(salts, ""))
	newSig := md5.Sum([]byte(str)) //转成加密编码
	// 将编码转换为字符串
	newArr := fmt.Sprintf("%x", newSig)
	//输出字符串字母都是小写，转换为大写
	data = strings.ToTitle(newArr)
	return data
}

func Md5(data string, salts ...string) string {
	return strings.ToLower(MD5(data, salts...))
}

// GenerateNo 创建单号
func GenerateNo(prefix, suffix string) string {
	randomNum := GetRandRange(0, 100000)
	timeNum := time.Now().Local().Format(FormatRawTime)
	return fmt.Sprintf("%s%s%06d%s", prefix, timeNum, randomNum, suffix)
}

// GetRawBody 获取原始body 并写回缓存
func GetRawBody(r *http.Request) string {
	buf, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return string(buf)
}

// ParseJsonToStruct 解析JSON到dest结构体 赋值类型尝试以dest的结构类型为主
func ParseJsonToStruct(byteStr []byte, dest interface{}) error {
	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return json.Unmarshal(byteStr, dest)
	}
	v = reflect.Indirect(v)
	jsonObject := gjson.ParseBytes(byteStr)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldTypeName := field.Type.Name()
		fieldKind := field.Type.Kind()
		if fieldKind == reflect.Ptr {
			fieldTypeName = field.Type.Elem().Name()
			fieldKind = field.Type.Elem().Kind()
		}
		jsonTag := field.Tag.Get("json")
		tagArr := strings.Split(jsonTag, ",")
		if len(tagArr) > 1 {
			jsonTag = tagArr[0]
		}
		if jsonTag == "" {
			jsonTag = field.Name
		}
		vField := v
		if vField.Kind() != reflect.Invalid {
			vField = v.Field(i)
		}
		if !vField.CanSet() {
			continue
		}

		jsonVal := jsonObject.Get(jsonTag)
		switch fieldKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vField.SetInt(jsonVal.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vField.SetUint(jsonVal.Uint())
		case reflect.Float32, reflect.Float64:
			vField.SetFloat(jsonVal.Float())
		case reflect.String:
			vField.SetString(jsonVal.String())
		case reflect.Bool:
			boolean := false
			if jsonVal.IsBool() {
				boolean = jsonVal.Bool()
			} else if !(jsonVal.String() == "" || jsonVal.Float() == 0 || len(jsonVal.Array()) == 0) {
				boolean = true
			}
			vField.SetBool(boolean)
		case reflect.Struct:
			if fieldTypeName == "Time" {
				val := jsonVal.String()
				timeVal := ParseDate(val)
				vField.Set(reflect.ValueOf(timeVal))
				break
			}
			if fieldTypeName == "NullTime" {
				newStruct := sql.NullTime{}
				json.Unmarshal([]byte(jsonVal.Raw), &newStruct)
				vField.Set(reflect.ValueOf(newStruct))
				break
			}

			// 匿名类型
			if field.Anonymous && !jsonVal.Exists() {
				jsonVal = jsonObject
			}

			newStruct := reflect.New(field.Type)
			if field.Type.Kind() == reflect.Ptr {
				newStruct = reflect.New(field.Type.Elem())
				vField.Set(newStruct)
			}
			_ = ParseJsonToStruct([]byte(jsonVal.Raw), newStruct.Interface())
		case reflect.Map:
			if len(jsonVal.Map()) == 0 {
				break
			}
			mapKeyKind := field.Type.Key().Kind()
			mapValueKind := vField.Type().Elem().Kind()
			if mapKeyKind != reflect.String {
				break
			}
			newMap := reflect.MakeMap(field.Type)
			vField.Set(newMap)
			jsonVal.ForEach(func(k, v gjson.Result) bool {
				mapK := reflect.ValueOf(v.String())
				switch mapValueKind {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					newMap.SetMapIndex(mapK, reflect.ValueOf(v.Int()))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					newMap.SetMapIndex(mapK, reflect.ValueOf(v.Uint()))
				case reflect.Float32, reflect.Float64:
					newMap.SetMapIndex(mapK, reflect.ValueOf(v.Float()))
				case reflect.String:
					newMap.SetMapIndex(mapK, reflect.ValueOf(v.String()))
				case reflect.Interface:
					newMap.SetMapIndex(mapK, reflect.ValueOf(v.Value()))
				}
				return true
			})
		case reflect.Slice:
			arr := jsonObject.Get(jsonTag).Array()
			newSliceLen := len(arr)
			var (
				fieldType reflect.Type
				newSlice  reflect.Value
				ptr       reflect.Value
			)
			if field.Type.Kind() == reflect.Ptr {
				fieldType = field.Type.Elem()
				ptr = reflect.New(fieldType)
				newSlice = reflect.MakeSlice(fieldType, newSliceLen, newSliceLen)
				ptr.Elem().Set(newSlice)
				vField.Set(ptr)
			} else {
				fieldType = field.Type
				newSlice = reflect.MakeSlice(fieldType, newSliceLen, newSliceLen)
				vField.Set(newSlice)
			}
			if newSliceLen == 0 {
				break
			}
			fieldKind = fieldType.Elem().Kind()
			if fieldKind == reflect.Ptr {
				fieldKind = fieldType.Elem().Elem().Kind()
			}
			for k, v := range arr {
				switch fieldKind {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					newSlice.Index(k).SetInt(v.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					newSlice.Index(k).SetUint(v.Uint())
				case reflect.Float32, reflect.Float64:
					newSlice.Index(k).SetFloat(v.Float())
				case reflect.String:
					newSlice.Index(k).SetString(v.String())
				case reflect.Bool:
					newSlice.Index(k).SetBool(v.Bool())
				case reflect.Struct:
					newStruct := reflect.New(fieldType.Elem())
					if fieldType.Elem().Kind() == reflect.Ptr {
						newStruct = reflect.New(fieldType.Elem().Elem())
						newSlice.Index(k).Set(newStruct)
					}
					_ = ParseJsonToStruct([]byte(v.Raw), newStruct.Interface())
				}
			}
		}
	}
	return nil
}

// StructToMap struct通过json转map
func StructToMap(source, target interface{}) error {
	arrByte, err := json.Marshal(source)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(arrByte, target); err != nil {
		return err
	}
	return nil
}

// IsMobile 判断手机号
func IsMobile(mobile string) bool {
	match, err := regexp.MatchString("^(0|86|17951)?(13[0-9]|15[012356789]|166|17[3678]|18[0-9]|14[57])[0-9]{8}$", mobile)
	if err != nil {
		return false
	}
	return match
}

// InitialUpper 首字母大写
func InitialUpper(text string) string {
	if len(text) == 0 {
		return ""
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
