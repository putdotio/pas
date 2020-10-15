package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Json struct{}

func newJson(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return Json{}, nil
}

func (j Json) ColumnType() string {
	return "json"
}

func (j Json) ConvertValue(v interface{}) (driver.Value, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("value is not json: %v", v)
	}
	return b, nil
}
