package request

import (
	"errors"
	"fiber-ee/internal/pkg/i18n"
	"fiber-ee/internal/pkg/response"
	"fiber-ee/internal/pkg/validator"
	"fmt"
	"regexp"
	"strings"

	"github.com/bytedance/sonic/decoder"
	validatorLib "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var fieldRegexp = regexp.MustCompile(`"([^"]+)"\s*:`)
var mismatchTypeRe = regexp.MustCompile(`Mismatch type (\w+) with value (\w+)`)

type Validate struct {
}

// NewValidate 创建 Validate 实例（用于 dig 注入）
func NewValidate() Validate {
	return Validate{}
}
func ExtractFieldFromMismatch(desc string) string {
	// 找到 ^ 的位置
	caret := strings.Index(desc, "^")
	if caret == -1 {
		return ""
	}

	// 只看 ^ 之前的内容
	left := desc[:caret]

	// 在左侧找最后一个 "field":
	matches := fieldRegexp.FindAllStringSubmatch(left, -1)
	if len(matches) == 0 {
		return ""
	}

	return matches[len(matches)-1][1]
}

func ExtractTypeMismatch(err error) (expect, actual string, ok bool) {
	var me *decoder.MismatchTypeError
	is := errors.As(err, &me)
	if !is {
		return "", "", false
	}
	m := mismatchTypeRe.FindStringSubmatch(me.Error())
	if len(m) != 3 {
		return "", "", false
	}
	return m[1], m[2], true
}

// BindAndValidate 统一处理请求 Body 绑定和校验 (JSON/Form)
func (v Validate) BindAndValidate(c fiber.Ctx, obj any) error {
	if err := c.Bind().Body(obj); err != nil {
		return v.handleError(c, err)
	}
	return nil
}

// BindQuery 统一处理 Query 参数绑定和校验
func (v Validate) BindQuery(c fiber.Ctx, obj any) error {
	if err := c.Bind().Query(obj); err != nil {
		return v.handleError(c, err)
	}
	return nil
}

// BindForm 统一处理 Form 参数绑定和校验 (支持 Multipart)
func (v Validate) BindForm(c fiber.Ctx, obj any) error {
	if err := c.Bind().Form(obj); err != nil {
		return v.handleError(c, err)
	}
	return nil
}

// Bind 自动根据 Content-Type 尝试绑定 Body 或 Query
func (v Validate) Bind(c fiber.Ctx, obj any) error {
	if err := c.Bind().All(obj); err != nil {
		return v.handleError(c, err)
	}
	return nil
}

// handleError 内部错误处理与翻译
func (v Validate) handleError(c fiber.Ctx, err error) error {
	// 检查是否是验证错误
	var validationErrors validatorLib.ValidationErrors
	if errors.As(err, &validationErrors) {
		vc := c.App().Config().StructValidator
		if cv, ok := vc.(*validator.CustomValidator); ok {
			errs := cv.TranslateError(err)
			return response.ErrParams.WithMsg(i18n.T(c, "params_error")).WithData(errs)
		}
		var errMsgs []string
		for _, e := range validationErrors {
			errMsgs = append(errMsgs, e.Error())
		}

		fmt.Println(validationErrors)
		return response.ErrParams.WithData(errMsgs)
	}

	// Sonic 语法错误
	var syntaxError decoder.SyntaxError
	if errors.As(err, &syntaxError) {
		return response.ErrParamsType.WithMsg(i18n.T(c, "json_syntax_error"))
	}

	// Sonic 类型不匹配错误
	var me *decoder.MismatchTypeError
	if errors.As(err, &me) {
		field := ExtractFieldFromMismatch(me.Description())
		expect, actual, _ := ExtractTypeMismatch(me)
		msg := i18n.TWithData(c, "field_type_error", map[string]any{
			"Field":  field,
			"Type":   expect,
			"Actual": actual,
		})
		return response.ErrParamsType.WithMsg(msg)
	}

	// 其他绑定错误
	if err.Error() != "" {
		return response.ErrParamsType.WithMsg(err.Error())
	}

	return response.ErrParams
}
