package i18n

import (
	"embed"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var LocaleFS embed.FS

var Bundle *i18n.Bundle

// Init 初始化 i18n Bundle，加载多语言文件
func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files := []string{
		"locales/active.en.toml",
		"locales/active.zh.toml",
	}

	for _, f := range files {
		if _, err := Bundle.LoadMessageFileFS(LocaleFS, f); err != nil {
			log.Fatalf("[i18n] failed to load locale file %s: %v", f, err)
		}
	}

	log.Printf("[i18n] Loaded locales: %v\n", Bundle.LanguageTags())
}
