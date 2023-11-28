package i18n

import "testing"

func TestI18n_Msg(t *testing.T) {
	Init("./lang")
	n := NewI18n()
	n.SetLanguage("zh_cn1")
	t.Log(n.Tfd("server.error", map[string]interface{}{}))
}
