package main

import (
	"encoding/json"
	"errors"
)

type User struct {
	ID         UserID     `json:"id"`
	Properties []Property `json:"properties"`
}

type UserID string

func (u *UserID) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if len(s) > 255 {
		return errors.New("user_id too big")
	}
	*u = UserID(s)
	return nil
}
