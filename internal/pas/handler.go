package pas

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	http.Handler
	analytics *Analytics
}

func NewHandler(analytics *Analytics) *Handler {
	h := &Handler{
		analytics: analytics,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", h.handleEvents)
	mux.HandleFunc("/api/users", h.handleUsers)
	h.Handler = mux
	return h
}

type Events struct {
	Events []Event `json:"events"`
}

func (s *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var events Events
	err := dec.Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = s.analytics.InsertEvents(events.Events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Users struct {
	Users []User `json:"users"`
}

func (s *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var users Users
	err := dec.Decode(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = s.analytics.UpdateUsers(users.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
