package types

import (
	"database/sql/driver"
	"fmt"
)

// Bool 自定义布尔类型，数据库存 int32 (0/1)，JSON 输出 true/false
type Bool int32

// MarshalJSON 序列化为 true/false
func (b Bool) MarshalJSON() ([]byte, error) {
	if b == 0 {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// UnmarshalJSON 从 true/false 或 0/1 解析
func (b *Bool) UnmarshalJSON(data []byte) error {
	s := string(data)
	switch s {
	case "true", "1":
		*b = 1
	default:
		*b = 0
	}
	return nil
}

// Value 实现 driver.Valuer
func (b Bool) Value() (driver.Value, error) {
	return int64(b), nil
}

// Scan 实现 sql.Scanner
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		*b = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		*b = Bool(v)
	case int32:
		*b = Bool(v)
	case int:
		*b = Bool(v)
	case bool:
		if v {
			*b = 1
		} else {
			*b = 0
		}
	default:
		return fmt.Errorf("cannot scan %T into Bool", value)
	}
	return nil
}

// IsTrue 判断是否为真
func (b Bool) IsTrue() bool {
	return b != 0
}
