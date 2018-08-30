package main

import "strings"

func insertEvent(e Event) (string, []interface{}) {
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
	return sb.String(), values
}

func createEventTable(e Event) string {
	var sb strings.Builder
	sb.WriteString("create table ")
	sb.WriteString(string(e.Name))
	sb.WriteString("(user_id varchar(255) not null, timestamp datetime not null")
	for _, p := range e.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
		sb.WriteRune(' ')
		sb.WriteString(p.DBType())
	}
	sb.WriteRune(')')
	return sb.String()
}

func alterEventTable(e Event, existingColumns map[string]struct{}) string {
	var sb strings.Builder
	sb.WriteString("alter table ")
	sb.WriteString(string(e.Name))
	for _, p := range e.Properties {
		_, ok := existingColumns[string(p.Name)]
		if !ok {
			sb.WriteString(" add column ")
			sb.WriteString(string(p.Name))
			sb.WriteRune(' ')
			sb.WriteString(p.DBType())
		}
	}
	return sb.String()
}

func existingEventColumns(table string) (map[string]struct{}, error) {
	rows, err := db.Query("select column_name from information_schema.columns where table_name = ? and column_name not in ('user_id', 'timestamp')", table)
	if err != nil {
		return nil, err
	}
	existingColumns := make(map[string]struct{})
	for rows.Next() {
		var col string
		err = rows.Scan(&col)
		if err != nil {
			return nil, err
		}
		existingColumns[col] = struct{}{}
	}
	return existingColumns, rows.Err()
}

func insertUser(u User) (string, []interface{}) {
	var sb strings.Builder
	sb.WriteString("insert into user")
	sb.WriteString("(id")
	for _, p := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
	}
	sb.WriteString(") values (?")
	for range u.Properties {
		sb.WriteString(",?")
	}
	sb.WriteRune(')')
	values := make([]interface{}, len(u.Properties)+1)
	values[0] = string(u.ID)
	for i := range u.Properties {
		values[i+1] = u.Properties[i].Value
	}
	return sb.String(), values
}

func createUserTable(u User) string {
	var sb strings.Builder
	sb.WriteString("create table user")
	sb.WriteString("(id varchar(255) not null")
	for _, p := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
		sb.WriteRune(' ')
		sb.WriteString(p.DBType())
	}
	sb.WriteRune(')')
	return sb.String()
}

func alterUserTable(u User, existingColumns map[string]struct{}) string {
	var sb strings.Builder
	sb.WriteString("alter table user ")
	for _, p := range u.Properties {
		_, ok := existingColumns[string(p.Name)]
		if !ok {
			sb.WriteString(" add column ")
			sb.WriteString(string(p.Name))
			sb.WriteRune(' ')
			sb.WriteString(p.DBType())
		}
	}
	return sb.String()
}

func existingUserColumns() (map[string]struct{}, error) {
	rows, err := db.Query("select column_name from information_schema.columns where table_name = user and column_name != id")
	if err != nil {
		return nil, err
	}
	existingColumns := make(map[string]struct{})
	for rows.Next() {
		var col string
		err = rows.Scan(&col)
		if err != nil {
			return nil, err
		}
		existingColumns[col] = struct{}{}
	}
	return existingColumns, rows.Err()
}
