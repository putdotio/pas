package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Double struct{}

func newDouble(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Double{}, nil
}

func (t Double) ColumnType() string {
	return "double(53, 2)"
}

func (t Double) ConvertValue(v interface{}) (driver.Value, error) {
	if i, ok := v.(json.Number); ok {
		return i.Float64()
	}
	return nil, fmt.Errorf("value is not double: %v", v)
}
