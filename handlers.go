package main

import (
	"encoding/json"
	"net/http"

	"github.com/putdotio/pas/internal/pas"
)

type Events struct {
	Events []pas.Event `json:"events"`
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var events Events
	err := dec.Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = analytics.InsertEvents(events.Events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Users struct {
	Users []pas.User `json:"users"`
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var users Users
	err := dec.Decode(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = analytics.UpdateUsers(users.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
