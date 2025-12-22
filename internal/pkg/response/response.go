package response

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// RespType 响应类型
type RespType struct {
	code int
	msg  string
	data any
}

// Response 响应格式结构
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// PageData 分页数据结构
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// 错误码规范:
// 0/200  成功
// 1xxx   通用业务错误
// 2xxx   参数相关错误
// 3xxx   数据操作错误
// 4xxx   认证授权错误
// 5xxx   系统级错误

var (
	// 成功
	OK = RespType{code: 200, msg: "操作成功"}

	// 通用业务错误 (1xxx)
	ErrFailed = RespType{code: 1000, msg: "操作失败"}

	// 参数错误 (2xxx)
	ErrParams     = RespType{code: 2001, msg: "参数校验错误"}
	ErrParamsType = RespType{code: 2002, msg: "参数类型错误"}
	ErrTooMany    = RespType{code: 2003, msg: "请求过于频繁"}

	// 数据操作错误 (3xxx)
	ErrQuery    = RespType{code: 3001, msg: "查询数据失败"}
	ErrCreate   = RespType{code: 3002, msg: "创建数据失败"}
	ErrUpdate   = RespType{code: 3003, msg: "更新数据失败"}
	ErrDelete   = RespType{code: 3004, msg: "删除数据失败"}
	ErrNotFound = RespType{code: 3005, msg: "数据不存在"}
	ErrExist    = RespType{code: 3006, msg: "数据已存在"}

	// 认证授权错误 (4xxx)
	ErrUnauthorized = RespType{code: 4001, msg: "未登录或登录已过期"}
	ErrTokenInvalid = RespType{code: 4002, msg: "token无效"}
	ErrForbidden    = RespType{code: 4003, msg: "无访问权限"}
	ErrLogin        = RespType{code: 4004, msg: "账号或密码错误"}

	// 系统错误 (5xxx)
	ErrInternal = RespType{code: 5000, msg: "系统错误"}
	ErrNotRoute = RespType{code: 5001, msg: "接口不存在"}
)

// Error 实现 error 接口
func (r RespType) Error() string {
	return strconv.Itoa(r.code) + ": " + r.msg
}

// WithMsg 自定义消息
func (r RespType) WithMsg(msg string) RespType {
	r.msg = msg
	return r
}

// WithData 附带数据
func (r RespType) WithData(data any) RespType {
	r.data = data
	return r
}

// Code 获取错误码
func (r RespType) Code() int { return r.code }

// Msg 获取消息
func (r RespType) Msg() string { return r.msg }

// Data 获取数据
func (r RespType) Data() any { return r.data }

// Success 成功响应
func Success(c fiber.Ctx, data any) error {
	return c.JSON(Response{
		Code: OK.code,
		Msg:  OK.msg,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c fiber.Ctx, resp RespType) error {
	return c.Status(httpStatus(resp.code)).JSON(Response{
		Code: resp.code,
		Msg:  resp.msg,
		Data: resp.data,
	})
}

// Error 错误处理（自动识别 RespType）
func Error(c fiber.Ctx, err error) error {
	if err == nil {
		return Success(c, nil)
	}

	var resp RespType
	if errors.As(err, &resp) {
		return Fail(c, resp)
	}

	return Fail(c, ErrInternal.WithMsg(err.Error()))
}

// Check 判断是否出现错误，并返回对应响应
func Check(c fiber.Ctx, err error) error {
	if err != nil {
		return Error(c, err)
	}
	return Success(c, nil)
}

// CheckWithData 判断是否出现错误，并返回对应响应（带data数据）
func CheckWithData(c fiber.Ctx, data any, err error) error {
	if err != nil {
		return Error(c, err)
	}
	return Success(c, data)
}

// SuccessPage 分页成功响应
func SuccessPage(c fiber.Ctx, list any, total int64, page, pageSize int) error {
	return Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// httpStatus 业务码转 HTTP 状态码
func httpStatus(code int) int {
	switch {
	case code == OK.code:
		return fiber.StatusOK
	case code >= 2000 && code < 3000:
		return fiber.StatusBadRequest
	case code >= 4000 && code < 4003:
		return fiber.StatusUnauthorized
	case code == 4003:
		return fiber.StatusForbidden
	case code >= 5000:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusOK
	}
}
