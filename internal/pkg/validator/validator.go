package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type CustomValidator struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

func (cv *CustomValidator) Validate(out any) error {
	if err := cv.Validator.Struct(out); err != nil {
		return err
	}
	return nil
}

// TranslateError 将验证错误转换为翻译后的字符串
func (cv *CustomValidator) TranslateError(err error) []string {
	var errs []string
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			errs = append(errs, e.Translate(cv.Trans))
		}
		return errs
	}
	return []string{err.Error()}
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// 注册函数以从 json, query, form tag 获取结构体字段名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 优先级：json > query > form
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		if name == "-" {
			return ""
		}
		return name
	})

	zhT := zh.New()
	uni := ut.New(zhT, zhT)

	trans, _ := uni.GetTranslator("zh")
	if err := zh_translations.RegisterDefaultTranslations(v, trans); err != nil {
		panic(err)
	}

	return &CustomValidator{
		Validator: v,
		Trans:     trans,
	}
}
