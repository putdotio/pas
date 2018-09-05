package pas_test

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

	"github.com/putdotio/pas/internal/pas"
)

const localDSN = "root@(127.0.0.1:3306)/test"

var handler *pas.Handler

func init() {
	db, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}

	analytics := pas.NewAnalytics(db)

	handler = pas.NewHandler(analytics, "")
}

func TestPostEvents(t *testing.T) {
	s := `{
		"events": [
		{"name": "test_done", "user_id": "1234", "timestamp": "2000-01-01T01:02:03Z", "properties": [
				{"name": "foo", "value": "bar", "type": "string"}
		]}]}
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
	s := `{
		"users": [
		{"id": "1234", "properties": [
				{"name": "foo", "value": 1, "type": "string"}
		]}]}
	`
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

	analytics := pas.NewAnalytics(db)

	handler := pas.NewHandler(analytics, secret)

	s0 := `{
		"events": [
		{"name": "test_done", "user_id": "1234", "user_hash": "%s", "timestamp": "2000-01-01T01:02:03Z", "properties": [
				{"name": "foo", "value": "bar", "type": "string"}
		]}]}
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
