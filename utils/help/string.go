package help

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type strToAny struct {
	str string
}

// ToAny 字符串转多个类型
func ToAny[T ~string | ~[]byte](value T) strToAny {
	return strToAny{
		str: string(value),
	}
}

func (a strToAny) Int() (res int) {
	res, _ = strconv.Atoi(a.str)
	return
}
func (a strToAny) Int32() (res int32) {
	intVal, _ := strconv.ParseInt(a.str, 10, 32)
	res = int32(intVal)
	return
}
func (a strToAny) Int64() (res int64) {
	res, _ = strconv.ParseInt(a.str, 10, 64)
	return
}
func (a strToAny) Float32() (res float32) {
	intVal, _ := strconv.ParseFloat(a.str, 32)
	res = float32(intVal)
	return
}
func (a strToAny) Float64() (res float64) {
	res, _ = strconv.ParseFloat(a.str, 64)
	return
}
func (a strToAny) Bool() (res bool) {
	res, _ = strconv.ParseBool(a.str)
	return
}
func (a strToAny) Slice() (res []string) {
	res = []string{a.str}
	return
}
func (a strToAny) Str() (res string) {
	res = a.str
	return
}
func (a strToAny) Byte() (res []byte) {
	res = []byte(a.str)
	return
}

// ToString 转字符串
func ToString(val interface{}) string {
	if val == nil {
		return ""
	}
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	valValue := reflect.ValueOf(val)
	valValue = reflect.Indirect(valValue)
	str := ""
	switch valValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.Itoa(int(valValue.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.Itoa(int(valValue.Uint()))
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(valValue.Float(), 'f', -1, 64)
	case reflect.Slice, reflect.Array:
		retBytes, _ := json.Marshal(valValue.Interface())
		if valType.Elem().Kind() == reflect.Uint8 {
			retBytes = valValue.Interface().([]byte)
		}
		str = string(retBytes)
	case reflect.Struct, reflect.Map:
		retBytes, _ := json.Marshal(valValue.Interface())
		str = string(retBytes)
	case reflect.String:
		str = valValue.String()
	default:
	}
	return str
}
