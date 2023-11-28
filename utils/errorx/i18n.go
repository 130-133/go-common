package errorx

//
//import (
//	"github.com/BurntSushi/toml"
//	"github.com/nicksnyder/go-i18n/v2/i18n"
//	"golang.org/x/text/language"
//)
//
//type I18n struct {
//	Lang       string
//	Bundle     *i18n.Bundle
//	Message    map[int]*i18n.Message
//	UnknownMsg string
//}
//
//func NewI18n(data map[int]*i18n.Message, i18nFile []string, lang string) *I18n {
//	bundle := i18n.NewBundle(language.English)
//	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
//	for _, v := range i18nFile {
//		bundle.MustLoadMessageFile(v)
//	}
//	return &I18n{
//		Bundle:     bundle,
//		Lang:       lang,
//		Message:    data,
//		UnknownMsg: "Unknown Error",
//	}
//}
//
//func (i *I18n) Msg(code int) string {
//	localize := i18n.NewLocalizer(i.Bundle, i.Lang)
//	message, ok := i.Message[code]
//	if !ok {
//		message = &i18n.Message{ID: "Unknown", Other: i.UnknownMsg}
//	}
//	msg, _ := localize.LocalizeMessage(message)
//	return msg
//}
