package help

import (
	"testing"
)

func TestInArray(t *testing.T) {
	result, err := InArray("11.10", []interface{}{"1", 12, 3, 11.10})
	if err != nil {
		t.Error(err)
	}
	t.Logf("result:%+v", result)
}

func TestStructToStruct(t *testing.T) {
	type Common struct {
		Co int
	}
	type A struct {
		Common
		AA   map[int]int
		AAA  string
		AAAA int
	}
	type B struct {
		Common
		AAA  string
		AAAA int
		BBB  []string
	}
	type C struct {
		AA   map[int]int
		AAA  string
		AAAA string
		CCC  []string
	}
	a := A{
		Common: Common{Co: 1},
		AA:     map[int]int{1: 11, 2: 22, 3: 33},
		AAA:    "11",
		AAAA:   2,
	}
	b := B{}
	c := C{}
	StructToStruct(&a, &b)
	StructToStruct(&a, &c)
	t.Logf("%+v", a)
	t.Logf("%+v", b)
	t.Logf("%+v", c)
}

func TestParseJsonToStruct(t *testing.T) {
	jsonStr := `{
		"a": "",
		"b": 2,
		"c": {
			"e": 1,
			"f": ["xx", "yy"]
		},
		"d": [4,5],
		"f": [{"n":1}, {"n":2}],
		"h": 12,
		"i": {
			"e": 2,
			"f": ["xx", "yy"]
		},
		"j": {
			"k": "kkk",
            "l": 333
        },
		"k": [{"n":1}, {"n":2}],
	}`
	type C struct {
		E string   `json:"e"`
		F []string `json:"f"`
	}
	type G struct {
		H int                    `json:"h"`
		J map[string]interface{} `json:"j"`
	}
	type M struct {
		N int `json:"n"`
	}
	type A struct {
		//A bool  `json:"a"`
		//B int   `json:"b"`
		//C *C    `json:"c"`
		//D []int `json:"d"`
		F []*M  `json:"f"`
		K *[]*M `json:"k"`
		//*G
		//I C `json:"i"` //无法复制，必须使用指针类型
	}
	a := A{
		//I: C{
		//	E: "1234",
		//	F: []string{"123"},
		//},
	}
	ParseJsonToStruct([]byte(jsonStr), &a)
	t.Logf("a:%#v", a)
	//t.Logf("g:%#v", a.G)
	//t.Logf("c:%#v", a.C)
	for _, ii := range a.F {
		t.Logf("f:%#v", ii)
	}
	for _, ii := range *a.K {
		t.Logf("k:%#v", ii)
	}

}
