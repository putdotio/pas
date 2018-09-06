package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Float struct{}

func newFloat(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Float{}, nil
}

func (t Float) ColumnType() string {
	return "float(23, 3)"
}

func (t Float) ConvertValue(v interface{}) (driver.Value, error) {
	switch v := v.(type) {
	case json.Number:
		return v.Float64()
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return nil, fmt.Errorf("value is not float: %v", v)
	}
}
