package help

import (
	"errors"
	"reflect"
)

// FilterSingle 判断info结构体 是否符合过滤条件Slice
func FilterSingle(info interface{}, filter interface{}) bool {
	if filter == nil {
		return true
	}
	if reflect.TypeOf(info).Kind() != reflect.Ptr {
		return true
	}
	var oks []bool
	tv := reflect.TypeOf(filter)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	vv := reflect.Indirect(reflect.ValueOf(filter))
	if vv.Kind() == reflect.Invalid {
		return true
	}
	targetT := reflect.TypeOf(info).Elem()
	if targetT.Kind() == reflect.Ptr {
		targetT = targetT.Elem()
	}
	targetV := reflect.Indirect(reflect.ValueOf(info))

	// 过滤每个字段
	for i := 0; i < tv.NumField(); i++ {
		filterFieldName := tv.Field(i).Name
		filterV := vv.FieldByName(filterFieldName)
		if !filterV.CanInterface() || filterV.Kind() != reflect.Slice {
			continue
		}
		targetField := targetV.FieldByName(filterFieldName)
		isOk := false
		if filterV.Len() == 0 {
			continue
		}

		// 判断每个字段的多个值
		for x := 0; x < filterV.Len(); x++ {
			if targetField.Interface() == filterV.Index(x).Interface() {
				isOk = true
				break
			}
		}
		oks = append(oks, isOk)
	}
	// 判断是否存在不符合要求
	for _, ok := range oks {
		if !ok {
			return false
		}
	}
	return true
}

// InArray 指定值判断是否被包含
// *注意字符串与浮点数对比
// "1.1" == 1.1
// "1" == 1.0
// "1.11" == 1.110
// "1.10" != 1.1
func InArray(val, arr interface{}) (bool, error) {
	if arr == nil {
		return false, nil
	}
	valType := reflect.TypeOf(val)
	arrType := reflect.TypeOf(arr)
	arrValue := reflect.ValueOf(arr)
	if arrType.Kind() == reflect.Ptr {
		arrType = arrType.Elem()
	}
	if arrType.Kind() != reflect.Slice {
		return false, errors.New("not Array type for 'arr'")
	}
	if valType.Kind() == reflect.Slice ||
		valType.Kind() == reflect.Array ||
		valType.Kind() == reflect.Map {
		return false, errors.New("not allow Slice|Array|Map type for 'val'")
	}
	compA := ToString(val)
	for i := 0; i < arrValue.Len(); i++ {
		compB := ToString(arrValue.Index(i).Interface())
		if compA == compB {
			return true, nil
		}
	}
	return false, nil
}
