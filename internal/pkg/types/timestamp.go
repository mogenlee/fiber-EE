package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// Timestamp 自定义时间戳类型，数据库存 int，JSON 输出日期格式
type Timestamp int64

const DateTimeFormat = "2006-01-02 15:04:05"

// MarshalJSON 序列化为日期字符串
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t == 0 {
		return []byte("null"), nil
	}
	tm := time.Unix(int64(t), 0)
	return []byte(fmt.Sprintf(`"%s"`, tm.Format(DateTimeFormat))), nil
}

// UnmarshalJSON 从日期字符串或时间戳解析
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` {
		*t = 0
		return nil
	}
	// 尝试解析为数字
	s = s[1 : len(s)-1] // 去掉引号
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		*t = Timestamp(ts)
		return nil
	}
	// 尝试解析为日期字符串
	tm, err := time.ParseInLocation(DateTimeFormat, s, time.Local)
	if err != nil {
		return err
	}
	*t = Timestamp(tm.Unix())
	return nil
}

// Value 实现 driver.Valuer，存入数据库
func (t Timestamp) Value() (driver.Value, error) {
	return int64(t), nil
}

// Scan 实现 sql.Scanner，从数据库读取
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*t = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		*t = Timestamp(v)
	case int32:
		*t = Timestamp(v)
	case int:
		*t = Timestamp(v)
	default:
		return fmt.Errorf("cannot scan %T into Timestamp", value)
	}
	return nil
}

// Time 转换为 time.Time
func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// Now 返回当前时间戳
func Now() Timestamp {
	return Timestamp(time.Now().Unix())
}
