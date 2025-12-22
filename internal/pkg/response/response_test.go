package response

import (
	"testing"
)

func TestRespTypeError(t *testing.T) {
	err := ErrParams
	expected := "2001: 参数校验错误"

	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestRespTypeWithMsg(t *testing.T) {
	err := ErrParams.WithMsg("自定义消息")

	if err.Msg() != "自定义消息" {
		t.Errorf("Expected '自定义消息', got '%s'", err.Msg())
	}

	// 原始错误不应被修改
	if ErrParams.Msg() != "参数校验错误" {
		t.Errorf("Original error was modified")
	}
}

func TestRespTypeWithData(t *testing.T) {
	data := map[string]string{"field": "error"}
	err := ErrParams.WithData(data)

	if err.Data() == nil {
		t.Error("Expected data, got nil")
	}
}

func TestRespTypeCodes(t *testing.T) {
	tests := []struct {
		resp RespType
		code int
	}{
		{OK, 200},
		{ErrFailed, 1000},
		{ErrParams, 2001},
		{ErrQuery, 3001},
		{ErrUnauthorized, 4001},
		{ErrInternal, 5000},
	}

	for _, tt := range tests {
		if tt.resp.Code() != tt.code {
			t.Errorf("Expected code %d for %s, got %d", tt.code, tt.resp.Msg(), tt.resp.Code())
		}
	}
}

func TestHttpStatus(t *testing.T) {
	tests := []struct {
		code       int
		httpStatus int
	}{
		{200, 200},
		{2001, 400},
		{4001, 401},
		{4003, 403},
		{5000, 500},
		{1000, 200},
	}

	for _, tt := range tests {
		status := httpStatus(tt.code)
		if status != tt.httpStatus {
			t.Errorf("Expected HTTP status %d for code %d, got %d", tt.httpStatus, tt.code, status)
		}
	}
}
