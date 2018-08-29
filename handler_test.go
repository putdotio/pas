package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const localDSN = "root@(127.0.0.1:3306)/test"

func init() {
	var err error
	db, err = sql.Open("mysql", localDSN)
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostEvents(t *testing.T) {
	s := `{
		"events": [
		{"name": "test", "user_id": "1234", "timestamp": "2000-01-01T01:02:03Z", "properties": [
				{"name": "foo", "value": "bar", "type": "string"}
		]}]}
	`
	var postBody = bytes.NewBufferString(s)
	req, err := http.NewRequest("POST", "/api/events", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleEvents)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Body.String())
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPostUsers(t *testing.T) {
	s := `{
		"users": [
		{"id": "1234", "properties": [
				{"name": "foo", "value": "bar", "type": "string"}
		]}]}
	`
	var postBody = bytes.NewBufferString(s)
	req, err := http.NewRequest("POST", "/api/users", postBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Body.String())
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
