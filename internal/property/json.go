package property

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON struct{}

func newJSON(args []string) (typeInterface, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("invalid arguments: %q", args)
	}
	return JSON{}, nil
}

func (j JSON) ColumnType() string {
	return "json"
}

func (j JSON) ConvertValue(v interface{}) (driver.Value, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("value is not json: %v", v)
	}
	return b, nil
}
