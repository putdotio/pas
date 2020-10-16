package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
)

type Integer struct{}

func newInteger(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Integer{}, nil
}

func (t Integer) ColumnType() string {
	return "int"
}

func (t Integer) ConvertValue(v interface{}) (driver.Value, error) {
	i, err := t.int64(v)
	if err != nil {
		return nil, err
	}
	if i > math.MaxInt32 {
		return nil, fmt.Errorf("number is greater then max integer value: %v", v)
	}
	return i, nil
}

func (t Integer) int64(v interface{}) (int64, error) {
	switch i := v.(type) {
	case json.Number:
		return i.Int64()
	case int:
		return int64(i), nil
	case int64:
		return i, nil
	case int32:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case uint:
		if i > math.MaxUint32 {
			return 0, fmt.Errorf("number is greater then max integer value: %v", v)
		}
		return int64(i), nil
	case uint64:
		if i > math.MaxInt64 {
			return 0, fmt.Errorf("number is greater then max integer value: %v", v)
		}
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint8:
		return int64(i), nil
	default:
		return 0, fmt.Errorf("value is not integer: %v", v)
	}
}
