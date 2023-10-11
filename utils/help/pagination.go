package help

import "reflect"

// Pagination 列表分页
func Pagination(list interface{}, currentPage, pageSize int64) {
	if pageSize <= 0 {
		return
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	if list == nil {
		return
	}
	t := reflect.TypeOf(list)
	v := reflect.Indirect(reflect.ValueOf(list))
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		return
	}
	if v.Len() == 0 {
		return
	}
	length := v.Len() - 1
	start := int((currentPage - 1) * pageSize)
	end := start + int(pageSize)
	if start > length {
		v.Set(reflect.MakeSlice(t, 0, 0))
		return
	}
	if end > length {
		end = length + 1
	}
	newV := v.Slice(start, end)
	v.Set(newV)
}
