package pas

import (
	"database/sql"
	"strings"
	"time"
)

type userInserter struct {
	u User
}

func (i userInserter) InsertSQL(t time.Time) (string, []interface{}) {
	u := i.u
	var sb strings.Builder
	sb.WriteString("insert into user")
	sb.WriteString("(id, timestamp")
	for _, p := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
	}
	sb.WriteString(") values (?, ?")
	for range u.Properties {
		sb.WriteString(",?")
	}
	sb.WriteString(") on duplicate key update timestamp = ?")
	for _, p := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
		sb.WriteString("=?")
	}
	values := make([]interface{}, 2*len(u.Properties)+3)
	values[0] = string(u.ID)
	values[1] = t
	for i := range u.Properties {
		values[i+2] = u.Properties[i].Value
	}
	values[len(u.Properties)+2] = t
	for i := range u.Properties {
		values[i+len(u.Properties)+3] = u.Properties[i].Value
	}
	return sb.String(), values
}

func (i userInserter) CreateTableSQL() string {
	u := i.u
	var sb strings.Builder
	sb.WriteString("create table user")
	sb.WriteString("(id varchar(255) not null, timestamp datetime not null")
	for _, p := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(p.Name))
		sb.WriteRune(' ')
		sb.WriteString(p.dbType())
	}
	sb.WriteString(", primary key (id), index idx_userid (id), index idx_timestamp(timestamp))")
	return sb.String()
}

func (i userInserter) ExistingColumns(db *sql.DB) (map[string]struct{}, error) {
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

func (i userInserter) AlterTableSQL(existingColumns map[string]struct{}) string {
	u := i.u
	var sb strings.Builder
	sb.WriteString("alter table user ")
	for _, p := range u.Properties {
		_, ok := existingColumns[string(p.Name)]
		if !ok {
			sb.WriteString(" add column ")
			sb.WriteString(string(p.Name))
			sb.WriteRune(' ')
			sb.WriteString(p.dbType())
		}
	}
	return sb.String()
}
