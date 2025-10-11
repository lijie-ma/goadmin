package util

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type DateTime time.Time

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if nil != err {
		return
	}
	time.Local = loc
}

func (t DateTime) IsZero() bool {
	return time.Time(t).IsZero()
}

// gin 绑定
func (t *DateTime) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	if len(param) == 10 {
		param += " 00:00:00"
	}

	parsed, err := time.Parse(time.DateTime, param)
	if err != nil {
		return err
	}

	*t = DateTime(parsed)
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(t).Format(time.DateTime))), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsed, err := time.Parse(time.DateTime, str)
	if err != nil {
		return err
	}
	*t = DateTime(parsed)
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

func (t *DateTime) Scan(value interface{}) error {
	if value == nil {
		*t = DateTime(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = DateTime(v)
	case []byte:
		return t.UnmarshalJSON(v)
	case string:
		// 如果是 CURRENT_TIMESTAMP 等特殊值，返回当前时间
		if v == "CURRENT_TIMESTAMP" {
			*t = Now()
			return nil
		}
		return t.UnmarshalJSON([]byte(v))
	case int64:
		*t = DateTime(time.Unix(v, 0))
	default:
		return fmt.Errorf("无法扫描类型%T到DateTime", value)
	}
	return nil
}

func (t DateTime) String() string {
	return time.Time(t).Format(time.DateTime)
}

func Now() DateTime {
	return DateTime(time.Now())
}

func ZNow() time.Time {
	return time.Now().In(time.Local)
}

func ZParse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}
