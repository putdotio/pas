package property

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Decimal struct {
	Precision int
	Scale     int
}

func newDecimal(args []string) (typeInterface, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	precision, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}
	scale, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	return Decimal{Precision: precision, Scale: scale}, nil
}

func (t Decimal) ColumnType() string {
	return fmt.Sprintf("decimal(%d,%d)", t.Precision, t.Scale)
}

func (t Decimal) ConvertValue(v interface{}) (driver.Value, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return nil, fmt.Errorf("value is not decimal: %v", v)
}
