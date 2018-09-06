package property

import (
	"database/sql/driver"
	"fmt"
)

type Date struct{}

func newDate(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Date{}, nil
}

func (t Date) ColumnType() string {
	return "date"
}

func (t Date) ConvertValue(v interface{}) (driver.Value, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return nil, fmt.Errorf("value is not date: %v", v)
}
