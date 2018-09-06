package user

import (
	"encoding/json"
	"errors"

	"github.com/putdotio/pas/internal/property"
)

type User struct {
	ID         ID                            `json:"id"`
	Hash       string                        `json:"hash"`
	Properties map[property.Name]interface{} `json:"properties"`
}

type ID string

func (u *ID) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if len(s) > 255 {
		return errors.New("user_id too big")
	}
	*u = ID(s)
	return nil
}
