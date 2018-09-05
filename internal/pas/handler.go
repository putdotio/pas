package pas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

type Handler struct {
	http.Handler
	analytics *Analytics
	secret    []byte
}

func NewHandler(analytics *Analytics, secret string) *Handler {
	h := &Handler{
		analytics: analytics,
		secret:    []byte(secret),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", h.handleEvents)
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/health", h.handleHealth)
	h.Handler = mux
	return h
}

type eventsRequest struct {
	Events []Event `json:"events"`
}

func (s *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var events eventsRequest
	err := dec.Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(s.secret) > 0 {
		hash := hmac.New(sha256.New, s.secret)
		for _, e := range events.Events {
			hash.Write([]byte(e.UserID))
			if hex.EncodeToString(hash.Sum(nil)) != e.UserHash {
				http.Error(w, "invalid user_hash", http.StatusBadRequest)
				return
			}
			hash.Reset()
		}
	}
	_, err = s.analytics.InsertEvents(events.Events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type usersRequest struct {
	Users []User `json:"users"`
}

func (s *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var users usersRequest
	err := dec.Decode(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(s.secret) > 0 {
		hash := hmac.New(sha256.New, s.secret)
		for _, u := range users.Users {
			hash.Write([]byte(u.ID))
			if hex.EncodeToString(hash.Sum(nil)) != u.Hash {
				http.Error(w, "invalid user_hash", http.StatusBadRequest)
				return
			}
			hash.Reset()
		}
	}
	_, err = s.analytics.UpdateUsers(users.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	err := s.analytics.Health()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
