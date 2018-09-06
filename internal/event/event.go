package event

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/putdotio/pas/internal/property"
	"github.com/putdotio/pas/internal/user"
)

type Event struct {
	UserID     user.ID                       `json:"user_id"`
	UserHash   string                        `json:"user_hash"`
	Timestamp  *time.Time                    `json:"timestamp"`
	Name       Name                          `json:"name"`
	Properties map[property.Name]interface{} `json:"properties"`
}

type Name string

var eventNameRegex = regexp.MustCompile(`[a-z0-9]+_[a-z0-9]+`)

func (n *Name) UnmarshalJSON(b []byte) error {
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
	*n = Name(s)
	return nil
}
