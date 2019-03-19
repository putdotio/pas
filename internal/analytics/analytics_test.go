package analytics_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"testing"
	"time"

	"github.com/putdotio/pas/internal/analytics"
	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/property"
	"github.com/putdotio/pas/internal/user"
	"github.com/stretchr/testify/assert"
)

const localDSN = "root@(127.0.0.1:3306)/test"
const secret = "foobar"

func TestInsertEvents(t *testing.T) {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	types := map[event.Name]property.Types{
		"page_viewed": property.Types{
			"foo":           property.Must(property.New("integer")),
			"type_string":   property.Must(property.New("string")),
			"type_boolean":  property.Must(property.New("boolean")),
			"type_float":    property.Must(property.New("float")),
			"type_decimal":  property.Must(property.New("decimal(6, 3)")),
			"type_datetime": property.Must(property.New("datetime")),
		},
	}

	analytics := analytics.New(db, secret, nil, types)

	_, err = db.Exec("drop table if exists page_viewed")
	if err != nil {
		t.Fatal(err)
	}

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	e := event.Event{
		UserID:    "1234",
		Timestamp: &ts,
		Name:      "page_viewed",
		Properties: map[property.Name]interface{}{
			"foo":           1,
			"type_string":   "test",
			"type_boolean":  true,
			"type_float":    123.456,
			"type_decimal":  "123.456",
			"type_datetime": "2010-02-03T01:02:03",
		},
	}
	events := []event.Event{e}

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
	types["page_viewed"]["bar"] = property.Must(property.New("string"))
	events[0].Properties["bar"] = "test"
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

	types := property.Types{
		"foo": property.Must(property.New("integer")),
	}
	analytics := analytics.New(db, secret, types, nil)

	_, err = db.Exec("drop table if exists user")
	if err != nil {
		t.Fatal(err)
	}

	u := user.User{
		ID:   "1234",
		Hash: generateUserHash("1234", secret),
		Properties: map[property.Name]interface{}{
			"foo": 1,
		},
	}
	users := []user.User{u}

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
	types["bar"] = property.Must(property.New("string"))
	users[0].Properties["bar"] = "test"
	n, err = analytics.UpdateUsers(users)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)
}

func TestAlias(t *testing.T) {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	types := map[event.Name]property.Types{
		"page_viewed": property.Types{
			"foo": property.Must(property.New("integer")),
		},
	}

	analytics := analytics.New(db, secret, nil, types)

	_, err = db.Exec("drop table if exists page_viewed")
	if err != nil {
		t.Fatal(err)
	}

	e := event.Event{
		UserID: "1",
		Name:   "page_viewed",
		Properties: map[property.Name]interface{}{
			"foo": 1,
		},
	}
	events := []event.Event{e}

	n, err := analytics.InsertEvents(events)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, 1)

	var userID string
	var isAnonymous bool
	row := db.QueryRow("select user_id, is_anonymous from page_viewed")
	err = row.Scan(&userID, &isAnonymous)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, userID, "1")
	assert.True(t, isAnonymous)

	err = analytics.Alias("1", "2", generateUserHash("2", secret))
	if err != nil {
		t.Fatal(err)
	}
	row = db.QueryRow("select user_id, is_anonymous from page_viewed")
	err = row.Scan(&userID, &isAnonymous)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, userID, "2")
	assert.False(t, isAnonymous)
}

func generateUserHash(userID, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	_, _ = hash.Write([]byte(userID))
	return hex.EncodeToString(hash.Sum(nil))
}
