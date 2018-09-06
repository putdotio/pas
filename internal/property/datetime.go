package property

import (
	"database/sql/driver"
	"fmt"
)

type DateTime struct{}

func newDateTime(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return DateTime{}, nil
}

func (t DateTime) ColumnType() string {
	return "datetime"
}

func (t DateTime) ConvertValue(v interface{}) (driver.Value, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return nil, fmt.Errorf("value is not datetime: %v", v)
}
