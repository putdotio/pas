package handler

import (
	"encoding/json"
	"net/http"

	"pas/internal/analytics"
	"pas/internal/event"
	"pas/internal/user"
)

type Handler struct {
	http.Handler
	analytics *analytics.Analytics
}

func New(analytics *analytics.Analytics) *Handler {
	h := &Handler{
		analytics: analytics,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", h.handleEvents)
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/alias", h.handleAlias)
	mux.HandleFunc("/health", h.handleHealth)
	h.Handler = mux
	return h
}

func successResponse(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

type eventsRequest struct {
	Events []event.Event `json:"events"`
}

func (s *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	var events eventsRequest
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
	successResponse(w)
}

type usersRequest struct {
	Users []user.User `json:"users"`
}

func (s *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	var users usersRequest
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
	successResponse(w)
}

func (s *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	err := s.analytics.Health()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	successResponse(w)
}

type aliasRequest struct {
	UserID     user.ID `json:"user_id"`
	UserHash   string  `json:"user_hash"`
	PreviousID user.ID `json:"previous_id"`
}

func (s *Handler) handleAlias(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req aliasRequest
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.analytics.Alias(req.PreviousID, req.UserID, req.UserHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	successResponse(w)
}
