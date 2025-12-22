package request

import (
	"encoding/json"
	"errors"
	"fiber-ee/internal/pkg/i18n"
	"fiber-ee/internal/pkg/response"
	"fiber-ee/internal/pkg/validator"

	validatorLib "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Validate struct {
}

// NewValidate 创建 Validate 实例（用于 dig 注入）
func NewValidate() Validate {
	return Validate{}
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
		// 没有翻译器，返回原始错误
		var errMsgs []string
		for _, e := range validationErrors {
			errMsgs = append(errMsgs, e.Error())
		}
		return response.ErrParams.WithData(errMsgs)
	}

	// JSON 类型错误处理
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		msg := i18n.TWithData(c, "field_type_error", map[string]any{
			"Field": unmarshalErr.Field,
			"Type":  unmarshalErr.Type.String(),
		})
		return response.ErrParamsType.WithMsg(msg)
	}

	// JSON 语法错误
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return response.ErrParamsType.WithMsg(i18n.T(c, "json_syntax_error"))
	}

	// 其他绑定错误
	if err.Error() != "" {
		return response.ErrParamsType.WithMsg(err.Error())
	}

	return response.ErrParams
}
