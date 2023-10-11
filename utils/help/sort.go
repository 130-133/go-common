package help

import (
	"reflect"
	"sort"
	"strings"
)

type Sort struct {
	Field string `json:"field"` // 字段名 对应 结构体字段名 非json名
	Asc   bool   `json:"asc"`   // 是否升序
}

// Sorter 列表排序（多列排序）
func Sorter(list interface{}, sorts []Sort) {
	if len(sorts) == 0 {
		return
	}
	if reflect.TypeOf(list).Kind() != reflect.Slice {
		return
	}
	sort.Slice(list, func(i, j int) bool {
		t := reflect.TypeOf(list)
		if t.Kind() != reflect.Slice {
			return false
		}
		v := reflect.ValueOf(list)

		u1 := reflect.Indirect(v.Index(i))
		u2 := reflect.Indirect(v.Index(j))
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		// 顺序排序
		for _, s := range sorts {
			// 首字母大写兼容
			field := InitialUpper(s.Field)
			uc1 := u1.FieldByName(field)
			uc2 := u2.FieldByName(field)
			if !uc1.IsValid() {
				continue
			}
			if uc1.Interface() == uc2.Interface() {
				continue
			}
			result := false
			switch {
			case uc1.Kind() >= reflect.Int && uc1.Kind() <= reflect.Int64:
				result = uc1.Int() < uc2.Int()
			case uc1.Kind() >= reflect.Uint && uc1.Kind() <= reflect.Uint64:
				result = uc1.Uint() < uc2.Uint()
			case uc1.Kind() >= reflect.Float32 && uc1.Kind() <= reflect.Float64:
				result = uc1.Float() < uc2.Float()
			case uc1.Kind() == reflect.String:
				result = uc1.String() < uc2.String()
			}
			if !s.Asc {
				return !result
			}
			return result
		}
		return false
	})
}

type SortKind string

const (
	Asc  SortKind = "asc"
	Desc SortKind = "desc"
)

func (k SortKind) Sort() SortKind {
	sort := strings.ToLower(strings.TrimSpace(string(k)))
	if sort == "desc" {
		return Desc
	}
	return Asc
}
