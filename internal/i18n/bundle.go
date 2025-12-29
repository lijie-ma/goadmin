package i18n

import (
	"embed"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/**/*.toml
var LocaleFS embed.FS

var Bundle *i18n.Bundle

// Init 初始化 i18n Bundle，加载多语言文件
func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	dirs := []string{"locales/common", "locales/module"}
	for _, dir := range dirs {
		files, _ := LocaleFS.ReadDir(dir)
		for _, f := range files {
			data, _ := LocaleFS.ReadFile(dir + "/" + f.Name())
			_, _ = Bundle.ParseMessageFileBytes(data, f.Name())
		}
	}

	log.Printf("[i18n] Loaded locales: %v\n", Bundle.LanguageTags())
}
