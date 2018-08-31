package pas_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/putdotio/pas/internal/pas"
	"github.com/stretchr/testify/assert"
)

func TestInsertEvents(t *testing.T) {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	analytics := pas.NewAnalytics(db)

	_, err = db.Exec("drop table if exists page_viewed")
	if err != nil {
		t.Fatal(err)
	}

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	e := pas.Event{
		UserID:    pas.UserID(1234),
		Timestamp: &ts,
		Name:      "page_viewed",
		Properties: []pas.Property{
			{
				Type:  pas.TypeInteger,
				Name:  "foo",
				Value: 1,
			},
			{
				Type:  pas.TypeString,
				Name:  "type_string",
				Value: "test",
			},
			{
				Type:  pas.TypeBoolean,
				Name:  "type_boolean",
				Value: true,
			},
			{
				Type:  pas.TypeFloat,
				Name:  "type_float",
				Value: 123.456,
			},
			{
				Type:  pas.TypeDecimal,
				Name:  "type_decimal",
				Value: 123.456,
			},
			{
				Type:  pas.TypeDateTime,
				Name:  "type_datetime",
				Value: "2010-02-03T01:02:03",
			},
		},
	}
	events := []pas.Event{e}

	// will create table
	n, err := analytics.InsertEvents(events)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)

	// will insert
	n, err = analytics.InsertEvents(events)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)

	// will add column
	events[0].Properties = append(events[0].Properties, pas.Property{
		Type:  pas.TypeString,
		Name:  "bar",
		Value: "test",
	})
	n, err = analytics.InsertEvents(events)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)
}

func TestUpdateUsers(t *testing.T) {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	analytics := pas.NewAnalytics(db)

	_, err = db.Exec("drop table if exists user")
	if err != nil {
		t.Fatal(err)
	}

	u := pas.User{
		ID: pas.UserID(1234),
		Properties: []pas.Property{
			{
				Type:  pas.TypeInteger,
				Name:  "foo",
				Value: 1,
			},
		},
	}
	users := []pas.User{u}

	// will create table
	n, err := analytics.UpdateUsers(users)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)

	// will insert
	n, err = analytics.UpdateUsers(users)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)

	// will add column
	users[0].Properties = append(users[0].Properties, pas.Property{
		Type:  pas.TypeString,
		Name:  "bar",
		Value: "test",
	})
	n, err = analytics.UpdateUsers(users)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)
}
