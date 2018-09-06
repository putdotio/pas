package property

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Type struct {
	typeInterface
}

func New(s string) (Type, error) {
	var t Type
	err := t.UnmarshalText([]byte(s))
	return t, err
}

func Must(t Type, err error) Type {
	if err != nil {
		panic(err)
	}
	return t
}

type typeInterface interface {
	ColumnType() string
	ConvertValue(v interface{}) (driver.Value, error)
}

var types = map[string]func([]string) (typeInterface, error){
	"string":      newString,
	"integer":     newInteger,
	"big_integer": newBigInteger,
	"float":       newFloat,
	"double":      newDouble,
	"boolean":     newBoolean,
	"date":        newDate,
	"datetime":    newDateTime,
	"decimal":     newDecimal,
}

func (t2 *Type) UnmarshalText(text []byte) error {
	s := string(text)
	s = strings.Replace(s, " ", "", -1)
	parts := strings.SplitN(s, "(", 2)
	t := parts[0]
	var args []string
	if len(parts) == 2 {
		arg := parts[1]
		if arg[len(arg)-1] != ')' {
			return fmt.Errorf("invalid type: %s", s)
		}
		arg = arg[:len(arg)-1]
		args = strings.Split(arg, ",")
	}
	f, ok := types[t]
	if !ok {
		return fmt.Errorf("unknown type: %s", t)
	}
	ty, err := f(args)
	if err != nil {
		return fmt.Errorf("invalid type: %s (%s)", t, err)
	}
	t2.typeInterface = ty
	return nil
}
