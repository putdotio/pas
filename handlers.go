package main

import (
	"encoding/json"
	"net/http"
	"strings"

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
	for _, e := range events.Events {
		var sb strings.Builder
		sb.WriteString("insert into ")
		sb.WriteString(string(e.Name))
		sb.WriteString("(user_id, timestamp")
		for _, p := range e.Properties {
			sb.WriteRune(',')
			sb.WriteString(string(p.Name))
		}
		sb.WriteString(") values (?, ?")
		for range e.Properties {
			sb.WriteString(",?")
		}
		sb.WriteRune(')')
		values := make([]interface{}, len(e.Properties)+2)
		values[0] = string(e.UserID)
		values[1] = e.Timestamp
		for i := range e.Properties {
			values[i+2] = e.Properties[i].Value
		}
		sql := sb.String()
		_, err = db.Exec(sql, values...)
		if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1146 {
			var cb strings.Builder
			cb.WriteString("create table ")
			cb.WriteString(string(e.Name))
			cb.WriteString("(user_id varchar(255) not null, timestamp datetime not null")
			for _, p := range e.Properties {
				cb.WriteRune(',')
				cb.WriteString(string(p.Name))
				cb.WriteRune(' ')
				cb.WriteString(p.DBType())
			}
			cb.WriteRune(')')
			_, err = db.Exec(cb.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(sql, values...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1054 {
			rows, err2 := db.Query("select column_name from information_schema.columns where table_name = ? and column_name not in ('user_id', 'timestamp')", string(e.Name))
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			existingColumns := make(map[string]struct{})
			for rows.Next() {
				var col string
				err = rows.Scan(&col)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				existingColumns[col] = struct{}{}
			}
			err = rows.Err()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var ab strings.Builder
			ab.WriteString("alter table ")
			ab.WriteString(string(e.Name))
			for _, p := range e.Properties {
				_, ok := existingColumns[string(p.Name)]
				if !ok {
					ab.WriteString(" add column ")
					ab.WriteString(string(p.Name))
					ab.WriteRune(' ')
					ab.WriteString(p.DBType())
				}
			}
			_, err = db.Exec(ab.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(sql, values...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
	for _, e := range users.Users {
		var sb strings.Builder
		sb.WriteString("insert into user")
		sb.WriteString("(id")
		for _, p := range e.Properties {
			sb.WriteRune(',')
			sb.WriteString(string(p.Name))
		}
		sb.WriteString(") values (?")
		for range e.Properties {
			sb.WriteString(",?")
		}
		sb.WriteRune(')')
		values := make([]interface{}, len(e.Properties)+1)
		values[0] = string(e.ID)
		for i := range e.Properties {
			values[i+1] = e.Properties[i].Value
		}
		sql := sb.String()
		_, err = db.Exec(sql, values...)
		if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1146 {
			var cb strings.Builder
			cb.WriteString("create table user")
			cb.WriteString("(id varchar(255) not null")
			for _, p := range e.Properties {
				cb.WriteRune(',')
				cb.WriteString(string(p.Name))
				cb.WriteRune(' ')
				cb.WriteString(p.DBType())
			}
			cb.WriteRune(')')
			_, err = db.Exec(cb.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(sql, values...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1054 {
			rows, err2 := db.Query("select column_name from information_schema.columns where table_name = user and column_name not in ('id')")
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			existingColumns := make(map[string]struct{})
			for rows.Next() {
				var col string
				err = rows.Scan(&col)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				existingColumns[col] = struct{}{}
			}
			err = rows.Err()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var ab strings.Builder
			ab.WriteString("alter table user ")
			for _, p := range e.Properties {
				_, ok := existingColumns[string(p.Name)]
				if !ok {
					ab.WriteString(" add column ")
					ab.WriteString(string(p.Name))
					ab.WriteRune(' ')
					ab.WriteString(p.DBType())
				}
			}
			_, err = db.Exec(ab.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(sql, values...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
