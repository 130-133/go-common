package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
)

type I18n struct {
	language  string
	localizer *i18n.Localizer
}

var bundle *i18n.Bundle

func NewI18n() *I18n {
	ins := &I18n{}
	ins.SetLanguage("en")
	return ins
}

// SetLanguage 设置当前会话的语言
func (n *I18n) SetLanguage(language string) {
	n.language = language
	n.localizer = i18n.NewLocalizer(bundle, language)
}

func Init(path string) {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	files := getFiles(path)
	if len(files) == 0 {
		panic("no i18n config file")
	}
	for _, v := range files {
		bundle.MustLoadMessageFile(v)
	}
}

func (n *I18n) Tfd(format string, value map[string]interface{}) string {
	res, err := n.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    format,
		TemplateData: value,
	})
	if err != nil || res == "" {
		return format
	}
	return res
}

func getFiles(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var res []string
	for _, file := range files {
		if file.IsDir() {
			res = append(res, getFiles(path+"/"+file.Name())...)
		} else {
			if file.Name()[len(file.Name())-5:] != ".toml" {
				continue
			}
			res = append(res, path+"/"+file.Name())
		}
	}
	return res
}
