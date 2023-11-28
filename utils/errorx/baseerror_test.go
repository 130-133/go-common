package errorx

import (
	"testing"
)

func TestNewError(t *testing.T) {
	//Init(Admin)
	//e := NewSystemError("错 错:错", 01)
	//t.Log(e)
	//trans := func(err error) error {
	//	return err
	//}
	//
	//t.Logf("%+v\n", ParseErr(trans(e)).Message())
	//t.Logf("%+v\n", ParseErr(errors.New("xx x :")).Message())
}

func TestNewCodeError(t *testing.T) {
	//Init(Admin, WithErrMsgMap(map[int]string{
	//	1: "this is test",
	//}))
	//e := NewSystemCodeError(2)
	//t.Log(e)
}

func TestParseErr(t *testing.T) {
	//var err error
	//err = NewParamError("a", 1)
	//s, _ := status.FromError(err)
	//t.Log(s.Err())
	//errP := ParseErr(s.Err())
	//t.Log(errP.Message())
	//t.Log(errP.Code())
	//t.Log(errP.Error())
	//println("")
	//errP2 := ParseErr(err)
	//t.Log(errP2.Message())
	//t.Log(errP2.Code())
	//t.Log(errP2.Error())
}

func TestI18n_Msg(t *testing.T) {
	//Init(Admin, WithLocalize(map[int]*i18n.Message{
	//	1: {ID: "test", Other: "测试"},
	//	//2: {ID: "test2", Other: "测试2"},
	//}, []string{"./i18n/zh_cn.toml"}, "zh_cn"), WithUnknownMsg("Known Error"))
	//
	//err := NewParamCodeError(2)
	//t.Log(err)
}
