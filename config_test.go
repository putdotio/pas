package main

import (
	"testing"

	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/property"

	"github.com/naoina/toml"
	"github.com/stretchr/testify/assert"
)

func TestConfigUnmarshalEvents(t *testing.T) {
	s := `
	[user]
	user_property = "integer"

	[events]

	[events.test_event]
	test_property = "string"

	[events.test_event2]
	test_property2 = "string"
	`
	var c Config
	err := toml.Unmarshal([]byte(s), &c)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, c.User[property.Name("user_property")].ColumnType(), "int")

	assert.Equal(t, len(c.Events), 2)
	assert.Equal(t, c.Events[event.Name("test_event")][property.Name("test_property")].ColumnType(), "varchar(2000)")
}
