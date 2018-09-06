package property

import (
	"encoding/json"
	"errors"
	"regexp"
)

type Name string

var propertyNameRegex = regexp.MustCompile(`[a-z]+[a-z_0-9 \-]*`)

func (n *Name) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if len(*n) > 255 {
		return errors.New("property name too big")
	}
	if !propertyNameRegex.MatchString(s) {
		return errors.New("invalid property name")
	}
	*n = Name(s)
	return nil
}
