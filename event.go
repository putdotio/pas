package main

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

type Event struct {
	UserID     UserID     `json:"user_id"`
	Timestamp  *time.Time `json:"timestamp"`
	Name       EventName  `json:"name"`
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

type EventName string

var eventNameRegex = regexp.MustCompile(`[a-z]+[a-z_0-9]*`)

func (n *EventName) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if len(*n) > 255 {
		return errors.New("event name too big")
	}
	if !eventNameRegex.MatchString(s) {
		return errors.New("invalid event name")
	}
	*n = EventName(s)
	return nil
}
