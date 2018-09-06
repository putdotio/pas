package handler_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/putdotio/pas/internal/analytics"
	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/handler"
	"github.com/putdotio/pas/internal/property"
)

const localDSN = "root@(127.0.0.1:3306)/test"

func TestPostEvents(t *testing.T) {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	events := map[event.Name]property.Types{
		"test_done": property.Types{
			"foo": property.Must(property.New("string")),
		},
	}
	analytics := analytics.New(db, nil, events)
	handler := handler.New(analytics, "")

	s := `{
		"events": [
		{"name": "test_done", "user_id": "1234", "timestamp": "2000-01-01T01:02:03Z", "properties": {
				"foo": "bar"
		}}]}
	`
	var postBody = bytes.NewBufferString(s)
	req, err := http.NewRequest("POST", "/api/events", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Log(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPostUsers(t *testing.T) {
	const secret = "foobar"
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user := property.Types{
		"foo": property.Must(property.New("integer")),
	}
	analytics := analytics.New(db, user, nil)
	handler := handler.New(analytics, secret)

	s := `{
		"users": [
		{"id": "1234", "hash": "%s", "properties": {
				"foo": 1
		}}]}
	`
	s = fmt.Sprintf(s, generateUserHash("1234", secret))
	var postBody = bytes.NewBufferString(s)
	req, err := http.NewRequest("POST", "/api/users", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Log(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUserHash(t *testing.T) {
	const secret = "foobar"

	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	events := map[event.Name]property.Types{
		"test_done": property.Types{
			"foo": property.Must(property.New("string")),
		},
	}
	analytics := analytics.New(db, nil, events)
	handler := handler.New(analytics, secret)

	s0 := `{
		"events": [
		{"name": "test_done", "user_id": "1234", "user_hash": "%s", "timestamp": "2000-01-01T01:02:03Z", "properties": {
				"foo": "bar"
		}}]}
	`

	// Test invalid secret
	s := fmt.Sprintf(s0, generateUserHash("1234", "invalid"))
	var postBody = bytes.NewBufferString(s)
	req, err := http.NewRequest("POST", "/api/events", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Log(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Test correct secret
	s = fmt.Sprintf(s0, generateUserHash("1234", secret))
	postBody = bytes.NewBufferString(s)
	req, err = http.NewRequest("POST", "/api/events", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Log(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func generateUserHash(userID, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(userID))
	return hex.EncodeToString(hash.Sum(nil))
}
