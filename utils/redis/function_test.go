package redis

import (
	"database/sql"
	"testing"
	"time"
)

func TestLoadResult_Unmarshal(t *testing.T) {
	r := NewRedisConn()
	c, _ := r.LoadHSetEx("a", "b", 10*time.Minute, func() ([]byte, error) {
		a := []byte("{\"id\":2,\"uin\":\"1052999485\",\"orderId\":0,\"courseId\":\"60a7a4df510bf7c91ef2287c\",\"courseType\":4,\"courseName\":\"年卡\",\"courseImg\":\"\",\"createdAt\":{\"Time\":\"2022-01-05T14:34:09+08:00\",\"Valid\":true},\"updatedAt\":{\"Time\":\"2022-01-05T14:34:09+08:00\",\"Valid\":true}}")
		return a, nil
	})

	type UserCourse struct {
		ID         int64        `gorm:"primaryKey;column:id" json:"id"`
		Uin        string       `gorm:"column:uin" json:"uin"`
		OrderId    int64        `gorm:"column:order_id" json:"orderId"`
		CourseId   string       `gorm:"column:course_id" json:"courseId"`     // 课程id
		CourseType int32        `gorm:"column:course_type" json:"courseType"` // 课程类型
		CourseName string       `gorm:"column:course_name" json:"courseName"` // 课程名称
		CourseImg  string       `gorm:"column:course_img" json:"courseImg"`   // 课程图片
		CreatedAt  sql.NullTime `gorm:"column:created_at" json:"createdAt"`
		UpdatedAt  sql.NullTime `gorm:"column:updated_at" json:"updatedAt"`
	}
	b := UserCourse{}
	c.Unmarshal(&b)
	t.Logf("%+v\n", b)
}

func TestMRedis_LoadHSetEx(t *testing.T) {
	r := NewRedisConn(WithPwd("123456"))
	r.LoadHSetEx("hset", "3", 10*time.Minute, func() ([]byte, error) {
		return []byte("1"), nil
	})
}
