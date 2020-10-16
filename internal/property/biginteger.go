package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
)

type BigInteger struct{}

func newBigInteger(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return BigInteger{}, nil
}

func (t BigInteger) ColumnType() string {
	return "bigint"
}

func (t BigInteger) ConvertValue(v interface{}) (driver.Value, error) {
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
			return nil, fmt.Errorf("number is greater then max big_integer value: %v", v)
		}
		return int64(i), nil
	case uint64:
		if i > math.MaxUint64 {
			return nil, fmt.Errorf("number is greater then max big_integer value: %v", v)
		}
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint8:
		return int64(i), nil
	default:
		return nil, fmt.Errorf("value is not integer: %v", v)
	}
}
