package property

import (
	"database/sql/driver"
	"fmt"
)

type Boolean struct{}

func newBoolean(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Boolean{}, nil
}

func (t Boolean) ColumnType() string {
	return "tinyint(1)"
}

func (t Boolean) ConvertValue(v interface{}) (driver.Value, error) {
	if i, ok := v.(bool); ok {
		return i, nil
	}
	return nil, fmt.Errorf("value is not boolean: %v", v)
}
