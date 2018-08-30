package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Events struct {
	Events []Event `json:"events"`
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var events Events
	err := dec.Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	for _, e := range events.Events {
		i := EventInserter{e}
		err = runInserter(i, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type Users struct {
	Users []User `json:"users"`
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var users Users
	err := dec.Decode(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	for _, u := range users.Users {
		i := UserInserter{u}
		err = runInserter(i, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
