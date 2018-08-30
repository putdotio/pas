package pas_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/putdotio/pas/internal/pas"
	"github.com/stretchr/testify/assert"
)

func TestInsertEvent(t *testing.T) {
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
