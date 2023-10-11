package help

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectIDs(ids []string) []primitive.ObjectID {
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objectId)
	}
	return objectIds
}

func StructToBson(model interface{}) bson.M {
	data := make(map[string]interface{})
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return data
	}
	vPtr := reflect.ValueOf(model)
	v := reflect.Indirect(vPtr)

	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		interVal := v.Field(i).Interface()
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.Field(i).Int() == 0 {
				continue
			}
		case reflect.Float32, reflect.Float64:
			if v.Field(i).Float() == 0 {
				continue
			}
		case reflect.Bool:
			if v.Field(i).Bool() == false {
				continue
			}
		default:
			timeType, ok := interVal.(time.Time)
			if interVal == "" || interVal == nil || (ok && timeType.IsZero()) {
				continue
			}
		}
		if t.Field(i).Type.Kind() == reflect.Slice && v.Field(i).Len() == 0 {
			continue
		}

		var value interface{}
		value = interVal
		if t.Field(i).Type.Kind() == reflect.Ptr && t.Field(i).Type.Elem().Kind() == reflect.Struct {
			if v.Field(i).IsNil() {
				continue
			}
			value = StructToBson(interVal)
		}
		// key
		tags := t.Field(i).Tag.Get("bson")
		tag := strings.SplitN(tags, ",", 2)
		if len(tag) == 0 {
			tag[0] = t.Field(i).Name
		}
		data[tag[0]] = value
	}
	return data
}
