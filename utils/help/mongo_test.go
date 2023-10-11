package help

import "testing"

func TestStructToBson(t *testing.T) {
	type C struct {
		CC string `bson:"cc" json:"cc"`
	}
	type B struct {
		BB string `bson:"bb" json:"bb"`
	}
	type A struct {
		B *B     `bson:"b" json:"b"`
		C *C     `bson:"c" json:"c"`
		A string `bson:"a" json:"a"`
	}
	b := &B{
		BB: "xxxx",
	}
	a := A{
		B: b,
		A: "aaaa",
	}
	c := StructToBson(&a)
	t.Log(c)
}
