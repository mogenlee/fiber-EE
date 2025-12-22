package i18n

import (
	contribi18n "github.com/gofiber/contrib/v3/i18n"
	"github.com/gofiber/fiber/v3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// New 创建 i18n 中间件
func New() fiber.Handler {
	return contribi18n.New(&contribi18n.Config{
		RootPath:        "./locales",
		AcceptLanguages: []language.Tag{language.Chinese, language.English},
		DefaultLanguage: language.Chinese,
		LangHandler: func(c fiber.Ctx, defaultLang string) string {
			// 解析 Accept-Language，匹配到 en 系列都返回 en
			accept := c.Get("Accept-Language")
			if accept == "" {
				return defaultLang
			}
			tags, _, _ := language.ParseAcceptLanguage(accept)
			matcher := language.NewMatcher([]language.Tag{language.Chinese, language.English})
			_, idx, _ := matcher.Match(tags...)
			switch idx {
			case 1:
				return "en"
			default:
				return "zh"
			}
		},
	})
}

// T 获取翻译文本（简单 key）
func T(c fiber.Ctx, key string) string {
	msg, err := contribi18n.Localize(c, key)
	if err != nil {
		return key
	}
	return msg
}

// TWithData 获取翻译文本（带模板数据）
func TWithData(c fiber.Ctx, key string, data map[string]any) string {
	msg, err := contribi18n.Localize(c, &i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		return key
	}
	return msg
}
