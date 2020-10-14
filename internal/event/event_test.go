package event

import (
	"encoding/json"
	"testing"

	"pas/internal/user"

	"github.com/stretchr/testify/assert"
)

func TestEventMarshal(t *testing.T) {
	e := Event{UserID: "test"}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(b), `"user_id":"test"`)
}

func TestEventUnmarshal(t *testing.T) {
	var e Event
	b := []byte(`{"user_id": "test"}`)
	err := json.Unmarshal(b, &e)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e.UserID, user.ID("test"))
}
