package property

import (
	"database/sql/driver"
	"fmt"
)

type String struct{}

func newString(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return String{}, nil
}

func (t String) ColumnType() string {
	return "varchar(2000)"
}

func (t String) ConvertValue(v interface{}) (driver.Value, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return nil, fmt.Errorf("value is not string: %v", v)
}
