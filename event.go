package main

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

const (
	TypeString     = PropertyType("string")
	TypeInteger    = PropertyType("integer")
	TypeBigInteger = PropertyType("big_integer")
	TypeFloat      = PropertyType("float")
	TypeDecimal    = PropertyType("decimal")
	TypeBoolean    = PropertyType("boolean")
	TypeDate       = PropertyType("date")
	TypeDateTime   = PropertyType("datetime")
)

var propertyTypes = map[PropertyType]struct{}{
	TypeString:     struct{}{},
	TypeInteger:    struct{}{},
	TypeBigInteger: struct{}{},
	TypeFloat:      struct{}{},
	TypeDecimal:    struct{}{},
	TypeBoolean:    struct{}{},
	TypeDate:       struct{}{},
	TypeDateTime:   struct{}{},
}

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

type Property struct {
	Type  PropertyType `json:"type"`
	Name  PropertyName `json:"name"`
	Value interface{}  `json:"value"`
}

type PropertyName string

var propertyNameRegex = regexp.MustCompile(`[a-z]+[a-z_0-9 \-]*`)

func (n *PropertyName) UnmarshalJSON(b []byte) error {
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
	*n = PropertyName(s)
	return nil
}

type PropertyType string

func (t *PropertyType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = PropertyType(s)
	_, ok := propertyTypes[*t]
	if !ok {
		return errors.New("unknown property type")
	}
	return nil
}
