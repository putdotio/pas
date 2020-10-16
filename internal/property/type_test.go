package property_test

import (
	"encoding/json"
	"testing"

	"github.com/putdotio/pas/internal/property"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var v struct {
		Type property.Type
	}
	s := `{"Type": "string"}`
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, v.Type.ColumnType(), "varchar(2000)")
}

func TestDecimal(t *testing.T) {
	var v struct {
		Type property.Type
	}
	s := `{"Type": "decimal(5, 2)"}`
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, v.Type.ColumnType(), "decimal(5,2)")
}
