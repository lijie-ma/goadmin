package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Middleware 自动解析 Accept-Language，并注入 Localizer 到 gin.Context
func Middleware() gin.HandlerFunc {
	matcher := language.NewMatcher(Bundle.LanguageTags())

	return func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		tags, _, _ := language.ParseAcceptLanguage(accept)
		tag, _, _ := matcher.Match(tags...)

		localizer := i18n.NewLocalizer(Bundle, tag.String())
		c.Set("localizer", localizer)

		c.Next()
	}
}
