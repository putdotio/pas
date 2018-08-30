package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
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
		if e.Timestamp == nil {
			e.Timestamp = &now
		}
		sql, values := insertEvent(e)
		_, err = db.Exec(sql, values...)
		if merr, ok := err.(*mysql.MySQLError); ok {
			if merr.Number == 1146 { // table doesn't exist
				_, err = db.Exec(createEventTable(e))
				if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1050 { // table already exists
					err = nil
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(sql, values...)
			} else if merr.Number == 1054 { // unknown column
				cols, err2 := existingEventColumns(string(e.Name))
				if err2 != nil {
					http.Error(w, err2.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(alterEventTable(e, cols))
				if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1060 { // duplicate column name
					err = nil
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(sql, values...)
			}
		}
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
	for _, u := range users.Users {
		sql, values := insertUser(u)
		_, err = db.Exec(sql, values...)
		if merr, ok := err.(*mysql.MySQLError); ok {
			if merr.Number == 1146 { // table doesn't exist
				_, err = db.Exec(createUserTable(u))
				if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1050 { // table already exists
					err = nil
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(sql, values...)
			} else if merr.Number == 1054 { // unknown column
				cols, err2 := existingUserColumns()
				if err2 != nil {
					http.Error(w, err2.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(alterUserTable(u, cols))
				if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1060 { // duplicate column name
					err = nil
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.Exec(sql, values...)
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
